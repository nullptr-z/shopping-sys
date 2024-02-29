package handler

import (
	"context"
	"goods-server/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// GetAllCategorysList implements proto.GoodsServer.
func (*GoodsService) GetAllCategorysList(context.Context, *emptypb.Empty) (*proto.CategoryListResponse, error) {
	panic("unimplemented")
}

// CreateCategory implements proto.GoodsServer.
func (*GoodsService) CreateCategory(context.Context, *proto.CategoryInfoRequest) (*proto.CategoryInfoResponse, error) {
	panic("unimplemented")
}

// DeleteCategory implements proto.GoodsServer.
func (*GoodsService) DeleteCategory(context.Context, *proto.DeleteCategoryRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateCategory implements proto.GoodsServer.
func (*GoodsService) UpdateCategory(context.Context, *proto.CategoryInfoRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// GetSubCategory implements proto.GoodsServer.
func (*GoodsService) GetSubCategory(context.Context, *proto.CategoryListRequest) (*proto.SubCategoryListResponse, error) {
	panic("unimplemented")
}
