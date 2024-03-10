package handler

import (
	"context"
	"goods-server/global"
	"goods-server/model"
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
// 现在用户提交订单有多个商品，你得批量查询商品的信息吧
func (s *GoodsService) BatchGetGoods(ctx context.Context, req *proto.BatchGoodsIdInfo) (*proto.GoodsListResponse, error) {
	goodsListResponse := &proto.GoodsListResponse{}
	var goods []model.Goods

	//调用where并不会真正执行sql 只是用来生成sql的 当调用find， first才会去执行sql，
	result := global.DB.Where(req.Id).Find(&goods)
	for _, good := range goods {
		goodsInfoResponse := ModelToResponse(good)
		goodsListResponse.Data = append(goodsListResponse.Data, &goodsInfoResponse)
	}
	goodsListResponse.Total = int32(result.RowsAffected)
	return goodsListResponse, nil
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

func ModelToResponse(goods model.Goods) proto.GoodsInfoResponse {
	return proto.GoodsInfoResponse{
		Id:              goods.ID,
		CategoryId:      goods.CategoryId,
		Name:            goods.Name,
		GoodsSn:         goods.GoodsSn,
		ClickNum:        goods.ClickNum,
		SoldNum:         goods.SoldNum,
		FavNum:          goods.FavNum,
		MarketPrice:     goods.MarketPrice,
		ShopPrice:       goods.ShopPrice,
		GoodsBrief:      goods.GoodsBrief,
		ShipFree:        goods.ShipFree,
		GoodsFrontImage: goods.GoodsFrontImage,
		IsNew:           goods.IsNew,
		IsHot:           goods.IsHot,
		OnSale:          goods.OnSale,
		DescImages:      goods.DescImages,
		Images:          goods.Images,
		Category: &proto.CategoryBriefInfoResponse{
			Id:   goods.Category.ID,
			Name: goods.Category.Name,
		},
		Brand: &proto.BrandInfoResponse{
			Id:   goods.Brands.ID,
			Name: goods.Brands.Name,
			Logo: goods.Brands.Logo,
		},
	}
}
