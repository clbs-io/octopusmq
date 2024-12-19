// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: queues_management.proto

package grpcmgmtpb

import (
	context "context"
	protobuf "github.com/clbs-io/octopusmq/api/protobuf"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ManagementService_CreateQueue_FullMethodName = "/io.clbs.octopusmq.grpc.management.ManagementService/CreateQueue"
	ManagementService_ResizeQueue_FullMethodName = "/io.clbs.octopusmq.grpc.management.ManagementService/ResizeQueue"
	ManagementService_DeleteQueue_FullMethodName = "/io.clbs.octopusmq.grpc.management.ManagementService/DeleteQueue"
	ManagementService_ListQueues_FullMethodName  = "/io.clbs.octopusmq.grpc.management.ManagementService/ListQueues"
	ManagementService_PauseQueue_FullMethodName  = "/io.clbs.octopusmq.grpc.management.ManagementService/PauseQueue"
	ManagementService_ResumeQueue_FullMethodName = "/io.clbs.octopusmq.grpc.management.ManagementService/ResumeQueue"
)

// ManagementServiceClient is the client API for ManagementService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The service represents management on the queues.
type ManagementServiceClient interface {
	// CreateQueue creates a new queue. The named queue starts immediately.
	// If the queue already exists, the operation fails.
	CreateQueue(ctx context.Context, in *protobuf.CreateQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ResizeQueue resizes the queue. If the queue does not exist, the operation fails. To shrink the queue, it has to be empty enought.
	ResizeQueue(ctx context.Context, in *protobuf.ResizeQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// DeleteQueue deletes the queue. If the queue does not exist, the operation fails.
	DeleteQueue(ctx context.Context, in *protobuf.DeleteQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ListQueues lists all queues defined on the server side.
	ListQueues(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.ListQueuesResponse, error)
	// PauseQueue pauses the queue. If the queue does not exist, the operation fails. If the queue is already paused, the operation does nothing.
	PauseQueue(ctx context.Context, in *protobuf.PauseQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// ResumeQueue resumes the queue. If the queue does not exist, the operation fails. If the queue is already running, the operation does nothing.
	ResumeQueue(ctx context.Context, in *protobuf.ResumeQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type managementServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewManagementServiceClient(cc grpc.ClientConnInterface) ManagementServiceClient {
	return &managementServiceClient{cc}
}

func (c *managementServiceClient) CreateQueue(ctx context.Context, in *protobuf.CreateQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ManagementService_CreateQueue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) ResizeQueue(ctx context.Context, in *protobuf.ResizeQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ManagementService_ResizeQueue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) DeleteQueue(ctx context.Context, in *protobuf.DeleteQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ManagementService_DeleteQueue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) ListQueues(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*protobuf.ListQueuesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.ListQueuesResponse)
	err := c.cc.Invoke(ctx, ManagementService_ListQueues_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) PauseQueue(ctx context.Context, in *protobuf.PauseQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ManagementService_PauseQueue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *managementServiceClient) ResumeQueue(ctx context.Context, in *protobuf.ResumeQueueRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, ManagementService_ResumeQueue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ManagementServiceServer is the server API for ManagementService service.
// All implementations must embed UnimplementedManagementServiceServer
// for forward compatibility.
//
// The service represents management on the queues.
type ManagementServiceServer interface {
	// CreateQueue creates a new queue. The named queue starts immediately.
	// If the queue already exists, the operation fails.
	CreateQueue(context.Context, *protobuf.CreateQueueRequest) (*emptypb.Empty, error)
	// ResizeQueue resizes the queue. If the queue does not exist, the operation fails. To shrink the queue, it has to be empty enought.
	ResizeQueue(context.Context, *protobuf.ResizeQueueRequest) (*emptypb.Empty, error)
	// DeleteQueue deletes the queue. If the queue does not exist, the operation fails.
	DeleteQueue(context.Context, *protobuf.DeleteQueueRequest) (*emptypb.Empty, error)
	// ListQueues lists all queues defined on the server side.
	ListQueues(context.Context, *emptypb.Empty) (*protobuf.ListQueuesResponse, error)
	// PauseQueue pauses the queue. If the queue does not exist, the operation fails. If the queue is already paused, the operation does nothing.
	PauseQueue(context.Context, *protobuf.PauseQueueRequest) (*emptypb.Empty, error)
	// ResumeQueue resumes the queue. If the queue does not exist, the operation fails. If the queue is already running, the operation does nothing.
	ResumeQueue(context.Context, *protobuf.ResumeQueueRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedManagementServiceServer()
}

// UnimplementedManagementServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedManagementServiceServer struct{}

func (UnimplementedManagementServiceServer) CreateQueue(context.Context, *protobuf.CreateQueueRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateQueue not implemented")
}
func (UnimplementedManagementServiceServer) ResizeQueue(context.Context, *protobuf.ResizeQueueRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResizeQueue not implemented")
}
func (UnimplementedManagementServiceServer) DeleteQueue(context.Context, *protobuf.DeleteQueueRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteQueue not implemented")
}
func (UnimplementedManagementServiceServer) ListQueues(context.Context, *emptypb.Empty) (*protobuf.ListQueuesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListQueues not implemented")
}
func (UnimplementedManagementServiceServer) PauseQueue(context.Context, *protobuf.PauseQueueRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PauseQueue not implemented")
}
func (UnimplementedManagementServiceServer) ResumeQueue(context.Context, *protobuf.ResumeQueueRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResumeQueue not implemented")
}
func (UnimplementedManagementServiceServer) mustEmbedUnimplementedManagementServiceServer() {}
func (UnimplementedManagementServiceServer) testEmbeddedByValue()                           {}

// UnsafeManagementServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ManagementServiceServer will
// result in compilation errors.
type UnsafeManagementServiceServer interface {
	mustEmbedUnimplementedManagementServiceServer()
}

func RegisterManagementServiceServer(s grpc.ServiceRegistrar, srv ManagementServiceServer) {
	// If the following call pancis, it indicates UnimplementedManagementServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ManagementService_ServiceDesc, srv)
}

func _ManagementService_CreateQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.CreateQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).CreateQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_CreateQueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).CreateQueue(ctx, req.(*protobuf.CreateQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_ResizeQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.ResizeQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).ResizeQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_ResizeQueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).ResizeQueue(ctx, req.(*protobuf.ResizeQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_DeleteQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.DeleteQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).DeleteQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_DeleteQueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).DeleteQueue(ctx, req.(*protobuf.DeleteQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_ListQueues_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).ListQueues(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_ListQueues_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).ListQueues(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_PauseQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.PauseQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).PauseQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_PauseQueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).PauseQueue(ctx, req.(*protobuf.PauseQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ManagementService_ResumeQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.ResumeQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ManagementServiceServer).ResumeQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ManagementService_ResumeQueue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ManagementServiceServer).ResumeQueue(ctx, req.(*protobuf.ResumeQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ManagementService_ServiceDesc is the grpc.ServiceDesc for ManagementService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ManagementService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "io.clbs.octopusmq.grpc.management.ManagementService",
	HandlerType: (*ManagementServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateQueue",
			Handler:    _ManagementService_CreateQueue_Handler,
		},
		{
			MethodName: "ResizeQueue",
			Handler:    _ManagementService_ResizeQueue_Handler,
		},
		{
			MethodName: "DeleteQueue",
			Handler:    _ManagementService_DeleteQueue_Handler,
		},
		{
			MethodName: "ListQueues",
			Handler:    _ManagementService_ListQueues_Handler,
		},
		{
			MethodName: "PauseQueue",
			Handler:    _ManagementService_PauseQueue_Handler,
		},
		{
			MethodName: "ResumeQueue",
			Handler:    _ManagementService_ResumeQueue_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "queues_management.proto",
}
