package handler

import (
	"context"
	"goods-server/proto"

	"google.golang.org/protobuf/types/known/emptypb"
)

// BannerList implements proto.GoodsServer.
func (*GoodsService) BannerList(context.Context, *emptypb.Empty) (*proto.BannerListResponse, error) {
	panic("unimplemented")
}

// CreateBanner implements proto.GoodsServer.
func (*GoodsService) CreateBanner(context.Context, *proto.BannerRequest) (*proto.BannerResponse, error) {
	panic("unimplemented")
}

// DeleteBanner implements proto.GoodsServer.
func (*GoodsService) DeleteBanner(context.Context, *proto.BannerRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}

// UpdateBanner implements proto.GoodsServer.
func (*GoodsService) UpdateBanner(context.Context, *proto.BannerRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}
