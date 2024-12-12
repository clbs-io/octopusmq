package client

import (
	"context"
	"fmt"
	"io"
	"sync"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcstoragepb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// make that shit thread safe
type StorageClient struct {
	stoClient grpcstoragepb.StorageServiceClient
	stream    grpc.BidiStreamingClient[pb.StorageRequest, pb.StorageResponse]
	corrid    uint64
	opts      []grpc.CallOption // like really?
	ctx       context.Context
	stoname   string
	closed    bool // till no close called, try to open at any cost
	logger    *zap.SugaredLogger
	lock      sync.Mutex
	errch     chan error
	corrmap   map[uint64]chan *pb.StorageResponse
}

func (c *Client) OpenStorage(ctx context.Context, name string, opts ...grpc.CallOption) (*StorageClient, error) {
	// opening on demand, on every attempt
	ret := &StorageClient{
		stoClient: c.stoConnect,
		closed:    false,
		stream:    nil,
		corrid:    0,
		opts:      opts,
		ctx:       ctx,
		stoname:   name,
		logger:    c.logger,
		errch:     make(chan error, 1),
		corrmap:   make(map[uint64]chan *pb.StorageResponse),
	}
	err := ret.open() // no locks yet
	if err != nil {
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
			c.errch <- err
			return
		}
		// put r according to correlation id to proper channel, this is damn slow, maybe integrate that into caller application instead of this mess?
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

func (c *StorageClient) open() (err error) {
	defer func() {
		if err != nil && c.stream != nil {
			_ = c.stream.CloseSend()
			c.stream = nil
		}
	}()

	c.stream, err = c.stoClient.StorageConnect(c.ctx, c.opts...)
	if err != nil {
		return // handle errors?
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
			return fmt.Errorf("setup failed: %v", st.Status)
		}
	} else {
		return fmt.Errorf("setup failed, invalid response type: %T", resp)
	}
	return
}

func (c *StorageClient) handlesend(cmd *pb.StorageRequest) (ch chan *pb.StorageResponse, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil, ErrStorageClientClosed
	}

	c.corrid++
	if _, ok := c.corrmap[c.corrid]; ok { // damn shouldnt happen, fucking so many map check, that shit couldnt be fast at all
		return nil, fmt.Errorf("correlation id collision: %d", c.corrid)
	}
	cmd.CorrelationId = c.corrid
	err = c.stream.Send(cmd)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
			case codes.Canceled:
				c.logger.Errorf("stream send error: %v, reduced to io.EOF", err)
				return nil, io.EOF
			}
		}
	} else {
		ch = make(chan *pb.StorageResponse)
		c.corrmap[c.corrid] = ch
	}
	return
}

func (c *StorageClient) handleresp(cmd *pb.StorageRequest) (*pb.StorageResponse, error) {
	ch, err := c.handlesend(cmd)
	if err != nil {
		return nil, err
	}
	select {
	case reqp, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("forcibly closed")
		}
		return reqp, nil
	case err, ok := <-c.errch:
		if !ok {
			return nil, fmt.Errorf("already in error") // already error
		}
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.DeadlineExceeded:
			case codes.Canceled:
				c.logger.Errorf("stream recv error: %v, reduced to io.EOF", err)
				return nil, io.EOF
			}
		}
		return nil, err
	}
}

func decodestoragestatus(cc *pb.StorageResponse_Status) error {
	switch cc.Status.Code {
	case pb.StatusCode_STATUS_CODE_OK:
		return fmt.Errorf("unexpectedly ok")
	case pb.StatusCode_STATUS_CODE_TIMEOUT:
		return ErrStorageTimeout
	}
	return fmt.Errorf("command error: %s, status: %d", cc.Status.Message, cc.Status.Code)
}

func (c *StorageClient) Close() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		return ErrStorageClientClosed
	}
	c.closed = true
	if c.stream != nil {
		err = c.stream.CloseSend() // closed with error... consider it just closed...
		c.stream = nil             // just screw it
	}

	for _, ch := range c.corrmap { // just close it, damn i have to figure out a bit different way how to handle this... i just dont like it
		close(ch)
	}
	return
}

func (c *StorageClient) Get(req *pb.StorageGetRequest) (*pb.StorageDataResponse, error) {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Get{
			Get: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		return nil, decodestoragestatus(cc)
	case *pb.StorageResponse_DataResponse:
		return cc.DataResponse, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", cc)
	}
}

func (c *StorageClient) Set(req *pb.StorageSetRequest) error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Set{
			Set: req,
		},
	})
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
		return fmt.Errorf("unexpected response type: %T", cc)
	}
}

func (c *StorageClient) Delete(req *pb.StorageDeleteRequest) error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Delete{
			Delete: req,
		},
	})
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
		return fmt.Errorf("unexpected response type: %T", cc)
	}
}

func (c *StorageClient) LockAny(req *pb.StorageLockAnyWithIdRequest) (*pb.StorageDataResponse, error) {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_LockAnyWithId{
			LockAnyWithId: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.StorageResponse_Status:
		return nil, decodestoragestatus(cc)
	case *pb.StorageResponse_DataResponse:
		return cc.DataResponse, nil
	default:
		return nil, fmt.Errorf("unexpected response type: %T", cc)
	}
}

func (c *StorageClient) ReleaseId(req *pb.StorageReleaseIdRequest) error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_ReleaseId{
			ReleaseId: req,
		},
	})
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
		return fmt.Errorf("unexpected response type: %T", cc)
	}
}

func (c *StorageClient) Noop() error {
	reqp, err := c.handleresp(&pb.StorageRequest{
		Command: &pb.StorageRequest_Noop{
			Noop: &pb.NoopRequest{},
		},
	})
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
		return fmt.Errorf("unexpected response type: %T", cc)
	}
}
