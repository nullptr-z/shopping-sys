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

func TestCreateCartList(uid int32) {
	fmt.Println("udi:", uid)
	res, err := client.CartItemList(context.Background(), &proto.UserInfo{Id: uid})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	fmt.Println("======================", len(res.Data))
	for _, goods := range res.Data {
		fmt.Println("=======================TestCreateCartList goods:", goods)
	}
}

func TestUpdateCartItem(id int32) {
	res, err := client.UpdateCartItem(context.Background(), &proto.CartItemRequest{
		Id:      id,
		Checked: true, // 勾选购物车的物品
	})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	fmt.Println("res: ", res)
}

func TestCreateOrder(uid int32) {
	res, err := client.CreateOrder(context.Background(), &proto.OrderRequest{
		UserId:  uid,
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

func TestOrderDetail(orderId int32) {
	res, err := client.OrderDetail(context.Background(), &proto.OrderRequest{Id: orderId})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	for _, g := range res.Goods {
		fmt.Println("g:", g)
	}
}

func TestOrderList() {
	res, err := client.OrderList(context.Background(), &proto.OrderFilterRequest{
		UserId: 10,
	})
	if err != nil {
		fmt.Println(" error:", err)
		return
	}
	for _, g := range res.Data {
		fmt.Println("g:", g)
	}
}

func main() {
	var uid int32 = 10
	// var goodsId int32 = 421
	// 购物车Tests
	// TestCreateCartItem(uid, 421, 1)
	// TestCreateCartItem(uid, 422, 2)
	// TestCreateCartList(uid)
	// TestUpdateCartItem(6)

	// 订单Tests
	TestCreateOrder(uid)
	// TestOrderDetail(30)
	// TestOrderList()
}
