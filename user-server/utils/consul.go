package utils

import (
	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func client() *api.Client {
	// 连接 Consul 客户端
	consulConfig := api.DefaultConfig()
	// Consul服务的地址
	consulConfig.Address = fmt.Sprintf(
		"%s:%d",
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
	)
	client, err := api.NewClient(consulConfig)
	if err != nil {
		zap.S().Warnw(err.Error())
	}

	return client
}

func Register(host string, port int32) {
	client := client()
	// 添加需要健康检查的微服务信息
	check := &api.AgentServiceCheck{
		// 检查 gRPC 服务
		GRPC:                           fmt.Sprintf("%s:%d", host, port), // 健康检查的地址
		Interval:                       "10s",                            // 健康检查的间隔时间
		Timeout:                        "2s",                             // 健康检查的超时时间
		DeregisterCriticalServiceAfter: "5s",
	}

	// 设置要注册的服务的信息
	registration := new(api.AgentServiceRegistration)
	registration.ID = viper.GetString("name")                                // 服务ID，唯一
	registration.Name = viper.GetString("name")                              // 服务名称
	registration.Port = 11001                                                // 服务端口
	registration.Tags = []string{"user", "login", "register", "web", "http"} // 可选标签
	registration.Address = host                                              // 服务地址
	registration.Check = check                                               // 如果不填写，默认健康的

	// 注册服务到Consul
	err := client.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Warnw(err.Error())
	}

	zap.S().Infof("注册服务！")
}
