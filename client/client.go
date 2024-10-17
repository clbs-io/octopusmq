package client

import (
	"context"

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
}

// NewClient creates a new client for the given address
// It panics if the client cannot be created
func NewClient(target string, grpcOptions ...grpc.DialOption) *Client {
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
	}
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
}

// Queue Management

func (c *Client) EnsureQueue(ctx context.Context, req *pb.CreateQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	res, err := c.mgmtClient.CreateQueue(ctx, req, opts...)

	if err == nil {
		return res, nil
	}

	st := status.Convert(err)

	// suppress the error and return
	if st.Code() == codes.AlreadyExists {
		return res, nil
	}

	return res, err
}

func (c *Client) CreateQueue(ctx context.Context, req *pb.CreateQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.mgmtClient.CreateQueue(ctx, req, opts...)
}

func (c *Client) DeleteQueue(ctx context.Context, req *pb.DeleteQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.mgmtClient.DeleteQueue(ctx, req, opts...)
}

func (c *Client) ListQueues(ctx context.Context, req *emptypb.Empty, opts ...grpc.CallOption) (*pb.ListQueuesResponse, error) {
	return c.mgmtClient.ListQueues(ctx, req, opts...)
}

func (c *Client) PauseQueue(ctx context.Context, req *pb.PauseQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.mgmtClient.PauseQueue(ctx, req, opts...)
}

func (c *Client) ResumeQueue(ctx context.Context, req *pb.ResumeQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return c.mgmtClient.ResumeQueue(ctx, req, opts...)
}
