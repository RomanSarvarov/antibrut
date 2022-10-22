// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: antibrut/v1/antibrut.proto

package antibrut

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Attempt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Attempt) Reset() {
	*x = Attempt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Attempt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Attempt) ProtoMessage() {}

func (x *Attempt) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Attempt.ProtoReflect.Descriptor instead.
func (*Attempt) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{0}
}

func (x *Attempt) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Attempt) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type TryRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attempt *Attempt `protobuf:"bytes,1,opt,name=attempt,proto3" json:"attempt,omitempty"`
}

func (x *TryRequest) Reset() {
	*x = TryRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TryRequest) ProtoMessage() {}

func (x *TryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TryRequest.ProtoReflect.Descriptor instead.
func (*TryRequest) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{1}
}

func (x *TryRequest) GetAttempt() *Attempt {
	if x != nil {
		return x.Attempt
	}
	return nil
}

type TryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok bool `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
}

func (x *TryResponse) Reset() {
	*x = TryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TryResponse) ProtoMessage() {}

func (x *TryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TryResponse.ProtoReflect.Descriptor instead.
func (*TryResponse) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{2}
}

func (x *TryResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

type ResetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attempt *Attempt `protobuf:"bytes,1,opt,name=attempt,proto3" json:"attempt,omitempty"`
}

func (x *ResetRequest) Reset() {
	*x = ResetRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResetRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetRequest) ProtoMessage() {}

func (x *ResetRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetRequest.ProtoReflect.Descriptor instead.
func (*ResetRequest) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{3}
}

func (x *ResetRequest) GetAttempt() *Attempt {
	if x != nil {
		return x.Attempt
	}
	return nil
}

type AddToWhiteListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attempt *Attempt `protobuf:"bytes,1,opt,name=attempt,proto3" json:"attempt,omitempty"`
}

func (x *AddToWhiteListRequest) Reset() {
	*x = AddToWhiteListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddToWhiteListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddToWhiteListRequest) ProtoMessage() {}

func (x *AddToWhiteListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddToWhiteListRequest.ProtoReflect.Descriptor instead.
func (*AddToWhiteListRequest) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{4}
}

func (x *AddToWhiteListRequest) GetAttempt() *Attempt {
	if x != nil {
		return x.Attempt
	}
	return nil
}

type DeleteFromWhiteListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attempt *Attempt `protobuf:"bytes,1,opt,name=attempt,proto3" json:"attempt,omitempty"`
}

func (x *DeleteFromWhiteListRequest) Reset() {
	*x = DeleteFromWhiteListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFromWhiteListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFromWhiteListRequest) ProtoMessage() {}

func (x *DeleteFromWhiteListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFromWhiteListRequest.ProtoReflect.Descriptor instead.
func (*DeleteFromWhiteListRequest) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{5}
}

func (x *DeleteFromWhiteListRequest) GetAttempt() *Attempt {
	if x != nil {
		return x.Attempt
	}
	return nil
}

type AddToBlackListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attempt *Attempt `protobuf:"bytes,1,opt,name=attempt,proto3" json:"attempt,omitempty"`
}

func (x *AddToBlackListRequest) Reset() {
	*x = AddToBlackListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddToBlackListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddToBlackListRequest) ProtoMessage() {}

func (x *AddToBlackListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddToBlackListRequest.ProtoReflect.Descriptor instead.
func (*AddToBlackListRequest) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{6}
}

func (x *AddToBlackListRequest) GetAttempt() *Attempt {
	if x != nil {
		return x.Attempt
	}
	return nil
}

type DeleteFromBlackListRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attempt *Attempt `protobuf:"bytes,1,opt,name=attempt,proto3" json:"attempt,omitempty"`
}

func (x *DeleteFromBlackListRequest) Reset() {
	*x = DeleteFromBlackListRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_antibrut_v1_antibrut_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteFromBlackListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteFromBlackListRequest) ProtoMessage() {}

func (x *DeleteFromBlackListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_antibrut_v1_antibrut_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteFromBlackListRequest.ProtoReflect.Descriptor instead.
func (*DeleteFromBlackListRequest) Descriptor() ([]byte, []int) {
	return file_antibrut_v1_antibrut_proto_rawDescGZIP(), []int{7}
}

func (x *DeleteFromBlackListRequest) GetAttempt() *Attempt {
	if x != nil {
		return x.Attempt
	}
	return nil
}

var File_antibrut_v1_antibrut_proto protoreflect.FileDescriptor

var file_antibrut_v1_antibrut_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x6e,
	0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x61, 0x6e,
	0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x33, 0x0a, 0x07, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x3c, 0x0a, 0x0a, 0x54,
	0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x61, 0x74, 0x74,
	0x65, 0x6d, 0x70, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x6e, 0x74,
	0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74,
	0x52, 0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x22, 0x1d, 0x0a, 0x0b, 0x54, 0x72, 0x79,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x22, 0x3e, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x61, 0x74, 0x74, 0x65,
	0x6d, 0x70, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x6e, 0x74, 0x69,
	0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52,
	0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x22, 0x47, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x54,
	0x6f, 0x57, 0x68, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x2e, 0x0a, 0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52, 0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x22, 0x4c, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x57,
	0x68, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x2e, 0x0a, 0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x14, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52, 0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x22,
	0x47, 0x0a, 0x15, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x4c, 0x69, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x61, 0x74, 0x74, 0x65,
	0x6d, 0x70, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x6e, 0x74, 0x69,
	0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52,
	0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x22, 0x4c, 0x0a, 0x1a, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x61, 0x74, 0x74, 0x65, 0x6d, 0x70,
	0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72,
	0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x52, 0x07, 0x61,
	0x74, 0x74, 0x65, 0x6d, 0x70, 0x74, 0x32, 0xd3, 0x03, 0x0a, 0x0f, 0x41, 0x6e, 0x74, 0x69, 0x42,
	0x72, 0x75, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x38, 0x0a, 0x03, 0x54, 0x72,
	0x79, 0x12, 0x17, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e,
	0x54, 0x72, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x6e, 0x74,
	0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x05, 0x52, 0x65, 0x73, 0x65, 0x74, 0x12, 0x19, 0x2e,
	0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x65,
	0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x12, 0x4c, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x57, 0x68, 0x69, 0x74, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x22, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x57, 0x68, 0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x56,
	0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x57, 0x68, 0x69, 0x74,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x27, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74,
	0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x72, 0x6f, 0x6d, 0x57, 0x68,
	0x69, 0x74, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x4c, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x42,
	0x6c, 0x61, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x22, 0x2e, 0x61, 0x6e, 0x74, 0x69, 0x62,
	0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x42, 0x6c, 0x61, 0x63,
	0x6b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x12, 0x56, 0x0a, 0x13, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x46, 0x72,
	0x6f, 0x6d, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x27, 0x2e, 0x61, 0x6e,
	0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x46, 0x72, 0x6f, 0x6d, 0x42, 0x6c, 0x61, 0x63, 0x6b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x0d, 0x5a, 0x0b,
	0x2e, 0x2f, 0x3b, 0x61, 0x6e, 0x74, 0x69, 0x62, 0x72, 0x75, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_antibrut_v1_antibrut_proto_rawDescOnce sync.Once
	file_antibrut_v1_antibrut_proto_rawDescData = file_antibrut_v1_antibrut_proto_rawDesc
)

func file_antibrut_v1_antibrut_proto_rawDescGZIP() []byte {
	file_antibrut_v1_antibrut_proto_rawDescOnce.Do(func() {
		file_antibrut_v1_antibrut_proto_rawDescData = protoimpl.X.CompressGZIP(file_antibrut_v1_antibrut_proto_rawDescData)
	})
	return file_antibrut_v1_antibrut_proto_rawDescData
}

var file_antibrut_v1_antibrut_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_antibrut_v1_antibrut_proto_goTypes = []interface{}{
	(*Attempt)(nil),                    // 0: antibrut.v1.Attempt
	(*TryRequest)(nil),                 // 1: antibrut.v1.TryRequest
	(*TryResponse)(nil),                // 2: antibrut.v1.TryResponse
	(*ResetRequest)(nil),               // 3: antibrut.v1.ResetRequest
	(*AddToWhiteListRequest)(nil),      // 4: antibrut.v1.AddToWhiteListRequest
	(*DeleteFromWhiteListRequest)(nil), // 5: antibrut.v1.DeleteFromWhiteListRequest
	(*AddToBlackListRequest)(nil),      // 6: antibrut.v1.AddToBlackListRequest
	(*DeleteFromBlackListRequest)(nil), // 7: antibrut.v1.DeleteFromBlackListRequest
	(*emptypb.Empty)(nil),              // 8: google.protobuf.Empty
}
var file_antibrut_v1_antibrut_proto_depIdxs = []int32{
	0,  // 0: antibrut.v1.TryRequest.attempt:type_name -> antibrut.v1.Attempt
	0,  // 1: antibrut.v1.ResetRequest.attempt:type_name -> antibrut.v1.Attempt
	0,  // 2: antibrut.v1.AddToWhiteListRequest.attempt:type_name -> antibrut.v1.Attempt
	0,  // 3: antibrut.v1.DeleteFromWhiteListRequest.attempt:type_name -> antibrut.v1.Attempt
	0,  // 4: antibrut.v1.AddToBlackListRequest.attempt:type_name -> antibrut.v1.Attempt
	0,  // 5: antibrut.v1.DeleteFromBlackListRequest.attempt:type_name -> antibrut.v1.Attempt
	1,  // 6: antibrut.v1.AntiBrutService.Try:input_type -> antibrut.v1.TryRequest
	3,  // 7: antibrut.v1.AntiBrutService.Reset:input_type -> antibrut.v1.ResetRequest
	4,  // 8: antibrut.v1.AntiBrutService.AddToWhiteList:input_type -> antibrut.v1.AddToWhiteListRequest
	5,  // 9: antibrut.v1.AntiBrutService.DeleteFromWhiteList:input_type -> antibrut.v1.DeleteFromWhiteListRequest
	6,  // 10: antibrut.v1.AntiBrutService.AddToBlackList:input_type -> antibrut.v1.AddToBlackListRequest
	7,  // 11: antibrut.v1.AntiBrutService.DeleteFromBlackList:input_type -> antibrut.v1.DeleteFromBlackListRequest
	2,  // 12: antibrut.v1.AntiBrutService.Try:output_type -> antibrut.v1.TryResponse
	8,  // 13: antibrut.v1.AntiBrutService.Reset:output_type -> google.protobuf.Empty
	8,  // 14: antibrut.v1.AntiBrutService.AddToWhiteList:output_type -> google.protobuf.Empty
	8,  // 15: antibrut.v1.AntiBrutService.DeleteFromWhiteList:output_type -> google.protobuf.Empty
	8,  // 16: antibrut.v1.AntiBrutService.AddToBlackList:output_type -> google.protobuf.Empty
	8,  // 17: antibrut.v1.AntiBrutService.DeleteFromBlackList:output_type -> google.protobuf.Empty
	12, // [12:18] is the sub-list for method output_type
	6,  // [6:12] is the sub-list for method input_type
	6,  // [6:6] is the sub-list for extension type_name
	6,  // [6:6] is the sub-list for extension extendee
	0,  // [0:6] is the sub-list for field type_name
}

func init() { file_antibrut_v1_antibrut_proto_init() }
func file_antibrut_v1_antibrut_proto_init() {
	if File_antibrut_v1_antibrut_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_antibrut_v1_antibrut_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Attempt); i {
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
		file_antibrut_v1_antibrut_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TryRequest); i {
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
		file_antibrut_v1_antibrut_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TryResponse); i {
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
		file_antibrut_v1_antibrut_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResetRequest); i {
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
		file_antibrut_v1_antibrut_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddToWhiteListRequest); i {
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
		file_antibrut_v1_antibrut_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFromWhiteListRequest); i {
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
		file_antibrut_v1_antibrut_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddToBlackListRequest); i {
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
		file_antibrut_v1_antibrut_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteFromBlackListRequest); i {
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
			RawDescriptor: file_antibrut_v1_antibrut_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_antibrut_v1_antibrut_proto_goTypes,
		DependencyIndexes: file_antibrut_v1_antibrut_proto_depIdxs,
		MessageInfos:      file_antibrut_v1_antibrut_proto_msgTypes,
	}.Build()
	File_antibrut_v1_antibrut_proto = out.File
	file_antibrut_v1_antibrut_proto_rawDesc = nil
	file_antibrut_v1_antibrut_proto_goTypes = nil
	file_antibrut_v1_antibrut_proto_depIdxs = nil
}
