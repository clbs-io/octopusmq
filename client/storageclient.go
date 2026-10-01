package client

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcstoragepb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StorageClient is a thread-safe client for key-value storage operations over a
// bidirectional stream.
//
// A StorageClient opened by Client.OpenStorage outlives its stream: the caller
// never opens it again. When the stream breaks, or its broker stops leading
// the cluster, the broker releases the keys locked on it. An operation that
// meets such a loss while the stream held no lock is repeated on a new
// stream, reopened through the client's connection, with a backoff for up to
// the reopen window (see WithReopenWindow); GetKeys starts its listing over. One that meets it while the
// stream held locks fails with ErrLeasesLost and is not repeated: the caller
// must treat those keys as released, and whether the operation took effect
// is unknown.
type StorageClient struct {
	stoClient grpcstoragepb.StorageServiceClient // nil: the client cannot reopen its stream
	opts      []grpc.CallOption
	ctx       context.Context
	cancel    context.CancelFunc
	stoname   string
	logger    *zap.SugaredLogger
	reopen    sync.Mutex // held by the one caller reopening the stream
	lock      sync.Mutex
	closed    bool
	corrid    uint64
	cur       *storageStream
	lost      *lostHandler
	window    time.Duration // the reopen window, see WithReopenWindow
}

// storageStream is one stream of a StorageClient and the calls waiting on it;
// its err, corrmap, held, dropped, leased and reported are guarded by the
// StorageClient's lock.
type storageStream struct {
	s       grpc.BidiStreamingClient[pb.StorageRequest, pb.StorageResponse]
	cancel  context.CancelFunc
	errch   chan error
	err     error
	corrmap map[uint64]chan *pb.StorageResponse
	// held is the ids of the keys locked on the stream and not yet released.
	held map[uint64]struct{}
	// dropped is set once the stream was lost, and leased then reports
	// whether it held locks; reported is set once an operation failed with
	// ErrLeasesLost for it.
	dropped, leased, reported bool
}

func newStorageStream(s grpc.BidiStreamingClient[pb.StorageRequest, pb.StorageResponse], cancel context.CancelFunc) *storageStream {
	return &storageStream{
		s:       s,
		cancel:  cancel,
		errch:   make(chan error, 1),
		corrmap: make(map[uint64]chan *pb.StorageResponse),
		held:    make(map[uint64]struct{}),
	}
}

// OpenStorage opens a bidirectional stream bound to the named storage.
// The stream is derived from ctx; call Close to release it.
func (c *Client) OpenStorage(ctx context.Context, name string, opts ...grpc.CallOption) (*StorageClient, error) {
	clientCtx, cancel := context.WithCancel(ctx)
	ret := &StorageClient{
		stoClient: c.stoConnect,
		opts:      opts,
		ctx:       clientCtx,
		cancel:    cancel,
		stoname:   name,
		logger:    c.logger,
		lost:      c.lost,
		window:    reopenWindowOf(opts),
	}
	ss, err := ret.open()
	if err != nil {
		cancel()
		return nil, err
	}
	ret.cur = ss
	go ret.receive(ss)
	return ret, nil
}

// receive dispatches the responses of ss to their callers until the stream
// ends, then reports why on ss.errch.
func (c *StorageClient) receive(ss *storageStream) {
	defer close(ss.errch)
	for {
		r, err := ss.s.Recv()
		if err != nil {
			broken := &brokenStreamError{cause: err}
			c.fail(ss, broken)
			ss.errch <- broken
			return
		}
		// Dispatch response to the waiting caller by correlation ID.
		c.lock.Lock()
		if ch, ok := ss.corrmap[r.CorrelationId]; ok {
			delete(ss.corrmap, r.CorrelationId)
			c.lock.Unlock()
			ch <- r
			close(ch)
		} else {
			c.lock.Unlock()
			c.logger.Errorf("unexpected correlation id: %d", r.CorrelationId)
		}
	}
}

// fail records the error that terminated ss and abandons every caller waiting
// on it. The callers are released by the closing of ss.errch, which reports
// the cause, so the abandoned channels are dropped rather than closed.
func (c *StorageClient) fail(ss *storageStream, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	ss.err = err
	clear(ss.corrmap)
}

// terminalerr reports the error that terminated ss.
func (c *StorageClient) terminalerr(ss *storageStream) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if ss.err == nil {
		return status.Error(codes.Internal, "stream closed")
	}
	return ss.err
}

// open opens a stream and binds it to the storage. A stream that fails before
// the broker answered the setup is reported as broken, so that it is retried.
func (c *StorageClient) open() (*storageStream, error) {
	ctx, cancel := context.WithCancel(c.ctx)
	s, err := c.stoClient.StorageConnect(ctx, c.opts...)
	if err != nil {
		cancel()
		return nil, &brokenStreamError{cause: err}
	}
	err = s.Send(&pb.StorageRequest{
		CorrelationId: 0,
		Command: &pb.StorageRequest_Setup{
			Setup: &pb.StorageSetupRequest{
				StorageName: c.stoname,
			},
		},
	})
	if err != nil {
		cancel()
		return nil, &brokenStreamError{cause: err}
	}
	resp, err := s.Recv()
	if err != nil {
		cancel()
		return nil, &brokenStreamError{cause: err}
	}
	st, ok := resp.Response.(*pb.StorageResponse_Status)
	if !ok {
		cancel()
		return nil, status.Errorf(codes.Internal, "setup failed, invalid response type: %T", resp)
	}
	if st.Status.Code != pb.StatusCode_STATUS_CODE_OK {
		cancel()
		return nil, decodestoragestatus(st)
	}
	return newStorageStream(s, cancel), nil
}

// reconnect replaces old, the stream an operation failed on, with a new one,
// unless another caller already did. The stream is opened outside the lock, so
// that Close can cancel an attempt the broker does not answer.
func (c *StorageClient) reconnect(old *storageStream) error {
	c.reopen.Lock()
	defer c.reopen.Unlock()

	c.lock.Lock()
	closed, replaced := c.closed, c.cur != old
	c.lock.Unlock()
	if closed {
		return ErrStorageClientClosed
	}
	if replaced {
		return nil
	}

	ss, err := c.open()
	if err != nil {
		return err
	}
	c.lock.Lock()
	if c.closed {
		c.lock.Unlock()
		ss.cancel()
		return ErrStorageClientClosed
	}
	c.cur = ss
	c.lock.Unlock()
	old.cancel()
	go c.receive(ss)
	return nil
}

func (c *StorageClient) handlesend(cmd *pb.StorageRequest, keepid bool) (*storageStream, chan *pb.StorageResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil, nil, ErrStorageClientClosed
	}
	ss := c.cur
	if ss.err != nil {
		return ss, nil, ss.err
	}

	var corrid uint64
	if keepid {
		corrid = cmd.CorrelationId
	} else {
		c.corrid++
		corrid = c.corrid
		cmd.CorrelationId = c.corrid
	}

	if _, ok := ss.corrmap[corrid]; ok {
		return ss, nil, status.Errorf(codes.Internal, "correlation id collision: %d", corrid)
	}
	if err := ss.s.Send(cmd); err != nil {
		return ss, nil, &brokenStreamError{cause: err}
	}
	ch := make(chan *pb.StorageResponse)
	ss.corrmap[corrid] = ch
	return ss, ch, nil
}

// roundtrip sends cmd on the current stream and waits for its response.
func (c *StorageClient) roundtrip(cmd *pb.StorageRequest, keepid bool) (*pb.StorageResponse, *storageStream, error) {
	ss, ch, err := c.handlesend(cmd, keepid)
	if err != nil {
		return nil, ss, err
	}
	select {
	case reqp, ok := <-ch:
		if !ok {
			return nil, ss, status.Error(codes.Canceled, "forcibly closed")
		}
		return reqp, ss, nil
	case err, ok := <-ss.errch:
		if !ok {
			return nil, ss, c.terminalerr(ss) // the failure was already reported to another caller
		}
		return nil, ss, err
	}
}

// handleresp runs cmd and returns its response. A command whose stream was
// lost, or which the broker refused while leadership moved, is repeated on a
// new stream while the lost stream held no lock, with a backoff for up to
// the reopen window (see WithReopenWindow); otherwise it fails with ErrLeasesLost. A command continuing a
// listing (keepid) belongs to its stream and is never repeated.
func (c *StorageClient) handleresp(cmd *pb.StorageRequest, keepid bool) (*pb.StorageResponse, error) {
	var wait time.Duration
	var window context.Context
	for {
		if !keepid {
			if err := c.ensure(); err != nil {
				return nil, err
			}
		}
		resp, ss, err := c.roundtrip(cmd, keepid)
		if err == nil {
			st, ok := resp.Response.(*pb.StorageResponse_Status)
			if !ok || st.Status.Code != pb.StatusCode_STATUS_CODE_LEADER_SWITCH || c.stoClient == nil {
				c.track(ss, cmd, resp)
				return resp, nil
			}
			err = &leaderSwitchError{msg: st.Status.Message}
		}
		if c.stoClient == nil || ss == nil || !retryable(err) {
			return nil, err
		}
		if c.drop(ss, err) {
			c.markreported(ss)
			return nil, &leasesLostError{cause: err}
		}
		if keepid {
			return nil, err
		}
		if window == nil {
			var cancel context.CancelFunc
			window, cancel = context.WithTimeout(c.ctx, c.window)
			defer cancel()
		}
		if !retrywait(window, &wait) {
			return nil, err
		}
	}
}

// track records on ss the key a lock took and forgets one a release ended.
func (c *StorageClient) track(ss *storageStream, cmd *pb.StorageRequest, resp *pb.StorageResponse) {
	c.lock.Lock()
	defer c.lock.Unlock()
	switch k := cmd.Command.(type) {
	case *pb.StorageRequest_LockAnyWithId:
		if d, ok := resp.Response.(*pb.StorageResponse_DataResponse); ok {
			ss.held[d.DataResponse.GetId()] = struct{}{}
		}
	case *pb.StorageRequest_ReleaseId:
		if st, ok := resp.Response.(*pb.StorageResponse_Status); ok {
			switch st.Status.Code {
			case pb.StatusCode_STATUS_CODE_OK, pb.StatusCode_STATUS_CODE_ITEM_NOT_FOUND: // released, or never this stream's
				delete(ss.held, k.ReleaseId.GetId())
			}
		}
	}
}

// ensure reopens the stream once it was dropped or broke, when the client
// can, retrying for a while (see reopening); the operation returns its error.
// A stream that broke while no operation was waiting on it still held locks
// the broker released: the first operation to find it fails with
// ErrLeasesLost.
func (c *StorageClient) ensure() error {
	c.lock.Lock()
	closed, ss := c.closed, c.cur
	dead := ss != nil && ss.err != nil
	var cause error
	if dead {
		cause = ss.err
	}
	c.lock.Unlock()
	if closed {
		return ErrStorageClientClosed
	}
	if !dead || c.stoClient == nil {
		return nil
	}
	if c.drop(ss, cause) && c.markreported(ss) {
		return &leasesLostError{cause: cause}
	}
	return reopening(c.ctx, c.window, func() error { return c.reconnect(ss) })
}

// markreported records that an operation reported the lost locks of ss, and
// reports whether it is the first.
func (c *StorageClient) markreported(ss *storageStream) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	first := !ss.reported
	ss.reported = true
	return first
}

// drop ends ss after err, so that the next operation reopens the stream; the
// broker releases the stream's locks once it closes. It reports whether the
// stream held locks, and reports the loss, once per stream, to the client's
// handler.
func (c *StorageClient) drop(ss *storageStream, err error) bool {
	c.lock.Lock()
	first := !ss.dropped
	if first {
		ss.dropped = true
		ss.leased = len(ss.held) > 0
		if ss.err == nil {
			ss.err = err
		}
	}
	leased := ss.leased
	c.lock.Unlock()
	ss.cancel()
	if first {
		c.logger.Warnf("storage %s: %v; holding locks: %v; the next operation reconnects", c.stoname, err, leased)
		c.lost.call(ConnectionLost{Storage: c.stoname, Err: err, LeasesLost: leased})
	}
	return leased
}

func decodestoragestatus(cc *pb.StorageResponse_Status) error {
	switch cc.Status.Code {
	case pb.StatusCode_STATUS_CODE_TIMEOUT:
		return ErrStorageTimeout
	case pb.StatusCode_STATUS_CODE_ITEM_NOT_FOUND:
		return ErrStorageKeyNotFound
	case pb.StatusCode_STATUS_CODE_STORAGE_NOT_FOUND:
		return ErrStorageNotFound
	case pb.StatusCode_STATUS_CODE_LEADER_SWITCH:
		return &leaderSwitchError{msg: cc.Status.Message}
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
	ss := c.cur
	err := ss.s.CloseSend()
	// Release callers still waiting for a response. The entries must leave the map
	// so that a late response cannot make the receiver send on a closed channel.
	for id, ch := range ss.corrmap {
		delete(ss.corrmap, id)
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

// GetKeys lists the keys matching req. A listing whose stream was lost
// midway starts over on a new stream, like any other operation (see
// StorageClient).
func (c *StorageClient) GetKeys(req *pb.StorageGetKeysRequest) ([][]byte, error) {
	var wait time.Duration
	var window context.Context
	for {
		keys, err := c.getkeys(req)
		if err == nil || c.stoClient == nil || !retryable(err) || errors.Is(err, ErrLeasesLost) {
			return keys, err
		}
		if window == nil {
			var cancel context.CancelFunc
			window, cancel = context.WithTimeout(c.ctx, c.window)
			defer cancel()
		}
		if !retrywait(window, &wait) {
			return nil, err
		}
	}
}

func (c *StorageClient) getkeys(req *pb.StorageGetKeysRequest) ([][]byte, error) {
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
