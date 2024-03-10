package handler

import (
	"context"
	"fmt"
	. "order-server/global"
	"order-server/model"
	. "order-server/proto"
	"order-server/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

// CreateOrder implements proto.OrderServer.
/* 新建订单
0. 从购物车中获取商品
1. 查询商品价格 - 调用商品微服务
2. 库存扣减 - 调用库存微服务
3. 订单基本信息表、订单商品表 - 当前服务数据库
4. 从购物车中删除已购买的商品
*/
func (*OrderService) CreateOrder(ctx context.Context, req *OrderRequest) (*OrderInfoResponse, error) {
	var goodsIds []int32
	var shopcarts []model.ShoppingCart
	// 查询购物车列表，获取用户选中的商品
	result := DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopcarts)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有选中任何商品")
	}
	// 获取用户选中的商品
	for _, shopcart := range shopcarts {
		goodsIds = append(goodsIds, shopcart.Goods)
	}

	// 调用商品微服务获取商品信息：价格
	goodsList, err := GoodsSvc.BatchGetGoods(ctx, &BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "购物车中商品已失效，请重新下单", err)
	}

	var goodsInfos = make([]*GoodsStockInfo, 0)
	var orderAmount float32
	var orderGoods = make([]*model.OrderGoods, 0)
	var shopMap = make(map[int32]int32, 0)
	for _, v := range shopcarts {
		shopMap[v.Goods] = v.Nums
	}

	for _, goods := range goodsList.Data {
		// 商品价格 * 数量=订单总金额
		orderAmount += goods.ShopPrice * float32(shopMap[goods.Id])
		// 订单中的商品
		orderGoods = append(orderGoods, &model.OrderGoods{
			// Order:      req.Id, 这里还没有生成订单
			Goods:      goods.Id,
			GoodsName:  goods.Name,
			GoodsImage: goods.GoodsFrontImage,
			GoodsPrice: goods.ShopPrice,
			Nums:       shopMap[goods.Id],
		})

		goodsInfos = append(goodsInfos, &GoodsStockInfo{
			GoodsId: goods.Id,
			Num:     shopMap[goods.Id],
		})
	}

	for _, v := range goodsInfos {
		fmt.Println("v:", v)

	}

	// 扣减库存，调用库存微服务
	_, err = StockSvc.Sell(ctx, &SellInfo{GoodsInfo: goodsInfos})
	if err != nil {
		return nil, status.Errorf(codes.ResourceExhausted, "库存不足，扣减库存失败", err.Error())
	}

	// 生成订单，包含订单基本信息表
	order := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      utils.GenerateOrderSn(req.UserId),
		OrderMount:   orderAmount,
		Address:      req.Address,
		Post:         req.Post,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
	}
	// #本地事务
	tx := DB.Begin()
	if result := tx.Save(&order); result.Error != nil {
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "订单生成失败")
	}

	fmt.Println("===============orderGoods:", len(orderGoods), order.ID)
	for _, o := range orderGoods {
		o.Order = int32(order.ID)
	}
	//  订单商品，批量每次插入100条
	if result := DB.CreateInBatches(orderGoods, 100); result.Error != nil {
		fmt.Println("result.Error:", result.Error)
		tx.Rollback()
		return nil, status.Errorf(codes.Internal, "订单商品生成失败")
	}
	tx.Commit()
	// ##本地事务结束

	// 从购物车，删除已购买的商品
	DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{})

	return &OrderInfoResponse{Id: int32(order.ID), OrderSn: order.OrderSn, Total: order.OrderMount}, nil
}

// OrderList implements proto.OrderServer.
func (*OrderService) OrderList(ctx context.Context, req *OrderFilterRequest) (*OrderListResponse, error) {
	var orderList []model.OrderInfo
	var total int64

	DB.Where(model.OrderInfo{User: req.UserId}).Count(&total)
	DB.Scopes(utils.Paginate(int(req.Pages), int(req.PageSize))).Find(&orderList)

	var resp OrderListResponse
	resp.Total = int32(total)
	for _, order := range orderList {
		resp.Data = append(resp.Data, order.IntoOrderInfoResponse())
	}

	return &resp, nil
}

// OrderDetail implements proto.OrderServer.
func (*OrderService) OrderDetail(ctx context.Context, req *OrderRequest) (*OrderInfoDetailResponse, error) {
	var orderInfo model.OrderInfo

	// 需要检查用户ID，防止爬虫查询其他用户的订单
	// 管理平台直接不传用户ID
	BaseModel := gorm.Model{ID: uint(req.Id)}
	ret := DB.Where(&model.OrderInfo{User: req.UserId, Model: BaseModel}).First(&orderInfo)
	if ret.Error != nil {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	var orderGoods []model.OrderGoods
	DB.Where("order_id = ?", orderInfo.ID).Find(&orderGoods)
	var goodsList []*OrderItemResponse
	for _, goods := range orderGoods {
		goodsList = append(goodsList, goods.IntoOrderItemResponse())
		// var orderItem IntoOrderItemResponse = &goods
		// goodsList = append(goodsList, orderItem.IntoOrderItemResponse())
	}

	return &OrderInfoDetailResponse{OrderInfo: orderInfo.IntoOrderInfoResponse(), Goods: goodsList}, nil
}

// UpdateOrderStatus implements proto.OrderServer.
func (*OrderService) UpdateOrderStatus(ctx context.Context, req *OrderStatus) (*emptypb.Empty, error) {
	// 支付时order_sn会传到支付宝，成功以后支付宝会回传order_sn
	err := DB.Model(&model.OrderInfo{}).Where("order_sn = ?", req.OrderSn).Update("status", req.Status)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}

	return &emptypb.Empty{}, nil
}
