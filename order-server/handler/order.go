package handler

import (
	"context"
	"encoding/json"
	"fmt"
	. "order-server/global"
	"order-server/model"
	"order-server/proto"
	. "order-server/proto"
	"order-server/utils"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
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
	o := OrderListener{}
	p, err := rocketmq.NewTransactionProducer(
		&o,
		producer.WithNameServer([]string{"192.168.1.107:9876"}),
	)
	if err != nil {
		zap.S().Errorf("生成producer失败 %s", err.Error())
		return nil, err
	}

	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败 %s", err.Error())
		return nil, err
	}
	defer p.Shutdown()

	// 一定要在这里生成好了传进去，因为ExecuteLocalTransaction，回查方法也需要这个数据，必须一致
	req.Sn = utils.GenerateOrderSn(req.UserId)
	messageStr, err := json.Marshal(req)
	if err != nil {
		zap.S().Errorf("订单序列化失败 %s", err.Error())
		return nil, err
	}

	_, err = p.SendMessageInTransaction(
		context.Background(),
		primitive.NewMessage("Order", messageStr),
	)
	if o.Code != codes.OK || err != nil {
		zap.S().Errorf("投递失败: ", o.Code, o.ErrorMsg)
		return nil, status.Error(o.Code, o.ErrorMsg)
	}

	// if err = p.Shutdown(); err != nil {
	// 	zap.S().Info("producer shutdown", err.Error())
	// 	return nil, err
	// }

	return &OrderInfoResponse{Id: int32(o.ID), OrderSn: o.OrderSn, Total: o.OrderMount}, nil
}

// OrderList implements proto.OrderServer.
func (*OrderService) OrderList(ctx context.Context, req *OrderFilterRequest) (*OrderListResponse, error) {
	var orderList []model.OrderInfo
	var total int64

	ret := DB.Model(&model.OrderInfo{}).Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	if ret.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取订单用户失败", ret.Error.Error())
	}
	ret = DB.Where(&model.OrderInfo{User: req.UserId}).Scopes(utils.Paginate(int(req.Pages), int(req.PageSize))).Find(&orderList)
	if ret.Error != nil {
		return nil, status.Errorf(codes.Internal, "获取订单列表失败", ret.Error.Error())
	}

	var resp OrderListResponse
	resp.Total = int32(total)
	for _, order := range orderList {
		fmt.Println("order:", order)
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
	DB.Where(&model.OrderGoods{Order: int32(orderInfo.ID)}).Find(&orderGoods)
	var goodsList []*OrderItemResponse
	for _, goods := range orderGoods {
		// goodsList = append(goodsList, goods.IntoOrderItemResponse())
		var orderItem IntoOrderItemResponse = &goods
		goodsList = append(goodsList, orderItem.IntoOrderItemResponse())
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

// 查询订单支付状态，如果支付成功，啥也不做，如果未支付，归还库存
func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgs {
		var req proto.OrderRequest
		json.Unmarshal(msg.Body, &req)

		zap.S().Info("获取到订单超时消息")

		var orderInfo model.OrderInfo
		ret := DB.Where(&model.OrderInfo{OrderSn: req.Sn}).First(&orderInfo)
		if ret.Error != nil {
			return consumer.ConsumeSuccess, nil
		}
		if orderInfo.Status != "TRADE_SUCCESS" {
			tx := DB.Begin()
			orderInfo.Status = "TRADE_CLOSED"
			tx.Save(&orderInfo)
			// 归还库存，发送一个归还订单的消息
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.1.107:9876"}))
			if err != nil {
				zap.S().Errorf("订单超市，生成producer失败 %s", err.Error())
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}

			if err = p.Start(); err != nil {
				zap.S().Errorf("订单超市，启动producer失败 %s", err.Error())
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			// defer p.Shutdown()

			message := primitive.NewMessage("Order", msg.Body)
			message.WithDelayTimeLevel(4)
			_, err = p.SendSync(context.Background(), message)
			if err != nil {
				zap.S().Errorf("延时消息发送失败 %s", err.Error())
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
			tx.Commit()
		}
	}

	return consumer.ConsumeSuccess, nil
}
