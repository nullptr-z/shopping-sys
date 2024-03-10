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

func TestCreateCartItem(udi, goodsid, num int32) {
	res, err := client.CreateCartItem(context.Background(), &proto.CartItemRequest{
		UserId:  udi,
		GoodsId: goodsid,
		Nums:    num,
	})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	fmt.Println("res: ", res)
}

func TestCreateCartList(udi int32) {
	res, err := client.CartItemList(context.Background(), &proto.UserInfo{Id: udi})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	for _, goods := range res.Data {
		fmt.Println("goods:", goods)
	}
}

func TestUpdateCartItem() {
	res, err := client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      4,
		Checked: true, // 勾选购物车的物品
	})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	fmt.Println("res: ", res)
}

func TestCreateOrder() {
	res, err := client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  10,
		Address: "北京市朝阳区",
		Name:    "测试",
		Mobile:  "12345678901",
		Post:    "请速速发货",
	})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	fmt.Println("res: ", res)
}

func main() {
	// 购物车Tests
	TestCreateCartItem(10, 422, 5)
	TestCreateCartItem(10, 421, 5)
	TestUpdateCartItem()
	TestCreateCartList(10)

	// 订单Tests
	TestCreateOrder()
}
