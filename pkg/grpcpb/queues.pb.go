// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.3
// source: queues.proto

package grpcpb

import (
	protobuf "github.com/clbs-io/octopusmq/api/protobuf"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_queues_proto protoreflect.FileDescriptor

var file_queues_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x71, 0x75, 0x65, 0x75, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x16,
	0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d,
	0x71, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x1a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x73, 0x0a, 0x0d, 0x51, 0x75, 0x65, 0x75, 0x65, 0x73,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x62, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x12, 0x28, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74,
	0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x29, 0x2e, 0x69,
	0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x28, 0x01, 0x30, 0x01, 0x42, 0x29, 0x5a, 0x27, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x62, 0x73, 0x2d, 0x69,
	0x6f, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x67, 0x72, 0x70, 0x63, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_queues_proto_goTypes = []any{
	(*protobuf.QueueRequest)(nil),  // 0: io.clbs.octopusmq.protobuf.QueueRequest
	(*protobuf.QueueResponse)(nil), // 1: io.clbs.octopusmq.protobuf.QueueResponse
}
var file_queues_proto_depIdxs = []int32{
	0, // 0: io.clbs.octopusmq.grpc.QueuesService.Connect:input_type -> io.clbs.octopusmq.protobuf.QueueRequest
	1, // 1: io.clbs.octopusmq.grpc.QueuesService.Connect:output_type -> io.clbs.octopusmq.protobuf.QueueResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_queues_proto_init() }
func file_queues_proto_init() {
	if File_queues_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_queues_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_queues_proto_goTypes,
		DependencyIndexes: file_queues_proto_depIdxs,
	}.Build()
	File_queues_proto = out.File
	file_queues_proto_rawDesc = nil
	file_queues_proto_goTypes = nil
	file_queues_proto_depIdxs = nil
}
