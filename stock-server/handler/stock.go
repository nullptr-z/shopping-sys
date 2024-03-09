package handler

import (
	"context"
	"fmt"
	"stock-server/global"
	"stock-server/model"
	. "stock-server/proto"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

// var mx sync.Mutex

/*
扣减库存
 1. 需要满足本地事务，三件商品要能同时扣减成功，其中一个失败全部撤回；数据一致性
 2. 数据库的一个应用场景：事务；数据一致性
*/
func (*StockService) Sell(ctx context.Context, req *SellInfo) (*empty.Empty, error) {
	var s model.Stock
	tx := global.DB.Begin()
	// 无法在分布式环境下使用，必须要在事务提交之后释放
	// mx.Lock()
	// defer mx.Unlock()
	for _, goodsInfo := range req.GoodsInfo {
		// 悲观锁
		// ret := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Stock{Goods: goodsInfo.GoodsId}).First(&s)
		for {
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
			if result := tx.Model(&model.Stock{}).Where("goods = ? and version = ?", goodsInfo.GoodsId, s.Version).Updates(model.Stock{
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
