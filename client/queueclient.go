package client

import (
	"context"
	"io"
	"sync"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcpb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QueueClient is a thread-safe client for queue operations over a bidirectional stream.
type QueueClient struct {
	svcClient grpcpb.QueuesServiceClient
	stream    grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse]
	corrid    uint64
	opts      []grpc.CallOption
	ctx       context.Context
	qname     string
	closed    bool
	logger    *zap.SugaredLogger
	lock      sync.Mutex
	errch     chan error
	corrmap   map[uint64]chan *pb.QueueResponse
}

func (c *Client) OpenQueue(ctx context.Context, name string, opts ...grpc.CallOption) (*QueueClient, error) {
	// opening on demand, on every attempt
	ret := &QueueClient{
		svcClient: c.queueConnect,
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

func (c *QueueClient) open() (err error) {
	c.stream, err = c.svcClient.Connect(c.ctx, c.opts...)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = c.stream.CloseSend()
		}
	}()

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
			return decodestatus(st)
		}
	} else {
		return status.Errorf(codes.Internal, "setup failed, invalid response type: %T", resp)
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
	if _, ok := c.corrmap[c.corrid]; ok {
		return nil, status.Errorf(codes.Internal, "correlation id collision: %d", c.corrid)
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
	select {
	case reqp, ok := <-ch:
		if !ok {
			return nil, status.Error(codes.Canceled, "forcibly closed")
		}
		return reqp, nil
	case err, ok := <-c.errch:
		if !ok {
			return nil, status.Error(codes.Internal, "already in error") // already error
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
	case pb.StatusCode_STATUS_CODE_TIMEOUT:
		return ErrQueueTimeout
	case pb.StatusCode_STATUS_CODE_QUEUE_NOT_FOUND:
		return ErrQueueNotFound
	}
	return status.Errorf(codes.Internal, "command error: %s, status: %d", cc.Status.Message, cc.Status.Code)
}

func (c *QueueClient) Close() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.closed {
		return ErrQueueClientClosed
	}
	c.closed = true
	err = c.stream.CloseSend()

	for _, ch := range c.corrmap {
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
		return nil, status.Errorf(codes.Internal, "unexpected response type: %T", cc)
	}
}

func (c *QueueClient) BatchEnqueue(req *pb.BatchEnqueueRequest) (*pb.BatchEnqueueResponse, error) {
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
