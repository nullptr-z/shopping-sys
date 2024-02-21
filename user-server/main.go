package main

import (
	"fmt"
	"net"
	"user-server/handler"
	"user-server/initialize"
	"user-server/proto"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// get ip:port from terminal； --help 会得到参数提示
// ip := flag.String("ip", "0.0.0.0", "ip address")
// port := flag.Int("port", 10001, "socket port number")
// flag.Parse()

func main() {
	// 初始化。注意顺序不能变
	initialize.ViperConfig()
	initialize.Logger()
	initialize.MySql()

	ip := viper.GetString("host")
	port := viper.GetInt32("port")
	address := fmt.Sprintf("%s:%d", ip, port)
	fmt.Print("Listen on http://", address, "\n")

	// 创建 gRPC 服务
	server := grpc.NewServer()
	// 注册 grpc 调用
	proto.RegisterUserServer(server, &handler.UserService{})
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprint("error:", err.Error()))
	}

	err = server.Serve(lis)
	if err != nil {
		panic(fmt.Sprint("Listen error of grpc service:", err.Error()))
	}
}
