package handler

import (
	"context"
	. "order-server/global"
	"order-server/model"
	. "order-server/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderService struct {
	UnimplementedOrderServer
}

// CartItemList implements proto.OrderServer.
// 获取用户的购物车列表
func (*OrderService) CartItemList(ctx context.Context, req *UserInfo) (*CartItemListResponse, error) {
	var shopCarts []model.ShoppingCart
	retDb := DB.Where(model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if retDb.Error != nil {
		return nil, retDb.Error
	}

	resp := &CartItemListResponse{}

	for _, cart := range shopCarts {
		resp.Data = append(resp.Data, &ShopCartInfoResponse{
			Id:      int32(cart.ID),
			UserId:  cart.User,
			GoodsId: cart.Goods,
			Nums:    cart.Nums,
			Checked: cart.Checked,
		})
	}
	resp.Total = int32(retDb.RowsAffected)

	return resp, nil
}

// CreateCartItem implements proto.OrderServer.
// 添加商品到购物车,如果不存在则添加，存在则添加数量
func (*OrderService) CreateCartItem(ctx context.Context, req *CartItemRequest) (*ShopCartInfoResponse, error) {
	var shopCart model.ShoppingCart
	serachRet := DB.Where(model.ShoppingCart{User: req.UserId, Goods: req.GoodsId}).First(&shopCart)

	if serachRet.RowsAffected == 1 {
		shopCart.Nums += req.Nums
	} else {
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}
	creteRet := DB.Save(&shopCart)
	if creteRet.Error != nil {
		return nil, creteRet.Error
	}
	return &ShopCartInfoResponse{Id: int32(shopCart.ID)}, nil
}

// UpdateCartItem implements proto.OrderServer.
func (*OrderService) UpdateCartItem(ctx context.Context, req *CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart
	serachRet := DB.First(&shopCart, req.Id)
	if serachRet.Error != nil {
		return nil, status.Errorf(codes.NotFound, "购物车不存在这个商品")
	}
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	shopCart.Checked = req.Checked

	creteRet := DB.Save(&shopCart)
	if creteRet.Error != nil {
		return nil, status.Errorf(codes.InvalidArgument, "服务器内部错误")
	}

	return &emptypb.Empty{}, nil
}

// DeleteCartItem implements proto.OrderServer.
func (*OrderService) DeleteCartItem(ctx context.Context, req *CartItemRequest) (*emptypb.Empty, error) {
	creteRet := DB.Delete(&model.ShoppingCart{}, req.Id)
	if creteRet.Error != nil {
		return nil, status.Errorf(codes.NotFound, "购物车不存在这个商品")
	}
	return &emptypb.Empty{}, nil
}

// mustEmbedUnimplementedOrderServer implements proto.OrderServer.
func (*OrderService) mustEmbedUnimplementedOrderServer() {
	panic("unimplemented")
}
