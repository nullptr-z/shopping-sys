package initialize

import (
	"api-http/user-web/global"
	"api-http/user-web/proto"
	"fmt"

	"github.com/hashicorp/consul/api"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// user-server
func filterServer(consul *api.Client, name string) (host string, port int) {
	servers, err := consul.Agent().ServicesWithFilter(fmt.Sprintf(`Service == "%s"`, name))
	if err != nil {
		zap.S().Errorw(" error:", err.Error())
		return
	}
	for _, s := range servers {
		host = s.Address
		port = s.Port
		return
	}
	return
}

func ConnConsul() (*api.Client, error) {
	consulConfig := api.DefaultConfig()
	// Consul服务的地址
	consulConfig.Address = fmt.Sprintf(
		"%s:%d",
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
	)
	consul, err := api.NewClient(consulConfig)
	if err != nil {
		zap.S().Errorw("[ConnConsul] Failed")
		return nil, err
	}
	zap.S().Info("Connect Consul success")

	return consul, nil
}

/*
	 TODO:
		1.RPC 下线或者改地址了，怎么办；依靠负载均衡
		2.连接池
*/
func ConnUserRpc() error {
	// 手动分配
	// consul, err := ConnConsul()
	// host, port := filterServer(consul, viper.GetString("userServer.name"))
	// 负载均衡，自动分配
	rpcAddress := fmt.Sprintf(
		"consul://%s:%d/%s?wait=14s",
		global.Conf.Consul.Host,
		global.Conf.Consul.Port,
		viper.GetString("userServer.name"),
	)
	zap.S().Info("user service rpcAddress:", rpcAddress)
	conn, err := grpc.Dial(
		rpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`), // 轮询策略,负载均衡
	)
	if err != nil {
		zap.S().Fatal("[HandlerGetList] connect rpc server of user failed", err)
		return err
	}
	// 创建 grpc client
	global.UserRpc = proto.NewUserClient(conn)

	fmt.Println("Connect User gRpc .......")
	return err
}
