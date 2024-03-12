package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"stock-server/global"
	"stock-server/model"
	. "stock-server/proto"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type StockService struct {
	UnimplementedStockServer
}

// set up stock
func (*StockService) SetStock(ctx context.Context, req *GoodsStockInfo) (*empty.Empty, error) {
	var s model.Stock
	global.DB.Where(&model.Stock{Goods: req.GoodsId}).First(&s)
	// 如果不存在直接添加这个商品
	if s.Goods == 0 {
		s.Goods = req.GoodsId
		s.CratedAt = time.Now()
		s.UpdatedAt = time.Now()
	}
	s.Stocks = req.Num
	fmt.Println("s:", s)
	global.DB.Save(&s)

	return &empty.Empty{}, nil
}

func (*StockService) InvDetail(ctx context.Context, req *GoodsStockInfo) (*GoodsStockInfo, error) {
	var s model.Stock
	ret := global.DB.Where(&model.Stock{Goods: req.GoodsId}).First(&s)
	if ret.RowsAffected == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
	}

	info := &GoodsStockInfo{
		GoodsId: s.Goods,
		Num:     s.Stocks,
	}

	return info, nil
}

/*
扣减库存
 1. 需要满足本地事务，三件商品要能同时扣减成功，其中一个失败全部撤回；数据一致性
 2. 数据库的一个应用场景：事务；数据一致性
    红锁，续命，只允许自己释放锁
*/
func (*StockService) Sell(ctx context.Context, req *SellInfo) (*empty.Empty, error) {
	var s model.Stock
	tx := global.DB.Begin()

	var goodsList []model.GoodsDetail

	for _, goodsInfo := range req.GoodsInfo {
		mutexname := fmt.Sprint("stockSell.", goodsInfo.GoodsId)
		mutex := global.Rdsync.NewMutex(mutexname)
		if err := mutex.Lock(); err != nil {
			zap.S().Error("Redis lock failed", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Redis lock failed")
		}
		// time.Sleep(10 * time.Second)
		ret := global.DB.Where(&model.Stock{Goods: goodsInfo.GoodsId}).First(&s)
		if ret.Error != nil {
			tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
			return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息，"+ret.Error.Error())
		}
		// 库存是否充足
		if s.Stocks < goodsInfo.Num {
			tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
			return nil, status.Errorf(codes.InvalidArgument, "库存不足")
		}
		// 所有条件满足，扣除库存
		// 并发是，可能出现数据不一致问题，其他地方扣减了库存；分布式锁
		s.Stocks -= goodsInfo.Num
		tx.Save(&s)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			zap.S().Error("Redis unlock failed", zap.Error(err))
			return nil, status.Errorf(codes.Internal, "Redis unlock failed")
		}

		goodsList = append(goodsList, model.GoodsDetail{Goods: goodsInfo.GoodsId, Num: goodsInfo.Num})
	}
	// 添加库存扣减记录
	sellDetail := model.StockSellReback{
		OrderSn:   req.OrderSn,
		Status:    1,
		GoodsList: goodsList,
	}
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
		return nil, status.Errorf(codes.Internal, "保存库存扣减记录失败")
	}
	tx.Commit() // 提交事务
	return &emptypb.Empty{}, nil
}

// 扣减库存，使用mysql做分布式锁
// 乐观锁，悲观锁
// var mx sync.Mutex // go 的互斥锁，无法在分布式环境下使用
func (*StockService) Sell_mysql_lock(ctx context.Context, req *SellInfo) (*empty.Empty, error) {
	var s model.Stock
	tx := global.DB.Begin()
	// 无法在分布式环境下使用，必须要在事务提交之后释放
	// mx.Lock()
	// defer mx.Unlock()
	for _, goodsInfo := range req.GoodsInfo {
		// 悲观锁
		// ret := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Stock{Goods: goodsInfo.GoodsId}).First(&s)
		for {
			fmt.Println(" ================goodsInfo.GoodsId:", goodsInfo.GoodsId)
			ret := global.DB.Where(&model.Stock{Goods: goodsInfo.GoodsId}).First(&s)
			if ret.RowsAffected == 0 {
				tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
				return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
			}
			// 库存是否充足
			if s.Stocks < goodsInfo.Num {
				tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
				return nil, status.Errorf(codes.InvalidArgument, "库存不足")
			}
			// 所有条件满足，扣除库存
			// 并发是，可能出现数据不一致问题，其他地方扣减了库存；分布式锁
			s.Stocks -= goodsInfo.Num
			// 乐观锁
			if result := tx.Model(&model.Stock{}).Select("stocks", "version").Where("goods = ? and version = ?", goodsInfo.GoodsId, s.Version).Updates(model.Stock{
				Stocks:  s.Stocks,
				Version: s.Version + 1,
			}); result.RowsAffected == 0 {
				zap.S().Info("resutl:", result)
			} else {
				break
			}
		}
		// tx.Save(&s)
	}
	tx.Commit() // 提交事务
	return &emptypb.Empty{}, nil
}

/*
	 库存归还
		1. 订单超时，自动归还
		2. 订单创建失败
		3. 手动归还，取消订单

		订单中商品需要同时满足事务性，全部归还成功，才能取消订单
*/
func (*StockService) Reback(ctx context.Context, req *SellInfo) (*empty.Empty, error) {
	var s model.Stock
	tx := global.DB.Begin()
	for _, goodsInfo := range req.GoodsInfo {
		ret := global.DB.First(&s, goodsInfo.GoodsId)
		if ret.RowsAffected == 0 {
			tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
			return nil, status.Errorf(codes.InvalidArgument, "没有找到库存信息")
		}
		s.Stocks += goodsInfo.Num
		// 添加事务
		tx.Save(&s)
	}
	tx.Commit() // 提交事务

	return &emptypb.Empty{}, nil
}

func (*StockService) mustEmbedUnimplementedStockServer() {}

/*
1. 需要知道是哪个订单，哪些商品，要归还多少。
2. 需要确保幂等性，不能因为网络等延迟导致的重复归还。新建一张表，设置订单为主键，保证每个订单只能归还一次
3. 订单中商品需要同时满足事务性，全部归还成功，才能取消订单
*/
func StockReback(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string `json:"sn"`
	}
	// 打印接收到的消息
	for _, msg := range msgs {
		var orderInfo OrderInfo
		fmt.Printf("orderInfo..............:%s", msg.Body)
		err := json.Unmarshal(msg.Body, &orderInfo)
		fmt.Println("orderInfo。。。。。。。。:", orderInfo)
		if err != nil {
			zap.S().Error("json.Unmarshal failed", err.Error())
			return consumer.ConsumeSuccess, nil // 数据有问题，直接丢弃
		}
		// 查询扣减记录
		tx := global.DB.Begin()
		var sellDetail model.StockSellReback
		result := tx.Where(&model.StockSellReback{
			OrderSn: orderInfo.OrderSn,
			Status:  1,
		}).First(&sellDetail)
		if result.RowsAffected == 0 {
			zap.S().Info("该库存已经归还过了")
			return consumer.ConsumeSuccess, nil
		}

		// 逐个归还商品
		for _, g := range sellDetail.GoodsList {
			result := tx.Model(&model.Stock{}).Where(&model.Stock{Goods: g.Goods}).Update("stocks", gorm.Expr("stocks + ?", g.Num))
			if result.RowsAffected == 0 {
				zap.S().Info("reback goods failed:", result.Error.Error())
				tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
				return consumer.ConsumeRetryLater, nil
			}
		}

		// 标记为已归还
		sellDetail.Status = 2
		result = tx.Model(&sellDetail).Where(&model.StockSellReback{OrderSn: orderInfo.OrderSn}).Update("status", 2)
		if result.RowsAffected == 0 {
			zap.S().Info("set sellDetail.Status failed:", result.Error.Error())
			tx.Rollback() // 事务回滚，如果之前的商品成功扣减了的话
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
	}

	// 返回成功到主服务器，表示接受到了消息
	return consumer.ConsumeSuccess, nil
}
