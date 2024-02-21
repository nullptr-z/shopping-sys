package main

import (
	"fmt"
	"net"
	"user-server/handler"
	"user-server/initialize"
	"user-server/proto"
	"user-server/utils"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	gprcCheck "google.golang.org/grpc/health/grpc_health_v1"
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

	host := viper.GetString("host")
	port := viper.GetInt32("port")
	address := fmt.Sprintf("%s:%d", host, port)
	fmt.Print("\n Listen on http://", address, "\n\n")

	// 创建 gRPC 服务
	server := grpc.NewServer()
	// 用于服务健康检查
	gprcCheck.RegisterHealthServer(server, health.NewServer())
	utils.Register(host, port)
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
