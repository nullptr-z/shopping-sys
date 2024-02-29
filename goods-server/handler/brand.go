// 分类

package handler

import (
	"context"
	"goods-server/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateBrand implements proto.GoodsServer.
func (*GoodsService) CreateBrand(context.Context, *proto.BrandRequest) (*proto.BrandInfoResponse, error) {
	panic("unimplemented")
}

// BrandList implements proto.GoodsServer.
func (*GoodsService) BrandList(context.Context, *proto.BrandFilterRequest) (*proto.BrandListResponse, error) {
	panic("unimplemented")
}

// DeleteBrand implements proto.GoodsServer.
func (*GoodsService) DeleteBrand(context.Context, *proto.BrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateBrand implements proto.GoodsServer.
func (*GoodsService) UpdateBrand(context.Context, *proto.BrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}
