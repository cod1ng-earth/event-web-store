// Code generated by protoc-gen-go. DO NOT EDIT.
// source: topic.proto

package catalog

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
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Price                int64    `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Longtext             string   `protobuf:"bytes,5,opt,name=longtext,proto3" json:"longtext,omitempty"`
	Category             string   `protobuf:"bytes,6,opt,name=category,proto3" json:"category,omitempty"`
	SmallImageURL        string   `protobuf:"bytes,7,opt,name=smallImageURL,proto3" json:"smallImageURL,omitempty"`
	LargeImageURL        string   `protobuf:"bytes,8,opt,name=largeImageURL,proto3" json:"largeImageURL,omitempty"`
	Disabled             bool     `protobuf:"varint,9,opt,name=disabled,proto3" json:"disabled,omitempty"`
	PimOffset            int64    `protobuf:"varint,10,opt,name=pimOffset,proto3" json:"pimOffset,omitempty"`
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

func (m *Product) GetPimOffset() int64 {
	if m != nil {
		return m.PimOffset
	}
	return 0
}

func init() {
	proto.RegisterType((*TopicMessage)(nil), "catalog.TopicMessage")
	proto.RegisterType((*Product)(nil), "catalog.Product")
}

func init() { proto.RegisterFile("topic.proto", fileDescriptor_7312ad0e4fa171e8) }

var fileDescriptor_7312ad0e4fa171e8 = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x90, 0xcd, 0x4a, 0xf4, 0x30,
	0x14, 0x86, 0xbf, 0x76, 0x7e, 0xda, 0x9e, 0x7e, 0x8a, 0x04, 0x17, 0x41, 0x5c, 0x94, 0xc1, 0x45,
	0x17, 0xd2, 0x85, 0xde, 0x81, 0xab, 0x11, 0x14, 0xa5, 0xe8, 0x05, 0x64, 0x92, 0x33, 0x21, 0x90,
	0x36, 0x21, 0x89, 0xa0, 0xd7, 0xe3, 0x8d, 0x4a, 0x93, 0x99, 0xea, 0xec, 0xf2, 0x3e, 0xef, 0x93,
	0x84, 0x73, 0xa0, 0x0e, 0xc6, 0x2a, 0xde, 0x59, 0x67, 0x82, 0x21, 0x05, 0x67, 0x81, 0x69, 0x23,
	0x37, 0x5b, 0xf8, 0xff, 0x36, 0xf1, 0x67, 0xf4, 0x9e, 0x49, 0x24, 0xb7, 0x50, 0x58, 0x67, 0xc4,
	0x07, 0x0f, 0x34, 0x6b, 0xb2, 0xb6, 0xbe, 0xbb, 0xe8, 0x0e, 0x6a, 0xf7, 0x9a, 0xf8, 0xf6, 0x5f,
	0x7f, 0x54, 0x1e, 0x00, 0xca, 0x21, 0x5d, 0xf4, 0x9b, 0xef, 0x1c, 0x8a, 0x83, 0x42, 0xce, 0x21,
	0x57, 0x22, 0x3e, 0x50, 0xf5, 0xb9, 0x12, 0xe4, 0x12, 0x56, 0xd6, 0x29, 0x8e, 0x34, 0x6f, 0xb2,
	0x76, 0xd1, 0xa7, 0x40, 0x08, 0x2c, 0x47, 0x36, 0x20, 0x5d, 0x44, 0x2f, 0x9e, 0x49, 0x03, 0xb5,
	0x40, 0xcf, 0x9d, 0xb2, 0x41, 0x99, 0x91, 0x2e, 0x63, 0xf5, 0x17, 0x91, 0x2b, 0x28, 0xb5, 0x19,
	0x65, 0xc0, 0xcf, 0x40, 0x57, 0xb1, 0x9e, 0xf3, 0xd4, 0x71, 0x16, 0x50, 0x1a, 0xf7, 0x45, 0xd7,
	0xa9, 0x3b, 0x66, 0x72, 0x03, 0x67, 0x7e, 0x60, 0x5a, 0x3f, 0x0e, 0x4c, 0xe2, 0x7b, 0xff, 0x44,
	0x8b, 0x28, 0x9c, 0xc2, 0xc9, 0xd2, 0xcc, 0x49, 0x9c, 0xad, 0x32, 0x59, 0x27, 0x70, 0xfa, 0x47,
	0x28, 0xcf, 0x76, 0x1a, 0x05, 0xad, 0x9a, 0xac, 0x2d, 0xfb, 0x39, 0x93, 0x6b, 0xa8, 0xac, 0x1a,
	0x5e, 0xf6, 0x7b, 0x8f, 0x81, 0x42, 0x9c, 0xf7, 0x17, 0xec, 0xd6, 0x71, 0xff, 0xf7, 0x3f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x2b, 0x47, 0x48, 0x0f, 0x8e, 0x01, 0x00, 0x00,
}
