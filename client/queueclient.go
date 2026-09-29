package client

import (
	"context"
	"sync"
	"time"
	"uuid"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NewRequestID returns a new idempotency key for use as EnqueueRequest.RequestId
// or BatchEnqueueRequest.RequestId. Reuse the same id when retrying a request
// after an error, including a stream error, so the server can recognize the
// retry within its dedup window instead of enqueuing the item again.
func NewRequestID() []byte {
	id := uuid.NewV7()
	return id[:]
}

// QueueClient is a thread-safe client for queue operations over a bidirectional
// stream.
//
// A QueueClient opened by Client.OpenQueue survives the loss of its stream and
// changes of the broker's raft leader: an operation whose stream broke, or that
// the broker refused while leadership moved, waits, reopens the stream through
// the client's connection and is repeated, until the context the QueueClient
// was opened with ends. Enqueue and BatchEnqueue requests without a request id
// get one first, so a repeated enqueue is applied once by a broker at feature
// level 2. Items pulled on a stream that broke are redelivered once their lease
// expires; from feature level 2 they can still be committed, requeued or
// deleted by id on the new stream.
type QueueClient struct {
	svcClient grpcpb.QueuesServiceClient // nil: the client cannot reopen its stream
	opts      []grpc.CallOption
	ctx       context.Context
	cancel    context.CancelFunc
	qname     string
	logger    *zap.SugaredLogger
	reopen    sync.Mutex // held by the one caller reopening the stream
	lock      sync.Mutex
	closed    bool
	corrid    uint64
	cur       *queueStream
}

// queueStream is one stream of a QueueClient and the calls waiting on it; its
// err and corrmap are guarded by the QueueClient's lock.
type queueStream struct {
	s       grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse]
	cancel  context.CancelFunc
	errch   chan error
	err     error
	corrmap map[uint64]chan *pb.QueueResponse
}

func newQueueStream(s grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse], cancel context.CancelFunc) *queueStream {
	return &queueStream{
		s:       s,
		cancel:  cancel,
		errch:   make(chan error, 1),
		corrmap: make(map[uint64]chan *pb.QueueResponse),
	}
}

// OpenQueue opens a bidirectional stream bound to the named queue.
// The stream is derived from ctx; call Close to release it.
func (c *Client) OpenQueue(ctx context.Context, name string, opts ...grpc.CallOption) (*QueueClient, error) {
	clientCtx, cancel := context.WithCancel(ctx)
	ret := &QueueClient{
		svcClient: c.queueConnect,
		opts:      opts,
		ctx:       clientCtx,
		cancel:    cancel,
		qname:     name,
		logger:    c.logger,
	}
	qs, err := ret.open()
	if err != nil {
		cancel()
		return nil, err
	}
	ret.cur = qs
	go ret.receive(qs)
	return ret, nil
}

// receive dispatches the responses of qs to their callers until the stream
// ends, then reports why on qs.errch.
func (c *QueueClient) receive(qs *queueStream) {
	defer close(qs.errch)
	for {
		r, err := qs.s.Recv()
		if err != nil {
			broken := &brokenStreamError{cause: err}
			c.fail(qs, broken)
			qs.errch <- broken
			return
		}
		// Dispatch response to the waiting caller by correlation ID.
		c.lock.Lock()
		if ch, ok := qs.corrmap[r.CorrelationId]; ok {
			delete(qs.corrmap, r.CorrelationId)
			c.lock.Unlock()
			ch <- r
			close(ch)
		} else {
			c.lock.Unlock()
			c.logger.Errorf("unexpected correlation id: %d", r.CorrelationId)
		}
	}
}

// fail records the error that terminated qs and abandons every caller waiting
// on it. The callers are released by the closing of qs.errch, which reports
// the cause, so the abandoned channels are dropped rather than closed.
func (c *QueueClient) fail(qs *queueStream, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	qs.err = err
	clear(qs.corrmap)
}

// terminalerr reports the error that terminated qs.
func (c *QueueClient) terminalerr(qs *queueStream) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if qs.err == nil {
		return status.Error(codes.Internal, "stream closed")
	}
	return qs.err
}

// open opens a stream and binds it to the queue. A stream that fails before
// the broker answered the setup is reported as broken, so that it is retried.
func (c *QueueClient) open() (*queueStream, error) {
	ctx, cancel := context.WithCancel(c.ctx)
	s, err := c.svcClient.Connect(ctx, c.opts...)
	if err != nil {
		cancel()
		return nil, &brokenStreamError{cause: err}
	}
	err = s.Send(&pb.QueueRequest{
		CorrelationId: 0,
		Command: &pb.QueueRequest_Setup{
			Setup: &pb.SetupRequest{
				QueueName: c.qname,
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
	st, ok := resp.Response.(*pb.QueueResponse_Status)
	if !ok {
		cancel()
		return nil, status.Errorf(codes.Internal, "setup failed, invalid response type: %T", resp)
	}
	if st.Status.Code != pb.StatusCode_STATUS_CODE_OK {
		cancel()
		return nil, decodestatus(st)
	}
	return newQueueStream(s, cancel), nil
}

// reconnect replaces old, the stream an operation failed on, with a new one,
// unless another caller already did. The stream is opened outside the lock, so
// that Close can cancel an attempt the broker does not answer.
func (c *QueueClient) reconnect(old *queueStream) error {
	c.reopen.Lock()
	defer c.reopen.Unlock()

	c.lock.Lock()
	closed, replaced := c.closed, c.cur != old
	c.lock.Unlock()
	if closed {
		return ErrQueueClientClosed
	}
	if replaced {
		return nil
	}

	qs, err := c.open()
	if err != nil {
		return err
	}
	c.lock.Lock()
	if c.closed {
		c.lock.Unlock()
		qs.cancel()
		return ErrQueueClientClosed
	}
	c.cur = qs
	c.lock.Unlock()
	old.cancel()
	go c.receive(qs)
	return nil
}

func (c *QueueClient) handlesend(cmd *pb.QueueRequest) (*queueStream, chan *pb.QueueResponse, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil, nil, ErrQueueClientClosed
	}
	qs := c.cur
	if qs.err != nil {
		return qs, nil, qs.err
	}

	c.corrid++
	if _, ok := qs.corrmap[c.corrid]; ok {
		return qs, nil, status.Errorf(codes.Internal, "correlation id collision: %d", c.corrid)
	}
	cmd.CorrelationId = c.corrid
	if err := qs.s.Send(cmd); err != nil {
		return qs, nil, &brokenStreamError{cause: err}
	}
	ch := make(chan *pb.QueueResponse)
	qs.corrmap[c.corrid] = ch
	return qs, ch, nil
}

// roundtrip sends cmd on the current stream and waits for its response.
func (c *QueueClient) roundtrip(cmd *pb.QueueRequest) (*pb.QueueResponse, *queueStream, error) {
	qs, ch, err := c.handlesend(cmd)
	if err != nil {
		return nil, qs, err
	}
	select {
	case reqp, ok := <-ch:
		if !ok {
			return nil, qs, status.Error(codes.Canceled, "forcibly closed")
		}
		return reqp, qs, nil
	case err, ok := <-qs.errch:
		if !ok {
			return nil, qs, c.terminalerr(qs) // the failure was already reported to another caller
		}
		return nil, qs, err
	}
}

// handleresp runs cmd and returns its response. A command whose stream broke,
// or that the broker refused while leadership moved, is repeated on a new
// stream until the client's context ends.
func (c *QueueClient) handleresp(cmd *pb.QueueRequest) (*pb.QueueResponse, error) {
	var wait time.Duration
	for {
		resp, qs, err := c.roundtrip(cmd)
		if err == nil {
			st, ok := resp.Response.(*pb.QueueResponse_Status)
			if !ok || st.Status.Code != pb.StatusCode_STATUS_CODE_LEADER_SWITCH || c.svcClient == nil {
				return resp, nil
			}
			err = &leaderSwitchError{msg: st.Status.Message}
		}
		if c.svcClient == nil || qs == nil || !retryable(err) {
			return nil, err
		}
		c.logger.Warnf("queue %s: %v; reconnecting", c.qname, err)
		for {
			if !retrywait(c.ctx, &wait) {
				return nil, err
			}
			rerr := c.reconnect(qs)
			if rerr == nil {
				break
			}
			if !retryable(rerr) {
				return nil, rerr
			}
			err = rerr
		}
	}
}

func decodestatus(cc *pb.QueueResponse_Status) error {
	switch cc.Status.Code {
	case pb.StatusCode_STATUS_CODE_TIMEOUT:
		return ErrQueueTimeout
	case pb.StatusCode_STATUS_CODE_QUEUE_NOT_FOUND:
		return ErrQueueNotFound
	case pb.StatusCode_STATUS_CODE_QUEUE_PAUSED:
		return ErrQueuePaused
	case pb.StatusCode_STATUS_CODE_LEADER_SWITCH:
		return &leaderSwitchError{msg: cc.Status.Message}
	}
	return status.Errorf(codes.Internal, "command error: %s, status: %d", cc.Status.Message, cc.Status.Code)
}

func (c *QueueClient) Close() error {
	c.lock.Lock()
	if c.closed {
		c.lock.Unlock()
		return ErrQueueClientClosed
	}
	c.closed = true
	qs := c.cur
	err := qs.s.CloseSend()
	// Release callers still waiting for a response. The entries must leave the map
	// so that a late response cannot make the receiver send on a closed channel.
	for id, ch := range qs.corrmap {
		delete(qs.corrmap, id)
		close(ch)
	}
	c.lock.Unlock()
	c.cancel()
	return err
}

// withRequestID gives an enqueue without a request id one, when the client
// repeats failed operations, so that a repeated enqueue is applied once.
func (c *QueueClient) withRequestID(id []byte) []byte {
	if len(id) == 0 && c.svcClient != nil {
		return NewRequestID()
	}
	return id
}

// Enqueue enqueues one item. A request without a request id gets one, set on
// req.
func (c *QueueClient) Enqueue(req *pb.EnqueueRequest) (*pb.EnqueueResponse, error) {
	req.RequestId = c.withRequestID(req.RequestId)
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_Enqueue{
			Enqueue: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_Enqueue:
		return cc.Enqueue, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

// BatchEnqueue enqueues items in one commit. A request without a request id
// gets one, set on req.
func (c *QueueClient) BatchEnqueue(req *pb.BatchEnqueueRequest) (*pb.BatchEnqueueResponse, error) {
	req.RequestId = c.withRequestID(req.RequestId)
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_BatchEnqueue{
			BatchEnqueue: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_BatchEnqueue:
		return cc.BatchEnqueue, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) CommitSingle(req *pb.CommitSingleRequest) (*pb.CommitSingleResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_CommitSingle{
			CommitSingle: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_CommitSingle:
		return cc.CommitSingle, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) Commit(req *pb.CommitRequest) (*pb.CommitResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_Commit{
			Commit: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_Commit:
		return cc.Commit, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) RequeueSingle(req *pb.RequeueSingleRequest) (*pb.RequeueSingleResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_RequeueSingle{
			RequeueSingle: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_RequeueSingle:
		return cc.RequeueSingle, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) Requeue(req *pb.RequeueRequest) (*pb.RequeueResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_Requeue{
			Requeue: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_Requeue:
		return cc.Requeue, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) DeleteSingle(req *pb.DeleteSingleRequest) (*pb.DeleteSingleResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_DeleteSingle{
			DeleteSingle: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_DeleteSingle:
		return cc.DeleteSingle, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) Delete(req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_Delete{
			Delete: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_Delete:
		return cc.Delete, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) Pull(req *pb.PullRequest) (*pb.PullResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_Pull{
			Pull: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_Pull:
		return cc.Pull, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) PullSingle(req *pb.PullSingleRequest) (*pb.PullSingleResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_PullSingle{
			PullSingle: req,
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_PullSingle:
		return cc.PullSingle, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) Noop() error {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_Noop{
			Noop: &pb.NoopRequest{},
		},
	})
	if err != nil {
		return err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		if cc.Status.Code == pb.StatusCode_STATUS_CODE_OK {
			return nil
		}
		return decodestatus(cc)
	default:
		return status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) GetInfo() (*pb.GetQueueInfoResponse, error) {
	reqp, err := c.handleresp(&pb.QueueRequest{
		Command: &pb.QueueRequest_GetQueueInfo{
			GetQueueInfo: &pb.GetQueueInfoRequest{},
		},
	})
	if err != nil {
		return nil, err
	}
	switch cc := reqp.Response.(type) {
	case *pb.QueueResponse_Status:
		return nil, decodestatus(cc)
	case *pb.QueueResponse_GetQueueInfo:
		return cc.GetQueueInfo, nil
	default:
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}
