package client

import (
	"context"
	"sync"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcstoragepb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StorageClient is a thread-safe client for key-value storage operations over a bidirectional stream.
type StorageClient struct {
	stoClient grpcstoragepb.StorageServiceClient
	stream    grpc.BidiStreamingClient[pb.StorageRequest, pb.StorageResponse]
	corrid    uint64
	opts      []grpc.CallOption
	ctx       context.Context
	cancel    context.CancelFunc
	stoname   string
	closed    bool
	streamErr error
	logger    *zap.SugaredLogger
	lock      sync.Mutex
	errch     chan error
	corrmap   map[uint64]chan *pb.StorageResponse
}

// OpenStorage opens a bidirectional stream bound to the named storage.
// The stream is derived from ctx; call Close to release it.
func (c *Client) OpenStorage(ctx context.Context, name string, opts ...grpc.CallOption) (*StorageClient, error) {
	// opening on demand, on every attempt
	streamCtx, cancel := context.WithCancel(ctx)
	ret := &StorageClient{
		stoClient: c.stoConnect,
		closed:    false,
		stream:    nil,
		corrid:    0,
		opts:      opts,
		ctx:       streamCtx,
		cancel:    cancel,
		stoname:   name,
		logger:    c.logger,
		errch:     make(chan error, 1),
		corrmap:   make(map[uint64]chan *pb.StorageResponse),
	}
	err := ret.open() // no locks yet
	if err != nil {
		cancel()
		return nil, err
	}
	go ret.receiver()
	return ret, nil
}

func (c *StorageClient) receiver() {
	defer close(c.errch)
	for {
		r, err := c.stream.Recv()
		if err != nil {
			c.fail(err)
			c.errch <- err
			return
		}
		// Dispatch response to the waiting caller by correlation ID.
		c.lock.Lock()
		if ch, ok := c.corrmap[r.CorrelationId]; ok {
			delete(c.corrmap, r.CorrelationId)
			c.lock.Unlock()
			ch <- r
			close(ch)
		} else {
			c.lock.Unlock()
			c.logger.Errorf("unexpected correlation id: %d", r.CorrelationId)
		}
	}
}

// fail records the error that terminated the stream and abandons every pending
// caller. The callers are released by the closing of errch, which reports the
// cause, so the abandoned channels are dropped rather than closed.
func (c *StorageClient) fail(err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.streamErr = err
	clear(c.corrmap)
}

// terminalerr reports the error that terminated the stream.
func (c *StorageClient) terminalerr() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.streamErr == nil {
		return status.Error(codes.Internal, "stream closed")
	}
	return reducestreamerr(c.streamErr)
}

func (c *StorageClient) open() (err error) {
	c.stream, err = c.stoClient.StorageConnect(c.ctx, c.opts...)
	if err != nil {
		return
	}

	// manually setup stream without receiver yet
	err = c.stream.Send(&pb.StorageRequest{
		CorrelationId: 0,
		Command: &pb.StorageRequest_Setup{
			Setup: &pb.StorageSetupRequest{
				StorageName: c.stoname,
			},
		},
	})
	if err != nil {
		return
	}
	var resp *pb.StorageResponse
	resp, err = c.stream.Recv()
	if err != nil {
		return
	}
	if st, ok := resp.Response.(*pb.StorageResponse_Status); ok {
		if st.Status.Code != pb.StatusCode_STATUS_CODE_OK {
			return decodestoragestatus(st)
		}
	} else {
		return status.Errorf(codes.Internal, "setup failed, invalid response type: %T", resp)
	}
	return
}

func (c *StorageClient) handlesend(cmd *pb.StorageRequest, keepid bool) (ch chan *pb.StorageResponse, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil, ErrStorageClientClosed
	}
	if c.streamErr != nil {
		return nil, reducestreamerr(c.streamErr)
	}

	var corrid uint64
	if keepid {
		corrid = cmd.CorrelationId
	} else {
		c.corrid++
		corrid = c.corrid
		cmd.CorrelationId = c.corrid
	}

	if _, ok := c.corrmap[corrid]; ok {
		return nil, status.Errorf(codes.Internal, "correlation id collision: %d", corrid)
	}
	err = c.stream.Send(cmd)
	if err != nil {
		if reduced := reducestreamerr(err); reduced != err {
			c.logger.Errorf("stream send error: %v, reduced to io.EOF", err)
			return nil, reduced
		}
		return nil, err
	}
	ch = make(chan *pb.StorageResponse)
	c.corrmap[corrid] = ch
	return ch, nil
}

func (c *StorageClient) handleresp(cmd *pb.StorageRequest, keepid bool) (*pb.StorageResponse, error) {
	ch, err := c.handlesend(cmd, keepid)
	if err != nil {
		return nil, err
	}
	select {
	case reqp, ok := <-ch:
		if !ok {
			return nil, status.Error(codes.Canceled, "forcibly closed")
		}
		return reqp, nil
	case err, ok := <-c.errch:
		if !ok {
			return nil, c.terminalerr() // the failure was already reported to another caller
		}
		if reduced := reducestreamerr(err); reduced != err {
			c.logger.Errorf("stream recv error: %v, reduced to io.EOF", err)
			return nil, reduced
		}
		return nil, err
	}
}

func decodestoragestatus(cc *pb.StorageResponse_Status) error {
	switch cc.Status.Code {
	case pb.StatusCode_STATUS_CODE_TIMEOUT:
		return ErrStorageTimeout
	case pb.StatusCode_STATUS_CODE_ITEM_NOT_FOUND:
		return ErrStorageKeyNotFound
	case pb.StatusCode_STATUS_CODE_STORAGE_NOT_FOUND:
		return ErrStorageNotFound
	}
	return status.Errorf(codes.Internal, "command error: %s, status: %d", cc.Status.Message, cc.Status.Code)
}

func (c *StorageClient) Close() error {
	c.lock.Lock()
	if c.closed {
		c.lock.Unlock()
		return ErrStorageClientClosed
	}
	c.closed = true
	err := c.stream.CloseSend()
	// Release callers still waiting for a response. The entries must leave the map
	// so that a late response cannot make the receiver send on a closed channel.
	for id, ch := range c.corrmap {
		delete(c.corrmap, id)
		close(ch)
	}
	c.lock.Unlock()
	c.cancel()
	return err
}

func (c *StorageClient) Get(req *pb.StorageGetRequest) (*pb.StorageDataResponse, error) {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Get{
			Get: req,
		},
	}, false)
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		return nil, decodestoragestatus(cc)
	case *pb.StorageResponse_DataResponse:
		return cc.DataResponse, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *StorageClient) GetKeys(req *pb.StorageGetKeysRequest) ([][]byte, error) {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_GetKeys{
			GetKeys: req,
		},
	}, false)

	if err != nil {
		return nil, err
	}
	ret := make([][]byte, 0)
	for {
		switch cc := reqp.Response.(type) {
		case *pb.StorageResponse_Status:
			return nil, decodestoragestatus(cc)
		case *pb.StorageResponse_GetKeysResponse:
			ret = append(ret, cc.GetKeysResponse.Keys...)
			if !cc.GetKeysResponse.More {
				return ret, nil
			}
		default:
			return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
		}
		reqp, err = c.handleresp(&pb.StorageRequest{
			CorrelationId: reqp.CorrelationId,
			Command:       &pb.StorageRequest_GetKeysNext{GetKeysNext: &pb.StorageGetKeysNextRequest{}},
		}, true)
		if err != nil {
			return nil, err
		}
	}
}

func (c *StorageClient) Set(req *pb.StorageSetRequest) error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Set{
			Set: req,
		},
	}, false)
	if err != nil {
		return err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		if cc.Status.Code == pb.StatusCode_STATUS_CODE_OK {
			return nil
		}
		return decodestoragestatus(cc)
	default:
		return status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *StorageClient) Delete(req *pb.StorageDeleteRequest) error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Delete{
			Delete: req,
		},
	}, false)
	if err != nil {
		return err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		if cc.Status.Code == pb.StatusCode_STATUS_CODE_OK {
			return nil
		}
		return decodestoragestatus(cc)
	default:
		return status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *StorageClient) LockAny(req *pb.StorageLockAnyWithIdRequest) (*pb.StorageDataResponse, error) {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_LockAnyWithId{
			LockAnyWithId: req,
		},
	}, false)
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		return nil, decodestoragestatus(cc)
	case *pb.StorageResponse_DataResponse:
		return cc.DataResponse, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *StorageClient) ReleaseId(req *pb.StorageReleaseIdRequest) error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_ReleaseId{
			ReleaseId: req,
		},
	}, false)
	if err != nil {
		return err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		if cc.Status.Code == pb.StatusCode_STATUS_CODE_OK {
			return nil
		}
		return decodestoragestatus(cc)
	default:
		return status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *StorageClient) Noop() error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Noop{
			Noop: &pb.NoopRequest{},
		},
	}, false)
	if err != nil {
		return err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		if cc.Status.Code == pb.StatusCode_STATUS_CODE_OK {
			return nil
		}
		return decodestoragestatus(cc)
	default:
		return status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *StorageClient) GetInfo() (*pb.StorageGetInfoResponse, error) {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_GetInfo{
			GetInfo: &pb.StorageGetInfoRequest{},
		},
	}, false)
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		return nil, decodestoragestatus(cc)
	case *pb.StorageResponse_GetInfoResponse:
		return cc.GetInfoResponse, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}
