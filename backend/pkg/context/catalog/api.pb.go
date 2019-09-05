// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

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

type CatalogRequest_Sorting int32

const (
	CatalogRequest_ID    CatalogRequest_Sorting = 0
	CatalogRequest_PRICE CatalogRequest_Sorting = 1
	CatalogRequest_NAME  CatalogRequest_Sorting = 2
)

var CatalogRequest_Sorting_name = map[int32]string{
	0: "ID",
	1: "PRICE",
	2: "NAME",
}
var CatalogRequest_Sorting_value = map[string]int32{
	"ID":    0,
	"PRICE": 1,
	"NAME":  2,
}

func (x CatalogRequest_Sorting) String() string {
	return proto.EnumName(CatalogRequest_Sorting_name, int32(x))
}
func (CatalogRequest_Sorting) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_api_231a699d8bb86116, []int{0, 0}
}

type CatalogRequest struct {
	Sorting              CatalogRequest_Sorting `protobuf:"varint,1,opt,name=sorting,proto3,enum=catalog.CatalogRequest_Sorting" json:"sorting,omitempty"`
	Prefix               string                 `protobuf:"bytes,2,opt,name=prefix,proto3" json:"prefix,omitempty"`
	Page                 int64                  `protobuf:"varint,3,opt,name=page,proto3" json:"page,omitempty"`
	ItemsPerPage         int64                  `protobuf:"varint,4,opt,name=itemsPerPage,proto3" json:"itemsPerPage,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *CatalogRequest) Reset()         { *m = CatalogRequest{} }
func (m *CatalogRequest) String() string { return proto.CompactTextString(m) }
func (*CatalogRequest) ProtoMessage()    {}
func (*CatalogRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_231a699d8bb86116, []int{0}
}
func (m *CatalogRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CatalogRequest.Unmarshal(m, b)
}
func (m *CatalogRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CatalogRequest.Marshal(b, m, deterministic)
}
func (dst *CatalogRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CatalogRequest.Merge(dst, src)
}
func (m *CatalogRequest) XXX_Size() int {
	return xxx_messageInfo_CatalogRequest.Size(m)
}
func (m *CatalogRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CatalogRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CatalogRequest proto.InternalMessageInfo

func (m *CatalogRequest) GetSorting() CatalogRequest_Sorting {
	if m != nil {
		return m.Sorting
	}
	return CatalogRequest_ID
}

func (m *CatalogRequest) GetPrefix() string {
	if m != nil {
		return m.Prefix
	}
	return ""
}

func (m *CatalogRequest) GetPage() int64 {
	if m != nil {
		return m.Page
	}
	return 0
}

func (m *CatalogRequest) GetItemsPerPage() int64 {
	if m != nil {
		return m.ItemsPerPage
	}
	return 0
}

type CatalogResponse struct {
	Request              *CatalogRequest    `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	Products             []*ProductResponse `protobuf:"bytes,9,rep,name=products,proto3" json:"products,omitempty"`
	TotalItems           int64              `protobuf:"varint,2,opt,name=totalItems,proto3" json:"totalItems,omitempty"`
	TotalPages           int64              `protobuf:"varint,3,opt,name=totalPages,proto3" json:"totalPages,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *CatalogResponse) Reset()         { *m = CatalogResponse{} }
func (m *CatalogResponse) String() string { return proto.CompactTextString(m) }
func (*CatalogResponse) ProtoMessage()    {}
func (*CatalogResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_231a699d8bb86116, []int{1}
}
func (m *CatalogResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CatalogResponse.Unmarshal(m, b)
}
func (m *CatalogResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CatalogResponse.Marshal(b, m, deterministic)
}
func (dst *CatalogResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CatalogResponse.Merge(dst, src)
}
func (m *CatalogResponse) XXX_Size() int {
	return xxx_messageInfo_CatalogResponse.Size(m)
}
func (m *CatalogResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CatalogResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CatalogResponse proto.InternalMessageInfo

func (m *CatalogResponse) GetRequest() *CatalogRequest {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *CatalogResponse) GetProducts() []*ProductResponse {
	if m != nil {
		return m.Products
	}
	return nil
}

func (m *CatalogResponse) GetTotalItems() int64 {
	if m != nil {
		return m.TotalItems
	}
	return 0
}

func (m *CatalogResponse) GetTotalPages() int64 {
	if m != nil {
		return m.TotalPages
	}
	return 0
}

type ProductRequest struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ProductRequest) Reset()         { *m = ProductRequest{} }
func (m *ProductRequest) String() string { return proto.CompactTextString(m) }
func (*ProductRequest) ProtoMessage()    {}
func (*ProductRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_231a699d8bb86116, []int{2}
}
func (m *ProductRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProductRequest.Unmarshal(m, b)
}
func (m *ProductRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProductRequest.Marshal(b, m, deterministic)
}
func (dst *ProductRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProductRequest.Merge(dst, src)
}
func (m *ProductRequest) XXX_Size() int {
	return xxx_messageInfo_ProductRequest.Size(m)
}
func (m *ProductRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProductRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProductRequest proto.InternalMessageInfo

func (m *ProductRequest) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type ProductResponse struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Price                int64    `protobuf:"varint,2,opt,name=price,proto3" json:"price,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
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

func (m *ProductResponse) Reset()         { *m = ProductResponse{} }
func (m *ProductResponse) String() string { return proto.CompactTextString(m) }
func (*ProductResponse) ProtoMessage()    {}
func (*ProductResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_api_231a699d8bb86116, []int{3}
}
func (m *ProductResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProductResponse.Unmarshal(m, b)
}
func (m *ProductResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProductResponse.Marshal(b, m, deterministic)
}
func (dst *ProductResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProductResponse.Merge(dst, src)
}
func (m *ProductResponse) XXX_Size() int {
	return xxx_messageInfo_ProductResponse.Size(m)
}
func (m *ProductResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProductResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProductResponse proto.InternalMessageInfo

func (m *ProductResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *ProductResponse) GetPrice() int64 {
	if m != nil {
		return m.Price
	}
	return 0
}

func (m *ProductResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProductResponse) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *ProductResponse) GetLongtext() string {
	if m != nil {
		return m.Longtext
	}
	return ""
}

func (m *ProductResponse) GetCategory() string {
	if m != nil {
		return m.Category
	}
	return ""
}

func (m *ProductResponse) GetSmallImageURL() string {
	if m != nil {
		return m.SmallImageURL
	}
	return ""
}

func (m *ProductResponse) GetLargeImageURL() string {
	if m != nil {
		return m.LargeImageURL
	}
	return ""
}

func (m *ProductResponse) GetDisabled() bool {
	if m != nil {
		return m.Disabled
	}
	return false
}

func init() {
	proto.RegisterType((*CatalogRequest)(nil), "catalog.CatalogRequest")
	proto.RegisterType((*CatalogResponse)(nil), "catalog.CatalogResponse")
	proto.RegisterType((*ProductRequest)(nil), "catalog.ProductRequest")
	proto.RegisterType((*ProductResponse)(nil), "catalog.ProductResponse")
	proto.RegisterEnum("catalog.CatalogRequest_Sorting", CatalogRequest_Sorting_name, CatalogRequest_Sorting_value)
}

func init() { proto.RegisterFile("api.proto", fileDescriptor_api_231a699d8bb86116) }

var fileDescriptor_api_231a699d8bb86116 = []byte{
	// 404 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0xaa, 0xd4, 0x30,
	0x14, 0xc6, 0x4d, 0xe7, 0x4f, 0x9b, 0x33, 0x3a, 0x77, 0x38, 0x88, 0x06, 0x17, 0x5a, 0x8a, 0x48,
	0x57, 0x03, 0x8e, 0x6e, 0x5c, 0xca, 0xf5, 0x2e, 0x06, 0x54, 0x86, 0x88, 0x0f, 0x90, 0xdb, 0xc6,
	0x12, 0xe8, 0x34, 0x31, 0xc9, 0x85, 0xeb, 0x43, 0xf8, 0x42, 0x6e, 0x7d, 0x31, 0x69, 0xd2, 0x09,
	0x33, 0x17, 0xdc, 0x9d, 0xf3, 0x7d, 0xbf, 0xe6, 0x7c, 0xe7, 0x50, 0xa0, 0xc2, 0xa8, 0xad, 0xb1,
	0xda, 0x6b, 0xcc, 0x1b, 0xe1, 0x45, 0xaf, 0xbb, 0xea, 0x2f, 0x81, 0xf5, 0x75, 0xac, 0xb9, 0xfc,
	0x79, 0x27, 0x9d, 0xc7, 0x0f, 0x90, 0x3b, 0x6d, 0xbd, 0x1a, 0x3a, 0x46, 0x4a, 0x52, 0xaf, 0x77,
	0xaf, 0xb6, 0x13, 0xbd, 0xbd, 0x24, 0xb7, 0xdf, 0x22, 0xc6, 0x4f, 0x3c, 0x3e, 0x83, 0xa5, 0xb1,
	0xf2, 0x87, 0xba, 0x67, 0x59, 0x49, 0x6a, 0xca, 0xa7, 0x0e, 0x11, 0xe6, 0x46, 0x74, 0x92, 0xcd,
	0x4a, 0x52, 0xcf, 0x78, 0xa8, 0xb1, 0x82, 0xc7, 0xca, 0xcb, 0xa3, 0x3b, 0x48, 0x7b, 0x18, 0xbd,
	0x79, 0xf0, 0x2e, 0xb4, 0xea, 0x0d, 0xe4, 0xd3, 0x0c, 0x5c, 0x42, 0xb6, 0xff, 0xb4, 0x79, 0x84,
	0x14, 0x16, 0x07, 0xbe, 0xbf, 0xbe, 0xd9, 0x10, 0x2c, 0x60, 0xfe, 0xf5, 0xe3, 0x97, 0x9b, 0x4d,
	0x56, 0xfd, 0x21, 0x70, 0x95, 0xb2, 0x39, 0xa3, 0x07, 0x27, 0xf1, 0x2d, 0xe4, 0x36, 0xe6, 0x0c,
	0x6b, 0xac, 0x76, 0xcf, 0xff, 0xb3, 0x06, 0x3f, 0x71, 0xf8, 0x1e, 0x0a, 0x63, 0x75, 0x7b, 0xd7,
	0x78, 0xc7, 0x68, 0x39, 0xab, 0x57, 0x3b, 0x96, 0xbe, 0x39, 0x44, 0xe3, 0xf4, 0x3c, 0x4f, 0x24,
	0xbe, 0x04, 0xf0, 0xda, 0x8b, 0x7e, 0x3f, 0x26, 0x0f, 0x8b, 0xcf, 0xf8, 0x99, 0x92, 0xfc, 0x71,
	0x23, 0x37, 0x9d, 0xe0, 0x4c, 0xa9, 0x4a, 0x58, 0xa7, 0xc7, 0x63, 0x8e, 0x35, 0x64, 0xaa, 0x0d,
	0xa9, 0x29, 0xcf, 0x54, 0x5b, 0xfd, 0xce, 0xe0, 0xea, 0xc1, 0xfc, 0x87, 0x0c, 0x3e, 0x85, 0x85,
	0xb1, 0xaa, 0x91, 0x53, 0x80, 0xd8, 0x8c, 0x87, 0x1f, 0xc4, 0x31, 0x1e, 0x9e, 0xf2, 0x50, 0x63,
	0x09, 0xab, 0x56, 0xba, 0xc6, 0x2a, 0xe3, 0x95, 0x1e, 0xc2, 0xdd, 0x29, 0x3f, 0x97, 0xf0, 0x05,
	0x14, 0xbd, 0x1e, 0x3a, 0x2f, 0xef, 0x3d, 0x5b, 0x04, 0x3b, 0xf5, 0xa3, 0xd7, 0x08, 0x2f, 0x3b,
	0x6d, 0x7f, 0xb1, 0x65, 0xf4, 0x4e, 0x3d, 0xbe, 0x86, 0x27, 0xee, 0x28, 0xfa, 0x7e, 0x7f, 0x14,
	0x9d, 0xfc, 0xce, 0x3f, 0xb3, 0x3c, 0x00, 0x97, 0xe2, 0x48, 0xf5, 0xc2, 0x76, 0x32, 0x51, 0x45,
	0xa4, 0x2e, 0xc4, 0x71, 0x4e, 0xab, 0x9c, 0xb8, 0xed, 0x65, 0xcb, 0x68, 0x49, 0xea, 0x82, 0xa7,
	0xfe, 0x76, 0x19, 0x7e, 0xe2, 0x77, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x4a, 0x53, 0x14, 0xd4,
	0xd1, 0x02, 0x00, 0x00,
}