// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v3.21.12
// source: stock.proto

package proto

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

type GoodsStockInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoodsId int32 `protobuf:"varint,1,opt,name=goodsId,proto3" json:"goodsId,omitempty"`
	Num     int32 `protobuf:"varint,2,opt,name=num,proto3" json:"num,omitempty"`
}

func (x *GoodsStockInfo) Reset() {
	*x = GoodsStockInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stock_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GoodsStockInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoodsStockInfo) ProtoMessage() {}

func (x *GoodsStockInfo) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoodsStockInfo.ProtoReflect.Descriptor instead.
func (*GoodsStockInfo) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{0}
}

func (x *GoodsStockInfo) GetGoodsId() int32 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

func (x *GoodsStockInfo) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

type SellInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GoodsInfo []*GoodsStockInfo `protobuf:"bytes,1,rep,name=goodsInfo,proto3" json:"goodsInfo,omitempty"`
	OrderSn   string            `protobuf:"bytes,2,opt,name=orderSn,proto3" json:"orderSn,omitempty"`
}

func (x *SellInfo) Reset() {
	*x = SellInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stock_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SellInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SellInfo) ProtoMessage() {}

func (x *SellInfo) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SellInfo.ProtoReflect.Descriptor instead.
func (*SellInfo) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{1}
}

func (x *SellInfo) GetGoodsInfo() []*GoodsStockInfo {
	if x != nil {
		return x.GoodsInfo
	}
	return nil
}

func (x *SellInfo) GetOrderSn() string {
	if x != nil {
		return x.OrderSn
	}
	return ""
}

var File_stock_proto protoreflect.FileDescriptor

var file_stock_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65,
	0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a, 0x0e, 0x47, 0x6f,
	0x6f, 0x64, 0x73, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07,
	0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x67,
	0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x03, 0x6e, 0x75, 0x6d, 0x22, 0x53, 0x0a, 0x08, 0x53, 0x65, 0x6c, 0x6c,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2d, 0x0a, 0x09, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x6e, 0x66,
	0x6f, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x09, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49,
	0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x53, 0x6e, 0x32, 0xc3, 0x01,
	0x0a, 0x05, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x33, 0x0a, 0x08, 0x53, 0x65, 0x74, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x12, 0x0f, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2d, 0x0a, 0x09,
	0x49, 0x6e, 0x76, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x0f, 0x2e, 0x47, 0x6f, 0x6f, 0x64,
	0x73, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0f, 0x2e, 0x47, 0x6f, 0x6f,
	0x64, 0x73, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x0a, 0x04, 0x53,
	0x65, 0x6c, 0x6c, 0x12, 0x09, 0x2e, 0x53, 0x65, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x2b, 0x0a, 0x06, 0x52, 0x65, 0x62, 0x61, 0x63, 0x6b,
	0x12, 0x09, 0x2e, 0x53, 0x65, 0x6c, 0x6c, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stock_proto_rawDescOnce sync.Once
	file_stock_proto_rawDescData = file_stock_proto_rawDesc
)

func file_stock_proto_rawDescGZIP() []byte {
	file_stock_proto_rawDescOnce.Do(func() {
		file_stock_proto_rawDescData = protoimpl.X.CompressGZIP(file_stock_proto_rawDescData)
	})
	return file_stock_proto_rawDescData
}

var file_stock_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_stock_proto_goTypes = []interface{}{
	(*GoodsStockInfo)(nil), // 0: GoodsStockInfo
	(*SellInfo)(nil),       // 1: SellInfo
	(*emptypb.Empty)(nil),  // 2: google.protobuf.Empty
}
var file_stock_proto_depIdxs = []int32{
	0, // 0: SellInfo.goodsInfo:type_name -> GoodsStockInfo
	0, // 1: Stock.SetStock:input_type -> GoodsStockInfo
	0, // 2: Stock.InvDetail:input_type -> GoodsStockInfo
	1, // 3: Stock.Sell:input_type -> SellInfo
	1, // 4: Stock.Reback:input_type -> SellInfo
	2, // 5: Stock.SetStock:output_type -> google.protobuf.Empty
	0, // 6: Stock.InvDetail:output_type -> GoodsStockInfo
	2, // 7: Stock.Sell:output_type -> google.protobuf.Empty
	2, // 8: Stock.Reback:output_type -> google.protobuf.Empty
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_stock_proto_init() }
func file_stock_proto_init() {
	if File_stock_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stock_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GoodsStockInfo); i {
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
		file_stock_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SellInfo); i {
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
			RawDescriptor: file_stock_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stock_proto_goTypes,
		DependencyIndexes: file_stock_proto_depIdxs,
		MessageInfos:      file_stock_proto_msgTypes,
	}.Build()
	File_stock_proto = out.File
	file_stock_proto_rawDesc = nil
	file_stock_proto_goTypes = nil
	file_stock_proto_depIdxs = nil
}
