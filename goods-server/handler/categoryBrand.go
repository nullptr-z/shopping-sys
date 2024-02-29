// 商品分类
package handler

import (
	"context"
	"goods-server/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// CategoryBrandList implements proto.GoodsServer.
func (*GoodsService) CategoryBrandList(context.Context, *proto.CategoryBrandFilterRequest) (*proto.CategoryBrandListResponse, error) {
	panic("unimplemented")
}

// CreateCategoryBrand implements proto.GoodsServer.
func (*GoodsService) CreateCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*proto.CategoryBrandResponse, error) {
	panic("unimplemented")
}

// DeleteCategoryBrand implements proto.GoodsServer.
func (*GoodsService) DeleteCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// GetCategoryBrandList implements proto.GoodsServer.
func (*GoodsService) GetCategoryBrandList(context.Context, *proto.CategoryInfoRequest) (*proto.BrandListResponse, error) {
	panic("unimplemented")
}

// UpdateCategoryBrand implements proto.GoodsServer.
func (*GoodsService) UpdateCategoryBrand(context.Context, *proto.CategoryBrandRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}
