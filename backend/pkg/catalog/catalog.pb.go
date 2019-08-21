// Code generated by protoc-gen-go. DO NOT EDIT.
// source: catalog.proto

package catalog

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

type CatalogPage struct {
	Products             []*Product `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
	TotalItems           int64      `protobuf:"varint,2,opt,name=totalItems,proto3" json:"totalItems,omitempty"`
	TotalPages           int64      `protobuf:"varint,3,opt,name=totalPages,proto3" json:"totalPages,omitempty"`
	CurrentPage          int64      `protobuf:"varint,4,opt,name=currentPage,proto3" json:"currentPage,omitempty"`
	SetPageTo            int64      `protobuf:"varint,5,opt,name=setPageTo,proto3" json:"setPageTo,omitempty"`
	Sorting              string     `protobuf:"bytes,6,opt,name=sorting,proto3" json:"sorting,omitempty"`
	Filtering            string     `protobuf:"bytes,7,opt,name=filtering,proto3" json:"filtering,omitempty"`
	ItemsPerPage         int64      `protobuf:"varint,8,opt,name=itemsPerPage,proto3" json:"itemsPerPage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *CatalogPage) Reset()         { *m = CatalogPage{} }
func (m *CatalogPage) String() string { return proto.CompactTextString(m) }
func (*CatalogPage) ProtoMessage()    {}
func (*CatalogPage) Descriptor() ([]byte, []int) {
	return fileDescriptor_catalog_c7372d41ac8a6f2e, []int{0}
}
func (m *CatalogPage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CatalogPage.Unmarshal(m, b)
}
func (m *CatalogPage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CatalogPage.Marshal(b, m, deterministic)
}
func (dst *CatalogPage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CatalogPage.Merge(dst, src)
}
func (m *CatalogPage) XXX_Size() int {
	return xxx_messageInfo_CatalogPage.Size(m)
}
func (m *CatalogPage) XXX_DiscardUnknown() {
	xxx_messageInfo_CatalogPage.DiscardUnknown(m)
}

var xxx_messageInfo_CatalogPage proto.InternalMessageInfo

func (m *CatalogPage) GetProducts() []*Product {
	if m != nil {
		return m.Products
	}
	return nil
}

func (m *CatalogPage) GetTotalItems() int64 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *CatalogPage) GetTotalPages() int64 {
	if m != nil {
		return m.TotalPages
	}
	return 0
}

func (m *CatalogPage) GetCurrentPage() int64 {
	if m != nil {
		return m.CurrentPage
	}
	return 0
}

func (m *CatalogPage) GetSetPageTo() int64 {
	if m != nil {
		return m.SetPageTo
	}
	return 0
}

func (m *CatalogPage) GetSorting() string {
	if m != nil {
		return m.Sorting
	}
	return ""
}

func (m *CatalogPage) GetFiltering() string {
	if m != nil {
		return m.Filtering
	}
	return ""
}

func (m *CatalogPage) GetItemsPerPage() int64 {
	if m != nil {
		return m.ItemsPerPage
	}
	return 0
}

type Product struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Price                int64    `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	Title                string   `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Description          string   `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	Longtext             string   `protobuf:"bytes,5,opt,name=longtext,proto3" json:"longtext,omitempty"`
	Category             string   `protobuf:"bytes,6,opt,name=category,proto3" json:"category,omitempty"`
	SmallImageURL        string   `protobuf:"bytes,7,opt,name=smallImageURL,proto3" json:"smallImageURL,omitempty"`
	LargeImageURL        string   `protobuf:"bytes,8,opt,name=largeImageURL,proto3" json:"largeImageURL,omitempty"`
	Disabled             bool     `protobuf:"varint,9,opt,name=disabled,proto3" json:"disabled,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Product) Reset()         { *m = Product{} }
func (m *Product) String() string { return proto.CompactTextString(m) }
func (*Product) ProtoMessage()    {}
func (*Product) Descriptor() ([]byte, []int) {
	return fileDescriptor_catalog_c7372d41ac8a6f2e, []int{1}
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

func (m *Product) GetTitle() string {
	if m != nil {
		return m.Title
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

func init() {
	proto.RegisterType((*CatalogPage)(nil), "catalog.CatalogPage")
	proto.RegisterType((*Product)(nil), "catalog.Product")
}

func init() { proto.RegisterFile("catalog.proto", fileDescriptor_catalog_c7372d41ac8a6f2e) }

var fileDescriptor_catalog_c7372d41ac8a6f2e = []byte{
	// 326 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0x92, 0x41, 0x6a, 0xc3, 0x30,
	0x10, 0x45, 0xb1, 0xd3, 0xc4, 0xf6, 0xa4, 0x29, 0x45, 0x74, 0x21, 0x4a, 0x29, 0x26, 0x74, 0xe1,
	0x45, 0xc9, 0xa2, 0x3d, 0x42, 0x57, 0x81, 0x2e, 0x82, 0x68, 0x0f, 0xa0, 0xd8, 0x53, 0x21, 0x50,
	0x2c, 0x23, 0x4d, 0xa0, 0x3d, 0x41, 0x8f, 0xd0, 0xeb, 0x16, 0xc9, 0x8e, 0xe3, 0xec, 0xf4, 0xdf,
	0xff, 0x33, 0xe3, 0x0f, 0x86, 0x55, 0x2d, 0x49, 0x1a, 0xab, 0x36, 0x9d, 0xb3, 0x64, 0x59, 0x36,
	0xc8, 0xf5, 0x5f, 0x0a, 0xcb, 0xb7, 0xfe, 0xbd, 0x93, 0x0a, 0xd9, 0x33, 0xe4, 0x9d, 0xb3, 0xcd,
	0xb1, 0x26, 0xcf, 0x93, 0x72, 0x56, 0x2d, 0x5f, 0x6e, 0x37, 0xa7, 0xd1, 0x5d, 0x6f, 0x88, 0x31,
	0xc1, 0x1e, 0x01, 0xc8, 0x92, 0x34, 0x5b, 0xc2, 0x83, 0xe7, 0x69, 0x99, 0x54, 0x33, 0x31, 0x21,
	0xa3, 0x1f, 0x56, 0x7b, 0x3e, 0x9b, 0xf8, 0x91, 0xb0, 0x12, 0x96, 0xf5, 0xd1, 0x39, 0x6c, 0x29,
	0x68, 0x7e, 0x15, 0x03, 0x53, 0xc4, 0x1e, 0xa0, 0xf0, 0x18, 0x9f, 0x1f, 0x96, 0xcf, 0xa3, 0x7f,
	0x06, 0x8c, 0x43, 0xe6, 0xad, 0x23, 0xdd, 0x2a, 0xbe, 0x28, 0x93, 0xaa, 0x10, 0x27, 0x19, 0xe6,
	0xbe, 0xb4, 0x21, 0x74, 0xc1, 0xcb, 0xa2, 0x77, 0x06, 0x6c, 0x0d, 0xd7, 0x3a, 0x7c, 0xe0, 0x0e,
	0x5d, 0x3c, 0x9c, 0xc7, 0xc5, 0x17, 0x6c, 0xfd, 0x9b, 0x42, 0x36, 0x34, 0x66, 0x37, 0x90, 0xea,
	0x86, 0x27, 0x71, 0x4d, 0xaa, 0x1b, 0x76, 0x07, 0xf3, 0xce, 0xe9, 0x1a, 0x87, 0xca, 0xbd, 0x08,
	0x94, 0x34, 0x19, 0x8c, 0x45, 0x0b, 0xd1, 0x8b, 0xd0, 0xb1, 0x41, 0x5f, 0x3b, 0xdd, 0x91, 0xb6,
	0x6d, 0xec, 0x58, 0x88, 0x29, 0x62, 0xf7, 0x90, 0x1b, 0xdb, 0x2a, 0xc2, 0x6f, 0x8a, 0x15, 0x0b,
	0x31, 0xea, 0xe0, 0xd5, 0x92, 0x50, 0x59, 0xf7, 0x33, 0x54, 0x1c, 0x35, 0x7b, 0x82, 0x95, 0x3f,
	0x48, 0x63, 0xb6, 0x07, 0xa9, 0xf0, 0x53, 0xbc, 0x0f, 0x3d, 0x2f, 0x61, 0x48, 0x19, 0xe9, 0x14,
	0x8e, 0xa9, 0xbc, 0x4f, 0x5d, 0xc0, 0x70, 0xa7, 0xd1, 0x5e, 0xee, 0x0d, 0x36, 0xbc, 0x28, 0x93,
	0x2a, 0x17, 0xa3, 0xde, 0x2f, 0xe2, 0x3f, 0xf3, 0xfa, 0x1f, 0x00, 0x00, 0xff, 0xff, 0x33, 0xbf,
	0x54, 0x1d, 0x44, 0x02, 0x00, 0x00,
}