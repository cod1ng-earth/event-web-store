// Code generated by protoc-gen-go. DO NOT EDIT.
// source: topic.proto

package pim

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type TopicMessage struct {
	// Types that are valid to be assigned to Messages:
	//	*TopicMessage_Product
	Messages             isTopicMessage_Messages `protobuf_oneof:"messages"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *TopicMessage) Reset()         { *m = TopicMessage{} }
func (m *TopicMessage) String() string { return proto.CompactTextString(m) }
func (*TopicMessage) ProtoMessage()    {}
func (*TopicMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_7312ad0e4fa171e8, []int{0}
}

func (m *TopicMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TopicMessage.Unmarshal(m, b)
}
func (m *TopicMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TopicMessage.Marshal(b, m, deterministic)
}
func (m *TopicMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TopicMessage.Merge(m, src)
}
func (m *TopicMessage) XXX_Size() int {
	return xxx_messageInfo_TopicMessage.Size(m)
}
func (m *TopicMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_TopicMessage.DiscardUnknown(m)
}

var xxx_messageInfo_TopicMessage proto.InternalMessageInfo

type isTopicMessage_Messages interface {
	isTopicMessage_Messages()
}

type TopicMessage_Product struct {
	Product *Product `protobuf:"bytes,1,opt,name=product,proto3,oneof"`
}

func (*TopicMessage_Product) isTopicMessage_Messages() {}

func (m *TopicMessage) GetMessages() isTopicMessage_Messages {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *TopicMessage) GetProduct() *Product {
	if x, ok := m.GetMessages().(*TopicMessage_Product); ok {
		return x.Product
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*TopicMessage) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TopicMessage_Product)(nil),
	}
}

type Product struct {
	InternalOffset       int64    `protobuf:"varint,99,opt,name=internalOffset,proto3" json:"internalOffset,omitempty"`
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Price                int64    `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Longtext             string   `protobuf:"bytes,5,opt,name=longtext,proto3" json:"longtext,omitempty"`
	Category             string   `protobuf:"bytes,6,opt,name=category,proto3" json:"category,omitempty"`
	SmallImageURL        string   `protobuf:"bytes,7,opt,name=smallImageURL,proto3" json:"smallImageURL,omitempty"`
	LargeImageURL        string   `protobuf:"bytes,8,opt,name=largeImageURL,proto3" json:"largeImageURL,omitempty"`
	Disabled             bool     `protobuf:"varint,9,opt,name=disabled,proto3" json:"disabled,omitempty"`
	Tax                  int64    `protobuf:"varint,10,opt,name=tax,proto3" json:"tax,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_7312ad0e4fa171e8, []int{1}
}

func (m *Product) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Product.Unmarshal(m, b)
}
func (m *Product) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Product.Marshal(b, m, deterministic)
}
func (m *Product) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Product.Merge(m, src)
}
func (m *Product) XXX_Size() int {
	return xxx_messageInfo_Product.Size(m)
}
func (m *Product) XXX_DiscardUnknown() {
	xxx_messageInfo_Product.DiscardUnknown(m)
}

var xxx_messageInfo_Product proto.InternalMessageInfo

func (m *Product) GetInternalOffset() int64 {
	if m != nil {
		return m.InternalOffset
	}
	return 0
}

func (m *Product) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Product) GetPrice() int64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Product) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Product) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *Product) GetLongtext() string {
	if m != nil {
		return m.Longtext
	}
	return ""
}

func (m *Product) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

func (m *Product) GetSmallImageURL() string {
	if m != nil {
		return m.SmallImageURL
	}
	return ""
}

func (m *Product) GetLargeImageURL() string {
	if m != nil {
		return m.LargeImageURL
	}
	return ""
}

func (m *Product) GetDisabled() bool {
	if m != nil {
		return m.Disabled
	}
	return false
}

func (m *Product) GetTax() int64 {
	if m != nil {
		return m.Tax
	}
	return 0
}

func init() {
	proto.RegisterType((*TopicMessage)(nil), "pim.TopicMessage")
	proto.RegisterType((*Product)(nil), "pim.Product")
}

func init() { proto.RegisterFile("topic.proto", fileDescriptor_7312ad0e4fa171e8) }

var fileDescriptor_7312ad0e4fa171e8 = []byte{
	// 277 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x91, 0xdf, 0x4a, 0xc3, 0x30,
	0x14, 0xc6, 0x6d, 0xbb, 0x3f, 0xdd, 0xd9, 0x1c, 0x72, 0xf0, 0x22, 0x78, 0x55, 0x86, 0x48, 0xaf,
	0x76, 0xa1, 0x6f, 0x20, 0x5e, 0x28, 0x28, 0x4a, 0xd1, 0x07, 0xc8, 0x92, 0xb3, 0x12, 0x48, 0x9a,
	0x90, 0x44, 0x98, 0x4f, 0xe7, 0xab, 0x49, 0x53, 0x57, 0x9c, 0x77, 0xe7, 0xfb, 0x7d, 0xbf, 0x70,
	0x08, 0x07, 0x96, 0xd1, 0x3a, 0x25, 0xb6, 0xce, 0xdb, 0x68, 0xb1, 0x70, 0xca, 0x6c, 0x1e, 0x60,
	0xf5, 0xde, 0xb3, 0x17, 0x0a, 0x81, 0xb7, 0x84, 0x35, 0xcc, 0x9d, 0xb7, 0xf2, 0x53, 0x44, 0x96,
	0x55, 0x59, 0xbd, 0xbc, 0x5d, 0x6d, 0x9d, 0x32, 0xdb, 0xb7, 0x81, 0x3d, 0x9e, 0x35, 0xc7, 0xfa,
	0x1e, 0xa0, 0x34, 0xc3, 0xa3, 0xb0, 0xf9, 0xce, 0x61, 0xfe, 0xab, 0xe0, 0x0d, 0xac, 0x55, 0x17,
	0xc9, 0x77, 0x5c, 0xbf, 0xee, 0xf7, 0x81, 0x22, 0x13, 0x55, 0x56, 0x17, 0xcd, 0x3f, 0x8a, 0x6b,
	0xc8, 0x95, 0x4c, 0x4b, 0x16, 0x4d, 0xae, 0x24, 0x5e, 0xc2, 0xd4, 0x79, 0x25, 0x88, 0xe5, 0x49,
	0x1f, 0x02, 0x22, 0x4c, 0x3a, 0x6e, 0x88, 0x15, 0xc9, 0x4b, 0x33, 0x56, 0xb0, 0x94, 0x14, 0x84,
	0x57, 0x2e, 0x2a, 0xdb, 0xb1, 0x49, 0xaa, 0xfe, 0x22, 0xbc, 0x82, 0x52, 0xdb, 0xae, 0x8d, 0x74,
	0x88, 0x6c, 0x9a, 0xea, 0x31, 0xf7, 0x9d, 0xe0, 0x91, 0x5a, 0xeb, 0xbf, 0xd8, 0x6c, 0xe8, 0x8e,
	0x19, 0xaf, 0xe1, 0x3c, 0x18, 0xae, 0xf5, 0x93, 0xe1, 0x2d, 0x7d, 0x34, 0xcf, 0x6c, 0x9e, 0x84,
	0x53, 0xd8, 0x5b, 0x9a, 0xfb, 0x96, 0x46, 0xab, 0x1c, 0xac, 0x13, 0xd8, 0xef, 0x91, 0x2a, 0xf0,
	0x9d, 0x26, 0xc9, 0x16, 0x55, 0x56, 0x97, 0xcd, 0x98, 0xf1, 0x02, 0x8a, 0xc8, 0x0f, 0x0c, 0xd2,
	0x4f, 0xfb, 0x71, 0x37, 0x4b, 0x37, 0xb9, 0xfb, 0x09, 0x00, 0x00, 0xff, 0xff, 0xf5, 0x52, 0x69,
	0xc9, 0xa2, 0x01, 0x00, 0x00,
}
