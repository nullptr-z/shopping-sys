package main

import (
	"context"
	"fmt"
	"stock-server/proto"
	"sync"

	"google.golang.org/grpc"
)

var client proto.StockClient

func init() {

	conn, err := grpc.Dial("192.168.1.104:10003", grpc.WithInsecure())
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	client = proto.NewStockClient(conn)

}

func TestSetStock(goodsId, Num int) {
	_, err := client.SetStock(context.Background(), &proto.GoodsStockInfo{GoodsId: int32(goodsId), Num: int32(Num)})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
}

func TestSellStock(goodsId, Num int, wg *sync.WaitGroup) {
	defer wg.Done()
	_, err := client.Sell(
		context.Background(),
		&proto.SellInfo{
			GoodsInfo: []*proto.GoodsStockInfo{
				{GoodsId: int32(goodsId), Num: int32(Num)},
			},
		},
	)
	if err != nil {
		fmt.Println(" error:", err)
		return
	}

}

func main() {
	// 给商品设置库存100个
	// for i := 421; i <= 840; i++ {
	// 	TestSetStock(i, 100)
	// }

	// 扣减库存
	// TestSellStock(421, 10)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go TestSellStock(421, 1, &wg)
	}
	wg.Wait()
}
