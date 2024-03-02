package main

import (
	"context"
	"fmt"
	"stock-server/proto"

	"google.golang.org/grpc"
)

func TestSetStock(goodsId, Num int) {
	conn, err := grpc.Dial("192.168.1.104:45911", grpc.WithInsecure())
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	client := proto.NewStockClient(conn)

	_, err = client.SetStock(context.Background(), &proto.GoodsStockInfo{GoodsId: int32(goodsId), Num: int32(Num)})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}

}

func main() {
	for i := 421; i <= 840; i++ {
		TestSetStock(i, 100)
	}
}
