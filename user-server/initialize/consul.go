package initialize

import (
	"fmt"
	"user-server/global"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

func ConnConsul() {
	// 连接 Consul 客户端
	consulConfig := api.DefaultConfig()
	// Consul服务的地址
	consulConfig.Address = fmt.Sprintf(
		"%s:%d",
		viper.GetString("consul.host"),
		viper.GetInt("consul.port"),
	)
	var err error
	global.Consul, err = api.NewClient(consulConfig)
	if err != nil {
		fmt.Println("Failed Connect Consul", err.Error())
		panic("Consul")
	}
	fmt.Println("Connect Consul .......")
}
