// Code generated by protoc-gen-go.
// source: proto.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	proto.proto

It has these top-level messages:
	ShareInfo
	PersonInfo
	FlatInfo
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type ShareInfo struct {
	Id          int64  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	ShareType   int64  `protobuf:"varint,2,opt,name=share_type,json=shareType" json:"share_type,omitempty"`
	ShareWith   string `protobuf:"bytes,3,opt,name=share_with,json=shareWith" json:"share_with,omitempty"`
	UidOwner    string `protobuf:"bytes,4,opt,name=uid_owner,json=uidOwner" json:"uid_owner,omitempty"`
	Parent      int64  `protobuf:"varint,5,opt,name=parent" json:"parent,omitempty"`
	ItemType    string `protobuf:"bytes,6,opt,name=item_type,json=itemType" json:"item_type,omitempty"`
	ItemSource  int64  `protobuf:"varint,7,opt,name=item_source,json=itemSource" json:"item_source,omitempty"`
	ItemTarget  string `protobuf:"bytes,8,opt,name=item_target,json=itemTarget" json:"item_target,omitempty"`
	FileSource  int64  `protobuf:"varint,9,opt,name=file_source,json=fileSource" json:"file_source,omitempty"`
	FileTarget  string `protobuf:"bytes,10,opt,name=file_target,json=fileTarget" json:"file_target,omitempty"`
	Permissions string `protobuf:"bytes,11,opt,name=permissions" json:"permissions,omitempty"`
	Stime       int64  `protobuf:"varint,12,opt,name=stime" json:"stime,omitempty"`
	Accepted    int64  `protobuf:"varint,13,opt,name=accepted" json:"accepted,omitempty"`
	Expiration  string `protobuf:"bytes,14,opt,name=expiration" json:"expiration,omitempty"`
	Token       string `protobuf:"bytes,15,opt,name=token" json:"token,omitempty"`
	MailSend    int64  `protobuf:"varint,16,opt,name=mail_send,json=mailSend" json:"mail_send,omitempty"`
}

func (m *ShareInfo) Reset()                    { *m = ShareInfo{} }
func (m *ShareInfo) String() string            { return proto1.CompactTextString(m) }
func (*ShareInfo) ProtoMessage()               {}
func (*ShareInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type PersonInfo struct {
	Login        string `protobuf:"bytes,1,opt,name=login" json:"login,omitempty"`
	Uid          int64  `protobuf:"varint,2,opt,name=uid" json:"uid,omitempty"`
	Department   string `protobuf:"bytes,3,opt,name=department" json:"department,omitempty"`
	Group        string `protobuf:"bytes,4,opt,name=group" json:"group,omitempty"`
	Organization string `protobuf:"bytes,5,opt,name=organization" json:"organization,omitempty"`
	Company      string `protobuf:"bytes,6,opt,name=company" json:"company,omitempty"`
	Office       string `protobuf:"bytes,7,opt,name=office" json:"office,omitempty"`
}

func (m *PersonInfo) Reset()                    { *m = PersonInfo{} }
func (m *PersonInfo) String() string            { return proto1.CompactTextString(m) }
func (*PersonInfo) ProtoMessage()               {}
func (*PersonInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type FlatInfo struct {
	ShareInfo  *ShareInfo  `protobuf:"bytes,1,opt,name=shareInfo" json:"shareInfo,omitempty"`
	OwnerInfo  *PersonInfo `protobuf:"bytes,2,opt,name=ownerInfo" json:"ownerInfo,omitempty"`
	ShareeInfo *PersonInfo `protobuf:"bytes,3,opt,name=shareeInfo" json:"shareeInfo,omitempty"`
}

func (m *FlatInfo) Reset()                    { *m = FlatInfo{} }
func (m *FlatInfo) String() string            { return proto1.CompactTextString(m) }
func (*FlatInfo) ProtoMessage()               {}
func (*FlatInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *FlatInfo) GetShareInfo() *ShareInfo {
	if m != nil {
		return m.ShareInfo
	}
	return nil
}

func (m *FlatInfo) GetOwnerInfo() *PersonInfo {
	if m != nil {
		return m.OwnerInfo
	}
	return nil
}

func (m *FlatInfo) GetShareeInfo() *PersonInfo {
	if m != nil {
		return m.ShareeInfo
	}
	return nil
}

func init() {
	proto1.RegisterType((*ShareInfo)(nil), "proto.ShareInfo")
	proto1.RegisterType((*PersonInfo)(nil), "proto.PersonInfo")
	proto1.RegisterType((*FlatInfo)(nil), "proto.FlatInfo")
}

func init() { proto1.RegisterFile("proto.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 452 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x52, 0x4d, 0x8e, 0x13, 0x3d,
	0x10, 0x55, 0xa7, 0xbf, 0xce, 0xa4, 0xab, 0xe7, 0x1b, 0x82, 0x35, 0x42, 0x16, 0x08, 0x88, 0xb2,
	0x9a, 0x55, 0x10, 0x70, 0x07, 0x24, 0x56, 0xa0, 0xce, 0x48, 0x2c, 0x23, 0x13, 0x57, 0x92, 0x12,
	0x69, 0xdb, 0x72, 0xbb, 0x35, 0x84, 0xc3, 0x70, 0x16, 0x4e, 0xc1, 0x79, 0x90, 0xcb, 0x9e, 0x74,
	0x58, 0xb0, 0x49, 0xf2, 0x7e, 0xea, 0x39, 0xf6, 0x2b, 0x68, 0x9c, 0xb7, 0xc1, 0xae, 0xf8, 0x53,
	0x54, 0xfc, 0xb5, 0xfc, 0x5d, 0x42, 0xbd, 0x3e, 0x28, 0x8f, 0x1f, 0xcd, 0xce, 0x8a, 0x1b, 0x98,
	0x90, 0x96, 0xc5, 0xa2, 0xb8, 0x2b, 0xdb, 0x09, 0x69, 0xf1, 0x12, 0xa0, 0x8f, 0xe2, 0x26, 0x9c,
	0x1c, 0xca, 0x09, 0xf3, 0x35, 0x33, 0xf7, 0x27, 0x87, 0xa3, 0xfc, 0x40, 0xe1, 0x20, 0xcb, 0x45,
	0x71, 0x57, 0x67, 0xf9, 0x0b, 0x85, 0x83, 0x78, 0x01, 0xf5, 0x40, 0x7a, 0x63, 0x1f, 0x0c, 0x7a,
	0xf9, 0x1f, 0xab, 0xb3, 0x81, 0xf4, 0xa7, 0x88, 0xc5, 0x33, 0x98, 0x3a, 0xe5, 0xd1, 0x04, 0x59,
	0x71, 0x6c, 0x46, 0x71, 0x88, 0x02, 0x76, 0xe9, 0xc4, 0x69, 0x1a, 0x8a, 0x04, 0x1f, 0xf8, 0x1a,
	0x1a, 0x16, 0x7b, 0x3b, 0xf8, 0x2d, 0xca, 0x2b, 0x9e, 0x84, 0x48, 0xad, 0x99, 0x39, 0x1b, 0x82,
	0xf2, 0x7b, 0x0c, 0x72, 0xc6, 0xf3, 0x6c, 0xb8, 0x67, 0x26, 0x1a, 0x76, 0x74, 0xc4, 0xc7, 0x84,
	0x3a, 0x25, 0x44, 0x6a, 0x4c, 0x60, 0x43, 0x4e, 0x80, 0x94, 0x10, 0xa9, 0x9c, 0xb0, 0x80, 0xc6,
	0xa1, 0xef, 0xa8, 0xef, 0xc9, 0x9a, 0x5e, 0x36, 0x6c, 0xb8, 0xa4, 0xc4, 0x2d, 0x54, 0x7d, 0xa0,
	0x0e, 0xe5, 0x35, 0xa7, 0x27, 0x20, 0x9e, 0xc3, 0x4c, 0x6d, 0xb7, 0xe8, 0x02, 0x6a, 0xf9, 0x3f,
	0x0b, 0x67, 0x2c, 0x5e, 0x01, 0xe0, 0x77, 0x47, 0x5e, 0x05, 0xb2, 0x46, 0xde, 0xa4, 0x33, 0x47,
	0x26, 0x26, 0x06, 0xfb, 0x0d, 0x8d, 0x7c, 0xc2, 0x52, 0x02, 0xf1, 0xa9, 0x3a, 0x45, 0xc7, 0x4d,
	0x8f, 0x46, 0xcb, 0x79, 0x8a, 0x8c, 0xc4, 0x1a, 0x8d, 0x5e, 0xfe, 0x2a, 0x00, 0x3e, 0xa3, 0xef,
	0xad, 0xe1, 0x66, 0x6f, 0xa1, 0x3a, 0xda, 0x3d, 0x19, 0x2e, 0xb7, 0x6e, 0x13, 0x10, 0x73, 0x28,
	0x07, 0xd2, 0xb9, 0xd8, 0xf8, 0x33, 0xfe, 0x13, 0x8d, 0x4e, 0xf9, 0xd0, 0xc5, 0x6a, 0x52, 0xa5,
	0x17, 0x4c, 0xcc, 0xd9, 0x7b, 0x3b, 0xb8, 0xdc, 0x67, 0x02, 0x62, 0x09, 0xd7, 0xd6, 0xef, 0x95,
	0xa1, 0x1f, 0xe9, 0x06, 0x15, 0x8b, 0x7f, 0x71, 0x42, 0xc2, 0xd5, 0xd6, 0x76, 0x4e, 0x99, 0x53,
	0xae, 0xf5, 0x11, 0xc6, 0x55, 0xb0, 0xbb, 0x1d, 0xe5, 0x42, 0xeb, 0x36, 0xa3, 0xe5, 0xcf, 0x02,
	0x66, 0x1f, 0x8e, 0x2a, 0xf0, 0x05, 0x56, 0x90, 0x36, 0x2b, 0x02, 0xbe, 0x44, 0xf3, 0x6e, 0x9e,
	0x56, 0x79, 0x75, 0xde, 0xdf, 0x76, 0xb4, 0x88, 0x37, 0x50, 0xf3, 0xe2, 0xb1, 0x7f, 0xc2, 0xfe,
	0xa7, 0xd9, 0x3f, 0x3e, 0x4b, 0x3b, 0x7a, 0xc4, 0xdb, 0xbc, 0xcc, 0xe9, 0x84, 0xf2, 0x5f, 0x13,
	0x17, 0xa6, 0xaf, 0x53, 0x56, 0xdf, 0xff, 0x09, 0x00, 0x00, 0xff, 0xff, 0xa8, 0x05, 0x88, 0x58,
	0x59, 0x03, 0x00, 0x00,
}
