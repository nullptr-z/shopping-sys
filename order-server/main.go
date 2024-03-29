package main

import (
	"flag"
	"fmt"
	"net"
	. "order-server/global"
	"order-server/handler"
	"order-server/initialize"
	"order-server/proto"
	"order-server/utils"
	"os"
	"os/signal"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	gprcCheck "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	// 初始化。注意顺序不能变
	initialize.ViperConfig()
	initialize.LoadNaocs()
	initialize.Logger()
	initialize.MySql()
	initialize.ConnRedis()
	initialize.ConnRedisSync()
	initialize.ConnRpcAll()
	initialize.ConnConsul()

	host_defualt := Conf.Host
	port_defualt := utils.GetFreePort()

	// get ip:port from terminal； --help 会得到参数提示
	host := flag.String("ip", host_defualt, "ip address")
	port := flag.Int("port", port_defualt, "socket port number")
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *host, *port)
	fmt.Print("\n Listen on http://", address, "\n\n")

	// 创建 gRPC 服务
	server := grpc.NewServer()
	// 用于服务健康检查
	gprcCheck.RegisterHealthServer(server, health.NewServer())
	utils.RegisterRpcInConsul(*host, int32(*port))
	// 注册 grpc 调用
	proto.RegisterOrderServer(server, &handler.OrderService{})
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Sprint("error:", err.Error()))
	}

	// 监听订单超时topic
	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.1.107:9876"}),
		consumer.WithGroupName("order_reback"),
	)
	if err := c.Subscribe("order_timeout", consumer.MessageSelector{}, handler.OrderTimeout); err != nil {
		fmt.Println("读取消息失败")
	}
	c.Start()
	defer c.Shutdown()

	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic(fmt.Sprint("Listen error of grpc service:", err.Error()))
		}
	}()
	shutdown()
}

func shutdown() {
	// 接受退出信号，善后工作
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGABRT, syscall.SIGALRM) // 添加SIGINT
	<-quit
	utils.DeRegister()
}
