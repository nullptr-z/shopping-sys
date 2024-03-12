package handler

import (
	"context"
	"encoding/json"
	"fmt"
	. "order-server/global"
	"order-server/model"
	"order-server/proto"
	. "order-server/proto"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

type OrderListener struct {
	model.OrderInfo
	Code     codes.Code
	ErrorMsg string
}

// 扣减库存>生成订单
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	// ------------前戏，准备各种数据，查询购物车,用户勾选的商品，商品价格--------------
	var req proto.OrderRequest
	err := json.Unmarshal(msg.Body, &req)
	if err != nil {
		o.Code = codes.Internal
		o.ErrorMsg = err.Error()
		zap.S().Error("json.Unmarshal error:", err.Error())
		return primitive.RollbackMessageState
	}

	var goodsIds []int32
	var shopcarts []model.ShoppingCart
	// 查询购物车列表，获取用户选中的商品
	result := DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Find(&shopcarts)
	if result.RowsAffected == 0 {
		return primitive.RollbackMessageState
	}
	// 获取用户选中的商品
	for _, shopcart := range shopcarts {
		goodsIds = append(goodsIds, shopcart.Goods)
	}

	// 调用商品微服务获取商品信息：价格
	goodsList, err := GoodsSvc.BatchGetGoods(context.Background(), &BatchGoodsIdInfo{Id: goodsIds})
	if err != nil {
		o.Code = codes.Internal
		o.ErrorMsg = err.Error()
		return primitive.RollbackMessageState
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

	//-------------主逻辑开始，首先扣减库存-----------------
	// 扣减库存，调用库存微服务
	_, err = StockSvc.Sell(context.Background(), &SellInfo{OrderSn: req.Sn, GoodsInfo: goodsInfos})
	if err != nil {
		// 仅处理库存扣减失败，根据状态码判断；
		// TODO: 单独处理网络问题，需要重试，并且Sell中需要加上重试逻辑，如果之前已经扣减成功了，直接返回
		o.Code = codes.ResourceExhausted
		o.ErrorMsg = "库存扣减失败，" + err.Error()
		return primitive.RollbackMessageState
	}

	// -----------------至此，库存扣减成功，接下来开始生成订单----------------------
	zap.S().Info("库存扣减成功")
	// o.Code = codes.Internal
	// return primitive.UnknowState

	// 生成订单，包含订单基本信息表
	order := model.OrderInfo{
		User:         req.UserId,
		OrderSn:      req.Sn,
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
		o.Code = codes.Internal
		o.ErrorMsg = "生成订单失败，" + result.Error.Error()
		return primitive.CommitMessageState // 生成订单失败，归还库存
	}

	for _, o := range orderGoods {
		o.Order = int32(order.ID)
	}
	//  订单商品，批量每次插入100条
	if result := DB.CreateInBatches(orderGoods, 100); result.Error != nil {
		tx.Rollback()
		o.Code = codes.Internal
		o.ErrorMsg = "订单商品插入失败，" + result.Error.Error()
		return primitive.CommitMessageState // 订单商品插入失败，归还库存
	}

	// 从购物车，删除已购买的商品
	if result := DB.Where(&model.ShoppingCart{User: req.UserId, Checked: true}).Delete(&model.ShoppingCart{}); result.Error != nil {
		tx.Rollback()
		o.Code = codes.Internal
		o.ErrorMsg = "删除购物车商品记录失败，" + result.Error.Error()
		return primitive.CommitMessageState // 删除购物车商品失败，归还库存
	}
	// transactionWithTimeoutCancel([]byte(req.Sn))

	// 延时消息发送，订单超时
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.1.107:9876"}))
	if err != nil {
		zap.S().Errorf("延迟，生成producer失败 %s", err.Error())
		tx.Rollback()
		o.Code = codes.Internal
		o.ErrorMsg = "延迟，生成producer失败" + err.Error()
		return primitive.CommitMessageState // 延时消息发送失败，归还库存
	}

	if err = p.Start(); err != nil {
		zap.S().Errorf("延迟，启动producer失败 %s", err.Error())
		tx.Rollback()
		o.Code = codes.Internal
		o.ErrorMsg = "延迟，启动producer失败" + err.Error()
		return primitive.CommitMessageState // 延时消息发送失败，归还库存
	}
	// defer p.Shutdown()

	message := primitive.NewMessage("order_timeout", msg.Body)
	message.WithDelayTimeLevel(5)
	_, err = p.SendSync(context.Background(), message)
	if err != nil {
		zap.S().Errorf("延时消息发送失败 %s", err.Error())
		tx.Rollback()
		o.Code = codes.Internal
		o.ErrorMsg = "延时消息发送失败" + err.Error()
		return primitive.CommitMessageState // 延时消息发送失败，归还库存
	}

	tx.Commit()
	// ##本地事务结束

	// Response需要的参数
	o.Code = codes.OK
	o.ID = order.ID
	o.OrderSn = order.OrderSn
	o.OrderMount = order.OrderMount

	//本地执行逻辑无缘无故失败 代码异常 宕机
	return primitive.RollbackMessageState
}

// 延时消息，订单超时，取消订单
func transactionWithTimeoutCancel(sn []byte) {
}

func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	fmt.Println("-----------触发回检-----------")
	var req proto.OrderRequest
	json.Unmarshal(msg.Body, &req)

	var order model.OrderInfo
	// 如何检查本地事务是否成功
	result := DB.Where(&model.OrderInfo{OrderSn: o.OrderSn}).First(&order)
	if result.Error != nil {
		return primitive.CommitMessageState
	}

	return primitive.RollbackMessageState
}
