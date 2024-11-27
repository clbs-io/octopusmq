package client

import (
	"context"

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
