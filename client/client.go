package client

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/clbs-io/octopusmq/pkg/grpcmgmtpb"
	"github.com/clbs-io/octopusmq/pkg/grpcstoragemgmtpb"
	"github.com/clbs-io/octopusmq/pkg/grpcstoragepb"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/clbs-io/octopusmq/api/protobuf"
	"github.com/clbs-io/octopusmq/pkg/grpcpb"
	"google.golang.org/grpc"
)

// Client manages a gRPC connection and provides access to queue and storage services.
// Call Close() when done to release the underlying connection.
type Client struct {
	queueConnect grpcpb.QueuesServiceClient
	mgmtClient   grpcmgmtpb.ManagementServiceClient
	stoClient    grpcstoragemgmtpb.StorageManagementServiceClient
	stoConnect   grpcstoragepb.StorageServiceClient
	grpcClient   *grpc.ClientConn
	logger       *zap.SugaredLogger
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
	stoClient := grpcstoragemgmtpb.NewStorageManagementServiceClient(grpcClient)
	stoConnect := grpcstoragepb.NewStorageServiceClient(grpcClient)

	return &Client{
		queueConnect: svcClient,
		grpcClient:   grpcClient,
		mgmtClient:   mgmtClient,
		stoClient:    stoClient,
		stoConnect:   stoConnect,
		logger:       logger,
	}
}

func (c *Client) Close() error {
	return c.grpcClient.Close()
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

// EnsureQueue creates the queue and reports success when it already exists.
func (c *Client) EnsureQueue(ctx context.Context, req *pb.CreateQueueRequest, opts ...grpc.CallOption) error {
	err := c.CreateQueue(ctx, req, opts...)
	if errors.Is(err, ErrQueueAlreadyExists) {
		return nil
	}
	return err
}

func (c *Client) ResizeQueue(ctx context.Context, req *pb.ResizeQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.ResizeQueue(ctx, req, opts...)
	return handledeferrors(err)
}

func (c *Client) DeleteQueue(ctx context.Context, req *pb.DeleteQueueRequest, opts ...grpc.CallOption) error {
	_, err := c.mgmtClient.DeleteQueue(ctx, req, opts...)
	return handledeferrors(err)
}

func (c *Client) ListQueues(ctx context.Context, opts ...grpc.CallOption) (*pb.ListQueuesResponse, error) {
	ret, err := c.mgmtClient.ListQueues(ctx, &emptypb.Empty{}, opts...)
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

func handledefsterrors(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if ok {
		switch st.Code() {
		case codes.NotFound:
			return ErrStorageNotFound
		case codes.AlreadyExists:
			return ErrStorageAlreadyExists
		}
	}
	return err
}

func (c *Client) CreateStorage(ctx context.Context, req *pb.CreateStorageRequest, opts ...grpc.CallOption) error {
	_, err := c.stoClient.CreateStorage(ctx, req, opts...)
	return handledefsterrors(err)
}

func (c *Client) DeleteStorage(ctx context.Context, req *pb.DeleteStorageRequest, opts ...grpc.CallOption) error {
	_, err := c.stoClient.DeleteStorage(ctx, req, opts...)
	return handledefsterrors(err)
}

func (c *Client) ListStorages(ctx context.Context, opts ...grpc.CallOption) (*pb.ListStoragesResponse, error) {
	ret, err := c.stoClient.ListStorages(ctx, &emptypb.Empty{}, opts...)
	return ret, handledefsterrors(err)
}

// EnsureStorage creates the storage and reports success when it already exists.
func (c *Client) EnsureStorage(ctx context.Context, req *pb.CreateStorageRequest, opts ...grpc.CallOption) error {
	err := c.CreateStorage(ctx, req, opts...)
	if errors.Is(err, ErrStorageAlreadyExists) {
		return nil
	}
	return err
}
