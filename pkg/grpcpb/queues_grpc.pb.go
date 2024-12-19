// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: queues.proto

package grpcpb

import (
	context "context"
	protobuf "github.com/clbs-io/octopusmq/api/protobuf"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	QueuesService_Connect_FullMethodName = "/io.clbs.octopusmq.grpc.QueuesService/Connect"
)

// QueuesServiceClient is the client API for QueuesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The service represents main operations on the queues.
type QueuesServiceClient interface {
	// Main command connection to the queue server. Inside this stream, all commands are sent.
	Connect(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[protobuf.QueueRequest, protobuf.QueueResponse], error)
}

type queuesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQueuesServiceClient(cc grpc.ClientConnInterface) QueuesServiceClient {
	return &queuesServiceClient{cc}
}

func (c *queuesServiceClient) Connect(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[protobuf.QueueRequest, protobuf.QueueResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &QueuesService_ServiceDesc.Streams[0], QueuesService_Connect_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[protobuf.QueueRequest, protobuf.QueueResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type QueuesService_ConnectClient = grpc.BidiStreamingClient[protobuf.QueueRequest, protobuf.QueueResponse]

// QueuesServiceServer is the server API for QueuesService service.
// All implementations must embed UnimplementedQueuesServiceServer
// for forward compatibility.
//
// The service represents main operations on the queues.
type QueuesServiceServer interface {
	// Main command connection to the queue server. Inside this stream, all commands are sent.
	Connect(grpc.BidiStreamingServer[protobuf.QueueRequest, protobuf.QueueResponse]) error
	mustEmbedUnimplementedQueuesServiceServer()
}

// UnimplementedQueuesServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueuesServiceServer struct{}

func (UnimplementedQueuesServiceServer) Connect(grpc.BidiStreamingServer[protobuf.QueueRequest, protobuf.QueueResponse]) error {
	return status.Errorf(codes.Unimplemented, "method Connect not implemented")
}
func (UnimplementedQueuesServiceServer) mustEmbedUnimplementedQueuesServiceServer() {}
func (UnimplementedQueuesServiceServer) testEmbeddedByValue()                       {}

// UnsafeQueuesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QueuesServiceServer will
// result in compilation errors.
type UnsafeQueuesServiceServer interface {
	mustEmbedUnimplementedQueuesServiceServer()
}

func RegisterQueuesServiceServer(s grpc.ServiceRegistrar, srv QueuesServiceServer) {
	// If the following call pancis, it indicates UnimplementedQueuesServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&QueuesService_ServiceDesc, srv)
}

func _QueuesService_Connect_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(QueuesServiceServer).Connect(&grpc.GenericServerStream[protobuf.QueueRequest, protobuf.QueueResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type QueuesService_ConnectServer = grpc.BidiStreamingServer[protobuf.QueueRequest, protobuf.QueueResponse]

// QueuesService_ServiceDesc is the grpc.ServiceDesc for QueuesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QueuesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "io.clbs.octopusmq.grpc.QueuesService",
	HandlerType: (*QueuesServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Connect",
			Handler:       _QueuesService_Connect_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "queues.proto",
}
