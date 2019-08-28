// Code generated by protoc-gen-go. DO NOT EDIT.
// source: checkout.proto

package checkout

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CheckoutMessages struct {
	// Types that are valid to be assigned to CheckoutMessage:
	//	*CheckoutMessages_ChangeProductQuantity
	//	*CheckoutMessages_StockCorrected
	//	*CheckoutMessages_Product
	//	*CheckoutMessages_OrderCart
	CheckoutMessage      isCheckoutMessages_CheckoutMessage `protobuf_oneof:"checkoutMessage"`
	XXX_NoUnkeyedLiteral struct{}                           `json:"-"`
	XXX_unrecognized     []byte                             `json:"-"`
	XXX_sizecache        int32                              `json:"-"`
}

func (m *CheckoutMessages) Reset()         { *m = CheckoutMessages{} }
func (m *CheckoutMessages) String() string { return proto.CompactTextString(m) }
func (*CheckoutMessages) ProtoMessage()    {}
func (*CheckoutMessages) Descriptor() ([]byte, []int) {
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{0}
}
func (m *CheckoutMessages) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckoutMessages.Unmarshal(m, b)
}
func (m *CheckoutMessages) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckoutMessages.Marshal(b, m, deterministic)
}
func (dst *CheckoutMessages) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckoutMessages.Merge(dst, src)
}
func (m *CheckoutMessages) XXX_Size() int {
	return xxx_messageInfo_CheckoutMessages.Size(m)
}
func (m *CheckoutMessages) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckoutMessages.DiscardUnknown(m)
}

var xxx_messageInfo_CheckoutMessages proto.InternalMessageInfo

type isCheckoutMessages_CheckoutMessage interface {
	isCheckoutMessages_CheckoutMessage()
}

type CheckoutMessages_ChangeProductQuantity struct {
	ChangeProductQuantity *ChangeProductQuantity `protobuf:"bytes,1,opt,name=changeProductQuantity,proto3,oneof"`
}

type CheckoutMessages_StockCorrected struct {
	StockCorrected *StockCorrected `protobuf:"bytes,2,opt,name=stockCorrected,proto3,oneof"`
}

type CheckoutMessages_Product struct {
	Product *Product `protobuf:"bytes,3,opt,name=product,proto3,oneof"`
}

type CheckoutMessages_OrderCart struct {
	OrderCart *OrderCart `protobuf:"bytes,4,opt,name=orderCart,proto3,oneof"`
}

func (*CheckoutMessages_ChangeProductQuantity) isCheckoutMessages_CheckoutMessage() {}

func (*CheckoutMessages_StockCorrected) isCheckoutMessages_CheckoutMessage() {}

func (*CheckoutMessages_Product) isCheckoutMessages_CheckoutMessage() {}

func (*CheckoutMessages_OrderCart) isCheckoutMessages_CheckoutMessage() {}

func (m *CheckoutMessages) GetCheckoutMessage() isCheckoutMessages_CheckoutMessage {
	if m != nil {
		return m.CheckoutMessage
	}
	return nil
}

func (m *CheckoutMessages) GetChangeProductQuantity() *ChangeProductQuantity {
	if x, ok := m.GetCheckoutMessage().(*CheckoutMessages_ChangeProductQuantity); ok {
		return x.ChangeProductQuantity
	}
	return nil
}

func (m *CheckoutMessages) GetStockCorrected() *StockCorrected {
	if x, ok := m.GetCheckoutMessage().(*CheckoutMessages_StockCorrected); ok {
		return x.StockCorrected
	}
	return nil
}

func (m *CheckoutMessages) GetProduct() *Product {
	if x, ok := m.GetCheckoutMessage().(*CheckoutMessages_Product); ok {
		return x.Product
	}
	return nil
}

func (m *CheckoutMessages) GetOrderCart() *OrderCart {
	if x, ok := m.GetCheckoutMessage().(*CheckoutMessages_OrderCart); ok {
		return x.OrderCart
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*CheckoutMessages) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CheckoutMessages_OneofMarshaler, _CheckoutMessages_OneofUnmarshaler, _CheckoutMessages_OneofSizer, []interface{}{
		(*CheckoutMessages_ChangeProductQuantity)(nil),
		(*CheckoutMessages_StockCorrected)(nil),
		(*CheckoutMessages_Product)(nil),
		(*CheckoutMessages_OrderCart)(nil),
	}
}

func _CheckoutMessages_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*CheckoutMessages)
	// checkoutMessage
	switch x := m.CheckoutMessage.(type) {
	case *CheckoutMessages_ChangeProductQuantity:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChangeProductQuantity); err != nil {
			return err
		}
	case *CheckoutMessages_StockCorrected:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.StockCorrected); err != nil {
			return err
		}
	case *CheckoutMessages_Product:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Product); err != nil {
			return err
		}
	case *CheckoutMessages_OrderCart:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.OrderCart); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("CheckoutMessages.CheckoutMessage has unexpected type %T", x)
	}
	return nil
}

func _CheckoutMessages_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*CheckoutMessages)
	switch tag {
	case 1: // checkoutMessage.changeProductQuantity
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChangeProductQuantity)
		err := b.DecodeMessage(msg)
		m.CheckoutMessage = &CheckoutMessages_ChangeProductQuantity{msg}
		return true, err
	case 2: // checkoutMessage.stockCorrected
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(StockCorrected)
		err := b.DecodeMessage(msg)
		m.CheckoutMessage = &CheckoutMessages_StockCorrected{msg}
		return true, err
	case 3: // checkoutMessage.product
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Product)
		err := b.DecodeMessage(msg)
		m.CheckoutMessage = &CheckoutMessages_Product{msg}
		return true, err
	case 4: // checkoutMessage.orderCart
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(OrderCart)
		err := b.DecodeMessage(msg)
		m.CheckoutMessage = &CheckoutMessages_OrderCart{msg}
		return true, err
	default:
		return false, nil
	}
}

func _CheckoutMessages_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*CheckoutMessages)
	// checkoutMessage
	switch x := m.CheckoutMessage.(type) {
	case *CheckoutMessages_ChangeProductQuantity:
		s := proto.Size(x.ChangeProductQuantity)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CheckoutMessages_StockCorrected:
		s := proto.Size(x.StockCorrected)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CheckoutMessages_Product:
		s := proto.Size(x.Product)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CheckoutMessages_OrderCart:
		s := proto.Size(x.OrderCart)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case nil:
	default:
		panic(fmt.Sprintf("proto: unexpected type %T in oneof", x))
	}
	return n
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
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{1}
}
func (m *Product) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Product.Unmarshal(m, b)
}
func (m *Product) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Product.Marshal(b, m, deterministic)
}
func (dst *Product) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Product.Merge(dst, src)
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

type Cart struct {
	Id                   string      `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Positions            []*Position `protobuf:"bytes,2,rep,name=positions,proto3" json:"positions,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Cart) Reset()         { *m = Cart{} }
func (m *Cart) String() string { return proto.CompactTextString(m) }
func (*Cart) ProtoMessage()    {}
func (*Cart) Descriptor() ([]byte, []int) {
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{2}
}
func (m *Cart) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Cart.Unmarshal(m, b)
}
func (m *Cart) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Cart.Marshal(b, m, deterministic)
}
func (dst *Cart) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Cart.Merge(dst, src)
}
func (m *Cart) XXX_Size() int {
	return xxx_messageInfo_Cart.Size(m)
}
func (m *Cart) XXX_DiscardUnknown() {
	xxx_messageInfo_Cart.DiscardUnknown(m)
}

var xxx_messageInfo_Cart proto.InternalMessageInfo

func (m *Cart) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Cart) GetPositions() []*Position {
	if m != nil {
		return m.Positions
	}
	return nil
}

type Position struct {
	ProductID            string   `protobuf:"bytes,1,opt,name=productID,proto3" json:"productID,omitempty"`
	Price                int64    `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	SmallImageURL        string   `protobuf:"bytes,4,opt,name=smallImageURL,proto3" json:"smallImageURL,omitempty"`
	Quantity             int64    `protobuf:"varint,6,opt,name=quantity,proto3" json:"quantity,omitempty"`
	InStock              bool     `protobuf:"varint,7,opt,name=inStock,proto3" json:"inStock,omitempty"`
	MoreInStock          bool     `protobuf:"varint,8,opt,name=moreInStock,proto3" json:"moreInStock,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Position) Reset()         { *m = Position{} }
func (m *Position) String() string { return proto.CompactTextString(m) }
func (*Position) ProtoMessage()    {}
func (*Position) Descriptor() ([]byte, []int) {
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{3}
}
func (m *Position) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Position.Unmarshal(m, b)
}
func (m *Position) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Position.Marshal(b, m, deterministic)
}
func (dst *Position) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Position.Merge(dst, src)
}
func (m *Position) XXX_Size() int {
	return xxx_messageInfo_Position.Size(m)
}
func (m *Position) XXX_DiscardUnknown() {
	xxx_messageInfo_Position.DiscardUnknown(m)
}

var xxx_messageInfo_Position proto.InternalMessageInfo

func (m *Position) GetProductID() string {
	if m != nil {
		return m.ProductID
	}
	return ""
}

func (m *Position) GetPrice() int64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *Position) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Position) GetSmallImageURL() string {
	if m != nil {
		return m.SmallImageURL
	}
	return ""
}

func (m *Position) GetQuantity() int64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func (m *Position) GetInStock() bool {
	if m != nil {
		return m.InStock
	}
	return false
}

func (m *Position) GetMoreInStock() bool {
	if m != nil {
		return m.MoreInStock
	}
	return false
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
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{4}
}
func (m *ChangeProductQuantity) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChangeProductQuantity.Unmarshal(m, b)
}
func (m *ChangeProductQuantity) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChangeProductQuantity.Marshal(b, m, deterministic)
}
func (dst *ChangeProductQuantity) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChangeProductQuantity.Merge(dst, src)
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
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{5}
}
func (m *OrderCart) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderCart.Unmarshal(m, b)
}
func (m *OrderCart) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderCart.Marshal(b, m, deterministic)
}
func (dst *OrderCart) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderCart.Merge(dst, src)
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

type OrderCartResonse struct {
	Successful           bool     `protobuf:"varint,1,opt,name=successful,proto3" json:"successful,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderCartResonse) Reset()         { *m = OrderCartResonse{} }
func (m *OrderCartResonse) String() string { return proto.CompactTextString(m) }
func (*OrderCartResonse) ProtoMessage()    {}
func (*OrderCartResonse) Descriptor() ([]byte, []int) {
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{6}
}
func (m *OrderCartResonse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderCartResonse.Unmarshal(m, b)
}
func (m *OrderCartResonse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderCartResonse.Marshal(b, m, deterministic)
}
func (dst *OrderCartResonse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderCartResonse.Merge(dst, src)
}
func (m *OrderCartResonse) XXX_Size() int {
	return xxx_messageInfo_OrderCartResonse.Size(m)
}
func (m *OrderCartResonse) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderCartResonse.DiscardUnknown(m)
}

var xxx_messageInfo_OrderCartResonse proto.InternalMessageInfo

func (m *OrderCartResonse) GetSuccessful() bool {
	if m != nil {
		return m.Successful
	}
	return false
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
	return fileDescriptor_checkout_5a6dea247f8c3a63, []int{7}
}
func (m *StockCorrected) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StockCorrected.Unmarshal(m, b)
}
func (m *StockCorrected) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StockCorrected.Marshal(b, m, deterministic)
}
func (dst *StockCorrected) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StockCorrected.Merge(dst, src)
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

func init() {
	proto.RegisterType((*CheckoutMessages)(nil), "checkout.CheckoutMessages")
	proto.RegisterType((*Product)(nil), "checkout.Product")
	proto.RegisterType((*Cart)(nil), "checkout.Cart")
	proto.RegisterType((*Position)(nil), "checkout.Position")
	proto.RegisterType((*ChangeProductQuantity)(nil), "checkout.ChangeProductQuantity")
	proto.RegisterType((*OrderCart)(nil), "checkout.OrderCart")
	proto.RegisterType((*OrderCartResonse)(nil), "checkout.OrderCartResonse")
	proto.RegisterType((*StockCorrected)(nil), "checkout.StockCorrected")
}

func init() { proto.RegisterFile("checkout.proto", fileDescriptor_checkout_5a6dea247f8c3a63) }

var fileDescriptor_checkout_5a6dea247f8c3a63 = []byte{
	// 467 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x54, 0xc1, 0x6e, 0xd3, 0x40,
	0x10, 0x8d, 0x9d, 0x34, 0xb1, 0x27, 0xc2, 0xa4, 0x03, 0x45, 0x2b, 0x84, 0x20, 0x32, 0x08, 0x45,
	0x48, 0x54, 0x28, 0xfd, 0x83, 0x86, 0x83, 0x23, 0x81, 0x0a, 0x8b, 0x10, 0x67, 0xb3, 0x5e, 0xa7,
	0xab, 0xda, 0xde, 0xb0, 0xbb, 0x3e, 0x70, 0xe5, 0x1f, 0xf8, 0x05, 0xbe, 0x87, 0x4f, 0x42, 0xdd,
	0xac, 0x63, 0xbb, 0x44, 0xca, 0xad, 0xb7, 0xcc, 0x9b, 0xf7, 0xde, 0xbe, 0x19, 0x4d, 0x0c, 0x11,
	0xbb, 0xe6, 0xec, 0x46, 0xd6, 0xe6, 0x7c, 0xab, 0xa4, 0x91, 0x18, 0x34, 0x75, 0xfc, 0xc7, 0x87,
	0xd9, 0xca, 0x15, 0x1f, 0xb9, 0xd6, 0xe9, 0x86, 0x6b, 0xfc, 0x06, 0x67, 0xec, 0x3a, 0xad, 0x36,
	0xfc, 0x93, 0x92, 0x59, 0xcd, 0xcc, 0xe7, 0x3a, 0xad, 0x8c, 0x30, 0x3f, 0x89, 0x37, 0xf7, 0x16,
	0xd3, 0xe5, 0x8b, 0xf3, 0xbd, 0xdd, 0xea, 0x10, 0x2d, 0x19, 0xd0, 0xc3, 0x7a, 0xbc, 0x84, 0x48,
	0x1b, 0xc9, 0x6e, 0x56, 0x52, 0x29, 0xce, 0x0c, 0xcf, 0x88, 0x6f, 0x1d, 0x49, 0xeb, 0xf8, 0xa5,
	0xd7, 0x4f, 0x06, 0xf4, 0x8e, 0x02, 0xdf, 0xc2, 0x64, 0xbb, 0xb3, 0x25, 0x43, 0x2b, 0x3e, 0x6d,
	0xc5, 0xee, 0xbd, 0x64, 0x40, 0x1b, 0x0e, 0x5e, 0x40, 0x28, 0x55, 0xc6, 0xd5, 0x2a, 0x55, 0x86,
	0x8c, 0xac, 0xe0, 0x51, 0x2b, 0xb8, 0x6a, 0x5a, 0xc9, 0x80, 0xb6, 0xbc, 0xcb, 0x53, 0x78, 0xc8,
	0xfa, 0x4b, 0x89, 0x7f, 0x7b, 0x30, 0x71, 0xf6, 0xf8, 0x0c, 0x42, 0x67, 0xbf, 0x7e, 0x6f, 0x77,
	0x12, 0xd2, 0x16, 0xc0, 0xc7, 0x70, 0xb2, 0x55, 0x82, 0x71, 0x3b, 0xdb, 0x90, 0xee, 0x0a, 0x44,
	0x18, 0x55, 0x69, 0xc9, 0x6d, 0xe6, 0x90, 0xda, 0xdf, 0xf8, 0x0a, 0x1e, 0xe8, 0x32, 0x2d, 0x8a,
	0x75, 0x99, 0x6e, 0xf8, 0x57, 0xfa, 0xc1, 0xe6, 0x0b, 0x69, 0x1f, 0xb4, 0xaf, 0x89, 0xf2, 0x2a,
	0xcf, 0x35, 0x37, 0xe4, 0xc4, 0x7a, 0xb6, 0x40, 0x9c, 0xc0, 0xe8, 0x36, 0x32, 0x46, 0xe0, 0x8b,
	0xcc, 0x85, 0xf1, 0x45, 0x86, 0xef, 0x20, 0xdc, 0x4a, 0x2d, 0x8c, 0x90, 0x95, 0x26, 0xfe, 0x7c,
	0xb8, 0x98, 0x2e, 0xb1, 0xb3, 0x28, 0xd7, 0xa2, 0x2d, 0x29, 0xfe, 0xeb, 0x41, 0xd0, 0xe0, 0xf7,
	0x3c, 0xe2, 0x53, 0x08, 0x7e, 0x34, 0x37, 0x36, 0xb6, 0x96, 0xfb, 0x1a, 0x09, 0x4c, 0x44, 0x65,
	0xaf, 0x82, 0x4c, 0xe6, 0xde, 0x22, 0xa0, 0x4d, 0x89, 0x73, 0x98, 0x96, 0x52, 0xf1, 0xb5, 0xeb,
	0x06, 0xb6, 0xdb, 0x85, 0x62, 0x01, 0x67, 0x07, 0x2f, 0x14, 0x9f, 0xc0, 0x98, 0xa5, 0xea, 0x76,
	0xb6, 0x5d, 0x58, 0x57, 0x1d, 0x19, 0xbb, 0x1b, 0xd3, 0xef, 0xc7, 0x8c, 0x5f, 0x42, 0xb8, 0x3f,
	0xa6, 0x8e, 0xbd, 0xd7, 0xb5, 0x8f, 0x97, 0x30, 0xdb, 0x93, 0x28, 0xd7, 0xb2, 0xd2, 0x1c, 0x9f,
	0x03, 0xe8, 0x9a, 0x31, 0xae, 0x75, 0x5e, 0x17, 0x96, 0x1f, 0xd0, 0x0e, 0x12, 0xff, 0xf2, 0x20,
	0xea, 0xff, 0x29, 0x8e, 0xa4, 0x7c, 0x0d, 0x51, 0x93, 0x6a, 0x37, 0xbc, 0xcb, 0x7a, 0x07, 0xc5,
	0x37, 0x30, 0xcb, 0xeb, 0x22, 0x17, 0x45, 0xc9, 0x2b, 0xe3, 0xce, 0x6b, 0x68, 0x99, 0xff, 0xe1,
	0xdf, 0xc7, 0xf6, 0xbb, 0x71, 0xf1, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xa0, 0x29, 0xd2, 0x75, 0x49,
	0x04, 0x00, 0x00,
}
