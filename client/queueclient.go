package client

import (
	"context"
	"fmt"
	"io"
	"sync"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// make that shit thread safe
type QueueClient struct {
	svcClient grpcpb.QueuesServiceClient
	stream    grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse]
	corrid    uint64
	opts      []grpc.CallOption // like really?
	ctx       context.Context
	qname     string
	closed    bool // till no close called, try to open at any cost
	logger    *zap.SugaredLogger
	lock      sync.Mutex
	errch     chan error
	corrmap   map[uint64]chan *pb.QueueResponse
}

func (c *Client) OpenQueue(ctx context.Context, name string, opts ...grpc.CallOption) (*QueueClient, error) {
	// opening on demand, on every attempt
	ret := &QueueClient{
		svcClient: c.svcClient,
		closed:    false,
		stream:    nil,
		corrid:    0,
		opts:      opts,
		ctx:       ctx,
		qname:     name,
		logger:    c.logger,
		errch:     make(chan error, 1),
		corrmap:   make(map[uint64]chan *pb.QueueResponse),
	}
	err := ret.open() // no locks yet
	if err != nil {
		return nil, err
	}
	go ret.receiver()
	return ret, nil
}

func (c *QueueClient) receiver() {
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
			ch <- r
		} else {
			c.logger.Errorf("unexpected correlation id: %d", r.CorrelationId)
		}
		c.lock.Unlock()
	}
}

func (c *QueueClient) open() (err error) {
	defer func() {
		if err != nil && c.stream != nil {
			_ = c.stream.CloseSend()
			c.stream = nil
		}
	}()

	c.stream, err = c.svcClient.Connect(c.ctx, c.opts...)
	if err != nil {
		return // handle errors?
	}
	// manually setup stream without receiver yet
	err = c.stream.Send(&pb.QueueRequest{
		CorrelationId: 0,
		Command: &pb.QueueRequest_Setup{
			Setup: &pb.SetupRequest{
				QueueName: c.qname,
			},
		},
	})
	if err != nil {
		return
	}
	var resp *pb.QueueResponse
	resp, err = c.stream.Recv()
	if err != nil {
		return
	}
	if st, ok := resp.Response.(*pb.QueueResponse_Status); ok {
		if st.Status.Code != pb.StatusCode_STATUS_CODE_OK {
			return fmt.Errorf("setup failed: %v", st.Status)
		}
	} else {
		return fmt.Errorf("setup failed, invalid response type: %T", resp)
	}
	return
}

func (c *QueueClient) handlesend(cmd *pb.QueueRequest) (ch chan *pb.QueueResponse, err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil, ErrQueueClientClosed
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
		ch = make(chan *pb.QueueResponse)
		c.corrmap[c.corrid] = ch
	}
	return
}

func (c *QueueClient) handleresp(cmd *pb.QueueRequest) (*pb.QueueResponse, error) {
	ch, err := c.handlesend(cmd)
	if err != nil {
		return nil, err
	}
	defer func() { // fucking hell
		c.lock.Lock()
		delete(c.corrmap, cmd.CorrelationId)
		c.lock.Unlock()
	}()
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

func decodestatus(cc *pb.QueueResponse_Status) error {
	switch cc.Status.Code {
	case pb.StatusCode_STATUS_CODE_OK:
		return fmt.Errorf("unexpectedly ok")
	case pb.StatusCode_STATUS_CODE_TIMEOUT:
		return ErrQueueTimeout
	}
	return fmt.Errorf("command error: %s, status: %d", cc.Status.Message, cc.Status.Code)
}

func (c *QueueClient) Close() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		return ErrQueueClientClosed
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

func (c *QueueClient) Enqueue(req *pb.EnqueueRequest) (*pb.EnqueueResponse, error) {
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
	}
}

// async motherfuckers
func (c *QueueClient) Pull(req *pb.PullRequest) (*pb.PullResponse, error) {
	cmd := &pb.QueueRequest{
		Command: &pb.QueueRequest_Pull{
			Pull: req,
		},
	}
	ch, err := c.handlesend(cmd)
	if err != nil {
		return nil, err
	}

	defer func() { // fucking hell
		c.lock.Lock()
		delete(c.corrmap, cmd.CorrelationId)
		c.lock.Unlock()
	}()

	select {
	case reqp, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("forcibly closed")
		}
		switch cc := reqp.Response.(type) {
		case *pb.QueueResponse_Status:
			return nil, decodestatus(cc)
		case *pb.QueueResponse_Pull:
			return cc.Pull, nil
		default:
			return nil, fmt.Errorf("unexpected response type: %T", cc)
		}
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

func (c *QueueClient) PullSingle(req *pb.PullSingleRequest) (*pb.PullSingleResponse, error) {
	cmd := &pb.QueueRequest{
		Command: &pb.QueueRequest_PullSingle{
			PullSingle: req,
		},
	}
	ch, err := c.handlesend(cmd)
	if err != nil {
		return nil, err
	}

	defer func() { // fucking hell
		c.lock.Lock()
		delete(c.corrmap, cmd.CorrelationId)
		c.lock.Unlock()
	}()

	select {
	case reqp, ok := <-ch:
		if !ok {
			return nil, fmt.Errorf("forcibly closed")
		}
		switch cc := reqp.Response.(type) {
		case *pb.QueueResponse_Status:
			return nil, decodestatus(cc)
		case *pb.QueueResponse_PullSingle:
			return cc.PullSingle, nil
		default:
			return nil, fmt.Errorf("unexpected response type: %T", cc)
		}
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
