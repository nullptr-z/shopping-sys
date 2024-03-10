// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: stock.proto

package proto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Stock_SetStock_FullMethodName  = "/Stock/SetStock"
	Stock_InvDetail_FullMethodName = "/Stock/InvDetail"
	Stock_Sell_FullMethodName      = "/Stock/Sell"
	Stock_Reback_FullMethodName    = "/Stock/Reback"
)

// StockClient is the client API for Stock service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StockClient interface {
	SetStock(ctx context.Context, in *GoodsStockInfo, opts ...grpc.CallOption) (*empty.Empty, error)
	InvDetail(ctx context.Context, in *GoodsStockInfo, opts ...grpc.CallOption) (*GoodsStockInfo, error)
	// 我们一般买东西的时候喜欢从购物车中去买，事务
	Sell(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*empty.Empty, error)
	Reback(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*empty.Empty, error)
}

type stockClient struct {
	cc grpc.ClientConnInterface
}

func NewStockClient(cc grpc.ClientConnInterface) StockClient {
	return &stockClient{cc}
}

func (c *stockClient) SetStock(ctx context.Context, in *GoodsStockInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Stock_SetStock_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) InvDetail(ctx context.Context, in *GoodsStockInfo, opts ...grpc.CallOption) (*GoodsStockInfo, error) {
	out := new(GoodsStockInfo)
	err := c.cc.Invoke(ctx, Stock_InvDetail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) Sell(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Stock_Sell_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stockClient) Reback(ctx context.Context, in *SellInfo, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Stock_Reback_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StockServer is the server API for Stock service.
// All implementations must embed UnimplementedStockServer
// for forward compatibility
type StockServer interface {
	SetStock(context.Context, *GoodsStockInfo) (*empty.Empty, error)
	InvDetail(context.Context, *GoodsStockInfo) (*GoodsStockInfo, error)
	// 我们一般买东西的时候喜欢从购物车中去买，事务
	Sell(context.Context, *SellInfo) (*empty.Empty, error)
	Reback(context.Context, *SellInfo) (*empty.Empty, error)
	mustEmbedUnimplementedStockServer()
}

// UnimplementedStockServer must be embedded to have forward compatible implementations.
type UnimplementedStockServer struct {
}

func (UnimplementedStockServer) SetStock(context.Context, *GoodsStockInfo) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetStock not implemented")
}
func (UnimplementedStockServer) InvDetail(context.Context, *GoodsStockInfo) (*GoodsStockInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InvDetail not implemented")
}
func (UnimplementedStockServer) Sell(context.Context, *SellInfo) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sell not implemented")
}
func (UnimplementedStockServer) Reback(context.Context, *SellInfo) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reback not implemented")
}
func (UnimplementedStockServer) mustEmbedUnimplementedStockServer() {}

// UnsafeStockServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StockServer will
// result in compilation errors.
type UnsafeStockServer interface {
	mustEmbedUnimplementedStockServer()
}

func RegisterStockServer(s grpc.ServiceRegistrar, srv StockServer) {
	s.RegisterService(&Stock_ServiceDesc, srv)
}

func _Stock_SetStock_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoodsStockInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).SetStock(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Stock_SetStock_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).SetStock(ctx, req.(*GoodsStockInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_InvDetail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GoodsStockInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).InvDetail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Stock_InvDetail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).InvDetail(ctx, req.(*GoodsStockInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_Sell_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).Sell(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Stock_Sell_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).Sell(ctx, req.(*SellInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _Stock_Reback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SellInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StockServer).Reback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Stock_Reback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StockServer).Reback(ctx, req.(*SellInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Stock_ServiceDesc is the grpc.ServiceDesc for Stock service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Stock_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Stock",
	HandlerType: (*StockServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetStock",
			Handler:    _Stock_SetStock_Handler,
		},
		{
			MethodName: "InvDetail",
			Handler:    _Stock_InvDetail_Handler,
		},
		{
			MethodName: "Sell",
			Handler:    _Stock_Sell_Handler,
		},
		{
			MethodName: "Reback",
			Handler:    _Stock_Reback_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "stock.proto",
}
