package main

import (
	"api-http/user-web/proto"
	"context"
	"fmt"
	"log"

	_ "github.com/mbobakov/grpc-consul-resolver"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"consul://127.0.0.1:8500/user-server?wait=14s",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	page := proto.PageInfo{Page: 1, PageSize: 4}
	fmt.Println("page:", &page)

	userRpc := proto.NewUserClient(conn)
	uerList, err := userRpc.GetUserList(
		context.Background(),
		&page,
	)
	for _, u := range uerList.Data {
		fmt.Println("user info:", u.Id, u.Mobile, u.NickName)
	}

}
