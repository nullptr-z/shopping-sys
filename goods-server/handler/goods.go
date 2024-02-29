package handler

import (
	"context"
	"goods-server/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

type GoodsService struct {
	proto.UnimplementedGoodsServer // 嵌入这个结构体
}

// mustEmbedUnimplementedGoodsServer implements proto.GoodsServer.
func (*GoodsService) mustEmbedUnimplementedGoodsServer() {
	panic("unimplemented")
}

// BatchGetGoods implements proto.GoodsServer.
func (*GoodsService) BatchGetGoods(context.Context, *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	panic("unimplemented")
}

// CreateGoods implements proto.GoodsServer.
func (*GoodsService) CreateGoods(context.Context, *proto.CreateGoodsInfo) (*proto.GoodsInfoResponse, error) {
	panic("unimplemented")
}

// DeleteGoods implements proto.GoodsServer.
func (*GoodsService) DeleteGoods(context.Context, *proto.DeleteGoodsInfo) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// GetGoodsDetail implements proto.GoodsServer.
func (*GoodsService) GetGoodsDetail(context.Context, *proto.GoodInfoRequest) (*proto.GoodsInfoResponse, error) {
	panic("unimplemented")
}

// GoodsList implements proto.GoodsServer.
func (*GoodsService) GoodsList(context.Context, *proto.GoodsFilterRequest) (*proto.GoodsListResponse, error) {
	panic("unimplemented")
}

// UpdateGoods implements proto.GoodsServer.
func (*GoodsService) UpdateGoods(context.Context, *proto.CreateGoodsInfo) (*emptypb.Empty, error) {
	panic("unimplemented")
}
