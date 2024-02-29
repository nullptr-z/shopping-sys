package utils

import (
	"fmt"
	"goods-server/global"

	"github.com/gofrs/uuid"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var id, _ = uuid.NewV4()
var ConsulId = fmt.Sprint(id)

func RegisterRpcInConsul(host string, port int32) {
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
	registration.Name = global.Conf.Name // 服务名称
	registration.ID = ConsulId           // 服务ID，唯一
	registration.Port = int(port)        // 服务端口
	registration.Tags = global.Conf.Tags // 可选标签
	registration.Address = host          // 服务地址
	registration.Check = check           // 如果不填写，默认健康的
	fmt.Println("--------------global.Conf.Tags:", global.Conf.Tags)

	// 注册服务到Consul
	err := global.Consul.Agent().ServiceRegister(registration)
	if err != nil {
		zap.S().Warnw("服务注册失败", err.Error())
		panic("服务注册失败")
	}

	zap.S().Infof("注册服务成功！")
}

func DeRegister() {
	fmt.Println("ConsulId:", ConsulId)
	err := global.Consul.Agent().ServiceDeregister(ConsulId)
	if err != nil {
		zap.S().Fatal("服务注销失败", err.Error())
		return
	}
	zap.S().Info("服务注销成功！")
}
