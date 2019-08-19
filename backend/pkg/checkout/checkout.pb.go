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

type CheckoutContext struct {
	// Types that are valid to be assigned to CheckoutContextMsg:
	//	*CheckoutContext_ChangeProductQuantity
	//	*CheckoutContext_Stock
	//	*CheckoutContext_Product
	//	*CheckoutContext_CartOrder
	CheckoutContextMsg   isCheckoutContext_CheckoutContextMsg `protobuf_oneof:"checkoutContextMsg"`
	XXX_NoUnkeyedLiteral struct{}                             `json:"-"`
	XXX_unrecognized     []byte                               `json:"-"`
	XXX_sizecache        int32                                `json:"-"`
}

func (m *CheckoutContext) Reset()         { *m = CheckoutContext{} }
func (m *CheckoutContext) String() string { return proto.CompactTextString(m) }
func (*CheckoutContext) ProtoMessage()    {}
func (*CheckoutContext) Descriptor() ([]byte, []int) {
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{0}
}
func (m *CheckoutContext) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckoutContext.Unmarshal(m, b)
}
func (m *CheckoutContext) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckoutContext.Marshal(b, m, deterministic)
}
func (dst *CheckoutContext) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckoutContext.Merge(dst, src)
}
func (m *CheckoutContext) XXX_Size() int {
	return xxx_messageInfo_CheckoutContext.Size(m)
}
func (m *CheckoutContext) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckoutContext.DiscardUnknown(m)
}

var xxx_messageInfo_CheckoutContext proto.InternalMessageInfo

type isCheckoutContext_CheckoutContextMsg interface {
	isCheckoutContext_CheckoutContextMsg()
}

type CheckoutContext_ChangeProductQuantity struct {
	ChangeProductQuantity *ChangeProductQuantity `protobuf:"bytes,1,opt,name=changeProductQuantity,proto3,oneof"`
}

type CheckoutContext_Stock struct {
	Stock *Stock `protobuf:"bytes,2,opt,name=stock,proto3,oneof"`
}

type CheckoutContext_Product struct {
	Product *Product `protobuf:"bytes,3,opt,name=product,proto3,oneof"`
}

type CheckoutContext_CartOrder struct {
	CartOrder *OrderCart `protobuf:"bytes,4,opt,name=cartOrder,proto3,oneof"`
}

func (*CheckoutContext_ChangeProductQuantity) isCheckoutContext_CheckoutContextMsg() {}

func (*CheckoutContext_Stock) isCheckoutContext_CheckoutContextMsg() {}

func (*CheckoutContext_Product) isCheckoutContext_CheckoutContextMsg() {}

func (*CheckoutContext_CartOrder) isCheckoutContext_CheckoutContextMsg() {}

func (m *CheckoutContext) GetCheckoutContextMsg() isCheckoutContext_CheckoutContextMsg {
	if m != nil {
		return m.CheckoutContextMsg
	}
	return nil
}

func (m *CheckoutContext) GetChangeProductQuantity() *ChangeProductQuantity {
	if x, ok := m.GetCheckoutContextMsg().(*CheckoutContext_ChangeProductQuantity); ok {
		return x.ChangeProductQuantity
	}
	return nil
}

func (m *CheckoutContext) GetStock() *Stock {
	if x, ok := m.GetCheckoutContextMsg().(*CheckoutContext_Stock); ok {
		return x.Stock
	}
	return nil
}

func (m *CheckoutContext) GetProduct() *Product {
	if x, ok := m.GetCheckoutContextMsg().(*CheckoutContext_Product); ok {
		return x.Product
	}
	return nil
}

func (m *CheckoutContext) GetCartOrder() *OrderCart {
	if x, ok := m.GetCheckoutContextMsg().(*CheckoutContext_CartOrder); ok {
		return x.CartOrder
	}
	return nil
}

// XXX_OneofFuncs is for the internal use of the proto package.
func (*CheckoutContext) XXX_OneofFuncs() (func(msg proto.Message, b *proto.Buffer) error, func(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error), func(msg proto.Message) (n int), []interface{}) {
	return _CheckoutContext_OneofMarshaler, _CheckoutContext_OneofUnmarshaler, _CheckoutContext_OneofSizer, []interface{}{
		(*CheckoutContext_ChangeProductQuantity)(nil),
		(*CheckoutContext_Stock)(nil),
		(*CheckoutContext_Product)(nil),
		(*CheckoutContext_CartOrder)(nil),
	}
}

func _CheckoutContext_OneofMarshaler(msg proto.Message, b *proto.Buffer) error {
	m := msg.(*CheckoutContext)
	// checkoutContextMsg
	switch x := m.CheckoutContextMsg.(type) {
	case *CheckoutContext_ChangeProductQuantity:
		b.EncodeVarint(1<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.ChangeProductQuantity); err != nil {
			return err
		}
	case *CheckoutContext_Stock:
		b.EncodeVarint(2<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Stock); err != nil {
			return err
		}
	case *CheckoutContext_Product:
		b.EncodeVarint(3<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.Product); err != nil {
			return err
		}
	case *CheckoutContext_CartOrder:
		b.EncodeVarint(4<<3 | proto.WireBytes)
		if err := b.EncodeMessage(x.CartOrder); err != nil {
			return err
		}
	case nil:
	default:
		return fmt.Errorf("CheckoutContext.CheckoutContextMsg has unexpected type %T", x)
	}
	return nil
}

func _CheckoutContext_OneofUnmarshaler(msg proto.Message, tag, wire int, b *proto.Buffer) (bool, error) {
	m := msg.(*CheckoutContext)
	switch tag {
	case 1: // checkoutContextMsg.changeProductQuantity
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(ChangeProductQuantity)
		err := b.DecodeMessage(msg)
		m.CheckoutContextMsg = &CheckoutContext_ChangeProductQuantity{msg}
		return true, err
	case 2: // checkoutContextMsg.stock
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Stock)
		err := b.DecodeMessage(msg)
		m.CheckoutContextMsg = &CheckoutContext_Stock{msg}
		return true, err
	case 3: // checkoutContextMsg.product
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(Product)
		err := b.DecodeMessage(msg)
		m.CheckoutContextMsg = &CheckoutContext_Product{msg}
		return true, err
	case 4: // checkoutContextMsg.cartOrder
		if wire != proto.WireBytes {
			return true, proto.ErrInternalBadWireType
		}
		msg := new(OrderCart)
		err := b.DecodeMessage(msg)
		m.CheckoutContextMsg = &CheckoutContext_CartOrder{msg}
		return true, err
	default:
		return false, nil
	}
}

func _CheckoutContext_OneofSizer(msg proto.Message) (n int) {
	m := msg.(*CheckoutContext)
	// checkoutContextMsg
	switch x := m.CheckoutContextMsg.(type) {
	case *CheckoutContext_ChangeProductQuantity:
		s := proto.Size(x.ChangeProductQuantity)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CheckoutContext_Stock:
		s := proto.Size(x.Stock)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CheckoutContext_Product:
		s := proto.Size(x.Product)
		n += 1 // tag and wire
		n += proto.SizeVarint(uint64(s))
		n += s
	case *CheckoutContext_CartOrder:
		s := proto.Size(x.CartOrder)
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
	Title                string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	SmallImageURL        string   `protobuf:"bytes,4,opt,name=smallImageURL,proto3" json:"smallImageURL,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{1}
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

func (m *Product) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *Product) GetSmallImageURL() string {
	if m != nil {
		return m.SmallImageURL
	}
	return ""
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
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{2}
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
	Title                string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
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
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{3}
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

func (m *Position) GetTitle() string {
	if m != nil {
		return m.Title
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
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{4}
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
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{5}
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
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{6}
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

type Stock struct {
	ProductID            string   `protobuf:"bytes,1,opt,name=productID,proto3" json:"productID,omitempty"`
	Quantity             int64    `protobuf:"varint,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Stock) Reset()         { *m = Stock{} }
func (m *Stock) String() string { return proto.CompactTextString(m) }
func (*Stock) ProtoMessage()    {}
func (*Stock) Descriptor() ([]byte, []int) {
	return fileDescriptor_checkout_f71b3428b3a3f15f, []int{7}
}
func (m *Stock) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Stock.Unmarshal(m, b)
}
func (m *Stock) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Stock.Marshal(b, m, deterministic)
}
func (dst *Stock) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Stock.Merge(dst, src)
}
func (m *Stock) XXX_Size() int {
	return xxx_messageInfo_Stock.Size(m)
}
func (m *Stock) XXX_DiscardUnknown() {
	xxx_messageInfo_Stock.DiscardUnknown(m)
}

var xxx_messageInfo_Stock proto.InternalMessageInfo

func (m *Stock) GetProductID() string {
	if m != nil {
		return m.ProductID
	}
	return ""
}

func (m *Stock) GetQuantity() int64 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

func init() {
	proto.RegisterType((*CheckoutContext)(nil), "checkout.CheckoutContext")
	proto.RegisterType((*Product)(nil), "checkout.Product")
	proto.RegisterType((*Cart)(nil), "checkout.Cart")
	proto.RegisterType((*Position)(nil), "checkout.Position")
	proto.RegisterType((*ChangeProductQuantity)(nil), "checkout.ChangeProductQuantity")
	proto.RegisterType((*OrderCart)(nil), "checkout.OrderCart")
	proto.RegisterType((*OrderCartResonse)(nil), "checkout.OrderCartResonse")
	proto.RegisterType((*Stock)(nil), "checkout.Stock")
}

func init() { proto.RegisterFile("checkout.proto", fileDescriptor_checkout_f71b3428b3a3f15f) }

var fileDescriptor_checkout_f71b3428b3a3f15f = []byte{
	// 417 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x53, 0xdd, 0x8a, 0xd4, 0x30,
	0x18, 0x9d, 0x76, 0x76, 0x66, 0xda, 0x6f, 0x70, 0x57, 0x3f, 0x77, 0xa5, 0x88, 0x68, 0x89, 0x82,
	0x73, 0xe3, 0x22, 0xb3, 0x4f, 0xa0, 0xdd, 0x8b, 0x16, 0x14, 0x35, 0x22, 0x5e, 0xd7, 0x34, 0xce,
	0x86, 0xed, 0x34, 0x35, 0x49, 0xc1, 0x7d, 0x08, 0x5f, 0xce, 0x27, 0x92, 0xa6, 0xe9, 0xcf, 0xc0,
	0x80, 0xe0, 0x85, 0x97, 0xe7, 0x7c, 0xe7, 0x3b, 0x39, 0xa7, 0x49, 0xe1, 0x94, 0xdd, 0x70, 0x76,
	0x2b, 0x1b, 0x73, 0x59, 0x2b, 0x69, 0x24, 0x06, 0x3d, 0x26, 0xbf, 0x7c, 0x38, 0x4b, 0x1c, 0x48,
	0x64, 0x65, 0xf8, 0x4f, 0x83, 0x5f, 0xe1, 0x82, 0xdd, 0xe4, 0xd5, 0x8e, 0x7f, 0x54, 0xb2, 0x68,
	0x98, 0xf9, 0xd4, 0xe4, 0x95, 0x11, 0xe6, 0x2e, 0xf2, 0x62, 0x6f, 0xb3, 0xde, 0x3e, 0xbb, 0x1c,
	0xdc, 0x92, 0x63, 0xb2, 0x74, 0x46, 0x8f, 0xef, 0xe3, 0x4b, 0x58, 0x68, 0x23, 0xd9, 0x6d, 0xe4,
	0x5b, 0xa3, 0xb3, 0xd1, 0xe8, 0x73, 0x4b, 0xa7, 0x33, 0xda, 0xcd, 0xf1, 0x15, 0xac, 0xea, 0x6e,
	0x37, 0x9a, 0x5b, 0xe9, 0x83, 0x51, 0xea, 0x4c, 0xd3, 0x19, 0xed, 0x35, 0x78, 0x05, 0x21, 0xcb,
	0x95, 0xf9, 0xa0, 0x0a, 0xae, 0xa2, 0x13, 0xbb, 0xf0, 0x70, 0x5c, 0xb0, 0x74, 0x92, 0xab, 0x76,
	0x65, 0xd4, 0xbd, 0x3d, 0x07, 0x64, 0x87, 0xc5, 0xdf, 0xeb, 0x1d, 0xb9, 0x83, 0x95, 0x3b, 0x00,
	0x9f, 0x40, 0xe8, 0x0e, 0xc8, 0xae, 0x6d, 0xf5, 0x90, 0x8e, 0x04, 0x9e, 0xc3, 0xa2, 0x56, 0x82,
	0x71, 0xdb, 0x65, 0x4e, 0x3b, 0xd0, 0xb2, 0x46, 0x98, 0x92, 0xdb, 0xd8, 0x21, 0xed, 0x00, 0xbe,
	0x80, 0x7b, 0x7a, 0x9f, 0x97, 0x65, 0xb6, 0xcf, 0x77, 0xfc, 0x0b, 0x7d, 0x67, 0x33, 0x86, 0xf4,
	0x90, 0x24, 0x29, 0x9c, 0xb4, 0x29, 0xf1, 0x14, 0x7c, 0x51, 0xb8, 0x03, 0x7d, 0x51, 0xe0, 0x6b,
	0x08, 0x6b, 0xa9, 0x85, 0x11, 0xb2, 0xd2, 0x91, 0x1f, 0xcf, 0x37, 0xeb, 0x2d, 0x4e, 0x3e, 0x87,
	0x1b, 0xd1, 0x51, 0x44, 0x7e, 0x7b, 0x10, 0xf4, 0xfc, 0xff, 0xae, 0x81, 0x8f, 0x21, 0xf8, 0xd1,
	0x3f, 0x98, 0xa5, 0x35, 0x1d, 0x30, 0x46, 0xb0, 0x12, 0x95, 0xbd, 0xeb, 0x68, 0x15, 0x7b, 0x9b,
	0x80, 0xf6, 0x10, 0x63, 0x58, 0xef, 0xa5, 0xe2, 0x99, 0x9b, 0x06, 0x76, 0x3a, 0xa5, 0x88, 0x80,
	0x8b, 0xa3, 0xcf, 0x0d, 0x1f, 0xc1, 0xb2, 0xbd, 0xd5, 0xec, 0xda, 0xa5, 0x75, 0xe8, 0x2f, 0xc5,
	0xa7, 0x31, 0xfd, 0xc3, 0x98, 0xe4, 0x39, 0x84, 0xc3, 0xa3, 0x99, 0xd8, 0x7b, 0x53, 0x7b, 0xb2,
	0x85, 0xfb, 0x83, 0x88, 0x72, 0x2d, 0x2b, 0xcd, 0xf1, 0x29, 0x80, 0x6e, 0x18, 0xe3, 0x5a, 0x7f,
	0x6f, 0x4a, 0xab, 0x0f, 0xe8, 0x84, 0x21, 0x6f, 0x60, 0xd1, 0xd5, 0xfd, 0xe7, 0x6c, 0xdf, 0x96,
	0xf6, 0x0f, 0xbe, 0xfa, 0x13, 0x00, 0x00, 0xff, 0xff, 0xf8, 0xe4, 0x22, 0x28, 0xd3, 0x03, 0x00,
	0x00,
}
