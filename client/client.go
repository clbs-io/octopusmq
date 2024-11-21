package client

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/clbs-io/octopusmq/pkg/grpcmgmtpb"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcpb"
	"google.golang.org/grpc"
)

// Client is a simple wrapper of grpcpb.QueuesServiceClient, since grpcpb package is internal
// When you are done with the client, you should call the Close() method to release the resources
type Client struct {
	svcClient  grpcpb.QueuesServiceClient
	mgmtClient grpcmgmtpb.ManagementServiceClient
	grpcClient *grpc.ClientConn
	logger     *zap.SugaredLogger
}

type QueueClient struct {
	svcClient grpcpb.QueuesServiceClient
	stream    grpc.BidiStreamingClient[pb.QueueRequest, pb.QueueResponse]
	corrid    uint64
	opts      []grpc.CallOption // like really?
	ctx       context.Context
	qname     string
	closed    bool // till no close called, try to open at any cost
	logger    *zap.SugaredLogger
}

// NewClient creates a new client for the given address
// It panics if the client cannot be created
func NewClient(target string, logger *zap.SugaredLogger, grpcOptions ...grpc.DialOption) *Client {
	grpcClient, grpcClientErr := grpc.NewClient(target, grpcOptions...)
	if grpcClientErr != nil {
		panic(grpcClientErr)
	}

	svcClient := grpcpb.NewQueuesServiceClient(grpcClient)
	mgmtClient := grpcmgmtpb.NewManagementServiceClient(grpcClient)

	return &Client{
		svcClient:  svcClient,
		grpcClient: grpcClient,
		mgmtClient: mgmtClient,
		logger:     logger,
	}
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
}

// Queue Management
func (c *Client) EnsureQueue(ctx context.Context, req *pb.CreateQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.CreateQueue(ctx, req, opts...)

	if err == nil {
		return nil
	}

	st := status.Convert(err)

	// suppress the error and return
	if st.Code() == codes.AlreadyExists {
		return nil
	}

	return err
}

func handledeferrors(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.NotFound:
			return ErrQueueNotFound
		case codes.AlreadyExists:
			return ErrQueueAlreadyExists
		case codes.FailedPrecondition:
			return ErrQueuePaused
		}
	}
	return err
}

func (c *Client) CreateQueue(ctx context.Context, req *pb.CreateQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.CreateQueue(ctx, req, opts...)
	return handledeferrors(err)
}

func (c *Client) ResizeQueue(ctx context.Context, req *pb.ResizeQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.ResizeQueue(ctx, req, opts...)
	return handledeferrors(err)
}

func (c *Client) DeleteQueue(ctx context.Context, req *pb.DeleteQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.DeleteQueue(ctx, req, opts...)
	return handledeferrors(err)
}

func (c *Client) ListQueues(ctx context.Context, req *emptypb.Empty, opts ...grpc.CallOption) (*pb.ListQueuesResponse, error) {
	ret, err := c.mgmtClient.ListQueues(ctx, req, opts...)
	return ret, handledeferrors(err)
}

func (c *Client) PauseQueue(ctx context.Context, req *pb.PauseQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.PauseQueue(ctx, req, opts...)
	return handledeferrors(err)
}

func (c *Client) ResumeQueue(ctx context.Context, req *pb.ResumeQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.ResumeQueue(ctx, req, opts...)
	return handledeferrors(err)
}

func (c *Client) OpenQueue(ctx context.Context, name string, opts ...grpc.CallOption) *QueueClient {
	// opening on demand, on every attempt
	return &QueueClient{
		svcClient: c.svcClient,
		closed:    false,
		stream:    nil,
		corrid:    0,
		opts:      opts,
		ctx:       ctx,
		qname:     name,
		logger:    c.logger,
	}
}

func (c *QueueClient) open() (err error) {
	if c.closed {
		return ErrQueueClientClosed
	}
	if c.stream != nil {
		return nil
	}
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
	var resp *pb.QueueResponse
	resp, err = c.handleresp(&pb.QueueRequest{
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
	if st, ok := resp.Response.(*pb.QueueResponse_Status); ok {
		if st.Status.Code != pb.StatusCode_STATUS_CODE_OK {
			return fmt.Errorf("setup failed: %v", st.Status)
		}
	} else {
		return fmt.Errorf("setup failed, invalid response type: %T", resp)
	}
	return
}

func (c *QueueClient) handleresp(cmd *pb.QueueRequest) (reqp *pb.QueueResponse, err error) {
	err = c.open()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = c.stream.CloseSend()
			c.stream = nil
		}
	}()

	cmd.CorrelationId = atomic.AddUint64(&c.corrid, 1)
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
		return nil, err
	}
	reqp, err = c.stream.Recv()
	if err != nil {
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
	if reqp.CorrelationId != cmd.CorrelationId {
		return nil, fmt.Errorf("unexpected correlation id: %d != %d", cmd.CorrelationId, reqp.CorrelationId)
	}
	return reqp, nil
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

func (c *QueueClient) Close() error {
	if c.closed {
		return ErrQueueClientClosed
	}
	c.closed = true
	if c.stream == nil {
		return nil
	}
	return c.stream.CloseSend() // closed with error... consider it just closed...
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
		return nil, fmt.Errorf("unexpected response type: %T", cc)
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
