package initialize

import (
	"fmt"
	"order-server/global"
	"order-server/proto"

	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

/*
	 TODO:
		1.RPC 下线或者改地址了，怎么办；依靠负载均衡
		2.连接池
*/
func ConnGoodsRpc() error {
	// 负载均衡，自动分配
	rpcAddress := fmt.Sprintf(
		"consul://%s:%d/%s?wait=14s",
		global.Conf.Consul.Host,
		global.Conf.Consul.Port,
		global.Conf.RpcName.Goods,
	)
	zap.S().Info("goods service rpcAddress:", rpcAddress)
	conn, err := grpc.Dial(
		rpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[ConnGoodsRpc] connect rpc server failed", err)
		return err
	}
	// 创建 grpc client
	global.GoodsSvc = proto.NewGoodsClient(conn)

	fmt.Println("Connect Goods gRpc .......")
	return err
}

func ConnStockRpc() error {
	rpcAddress := fmt.Sprintf(
		"consul://%s:%d/%s?wait=14s",
		global.Conf.Consul.Host,
		global.Conf.Consul.Port,
		"stock-server",
	)
	zap.S().Info("stock service rpcAddress:", rpcAddress)
	// 负载均衡，自动分配
	conn, err := grpc.Dial(
		rpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		zap.S().Fatal("[ConnStockRpc] connect rpc server failed", err)
		return err
	}

	// 创建 grpc client
	global.StockSvc = proto.NewStockClient(conn)

	fmt.Println("Connect Stock gRpc .......")
	return err
}

func ConnRpcAll() {
	ConnGoodsRpc()
	ConnStockRpc()
}
