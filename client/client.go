package client

import (
	"context"
	"errors"

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

func (c *Client) Push(ctx context.Context, req *pb.EnqueueRequest, opts ...grpc.CallOption) (*pb.EnqueueResponse, error) {
	return c.svcClient.Push(ctx, req, opts...)
}

func (c *Client) Prepend(ctx context.Context, req *pb.EnqueueRequest, opts ...grpc.CallOption) (*pb.EnqueueResponse, error) {
	return c.svcClient.Prepend(ctx, req, opts...)
}

func (c *Client) Pull(ctx context.Context, req *pb.PullRequest, opts ...grpc.CallOption) (*pb.PullResponse, error) {
	return c.svcClient.Pull(ctx, req, opts...)
}

func (c *Client) PullSingle(ctx context.Context, req *pb.PullSingleRequest, opts ...grpc.CallOption) (*pb.PullSingleResponse, error) {
	return c.svcClient.PullSingle(ctx, req, opts...)
}

func (c *Client) CommitSingle(ctx context.Context, req *pb.CommitSingleRequest, opts ...grpc.CallOption) (*pb.CommitSingleResponse, error) {
	return c.svcClient.CommitSingle(ctx, req, opts...)
}

func (c *Client) Commit(ctx context.Context, req *pb.CommitRequest, opts ...grpc.CallOption) (*pb.CommitResponse, error) {
	return c.svcClient.Commit(ctx, req, opts...)
}

func (c *Client) RequeueSingle(ctx context.Context, req *pb.RequeueSingleRequest, opts ...grpc.CallOption) (*pb.RequeueSingleResponse, error) {
	return c.svcClient.RequeueSingle(ctx, req, opts...)
}

func (c *Client) Requeue(ctx context.Context, req *pb.RequeueRequest, opts ...grpc.CallOption) (*pb.RequeueResponse, error) {
	return c.svcClient.Requeue(ctx, req, opts...)
}

func (c *Client) DeleteSingle(ctx context.Context, req *pb.DeleteSingleRequest, opts ...grpc.CallOption) (*pb.DeleteSingleResponse, error) {
	return c.svcClient.DeleteSingle(ctx, req, opts...)
}

func (c *Client) Delete(ctx context.Context, req *pb.DeleteRequest, opts ...grpc.CallOption) (*pb.DeleteResponse, error) {
	return c.svcClient.Delete(ctx, req, opts...)
}

func (c *Client) Subscribe(ctx context.Context, opts ...grpc.CallOption) (*Subscriber, error) {
	panic("not implemented")
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
}

// Queue Management

func (c *Client) EnsureQueue(ctx context.Context, req *pb.CreateQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	res, err := c.mgmtClient.CreateQueue(ctx, req, opts...)
	if errors.Is(err, ErrQueueAlreadyExists) {
		// we suppress the error and return
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
