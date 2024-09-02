// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.3
// source: raftcommands.proto

package grpcraftcmdpb

import (
	protobuf "github.com/clbs-io/octopusmq/api/protobuf"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RaftCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Command:
	//
	//	*RaftCommand_Createqueue
	//	*RaftCommand_Deletequeue
	//	*RaftCommand_Pauseresumequeue
	//	*RaftCommand_Purgequeue
	//	*RaftCommand_Storagecommands
	//	*RaftCommand_Storagedeleteids
	Command isRaftCommand_Command `protobuf_oneof:"command"`
}

func (x *RaftCommand) Reset() {
	*x = RaftCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftCommand) ProtoMessage() {}

func (x *RaftCommand) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftCommand.ProtoReflect.Descriptor instead.
func (*RaftCommand) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{0}
}

func (m *RaftCommand) GetCommand() isRaftCommand_Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func (x *RaftCommand) GetCreatequeue() *RaftCreateQueue {
	if x, ok := x.GetCommand().(*RaftCommand_Createqueue); ok {
		return x.Createqueue
	}
	return nil
}

func (x *RaftCommand) GetDeletequeue() *protobuf.DeleteQueueRequest {
	if x, ok := x.GetCommand().(*RaftCommand_Deletequeue); ok {
		return x.Deletequeue
	}
	return nil
}

func (x *RaftCommand) GetPauseresumequeue() *RaftPauseResumeQueue {
	if x, ok := x.GetCommand().(*RaftCommand_Pauseresumequeue); ok {
		return x.Pauseresumequeue
	}
	return nil
}

func (x *RaftCommand) GetPurgequeue() *RaftPurgeQueue {
	if x, ok := x.GetCommand().(*RaftCommand_Purgequeue); ok {
		return x.Purgequeue
	}
	return nil
}

func (x *RaftCommand) GetStoragecommands() *RaftStorageCommands {
	if x, ok := x.GetCommand().(*RaftCommand_Storagecommands); ok {
		return x.Storagecommands
	}
	return nil
}

func (x *RaftCommand) GetStoragedeleteids() *RaftStorageDeleteIds {
	if x, ok := x.GetCommand().(*RaftCommand_Storagedeleteids); ok {
		return x.Storagedeleteids
	}
	return nil
}

type isRaftCommand_Command interface {
	isRaftCommand_Command()
}

type RaftCommand_Createqueue struct {
	Createqueue *RaftCreateQueue `protobuf:"bytes,1,opt,name=createqueue,proto3,oneof"`
}

type RaftCommand_Deletequeue struct {
	Deletequeue *protobuf.DeleteQueueRequest `protobuf:"bytes,2,opt,name=deletequeue,proto3,oneof"`
}

type RaftCommand_Pauseresumequeue struct {
	Pauseresumequeue *RaftPauseResumeQueue `protobuf:"bytes,3,opt,name=pauseresumequeue,proto3,oneof"`
}

type RaftCommand_Purgequeue struct {
	Purgequeue *RaftPurgeQueue `protobuf:"bytes,4,opt,name=purgequeue,proto3,oneof"`
}

type RaftCommand_Storagecommands struct {
	// internal storage commands
	Storagecommands *RaftStorageCommands `protobuf:"bytes,5,opt,name=storagecommands,proto3,oneof"`
}

type RaftCommand_Storagedeleteids struct {
	Storagedeleteids *RaftStorageDeleteIds `protobuf:"bytes,6,opt,name=storagedeleteids,proto3,oneof"`
}

func (*RaftCommand_Createqueue) isRaftCommand_Command() {}

func (*RaftCommand_Deletequeue) isRaftCommand_Command() {}

func (*RaftCommand_Pauseresumequeue) isRaftCommand_Command() {}

func (*RaftCommand_Purgequeue) isRaftCommand_Command() {}

func (*RaftCommand_Storagecommands) isRaftCommand_Command() {}

func (*RaftCommand_Storagedeleteids) isRaftCommand_Command() {}

type RaftCreateQueue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Queueid uint32                       `protobuf:"varint,1,opt,name=queueid,proto3" json:"queueid,omitempty"`
	Queue   *protobuf.CreateQueueRequest `protobuf:"bytes,2,opt,name=queue,proto3" json:"queue,omitempty"`
}

func (x *RaftCreateQueue) Reset() {
	*x = RaftCreateQueue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftCreateQueue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftCreateQueue) ProtoMessage() {}

func (x *RaftCreateQueue) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftCreateQueue.ProtoReflect.Descriptor instead.
func (*RaftCreateQueue) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{1}
}

func (x *RaftCreateQueue) GetQueueid() uint32 {
	if x != nil {
		return x.Queueid
	}
	return 0
}

func (x *RaftCreateQueue) GetQueue() *protobuf.CreateQueueRequest {
	if x != nil {
		return x.Queue
	}
	return nil
}

type RaftPauseResumeQueue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueName string `protobuf:"bytes,1,opt,name=queueName,proto3" json:"queueName,omitempty"`
	Pause     bool   `protobuf:"varint,2,opt,name=pause,proto3" json:"pause,omitempty"`
}

func (x *RaftPauseResumeQueue) Reset() {
	*x = RaftPauseResumeQueue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftPauseResumeQueue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftPauseResumeQueue) ProtoMessage() {}

func (x *RaftPauseResumeQueue) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftPauseResumeQueue.ProtoReflect.Descriptor instead.
func (*RaftPauseResumeQueue) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{2}
}

func (x *RaftPauseResumeQueue) GetQueueName() string {
	if x != nil {
		return x.QueueName
	}
	return ""
}

func (x *RaftPauseResumeQueue) GetPause() bool {
	if x != nil {
		return x.Pause
	}
	return false
}

type RaftPurgeQueue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	All       bool     `protobuf:"varint,1,opt,name=all,proto3" json:"all,omitempty"`
	QueueName []string `protobuf:"bytes,2,rep,name=queueName,proto3" json:"queueName,omitempty"`
}

func (x *RaftPurgeQueue) Reset() {
	*x = RaftPurgeQueue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftPurgeQueue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftPurgeQueue) ProtoMessage() {}

func (x *RaftPurgeQueue) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftPurgeQueue.ProtoReflect.Descriptor instead.
func (*RaftPurgeQueue) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{3}
}

func (x *RaftPurgeQueue) GetAll() bool {
	if x != nil {
		return x.All
	}
	return false
}

func (x *RaftPurgeQueue) GetQueueName() []string {
	if x != nil {
		return x.QueueName
	}
	return nil
}

type RaftStorageCommands struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QueueName string                `protobuf:"bytes,1,opt,name=queueName,proto3" json:"queueName,omitempty"`
	Id        uint32                `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Commands  []*RaftStorageCommand `protobuf:"bytes,3,rep,name=commands,proto3" json:"commands,omitempty"`
}

func (x *RaftStorageCommands) Reset() {
	*x = RaftStorageCommands{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftStorageCommands) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftStorageCommands) ProtoMessage() {}

func (x *RaftStorageCommands) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftStorageCommands.ProtoReflect.Descriptor instead.
func (*RaftStorageCommands) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{4}
}

func (x *RaftStorageCommands) GetQueueName() string {
	if x != nil {
		return x.QueueName
	}
	return ""
}

func (x *RaftStorageCommands) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RaftStorageCommands) GetCommands() []*RaftStorageCommand {
	if x != nil {
		return x.Commands
	}
	return nil
}

type RaftStorageDeleteIds struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ids []uint64 `protobuf:"varint,1,rep,packed,name=ids,proto3" json:"ids,omitempty"`
}

func (x *RaftStorageDeleteIds) Reset() {
	*x = RaftStorageDeleteIds{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftStorageDeleteIds) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftStorageDeleteIds) ProtoMessage() {}

func (x *RaftStorageDeleteIds) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftStorageDeleteIds.ProtoReflect.Descriptor instead.
func (*RaftStorageDeleteIds) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{5}
}

func (x *RaftStorageDeleteIds) GetIds() []uint64 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type RaftStorageCommand struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Command:
	//
	//	*RaftStorageCommand_Delete
	//	*RaftStorageCommand_Insert
	//	*RaftStorageCommand_Flags
	Command isRaftStorageCommand_Command `protobuf_oneof:"command"`
}

func (x *RaftStorageCommand) Reset() {
	*x = RaftStorageCommand{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftStorageCommand) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftStorageCommand) ProtoMessage() {}

func (x *RaftStorageCommand) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftStorageCommand.ProtoReflect.Descriptor instead.
func (*RaftStorageCommand) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{6}
}

func (m *RaftStorageCommand) GetCommand() isRaftStorageCommand_Command {
	if m != nil {
		return m.Command
	}
	return nil
}

func (x *RaftStorageCommand) GetDelete() *RaftStorageDelete {
	if x, ok := x.GetCommand().(*RaftStorageCommand_Delete); ok {
		return x.Delete
	}
	return nil
}

func (x *RaftStorageCommand) GetInsert() *RaftStorageInsert {
	if x, ok := x.GetCommand().(*RaftStorageCommand_Insert); ok {
		return x.Insert
	}
	return nil
}

func (x *RaftStorageCommand) GetFlags() *RaftStorageFlags {
	if x, ok := x.GetCommand().(*RaftStorageCommand_Flags); ok {
		return x.Flags
	}
	return nil
}

type isRaftStorageCommand_Command interface {
	isRaftStorageCommand_Command()
}

type RaftStorageCommand_Delete struct {
	Delete *RaftStorageDelete `protobuf:"bytes,1,opt,name=delete,proto3,oneof"`
}

type RaftStorageCommand_Insert struct {
	Insert *RaftStorageInsert `protobuf:"bytes,2,opt,name=insert,proto3,oneof"`
}

type RaftStorageCommand_Flags struct {
	Flags *RaftStorageFlags `protobuf:"bytes,3,opt,name=flags,proto3,oneof"`
}

func (*RaftStorageCommand_Delete) isRaftStorageCommand_Command() {}

func (*RaftStorageCommand_Insert) isRaftStorageCommand_Command() {}

func (*RaftStorageCommand_Flags) isRaftStorageCommand_Command() {}

type RaftStorageDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RaftStorageDelete) Reset() {
	*x = RaftStorageDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftStorageDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftStorageDelete) ProtoMessage() {}

func (x *RaftStorageDelete) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftStorageDelete.ProtoReflect.Descriptor instead.
func (*RaftStorageDelete) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{7}
}

func (x *RaftStorageDelete) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type RaftStorageInsert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Priority uint32 `protobuf:"varint,2,opt,name=priority,proto3" json:"priority,omitempty"`
	Flags    uint32 `protobuf:"varint,3,opt,name=flags,proto3" json:"flags,omitempty"`
	Hash     uint64 `protobuf:"varint,4,opt,name=hash,proto3" json:"hash,omitempty"`
	Key      []byte `protobuf:"bytes,5,opt,name=key,proto3" json:"key,omitempty"`
	Value    []byte `protobuf:"bytes,6,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *RaftStorageInsert) Reset() {
	*x = RaftStorageInsert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftStorageInsert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftStorageInsert) ProtoMessage() {}

func (x *RaftStorageInsert) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftStorageInsert.ProtoReflect.Descriptor instead.
func (*RaftStorageInsert) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{8}
}

func (x *RaftStorageInsert) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RaftStorageInsert) GetPriority() uint32 {
	if x != nil {
		return x.Priority
	}
	return 0
}

func (x *RaftStorageInsert) GetFlags() uint32 {
	if x != nil {
		return x.Flags
	}
	return 0
}

func (x *RaftStorageInsert) GetHash() uint64 {
	if x != nil {
		return x.Hash
	}
	return 0
}

func (x *RaftStorageInsert) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *RaftStorageInsert) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type RaftStorageFlags struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Flags uint32 `protobuf:"varint,2,opt,name=flags,proto3" json:"flags,omitempty"`
	Mask  uint32 `protobuf:"varint,3,opt,name=mask,proto3" json:"mask,omitempty"`
}

func (x *RaftStorageFlags) Reset() {
	*x = RaftStorageFlags{}
	if protoimpl.UnsafeEnabled {
		mi := &file_raftcommands_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RaftStorageFlags) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RaftStorageFlags) ProtoMessage() {}

func (x *RaftStorageFlags) ProtoReflect() protoreflect.Message {
	mi := &file_raftcommands_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RaftStorageFlags.ProtoReflect.Descriptor instead.
func (*RaftStorageFlags) Descriptor() ([]byte, []int) {
	return file_raftcommands_proto_rawDescGZIP(), []int{9}
}

func (x *RaftStorageFlags) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *RaftStorageFlags) GetFlags() uint32 {
	if x != nil {
		return x.Flags
	}
	return 0
}

func (x *RaftStorageFlags) GetMask() uint32 {
	if x != nil {
		return x.Mask
	}
	return 0
}

var File_raftcommands_proto protoreflect.FileDescriptor

var file_raftcommands_proto_rawDesc = []byte{
	0x0a, 0x12, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63,
	0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x72, 0x61, 0x66,
	0x74, 0x63, 0x6d, 0x64, 0x1a, 0x0e, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbc, 0x04, 0x0a, 0x0b, 0x52, 0x61, 0x66, 0x74, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x53, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2f, 0x2e, 0x69, 0x6f, 0x2e, 0x63,
	0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x67, 0x72,
	0x70, 0x63, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61, 0x66, 0x74, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0b, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x52, 0x0a, 0x0b, 0x64, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e,
	0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73,
	0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x48, 0x00,
	0x52, 0x0b, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x12, 0x62, 0x0a,
	0x10, 0x70, 0x61, 0x75, 0x73, 0x65, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62,
	0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x67, 0x72, 0x70, 0x63,
	0x2e, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61, 0x66, 0x74, 0x50, 0x61, 0x75,
	0x73, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x48, 0x00, 0x52,
	0x10, 0x70, 0x61, 0x75, 0x73, 0x65, 0x72, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x71, 0x75, 0x65, 0x75,
	0x65, 0x12, 0x50, 0x0a, 0x0a, 0x70, 0x75, 0x72, 0x67, 0x65, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e,
	0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x72,
	0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61, 0x66, 0x74, 0x50, 0x75, 0x72, 0x67, 0x65,
	0x51, 0x75, 0x65, 0x75, 0x65, 0x48, 0x00, 0x52, 0x0a, 0x70, 0x75, 0x72, 0x67, 0x65, 0x71, 0x75,
	0x65, 0x75, 0x65, 0x12, 0x5f, 0x0a, 0x0f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x69,
	0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61,
	0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x73, 0x48, 0x00, 0x52, 0x0f, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x73, 0x12, 0x62, 0x0a, 0x10, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x69, 0x64, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x34,
	0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73,
	0x6d, 0x71, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e,
	0x52, 0x61, 0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74,
	0x65, 0x49, 0x64, 0x73, 0x48, 0x00, 0x52, 0x10, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x69, 0x64, 0x73, 0x42, 0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x22, 0x71, 0x0a, 0x0f, 0x52, 0x61, 0x66, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x71, 0x75, 0x65, 0x75, 0x65, 0x69, 0x64,
	0x12, 0x44, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2e, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75,
	0x73, 0x6d, 0x71, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52,
	0x05, 0x71, 0x75, 0x65, 0x75, 0x65, 0x22, 0x4a, 0x0a, 0x14, 0x52, 0x61, 0x66, 0x74, 0x50, 0x61,
	0x75, 0x73, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6d, 0x65, 0x51, 0x75, 0x65, 0x75, 0x65, 0x12, 0x1c,
	0x0a, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x70, 0x61, 0x75, 0x73, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x70, 0x61, 0x75,
	0x73, 0x65, 0x22, 0x40, 0x0a, 0x0e, 0x52, 0x61, 0x66, 0x74, 0x50, 0x75, 0x72, 0x67, 0x65, 0x51,
	0x75, 0x65, 0x75, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x6c, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x03, 0x61, 0x6c, 0x6c, 0x12, 0x1c, 0x0a, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x09, 0x71, 0x75, 0x65, 0x75, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x22, 0x93, 0x01, 0x0a, 0x13, 0x52, 0x61, 0x66, 0x74, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x12, 0x1c, 0x0a, 0x09,
	0x71, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x71, 0x75, 0x65, 0x75, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x4e, 0x0a, 0x08, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x32, 0x2e, 0x69,
	0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71,
	0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61,
	0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x73, 0x22, 0x28, 0x0a, 0x14, 0x52, 0x61,
	0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x49,
	0x64, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52,
	0x03, 0x69, 0x64, 0x73, 0x22, 0x83, 0x02, 0x0a, 0x12, 0x52, 0x61, 0x66, 0x74, 0x53, 0x74, 0x6f,
	0x72, 0x61, 0x67, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x4b, 0x0a, 0x06, 0x64,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x69, 0x6f,
	0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e,
	0x67, 0x72, 0x70, 0x63, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61, 0x66,
	0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x48, 0x00,
	0x52, 0x06, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x4b, 0x0a, 0x06, 0x69, 0x6e, 0x73, 0x65,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c,
	0x62, 0x73, 0x2e, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x67, 0x72, 0x70,
	0x63, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61, 0x66, 0x74, 0x53, 0x74,
	0x6f, 0x72, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x48, 0x00, 0x52, 0x06, 0x69,
	0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x48, 0x0a, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x69, 0x6f, 0x2e, 0x63, 0x6c, 0x62, 0x73, 0x2e, 0x6f,
	0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d, 0x71, 0x2e, 0x67, 0x72, 0x70, 0x63, 0x2e, 0x72, 0x61,
	0x66, 0x74, 0x63, 0x6d, 0x64, 0x2e, 0x52, 0x61, 0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x48, 0x00, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x42,
	0x09, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x23, 0x0a, 0x11, 0x52, 0x61,
	0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x22,
	0x91, 0x01, 0x0a, 0x11, 0x52, 0x61, 0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x49,
	0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x70, 0x72, 0x69, 0x6f, 0x72, 0x69, 0x74,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x61, 0x73, 0x68, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x68, 0x61, 0x73, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0x4c, 0x0a, 0x10, 0x52, 0x61, 0x66, 0x74, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x46, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x12, 0x0a,
	0x04, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x6d, 0x61, 0x73,
	0x6b, 0x42, 0x30, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x63, 0x6c, 0x62, 0x73, 0x2d, 0x69, 0x6f, 0x2f, 0x6f, 0x63, 0x74, 0x6f, 0x70, 0x75, 0x73, 0x6d,
	0x71, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x72, 0x61, 0x66, 0x74, 0x63, 0x6d,
	0x64, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_raftcommands_proto_rawDescOnce sync.Once
	file_raftcommands_proto_rawDescData = file_raftcommands_proto_rawDesc
)

func file_raftcommands_proto_rawDescGZIP() []byte {
	file_raftcommands_proto_rawDescOnce.Do(func() {
		file_raftcommands_proto_rawDescData = protoimpl.X.CompressGZIP(file_raftcommands_proto_rawDescData)
	})
	return file_raftcommands_proto_rawDescData
}

var file_raftcommands_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_raftcommands_proto_goTypes = []any{
	(*RaftCommand)(nil),                 // 0: io.clbs.octopusmq.grpc.raftcmd.RaftCommand
	(*RaftCreateQueue)(nil),             // 1: io.clbs.octopusmq.grpc.raftcmd.RaftCreateQueue
	(*RaftPauseResumeQueue)(nil),        // 2: io.clbs.octopusmq.grpc.raftcmd.RaftPauseResumeQueue
	(*RaftPurgeQueue)(nil),              // 3: io.clbs.octopusmq.grpc.raftcmd.RaftPurgeQueue
	(*RaftStorageCommands)(nil),         // 4: io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommands
	(*RaftStorageDeleteIds)(nil),        // 5: io.clbs.octopusmq.grpc.raftcmd.RaftStorageDeleteIds
	(*RaftStorageCommand)(nil),          // 6: io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommand
	(*RaftStorageDelete)(nil),           // 7: io.clbs.octopusmq.grpc.raftcmd.RaftStorageDelete
	(*RaftStorageInsert)(nil),           // 8: io.clbs.octopusmq.grpc.raftcmd.RaftStorageInsert
	(*RaftStorageFlags)(nil),            // 9: io.clbs.octopusmq.grpc.raftcmd.RaftStorageFlags
	(*protobuf.DeleteQueueRequest)(nil), // 10: io.clbs.octopusmq.protobuf.DeleteQueueRequest
	(*protobuf.CreateQueueRequest)(nil), // 11: io.clbs.octopusmq.protobuf.CreateQueueRequest
}
var file_raftcommands_proto_depIdxs = []int32{
	1,  // 0: io.clbs.octopusmq.grpc.raftcmd.RaftCommand.createqueue:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftCreateQueue
	10, // 1: io.clbs.octopusmq.grpc.raftcmd.RaftCommand.deletequeue:type_name -> io.clbs.octopusmq.protobuf.DeleteQueueRequest
	2,  // 2: io.clbs.octopusmq.grpc.raftcmd.RaftCommand.pauseresumequeue:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftPauseResumeQueue
	3,  // 3: io.clbs.octopusmq.grpc.raftcmd.RaftCommand.purgequeue:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftPurgeQueue
	4,  // 4: io.clbs.octopusmq.grpc.raftcmd.RaftCommand.storagecommands:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommands
	5,  // 5: io.clbs.octopusmq.grpc.raftcmd.RaftCommand.storagedeleteids:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftStorageDeleteIds
	11, // 6: io.clbs.octopusmq.grpc.raftcmd.RaftCreateQueue.queue:type_name -> io.clbs.octopusmq.protobuf.CreateQueueRequest
	6,  // 7: io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommands.commands:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommand
	7,  // 8: io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommand.delete:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftStorageDelete
	8,  // 9: io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommand.insert:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftStorageInsert
	9,  // 10: io.clbs.octopusmq.grpc.raftcmd.RaftStorageCommand.flags:type_name -> io.clbs.octopusmq.grpc.raftcmd.RaftStorageFlags
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_raftcommands_proto_init() }
func file_raftcommands_proto_init() {
	if File_raftcommands_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_raftcommands_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*RaftCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*RaftCreateQueue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*RaftPauseResumeQueue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*RaftPurgeQueue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*RaftStorageCommands); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*RaftStorageDeleteIds); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*RaftStorageCommand); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*RaftStorageDelete); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*RaftStorageInsert); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_raftcommands_proto_msgTypes[9].Exporter = func(v any, i int) any {
			switch v := v.(*RaftStorageFlags); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_raftcommands_proto_msgTypes[0].OneofWrappers = []any{
		(*RaftCommand_Createqueue)(nil),
		(*RaftCommand_Deletequeue)(nil),
		(*RaftCommand_Pauseresumequeue)(nil),
		(*RaftCommand_Purgequeue)(nil),
		(*RaftCommand_Storagecommands)(nil),
		(*RaftCommand_Storagedeleteids)(nil),
	}
	file_raftcommands_proto_msgTypes[6].OneofWrappers = []any{
		(*RaftStorageCommand_Delete)(nil),
		(*RaftStorageCommand_Insert)(nil),
		(*RaftStorageCommand_Flags)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_raftcommands_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_raftcommands_proto_goTypes,
		DependencyIndexes: file_raftcommands_proto_depIdxs,
		MessageInfos:      file_raftcommands_proto_msgTypes,
	}.Build()
	File_raftcommands_proto = out.File
	file_raftcommands_proto_rawDesc = nil
	file_raftcommands_proto_goTypes = nil
	file_raftcommands_proto_depIdxs = nil
}
