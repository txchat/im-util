// protoc -I=. -I=$GOPATH/src --go_out=plugins=grpc:. *.proto

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.1
// source: device.proto

package device

import (
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

type ConnectReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConnectReq) Reset() {
	*x = ConnectReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectReq) ProtoMessage() {}

func (x *ConnectReq) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectReq.ProtoReflect.Descriptor instead.
func (*ConnectReq) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{0}
}

type ConnectReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ConnectReply) Reset() {
	*x = ConnectReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectReply) ProtoMessage() {}

func (x *ConnectReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectReply.ProtoReflect.Descriptor instead.
func (*ConnectReply) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{1}
}

type ReConnectReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReConnectReq) Reset() {
	*x = ReConnectReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReConnectReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReConnectReq) ProtoMessage() {}

func (x *ReConnectReq) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReConnectReq.ProtoReflect.Descriptor instead.
func (*ReConnectReq) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{2}
}

type ReConnectReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *ReConnectReply) Reset() {
	*x = ReConnectReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReConnectReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReConnectReply) ProtoMessage() {}

func (x *ReConnectReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReConnectReply.ProtoReflect.Descriptor instead.
func (*ReConnectReply) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{3}
}

type DisConnectReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DisConnectReq) Reset() {
	*x = DisConnectReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DisConnectReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DisConnectReq) ProtoMessage() {}

func (x *DisConnectReq) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DisConnectReq.ProtoReflect.Descriptor instead.
func (*DisConnectReq) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{4}
}

type DisConnectReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *DisConnectReply) Reset() {
	*x = DisConnectReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DisConnectReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DisConnectReply) ProtoMessage() {}

func (x *DisConnectReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DisConnectReply.ProtoReflect.Descriptor instead.
func (*DisConnectReply) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{5}
}

type FocusReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelType int32  `protobuf:"varint,1,opt,name=channelType,proto3" json:"channelType,omitempty"`
	Target      string `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
}

func (x *FocusReq) Reset() {
	*x = FocusReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FocusReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FocusReq) ProtoMessage() {}

func (x *FocusReq) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FocusReq.ProtoReflect.Descriptor instead.
func (*FocusReq) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{6}
}

func (x *FocusReq) GetChannelType() int32 {
	if x != nil {
		return x.ChannelType
	}
	return 0
}

func (x *FocusReq) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

type FocusReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *FocusReply) Reset() {
	*x = FocusReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FocusReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FocusReply) ProtoMessage() {}

func (x *FocusReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FocusReply.ProtoReflect.Descriptor instead.
func (*FocusReply) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{7}
}

type InputReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Text string `protobuf:"bytes,1,opt,name=text,proto3" json:"text,omitempty"`
}

func (x *InputReq) Reset() {
	*x = InputReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputReq) ProtoMessage() {}

func (x *InputReq) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputReq.ProtoReflect.Descriptor instead.
func (*InputReq) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{8}
}

func (x *InputReq) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

type InputReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err string `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *InputReply) Reset() {
	*x = InputReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InputReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InputReply) ProtoMessage() {}

func (x *InputReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InputReply.ProtoReflect.Descriptor instead.
func (*InputReply) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{9}
}

func (x *InputReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type OutputReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OutputReq) Reset() {
	*x = OutputReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutputReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutputReq) ProtoMessage() {}

func (x *OutputReq) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutputReq.ProtoReflect.Descriptor instead.
func (*OutputReq) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{10}
}

type OutputReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg         string `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	ChannelType int32  `protobuf:"varint,2,opt,name=channelType,proto3" json:"channelType,omitempty"`
	Target      string `protobuf:"bytes,3,opt,name=target,proto3" json:"target,omitempty"`
}

func (x *OutputReply) Reset() {
	*x = OutputReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_device_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutputReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutputReply) ProtoMessage() {}

func (x *OutputReply) ProtoReflect() protoreflect.Message {
	mi := &file_device_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutputReply.ProtoReflect.Descriptor instead.
func (*OutputReply) Descriptor() ([]byte, []int) {
	return file_device_proto_rawDescGZIP(), []int{11}
}

func (x *OutputReply) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *OutputReply) GetChannelType() int32 {
	if x != nil {
		return x.ChannelType
	}
	return 0
}

func (x *OutputReply) GetTarget() string {
	if x != nil {
		return x.Target
	}
	return ""
}

var File_device_proto protoreflect.FileDescriptor

var file_device_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x22, 0x0c, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x22, 0x0e, 0x0a, 0x0c, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x0e, 0x0a, 0x0c, 0x52, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x71, 0x22, 0x10, 0x0a, 0x0e, 0x52, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x0f, 0x0a, 0x0d, 0x44, 0x69, 0x73, 0x43, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x22, 0x11, 0x0a, 0x0f, 0x44, 0x69, 0x73, 0x43, 0x6f,
	0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x44, 0x0a, 0x08, 0x46, 0x6f,
	0x63, 0x75, 0x73, 0x52, 0x65, 0x71, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x68, 0x61,
	0x6e, 0x6e, 0x65, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x22, 0x0c, 0x0a, 0x0a, 0x46, 0x6f, 0x63, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x1e,
	0x0a, 0x08, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65,
	0x78, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x22, 0x1e,
	0x0a, 0x0a, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x22, 0x0b,
	0x0a, 0x09, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x71, 0x22, 0x59, 0x0a, 0x0b, 0x4f,
	0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73,
	0x67, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x20, 0x0a, 0x0b,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0b, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16,
	0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x32, 0xca, 0x02, 0x0a, 0x06, 0x44, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x33, 0x0a, 0x07, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12, 0x12, 0x2e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71,
	0x1a, 0x14, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x39, 0x0a, 0x09, 0x52, 0x65, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x12, 0x14, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x43,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x16, 0x2e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x52, 0x65, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x3c, 0x0a, 0x0a, 0x44, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x12,
	0x15, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x44, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x2e,
	0x44, 0x69, 0x73, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x2d, 0x0a, 0x05, 0x46, 0x6f, 0x63, 0x75, 0x73, 0x12, 0x10, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x2e, 0x46, 0x6f, 0x63, 0x75, 0x73, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x64, 0x65, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x46, 0x6f, 0x63, 0x75, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x2f,
	0x0a, 0x05, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x10, 0x2e, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65,
	0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x28, 0x01, 0x12,
	0x32, 0x0a, 0x06, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x11, 0x2e, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x13, 0x2e, 0x64,
	0x65, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x30, 0x01, 0x42, 0x2a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x78, 0x63, 0x68, 0x61, 0x74, 0x2f, 0x69, 0x6d, 0x2d, 0x75, 0x74, 0x69, 0x6c,
	0x2f, 0x76, 0x69, 0x72, 0x74, 0x75, 0x61, 0x6c, 0x2f, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_device_proto_rawDescOnce sync.Once
	file_device_proto_rawDescData = file_device_proto_rawDesc
)

func file_device_proto_rawDescGZIP() []byte {
	file_device_proto_rawDescOnce.Do(func() {
		file_device_proto_rawDescData = protoimpl.X.CompressGZIP(file_device_proto_rawDescData)
	})
	return file_device_proto_rawDescData
}

var file_device_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_device_proto_goTypes = []interface{}{
	(*ConnectReq)(nil),      // 0: device.ConnectReq
	(*ConnectReply)(nil),    // 1: device.ConnectReply
	(*ReConnectReq)(nil),    // 2: device.ReConnectReq
	(*ReConnectReply)(nil),  // 3: device.ReConnectReply
	(*DisConnectReq)(nil),   // 4: device.DisConnectReq
	(*DisConnectReply)(nil), // 5: device.DisConnectReply
	(*FocusReq)(nil),        // 6: device.FocusReq
	(*FocusReply)(nil),      // 7: device.FocusReply
	(*InputReq)(nil),        // 8: device.InputReq
	(*InputReply)(nil),      // 9: device.InputReply
	(*OutputReq)(nil),       // 10: device.OutputReq
	(*OutputReply)(nil),     // 11: device.OutputReply
}
var file_device_proto_depIdxs = []int32{
	0,  // 0: device.Device.Connect:input_type -> device.ConnectReq
	2,  // 1: device.Device.ReConnect:input_type -> device.ReConnectReq
	4,  // 2: device.Device.DisConnect:input_type -> device.DisConnectReq
	6,  // 3: device.Device.Focus:input_type -> device.FocusReq
	8,  // 4: device.Device.Input:input_type -> device.InputReq
	10, // 5: device.Device.Output:input_type -> device.OutputReq
	1,  // 6: device.Device.Connect:output_type -> device.ConnectReply
	3,  // 7: device.Device.ReConnect:output_type -> device.ReConnectReply
	5,  // 8: device.Device.DisConnect:output_type -> device.DisConnectReply
	7,  // 9: device.Device.Focus:output_type -> device.FocusReply
	9,  // 10: device.Device.Input:output_type -> device.InputReply
	11, // 11: device.Device.Output:output_type -> device.OutputReply
	6,  // [6:12] is the sub-list for method output_type
	0,  // [0:6] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_device_proto_init() }
func file_device_proto_init() {
	if File_device_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_device_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectReq); i {
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
		file_device_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectReply); i {
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
		file_device_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReConnectReq); i {
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
		file_device_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReConnectReply); i {
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
		file_device_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DisConnectReq); i {
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
		file_device_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DisConnectReply); i {
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
		file_device_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FocusReq); i {
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
		file_device_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FocusReply); i {
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
		file_device_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputReq); i {
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
		file_device_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InputReply); i {
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
		file_device_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutputReq); i {
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
		file_device_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutputReply); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_device_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_device_proto_goTypes,
		DependencyIndexes: file_device_proto_depIdxs,
		MessageInfos:      file_device_proto_msgTypes,
	}.Build()
	File_device_proto = out.File
	file_device_proto_rawDesc = nil
	file_device_proto_goTypes = nil
	file_device_proto_depIdxs = nil
}