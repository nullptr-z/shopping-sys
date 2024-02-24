package main

import (
	"fmt"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"

	// "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func main() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: "127.0.0.1", // Nacos的服务地址
			Port:   8848,        // Nacos的服务端口
			// ContextPath: "", // Nacos的ContextPath，默认/nacos，在2.0中不需要设置
			// Scheme:      "", // Nacos的服务地址前缀，默认http，在2.0中不需要设置
			// GrpcPort:    "", // Nacos的 grpc 服务端口, 默认为 服务端口+1000, 不是必填
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         "user", // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos/log",
		CacheDir:            "nacos/cache",
		LogLevel:            "debug",
	}

	var err error
	// 创建动态配置客户端
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		panic(fmt.Sprint("CreateConfigClient error: ", err.Error()))
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: "user-server",
		Group:  "dev",
	})
	if err != nil {
		panic(fmt.Sprint("GetConfig error:", err.Error()))
	}
	configClient.ListenConfig(vo.ConfigParam{
		DataId: "user-server",
		Group:  "dev",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("nacos listening on config change")
			fmt.Println("namespace, group, dataId, data string:", namespace, group, dataId, data)
		},
	})

	fmt.Println("content:", content)
	time.Sleep(1000 * time.Second)
}
