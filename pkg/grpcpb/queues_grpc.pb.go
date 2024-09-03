// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
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
	QueuesService_Push_FullMethodName          = "/io.clbs.octopusmq.grpc.QueuesService/Push"
	QueuesService_Prepend_FullMethodName       = "/io.clbs.octopusmq.grpc.QueuesService/Prepend"
	QueuesService_PullSingle_FullMethodName    = "/io.clbs.octopusmq.grpc.QueuesService/PullSingle"
	QueuesService_Pull_FullMethodName          = "/io.clbs.octopusmq.grpc.QueuesService/Pull"
	QueuesService_CommitSingle_FullMethodName  = "/io.clbs.octopusmq.grpc.QueuesService/CommitSingle"
	QueuesService_Commit_FullMethodName        = "/io.clbs.octopusmq.grpc.QueuesService/Commit"
	QueuesService_RequeueSingle_FullMethodName = "/io.clbs.octopusmq.grpc.QueuesService/RequeueSingle"
	QueuesService_Requeue_FullMethodName       = "/io.clbs.octopusmq.grpc.QueuesService/Requeue"
	QueuesService_DeleteSingle_FullMethodName  = "/io.clbs.octopusmq.grpc.QueuesService/DeleteSingle"
	QueuesService_Delete_FullMethodName        = "/io.clbs.octopusmq.grpc.QueuesService/Delete"
	QueuesService_Subscribe_FullMethodName     = "/io.clbs.octopusmq.grpc.QueuesService/Subscribe"
)

// QueuesServiceClient is the client API for QueuesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The service represents main operations on the queues.
type QueuesServiceClient interface {
	// Push enques an item to the queue, if the queue is not full (queue has max size) or unless the timeout expires.
	Push(ctx context.Context, in *protobuf.EnqueueRequest, opts ...grpc.CallOption) (*protobuf.EnqueueResponse, error)
	// Prepend adds an item to the head of the queue, if the queue is not full (queue has max size) or unless the timeout expires.
	Prepend(ctx context.Context, in *protobuf.EnqueueRequest, opts ...grpc.CallOption) (*protobuf.EnqueueResponse, error)
	// PullSingle pulls a single item from the queue.
	PullSingle(ctx context.Context, in *protobuf.PullSingleRequest, opts ...grpc.CallOption) (*protobuf.PullSingleResponse, error)
	// Pull pulls multiple items from the queue.
	Pull(ctx context.Context, in *protobuf.PullRequest, opts ...grpc.CallOption) (*protobuf.PullResponse, error)
	// CommitSingle commits a single item.
	CommitSingle(ctx context.Context, in *protobuf.CommitSingleRequest, opts ...grpc.CallOption) (*protobuf.CommitSingleResponse, error)
	// Commit commits multiple items.
	Commit(ctx context.Context, in *protobuf.CommitRequest, opts ...grpc.CallOption) (*protobuf.CommitResponse, error)
	// RequeueSingle requeues a single item.
	RequeueSingle(ctx context.Context, in *protobuf.RequeueSingleRequest, opts ...grpc.CallOption) (*protobuf.RequeueSingleResponse, error)
	// Requeue requeues multiple items.
	Requeue(ctx context.Context, in *protobuf.RequeueRequest, opts ...grpc.CallOption) (*protobuf.RequeueResponse, error)
	// DeleteSingle deletes a single item.
	DeleteSingle(ctx context.Context, in *protobuf.DeleteSingleRequest, opts ...grpc.CallOption) (*protobuf.DeleteSingleResponse, error)
	// Delete deletes multiple items.
	Delete(ctx context.Context, in *protobuf.DeleteRequest, opts ...grpc.CallOption) (*protobuf.DeleteResponse, error)
	// Subscribe subscribes to the queue, currently doing pulling based on signal from client, see SubscribeRequest for more details.
	Subscribe(ctx context.Context, in *protobuf.SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[protobuf.SubscribeResponse], error)
}

type queuesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewQueuesServiceClient(cc grpc.ClientConnInterface) QueuesServiceClient {
	return &queuesServiceClient{cc}
}

func (c *queuesServiceClient) Push(ctx context.Context, in *protobuf.EnqueueRequest, opts ...grpc.CallOption) (*protobuf.EnqueueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.EnqueueResponse)
	err := c.cc.Invoke(ctx, QueuesService_Push_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) Prepend(ctx context.Context, in *protobuf.EnqueueRequest, opts ...grpc.CallOption) (*protobuf.EnqueueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.EnqueueResponse)
	err := c.cc.Invoke(ctx, QueuesService_Prepend_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) PullSingle(ctx context.Context, in *protobuf.PullSingleRequest, opts ...grpc.CallOption) (*protobuf.PullSingleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.PullSingleResponse)
	err := c.cc.Invoke(ctx, QueuesService_PullSingle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) Pull(ctx context.Context, in *protobuf.PullRequest, opts ...grpc.CallOption) (*protobuf.PullResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.PullResponse)
	err := c.cc.Invoke(ctx, QueuesService_Pull_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) CommitSingle(ctx context.Context, in *protobuf.CommitSingleRequest, opts ...grpc.CallOption) (*protobuf.CommitSingleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.CommitSingleResponse)
	err := c.cc.Invoke(ctx, QueuesService_CommitSingle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) Commit(ctx context.Context, in *protobuf.CommitRequest, opts ...grpc.CallOption) (*protobuf.CommitResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.CommitResponse)
	err := c.cc.Invoke(ctx, QueuesService_Commit_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) RequeueSingle(ctx context.Context, in *protobuf.RequeueSingleRequest, opts ...grpc.CallOption) (*protobuf.RequeueSingleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.RequeueSingleResponse)
	err := c.cc.Invoke(ctx, QueuesService_RequeueSingle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) Requeue(ctx context.Context, in *protobuf.RequeueRequest, opts ...grpc.CallOption) (*protobuf.RequeueResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.RequeueResponse)
	err := c.cc.Invoke(ctx, QueuesService_Requeue_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) DeleteSingle(ctx context.Context, in *protobuf.DeleteSingleRequest, opts ...grpc.CallOption) (*protobuf.DeleteSingleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.DeleteSingleResponse)
	err := c.cc.Invoke(ctx, QueuesService_DeleteSingle_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) Delete(ctx context.Context, in *protobuf.DeleteRequest, opts ...grpc.CallOption) (*protobuf.DeleteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(protobuf.DeleteResponse)
	err := c.cc.Invoke(ctx, QueuesService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queuesServiceClient) Subscribe(ctx context.Context, in *protobuf.SubscribeRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[protobuf.SubscribeResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &QueuesService_ServiceDesc.Streams[0], QueuesService_Subscribe_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[protobuf.SubscribeRequest, protobuf.SubscribeResponse]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type QueuesService_SubscribeClient = grpc.ServerStreamingClient[protobuf.SubscribeResponse]

// QueuesServiceServer is the server API for QueuesService service.
// All implementations must embed UnimplementedQueuesServiceServer
// for forward compatibility.
//
// The service represents main operations on the queues.
type QueuesServiceServer interface {
	// Push enques an item to the queue, if the queue is not full (queue has max size) or unless the timeout expires.
	Push(context.Context, *protobuf.EnqueueRequest) (*protobuf.EnqueueResponse, error)
	// Prepend adds an item to the head of the queue, if the queue is not full (queue has max size) or unless the timeout expires.
	Prepend(context.Context, *protobuf.EnqueueRequest) (*protobuf.EnqueueResponse, error)
	// PullSingle pulls a single item from the queue.
	PullSingle(context.Context, *protobuf.PullSingleRequest) (*protobuf.PullSingleResponse, error)
	// Pull pulls multiple items from the queue.
	Pull(context.Context, *protobuf.PullRequest) (*protobuf.PullResponse, error)
	// CommitSingle commits a single item.
	CommitSingle(context.Context, *protobuf.CommitSingleRequest) (*protobuf.CommitSingleResponse, error)
	// Commit commits multiple items.
	Commit(context.Context, *protobuf.CommitRequest) (*protobuf.CommitResponse, error)
	// RequeueSingle requeues a single item.
	RequeueSingle(context.Context, *protobuf.RequeueSingleRequest) (*protobuf.RequeueSingleResponse, error)
	// Requeue requeues multiple items.
	Requeue(context.Context, *protobuf.RequeueRequest) (*protobuf.RequeueResponse, error)
	// DeleteSingle deletes a single item.
	DeleteSingle(context.Context, *protobuf.DeleteSingleRequest) (*protobuf.DeleteSingleResponse, error)
	// Delete deletes multiple items.
	Delete(context.Context, *protobuf.DeleteRequest) (*protobuf.DeleteResponse, error)
	// Subscribe subscribes to the queue, currently doing pulling based on signal from client, see SubscribeRequest for more details.
	Subscribe(*protobuf.SubscribeRequest, grpc.ServerStreamingServer[protobuf.SubscribeResponse]) error
	mustEmbedUnimplementedQueuesServiceServer()
}

// UnimplementedQueuesServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedQueuesServiceServer struct{}

func (UnimplementedQueuesServiceServer) Push(context.Context, *protobuf.EnqueueRequest) (*protobuf.EnqueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Push not implemented")
}
func (UnimplementedQueuesServiceServer) Prepend(context.Context, *protobuf.EnqueueRequest) (*protobuf.EnqueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Prepend not implemented")
}
func (UnimplementedQueuesServiceServer) PullSingle(context.Context, *protobuf.PullSingleRequest) (*protobuf.PullSingleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PullSingle not implemented")
}
func (UnimplementedQueuesServiceServer) Pull(context.Context, *protobuf.PullRequest) (*protobuf.PullResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Pull not implemented")
}
func (UnimplementedQueuesServiceServer) CommitSingle(context.Context, *protobuf.CommitSingleRequest) (*protobuf.CommitSingleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitSingle not implemented")
}
func (UnimplementedQueuesServiceServer) Commit(context.Context, *protobuf.CommitRequest) (*protobuf.CommitResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedQueuesServiceServer) RequeueSingle(context.Context, *protobuf.RequeueSingleRequest) (*protobuf.RequeueSingleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequeueSingle not implemented")
}
func (UnimplementedQueuesServiceServer) Requeue(context.Context, *protobuf.RequeueRequest) (*protobuf.RequeueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Requeue not implemented")
}
func (UnimplementedQueuesServiceServer) DeleteSingle(context.Context, *protobuf.DeleteSingleRequest) (*protobuf.DeleteSingleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSingle not implemented")
}
func (UnimplementedQueuesServiceServer) Delete(context.Context, *protobuf.DeleteRequest) (*protobuf.DeleteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedQueuesServiceServer) Subscribe(*protobuf.SubscribeRequest, grpc.ServerStreamingServer[protobuf.SubscribeResponse]) error {
	return status.Errorf(codes.Unimplemented, "method Subscribe not implemented")
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

func _QueuesService_Push_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.EnqueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).Push(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_Push_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).Push(ctx, req.(*protobuf.EnqueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_Prepend_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.EnqueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).Prepend(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_Prepend_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).Prepend(ctx, req.(*protobuf.EnqueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_PullSingle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.PullSingleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).PullSingle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_PullSingle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).PullSingle(ctx, req.(*protobuf.PullSingleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_Pull_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.PullRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).Pull(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_Pull_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).Pull(ctx, req.(*protobuf.PullRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_CommitSingle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.CommitSingleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).CommitSingle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_CommitSingle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).CommitSingle(ctx, req.(*protobuf.CommitSingleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.CommitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_Commit_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).Commit(ctx, req.(*protobuf.CommitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_RequeueSingle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.RequeueSingleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).RequeueSingle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_RequeueSingle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).RequeueSingle(ctx, req.(*protobuf.RequeueSingleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_Requeue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.RequeueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).Requeue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_Requeue_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).Requeue(ctx, req.(*protobuf.RequeueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_DeleteSingle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.DeleteSingleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).DeleteSingle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_DeleteSingle_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).DeleteSingle(ctx, req.(*protobuf.DeleteSingleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protobuf.DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueuesServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: QueuesService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueuesServiceServer).Delete(ctx, req.(*protobuf.DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QueuesService_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(protobuf.SubscribeRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(QueuesServiceServer).Subscribe(m, &grpc.GenericServerStream[protobuf.SubscribeRequest, protobuf.SubscribeResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type QueuesService_SubscribeServer = grpc.ServerStreamingServer[protobuf.SubscribeResponse]

// QueuesService_ServiceDesc is the grpc.ServiceDesc for QueuesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QueuesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "io.clbs.octopusmq.grpc.QueuesService",
	HandlerType: (*QueuesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Push",
			Handler:    _QueuesService_Push_Handler,
		},
		{
			MethodName: "Prepend",
			Handler:    _QueuesService_Prepend_Handler,
		},
		{
			MethodName: "PullSingle",
			Handler:    _QueuesService_PullSingle_Handler,
		},
		{
			MethodName: "Pull",
			Handler:    _QueuesService_Pull_Handler,
		},
		{
			MethodName: "CommitSingle",
			Handler:    _QueuesService_CommitSingle_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _QueuesService_Commit_Handler,
		},
		{
			MethodName: "RequeueSingle",
			Handler:    _QueuesService_RequeueSingle_Handler,
		},
		{
			MethodName: "Requeue",
			Handler:    _QueuesService_Requeue_Handler,
		},
		{
			MethodName: "DeleteSingle",
			Handler:    _QueuesService_DeleteSingle_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _QueuesService_Delete_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Subscribe",
			Handler:       _QueuesService_Subscribe_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "queues.proto",
}
