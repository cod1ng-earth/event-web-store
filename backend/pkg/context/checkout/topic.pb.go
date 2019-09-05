// Code generated by protoc-gen-go. DO NOT EDIT.
// source: topic.proto

package checkout

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
	//	*TopicMessage_ChangeProductQuantity
	//	*TopicMessage_StockCorrected
	//	*TopicMessage_Product
	//	*TopicMessage_OrderCart
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

type TopicMessage_ChangeProductQuantity struct {
	ChangeProductQuantity *ChangeProductQuantity `protobuf:"bytes,1,opt,name=changeProductQuantity,proto3,oneof"`
}

type TopicMessage_StockCorrected struct {
	StockCorrected *StockCorrected `protobuf:"bytes,2,opt,name=stockCorrected,proto3,oneof"`
}

type TopicMessage_Product struct {
	Product *Product `protobuf:"bytes,3,opt,name=product,proto3,oneof"`
}

type TopicMessage_OrderCart struct {
	OrderCart *OrderCart `protobuf:"bytes,4,opt,name=orderCart,proto3,oneof"`
}

func (*TopicMessage_ChangeProductQuantity) isTopicMessage_Messages() {}

func (*TopicMessage_StockCorrected) isTopicMessage_Messages() {}

func (*TopicMessage_Product) isTopicMessage_Messages() {}

func (*TopicMessage_OrderCart) isTopicMessage_Messages() {}

func (m *TopicMessage) GetMessages() isTopicMessage_Messages {
	if m != nil {
		return m.Messages
	}
	return nil
}

func (m *TopicMessage) GetChangeProductQuantity() *ChangeProductQuantity {
	if x, ok := m.GetMessages().(*TopicMessage_ChangeProductQuantity); ok {
		return x.ChangeProductQuantity
	}
	return nil
}

func (m *TopicMessage) GetStockCorrected() *StockCorrected {
	if x, ok := m.GetMessages().(*TopicMessage_StockCorrected); ok {
		return x.StockCorrected
	}
	return nil
}

func (m *TopicMessage) GetProduct() *Product {
	if x, ok := m.GetMessages().(*TopicMessage_Product); ok {
		return x.Product
	}
	return nil
}

func (m *TopicMessage) GetOrderCart() *OrderCart {
	if x, ok := m.GetMessages().(*TopicMessage_OrderCart); ok {
		return x.OrderCart
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*TopicMessage) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*TopicMessage_ChangeProductQuantity)(nil),
		(*TopicMessage_StockCorrected)(nil),
		(*TopicMessage_Product)(nil),
		(*TopicMessage_OrderCart)(nil),
	}
}

type ChangeProductQuantity struct {
	CartID               string   `protobuf:"bytes,3,opt,name=cartID,proto3" json:"cartID,omitempty"`
	ProductID            string   `protobuf:"bytes,1,opt,name=productID,proto3" json:"productID,omitempty"`
	Quantity             int64    `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ChangeProductQuantity) Reset()         { *m = ChangeProductQuantity{} }
func (m *ChangeProductQuantity) String() string { return proto.CompactTextString(m) }
func (*ChangeProductQuantity) ProtoMessage()    {}
func (*ChangeProductQuantity) Descriptor() ([]byte, []int) {
	return fileDescriptor_7312ad0e4fa171e8, []int{1}
}

func (m *ChangeProductQuantity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangeProductQuantity.Unmarshal(m, b)
}
func (m *ChangeProductQuantity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangeProductQuantity.Marshal(b, m, deterministic)
}
func (m *ChangeProductQuantity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeProductQuantity.Merge(m, src)
}
func (m *ChangeProductQuantity) XXX_Size() int {
	return xxx_messageInfo_ChangeProductQuantity.Size(m)
}
func (m *ChangeProductQuantity) XXX_DiscardUnknown() {
	xxx_messageInfo_ChangeProductQuantity.DiscardUnknown(m)
}

var xxx_messageInfo_ChangeProductQuantity proto.InternalMessageInfo

func (m *ChangeProductQuantity) GetCartID() string {
	if m != nil {
		return m.CartID
	}
	return ""
}

func (m *ChangeProductQuantity) GetProductID() string {
	if m != nil {
		return m.ProductID
	}
	return ""
}

func (m *ChangeProductQuantity) GetQuantity() int64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type StockCorrected struct {
	ProductID            string   `protobuf:"bytes,1,opt,name=productID,proto3" json:"productID,omitempty"`
	QuantityChange       int64    `protobuf:"varint,2,opt,name=quantityChange,proto3" json:"quantityChange,omitempty"`
	FulfilmentOffset     int64    `protobuf:"varint,3,opt,name=fulfilmentOffset,proto3" json:"fulfilmentOffset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StockCorrected) Reset()         { *m = StockCorrected{} }
func (m *StockCorrected) String() string { return proto.CompactTextString(m) }
func (*StockCorrected) ProtoMessage()    {}
func (*StockCorrected) Descriptor() ([]byte, []int) {
	return fileDescriptor_7312ad0e4fa171e8, []int{2}
}

func (m *StockCorrected) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StockCorrected.Unmarshal(m, b)
}
func (m *StockCorrected) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StockCorrected.Marshal(b, m, deterministic)
}
func (m *StockCorrected) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StockCorrected.Merge(m, src)
}
func (m *StockCorrected) XXX_Size() int {
	return xxx_messageInfo_StockCorrected.Size(m)
}
func (m *StockCorrected) XXX_DiscardUnknown() {
	xxx_messageInfo_StockCorrected.DiscardUnknown(m)
}

var xxx_messageInfo_StockCorrected proto.InternalMessageInfo

func (m *StockCorrected) GetProductID() string {
	if m != nil {
		return m.ProductID
	}
	return ""
}

func (m *StockCorrected) GetQuantityChange() int64 {
	if m != nil {
		return m.QuantityChange
	}
	return 0
}

func (m *StockCorrected) GetFulfilmentOffset() int64 {
	if m != nil {
		return m.FulfilmentOffset
	}
	return 0
}

type Product struct {
	ProductID            string   `protobuf:"bytes,1,opt,name=productID,proto3" json:"productID,omitempty"`
	Price                int64    `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	SmallImageURL        string   `protobuf:"bytes,4,opt,name=smallImageURL,proto3" json:"smallImageURL,omitempty"`
	PimOffset            int64    `protobuf:"varint,5,opt,name=pimOffset,proto3" json:"pimOffset,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_7312ad0e4fa171e8, []int{3}
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

func (m *Product) GetProductID() string {
	if m != nil {
		return m.ProductID
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

func (m *Product) GetSmallImageURL() string {
	if m != nil {
		return m.SmallImageURL
	}
	return ""
}

func (m *Product) GetPimOffset() int64 {
	if m != nil {
		return m.PimOffset
	}
	return 0
}

type OrderCart struct {
	CartID               string   `protobuf:"bytes,1,opt,name=cartID,proto3" json:"cartID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderCart) Reset()         { *m = OrderCart{} }
func (m *OrderCart) String() string { return proto.CompactTextString(m) }
func (*OrderCart) ProtoMessage()    {}
func (*OrderCart) Descriptor() ([]byte, []int) {
	return fileDescriptor_7312ad0e4fa171e8, []int{4}
}

func (m *OrderCart) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderCart.Unmarshal(m, b)
}
func (m *OrderCart) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderCart.Marshal(b, m, deterministic)
}
func (m *OrderCart) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderCart.Merge(m, src)
}
func (m *OrderCart) XXX_Size() int {
	return xxx_messageInfo_OrderCart.Size(m)
}
func (m *OrderCart) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderCart.DiscardUnknown(m)
}

var xxx_messageInfo_OrderCart proto.InternalMessageInfo

func (m *OrderCart) GetCartID() string {
	if m != nil {
		return m.CartID
	}
	return ""
}

func init() {
	proto.RegisterType((*TopicMessage)(nil), "checkout.TopicMessage")
	proto.RegisterType((*ChangeProductQuantity)(nil), "checkout.ChangeProductQuantity")
	proto.RegisterType((*StockCorrected)(nil), "checkout.StockCorrected")
	proto.RegisterType((*Product)(nil), "checkout.Product")
	proto.RegisterType((*OrderCart)(nil), "checkout.OrderCart")
}

func init() { proto.RegisterFile("topic.proto", fileDescriptor_7312ad0e4fa171e8) }

var fileDescriptor_7312ad0e4fa171e8 = []byte{
	// 361 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0x4d, 0x4f, 0xc2, 0x40,
	0x10, 0x6d, 0xf9, 0xee, 0xa0, 0x44, 0x57, 0x31, 0x8d, 0x31, 0xd1, 0x54, 0x63, 0x8c, 0x89, 0x1c,
	0xe4, 0x1f, 0x00, 0x87, 0x92, 0x68, 0xd0, 0x55, 0xe3, 0x79, 0x5d, 0xb6, 0xd0, 0xd0, 0x76, 0xeb,
	0x76, 0x7b, 0xf0, 0xea, 0x7f, 0xf0, 0xe4, 0x9f, 0x35, 0x6c, 0xb7, 0x14, 0xb0, 0x09, 0xb7, 0xce,
	0x9b, 0x79, 0xef, 0xcd, 0xbc, 0x2d, 0xb4, 0x25, 0x8f, 0x7d, 0xda, 0x8b, 0x05, 0x97, 0x1c, 0xb5,
	0xe8, 0x9c, 0xd1, 0x05, 0x4f, 0xa5, 0xf3, 0x5b, 0x81, 0xbd, 0xd7, 0x65, 0xe7, 0x91, 0x25, 0x09,
	0x99, 0x31, 0xf4, 0x0e, 0x5d, 0x3a, 0x27, 0xd1, 0x8c, 0x3d, 0x09, 0x3e, 0x4d, 0xa9, 0x7c, 0x4e,
	0x49, 0x24, 0x7d, 0xf9, 0x65, 0x9b, 0x17, 0xe6, 0x4d, 0xfb, 0xfe, 0xbc, 0x97, 0x53, 0x7b, 0xc3,
	0xb2, 0x31, 0xd7, 0xc0, 0xe5, 0x7c, 0x34, 0x80, 0x4e, 0x22, 0x39, 0x5d, 0x0c, 0xb9, 0x10, 0x8c,
	0x4a, 0x36, 0xb5, 0x2b, 0x4a, 0xd1, 0x2e, 0x14, 0x5f, 0x36, 0xfa, 0xae, 0x81, 0xb7, 0x18, 0xe8,
	0x0e, 0x9a, 0x71, 0x26, 0x6b, 0x57, 0x15, 0xf9, 0xb0, 0x20, 0x6b, 0x3f, 0xd7, 0xc0, 0xf9, 0x0c,
	0xea, 0x83, 0xc5, 0xc5, 0x94, 0x89, 0x21, 0x11, 0xd2, 0xae, 0x29, 0xc2, 0x51, 0x41, 0x98, 0xe4,
	0x2d, 0xd7, 0xc0, 0xc5, 0xdc, 0x00, 0xa0, 0x15, 0x66, 0x59, 0x24, 0x8e, 0x0f, 0xdd, 0xd2, 0x2b,
	0xd1, 0x09, 0x34, 0x28, 0x11, 0x72, 0x3c, 0x52, 0x7b, 0x58, 0x58, 0x57, 0xe8, 0x0c, 0x2c, 0x6d,
	0x3e, 0x1e, 0xa9, 0xc4, 0x2c, 0x5c, 0x00, 0xe8, 0x14, 0x5a, 0x9f, 0x79, 0x9c, 0xcb, 0xe3, 0xab,
	0x78, 0x55, 0x3b, 0xdf, 0x26, 0x74, 0x36, 0xef, 0xdf, 0x21, 0x76, 0x0d, 0x9d, 0x9c, 0x9c, 0xed,
	0xa8, 0x25, 0xb7, 0x50, 0x74, 0x0b, 0x07, 0x5e, 0x1a, 0x78, 0x7e, 0x10, 0xb2, 0x48, 0x4e, 0x3c,
	0x2f, 0x61, 0x59, 0x78, 0x55, 0xfc, 0x0f, 0x77, 0x7e, 0x4c, 0x68, 0xea, 0x53, 0x77, 0xb8, 0x1f,
	0x43, 0x3d, 0x16, 0x3e, 0xcd, 0x4d, 0xb3, 0x02, 0x21, 0xa8, 0x45, 0x24, 0x64, 0x3a, 0x14, 0xf5,
	0x8d, 0xae, 0x60, 0x3f, 0x09, 0x49, 0x10, 0x8c, 0x43, 0x32, 0x63, 0x6f, 0xf8, 0x41, 0x3d, 0x84,
	0x85, 0x37, 0x41, 0xe5, 0xe6, 0x87, 0x7a, 0xbd, 0xba, 0xd2, 0x2c, 0x00, 0xe7, 0x12, 0xac, 0xd5,
	0x6b, 0xad, 0x65, 0x6f, 0xae, 0x67, 0xff, 0xd1, 0x50, 0xff, 0x76, 0xff, 0x2f, 0x00, 0x00, 0xff,
	0xff, 0x45, 0xae, 0xae, 0x28, 0xea, 0x02, 0x00, 0x00,
}
