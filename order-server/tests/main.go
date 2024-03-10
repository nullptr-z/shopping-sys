package main

import (
	"context"
	"fmt"
	"order-server/proto"

	"google.golang.org/grpc"
)

var client proto.OrderClient

func init() {
	conn, err := grpc.Dial("192.168.1.104:10004", grpc.WithInsecure())
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	client = proto.NewOrderClient(conn)

}

func TestCreateCartItem() {
	res, err := client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  10,
		GoodsId: 421,
		Nums:    5,
	})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	fmt.Println("res: ", res)
}

func main() {
	TestCreateCartItem()
}
