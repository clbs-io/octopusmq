// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.2
// source: queues_management.proto

package grpcmgmtpb

import (
	protobuf "github.com/clbs-io/octopusmq/api/protobuf"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_queues_management_proto protoreflect.FileDescriptor

var file_queues_management_proto_rawDesc = []byte{
	0x0a, 0x17, 0x71, 0x75, 0x65, 0x75, 0x65, 0x73, 0x5f, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x21, 0x69, 0x6f, 0x2e, 0x63, 0x6c,
	0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x32, 0x9a, 0x04, 0x0a, 0x11, 0x4d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x55, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x2e,
	0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73,
	0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x55, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x69, 0x7a, 0x65,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x2e, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e,
	0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x52, 0x65, 0x73, 0x69, 0x7a, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x55, 0x0a,
	0x0b, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x2e, 0x2e, 0x69,
	0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x54, 0x0a, 0x0a, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x2e, 0x2e, 0x69, 0x6f, 0x2e,
	0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x51, 0x75, 0x65, 0x75,
	0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x53, 0x0a, 0x0a, 0x50, 0x61,
	0x75, 0x73, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x2d, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c,
	0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x50, 0x61, 0x75, 0x73, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12,
	0x55, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x2e,
	0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73,
	0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x52, 0x65, 0x73, 0x75,
	0x6d, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6c, 0x62, 0x73, 0x2d, 0x69, 0x6f, 0x2f, 0x6f, 0x63, 0x74,
	0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x6d,
	0x67, 0x6d, 0x74, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_queues_management_proto_goTypes = []any{
	(*protobuf.CreateQueueRequest)(nil), // 0: io.clbs.octopusmq.protobuf.CreateQueueRequest
	(*protobuf.ResizeQueueRequest)(nil), // 1: io.clbs.octopusmq.protobuf.ResizeQueueRequest
	(*protobuf.DeleteQueueRequest)(nil), // 2: io.clbs.octopusmq.protobuf.DeleteQueueRequest
	(*emptypb.Empty)(nil),               // 3: google.protobuf.Empty
	(*protobuf.PauseQueueRequest)(nil),  // 4: io.clbs.octopusmq.protobuf.PauseQueueRequest
	(*protobuf.ResumeQueueRequest)(nil), // 5: io.clbs.octopusmq.protobuf.ResumeQueueRequest
	(*protobuf.ListQueuesResponse)(nil), // 6: io.clbs.octopusmq.protobuf.ListQueuesResponse
}
var file_queues_management_proto_depIdxs = []int32{
	0, // 0: io.clbs.octopusmq.grpc.management.ManagementService.CreateQueue:input_type -> io.clbs.octopusmq.protobuf.CreateQueueRequest
	1, // 1: io.clbs.octopusmq.grpc.management.ManagementService.ResizeQueue:input_type -> io.clbs.octopusmq.protobuf.ResizeQueueRequest
	2, // 2: io.clbs.octopusmq.grpc.management.ManagementService.DeleteQueue:input_type -> io.clbs.octopusmq.protobuf.DeleteQueueRequest
	3, // 3: io.clbs.octopusmq.grpc.management.ManagementService.ListQueues:input_type -> google.protobuf.Empty
	4, // 4: io.clbs.octopusmq.grpc.management.ManagementService.PauseQueue:input_type -> io.clbs.octopusmq.protobuf.PauseQueueRequest
	5, // 5: io.clbs.octopusmq.grpc.management.ManagementService.ResumeQueue:input_type -> io.clbs.octopusmq.protobuf.ResumeQueueRequest
	3, // 6: io.clbs.octopusmq.grpc.management.ManagementService.CreateQueue:output_type -> google.protobuf.Empty
	3, // 7: io.clbs.octopusmq.grpc.management.ManagementService.ResizeQueue:output_type -> google.protobuf.Empty
	3, // 8: io.clbs.octopusmq.grpc.management.ManagementService.DeleteQueue:output_type -> google.protobuf.Empty
	6, // 9: io.clbs.octopusmq.grpc.management.ManagementService.ListQueues:output_type -> io.clbs.octopusmq.protobuf.ListQueuesResponse
	3, // 10: io.clbs.octopusmq.grpc.management.ManagementService.PauseQueue:output_type -> google.protobuf.Empty
	3, // 11: io.clbs.octopusmq.grpc.management.ManagementService.ResumeQueue:output_type -> google.protobuf.Empty
	6, // [6:12] is the sub-list for method output_type
	0, // [0:6] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_queues_management_proto_init() }
func file_queues_management_proto_init() {
	if File_queues_management_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_queues_management_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_queues_management_proto_goTypes,
		DependencyIndexes: file_queues_management_proto_depIdxs,
	}.Build()
	File_queues_management_proto = out.File
	file_queues_management_proto_rawDesc = nil
	file_queues_management_proto_goTypes = nil
	file_queues_management_proto_depIdxs = nil
}
