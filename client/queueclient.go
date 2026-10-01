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
// A QueueClient opened by Client.OpenQueue outlives its stream: the caller
// never opens it again. When the stream breaks, or its broker stops leading
// the cluster, the broker frees the items pulled on it and not yet committed,
// requeued or deleted. An operation that meets such a loss while the stream
// held no item is repeated on a new stream, reopened through the client's
// connection, with a backoff for up to the reopen window (see WithReopenWindow): nothing the caller holds
// is affected. One that meets it while the stream held items fails with
// ErrLeasesLost and is not repeated: the caller must treat those items as
// free, and whether the operation took effect is unknown. Enqueue and
// BatchEnqueue requests without a request id get one, set on the request, so
// a repeated enqueue is applied once by a broker at feature level 2; below
// it, a repeat may enqueue a duplicate.
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
	lost      *lostHandler
	window    time.Duration // the reopen window, see WithReopenWindow
}

// queueStream is one stream of a QueueClient and the calls waiting on it; its
// err, corrmap, held, dropped and leased are guarded by the QueueClient's
// lock.
type queueStream struct {
	s       grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse]
	cancel  context.CancelFunc
	errch   chan error
	err     error
	corrmap map[uint64]chan *pb.QueueResponse
	// held is the items pulled on the stream and not yet committed,
	// requeued or deleted.
	held map[uint64]struct{}
	// dropped is set once the stream was lost, and leased then reports
	// whether it held items; reported is set once an operation failed with
	// ErrLeasesLost for it.
	dropped, leased, reported bool
}

func newQueueStream(s grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse], cancel context.CancelFunc) *queueStream {
	return &queueStream{
		s:       s,
		cancel:  cancel,
		errch:   make(chan error, 1),
		corrmap: make(map[uint64]chan *pb.QueueResponse),
		held:    make(map[uint64]struct{}),
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
		lost:      c.lost,
		window:    reopenWindowOf(opts),
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

// handleresp runs cmd and returns its response. A command whose stream was
// lost, or which the broker refused while leadership moved, is repeated on a
// new stream while the lost stream held no item, with a backoff for up to
// the reopen window (see WithReopenWindow); otherwise it fails with ErrLeasesLost.
func (c *QueueClient) handleresp(cmd *pb.QueueRequest) (*pb.QueueResponse, error) {
	var wait time.Duration
	var window context.Context
	for {
		if err := c.ensure(); err != nil {
			return nil, err
		}
		resp, qs, err := c.roundtrip(cmd)
		if err == nil {
			st, ok := resp.Response.(*pb.QueueResponse_Status)
			if !ok || st.Status.Code != pb.StatusCode_STATUS_CODE_LEADER_SWITCH || c.svcClient == nil {
				c.track(qs, cmd, resp)
				return resp, nil
			}
			err = &leaderSwitchError{msg: st.Status.Message}
		}
		if c.svcClient == nil || qs == nil || !retryable(err) {
			return nil, err
		}
		if c.drop(qs, err) {
			c.markreported(qs)
			return nil, &leasesLostError{cause: err}
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

// track records on qs the items a pull handed out and forgets those an
// answered commit, requeue or delete settled.
func (c *QueueClient) track(qs *queueStream, cmd *pb.QueueRequest, resp *pb.QueueResponse) {
	c.lock.Lock()
	defer c.lock.Unlock()
	switch r := resp.Response.(type) {
	case *pb.QueueResponse_Pull:
		for _, it := range r.Pull.GetItems() {
			qs.held[it.GetId()] = struct{}{}
		}
		return
	case *pb.QueueResponse_PullSingle:
		if it := r.PullSingle.GetItem(); it != nil {
			qs.held[it.GetId()] = struct{}{}
		}
		return
	case *pb.QueueResponse_Status:
		return // the command failed: nothing was settled
	}
	switch k := cmd.Command.(type) {
	case *pb.QueueRequest_CommitSingle:
		delete(qs.held, k.CommitSingle.GetId())
	case *pb.QueueRequest_Commit:
		for _, id := range k.Commit.GetIds() {
			delete(qs.held, id)
		}
	case *pb.QueueRequest_RequeueSingle:
		delete(qs.held, k.RequeueSingle.GetItem().GetId())
	case *pb.QueueRequest_Requeue:
		for _, it := range k.Requeue.GetItems() {
			delete(qs.held, it.GetId())
		}
	case *pb.QueueRequest_DeleteSingle:
		delete(qs.held, k.DeleteSingle.GetId())
	case *pb.QueueRequest_Delete:
		for _, id := range k.Delete.GetIds() {
			delete(qs.held, id)
		}
	}
}

// ensure reopens the stream once it was dropped or broke, when the client
// can, retrying for a while (see reopening); the operation returns its error.
// A stream that broke while no operation was waiting on it still held items
// the broker freed: the first operation to find it fails with ErrLeasesLost.
func (c *QueueClient) ensure() error {
	c.lock.Lock()
	closed, qs := c.closed, c.cur
	dead := qs != nil && qs.err != nil
	var cause error
	if dead {
		cause = qs.err
	}
	c.lock.Unlock()
	if closed {
		return ErrQueueClientClosed
	}
	if !dead || c.svcClient == nil {
		return nil
	}
	if c.drop(qs, cause) && c.markreported(qs) {
		return &leasesLostError{cause: cause}
	}
	return reopening(c.ctx, c.window, func() error { return c.reconnect(qs) })
}

// markreported records that an operation reported the lost leases of qs, and
// reports whether it is the first.
func (c *QueueClient) markreported(qs *queueStream) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	first := !qs.reported
	qs.reported = true
	return first
}

// drop ends qs after err, so that the next operation reopens the stream; the
// broker frees what the stream held once it closes. It reports whether the
// stream held items, and reports the loss, once per stream, to the client's
// handler.
func (c *QueueClient) drop(qs *queueStream, err error) bool {
	c.lock.Lock()
	first := !qs.dropped
	if first {
		qs.dropped = true
		qs.leased = len(qs.held) > 0
		if qs.err == nil {
			qs.err = err
		}
	}
	leased := qs.leased
	c.lock.Unlock()
	qs.cancel()
	if first {
		c.logger.Warnf("queue %s: %v; holding items: %v; the next operation reconnects", c.qname, err, leased)
		c.lost.call(ConnectionLost{Queue: c.qname, Err: err, LeasesLost: leased})
	}
	return leased
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
// can reopen its stream, so that the caller may repeat a failed enqueue and
// have it applied once.
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
