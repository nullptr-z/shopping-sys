package handler

import (
	"context"
	"stock-server/global"
	"stock-server/model"
	. "stock-server/proto"

	"github.com/golang/protobuf/ptypes/empty"
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
	global.DB.First(&s, req.GoodsId)
	s.Goods = req.GoodsId
	s.Stocks = req.Num
	global.DB.Save(&s)

	return &empty.Empty{}, nil
}
func (*StockService) InvDetail(ctx context.Context, req *GoodsStockInfo) (*GoodsStockInfo, error) {
	var s model.Stock
	ret := global.DB.First(&s, req.GoodsId)
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
*/
func (*StockService) Sell(ctx context.Context, req *SellInfo) (*empty.Empty, error) {
	var s model.Stock
	tx := global.DB.Begin()
	for _, goodsInfo := range req.GoodsInfo {
		ret := global.DB.First(&s, goodsInfo.GoodsId)
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
		// 添加事务
		tx.Save(&s)
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
