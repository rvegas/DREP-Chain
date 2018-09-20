// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bean/node.proto

package bean

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import mycrypto "BlockChainTest/mycrypto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type PeerInfo struct {
	Pk                   *mycrypto.Point `protobuf:"bytes,1,opt,name=pk,proto3" json:"pk,omitempty"`
	Ip                   string          `protobuf:"bytes,2,opt,name=ip,proto3" json:"ip,omitempty"`
	Port                 int32           `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *PeerInfo) Reset()         { *m = PeerInfo{} }
func (m *PeerInfo) String() string { return proto.CompactTextString(m) }
func (*PeerInfo) ProtoMessage()    {}
func (*PeerInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_node_7b7652e688675472, []int{0}
}
func (m *PeerInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerInfo.Unmarshal(m, b)
}
func (m *PeerInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerInfo.Marshal(b, m, deterministic)
}
func (dst *PeerInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerInfo.Merge(dst, src)
}
func (m *PeerInfo) XXX_Size() int {
	return xxx_messageInfo_PeerInfo.Size(m)
}
func (m *PeerInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerInfo.DiscardUnknown(m)
}

var xxx_messageInfo_PeerInfo proto.InternalMessageInfo

func (m *PeerInfo) GetPk() *mycrypto.Point {
	if m != nil {
		return m.Pk
	}
	return nil
}

func (m *PeerInfo) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *PeerInfo) GetPort() int32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type PeerInfoList struct {
	List                 []*PeerInfo `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PeerInfoList) Reset()         { *m = PeerInfoList{} }
func (m *PeerInfoList) String() string { return proto.CompactTextString(m) }
func (*PeerInfoList) ProtoMessage()    {}
func (*PeerInfoList) Descriptor() ([]byte, []int) {
	return fileDescriptor_node_7b7652e688675472, []int{1}
}
func (m *PeerInfoList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerInfoList.Unmarshal(m, b)
}
func (m *PeerInfoList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerInfoList.Marshal(b, m, deterministic)
}
func (dst *PeerInfoList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerInfoList.Merge(dst, src)
}
func (m *PeerInfoList) XXX_Size() int {
	return xxx_messageInfo_PeerInfoList.Size(m)
}
func (m *PeerInfoList) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerInfoList.DiscardUnknown(m)
}

var xxx_messageInfo_PeerInfoList proto.InternalMessageInfo

func (m *PeerInfoList) GetList() []*PeerInfo {
	if m != nil {
		return m.List
	}
	return nil
}

type BlockReq struct {
	Height               int64    `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockReq) Reset()         { *m = BlockReq{} }
func (m *BlockReq) String() string { return proto.CompactTextString(m) }
func (*BlockReq) ProtoMessage()    {}
func (*BlockReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_node_7b7652e688675472, []int{2}
}
func (m *BlockReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockReq.Unmarshal(m, b)
}
func (m *BlockReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockReq.Marshal(b, m, deterministic)
}
func (dst *BlockReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockReq.Merge(dst, src)
}
func (m *BlockReq) XXX_Size() int {
	return xxx_messageInfo_BlockReq.Size(m)
}
func (m *BlockReq) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockReq.DiscardUnknown(m)
}

var xxx_messageInfo_BlockReq proto.InternalMessageInfo

func (m *BlockReq) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

type BlockResp struct {
	Height               int64    `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
	Blocks               []*Block `protobuf:"bytes,3,rep,name=blocks,proto3" json:"blocks,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BlockResp) Reset()         { *m = BlockResp{} }
func (m *BlockResp) String() string { return proto.CompactTextString(m) }
func (*BlockResp) ProtoMessage()    {}
func (*BlockResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_node_7b7652e688675472, []int{3}
}
func (m *BlockResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BlockResp.Unmarshal(m, b)
}
func (m *BlockResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BlockResp.Marshal(b, m, deterministic)
}
func (dst *BlockResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockResp.Merge(dst, src)
}
func (m *BlockResp) XXX_Size() int {
	return xxx_messageInfo_BlockResp.Size(m)
}
func (m *BlockResp) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockResp.DiscardUnknown(m)
}

var xxx_messageInfo_BlockResp proto.InternalMessageInfo

func (m *BlockResp) GetHeight() int64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockResp) GetBlocks() []*Block {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func init() {
	proto.RegisterType((*PeerInfo)(nil), "bean.peer_info")
	proto.RegisterType((*PeerInfoList)(nil), "bean.peer_info_list")
	proto.RegisterType((*BlockReq)(nil), "bean.block_req")
	proto.RegisterType((*BlockResp)(nil), "bean.block_resp")
}

func init() { proto.RegisterFile("bean/node.proto", fileDescriptor_node_7b7652e688675472) }

var fileDescriptor_node_7b7652e688675472 = []byte{
	// 228 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8f, 0xc1, 0x6e, 0x83, 0x30,
	0x10, 0x44, 0x65, 0x43, 0x91, 0xd8, 0x48, 0x89, 0xb4, 0x87, 0xca, 0xca, 0xa5, 0x08, 0x2e, 0x9c,
	0x88, 0xd4, 0xaa, 0x1f, 0xd1, 0x5b, 0xe5, 0x1f, 0x40, 0x90, 0xb8, 0x8d, 0x45, 0xca, 0x6e, 0x8d,
	0x2f, 0xfc, 0x7d, 0xc4, 0x92, 0x44, 0x39, 0xd9, 0x9e, 0x79, 0xde, 0x99, 0x85, 0x5d, 0xef, 0xba,
	0xf1, 0x30, 0xd2, 0xc9, 0x35, 0x1c, 0x28, 0x12, 0xa6, 0x8b, 0xb0, 0xc7, 0xbf, 0xf9, 0x18, 0x66,
	0x8e, 0x74, 0x18, 0xdc, 0xbc, 0x3a, 0xfb, 0x15, 0x3d, 0x75, 0xb1, 0x5b, 0x85, 0xf2, 0x1b, 0x72,
	0x76, 0x2e, 0xb4, 0x7e, 0xfc, 0x21, 0x7c, 0x03, 0xcd, 0x83, 0x51, 0x85, 0xaa, 0x37, 0xef, 0xbb,
	0xe6, 0xfe, 0xbd, 0x61, 0xf2, 0x63, 0xb4, 0x9a, 0x07, 0xdc, 0x82, 0xf6, 0x6c, 0x74, 0xa1, 0xea,
	0xdc, 0x6a, 0xcf, 0x88, 0x90, 0x32, 0x85, 0x68, 0x92, 0x42, 0xd5, 0x2f, 0x56, 0xee, 0xe5, 0x27,
	0x6c, 0x1f, 0x13, 0xdb, 0x8b, 0x9f, 0x22, 0x56, 0x90, 0x2e, 0xa7, 0x51, 0x45, 0x22, 0x83, 0x97,
	0x0e, 0xcd, 0x83, 0xb1, 0x62, 0x96, 0x15, 0xe4, 0xfd, 0x85, 0x8e, 0x43, 0x1b, 0xdc, 0x3f, 0xbe,
	0x42, 0x76, 0x76, 0xfe, 0xf7, 0x1c, 0xa5, 0x4c, 0x62, 0x6f, 0xaf, 0xf2, 0x0b, 0xe0, 0x0e, 0x4d,
	0xfc, 0x44, 0xe9, 0x67, 0x0a, 0x2b, 0xc8, 0x84, 0x9a, 0x4c, 0x22, 0x89, 0x9b, 0x35, 0x51, 0x34,
	0x7b, 0xb3, 0xfa, 0x4c, 0xf6, 0xff, 0xb8, 0x06, 0x00, 0x00, 0xff, 0xff, 0x00, 0x37, 0x40, 0xf0,
	0x3d, 0x01, 0x00, 0x00,
}
