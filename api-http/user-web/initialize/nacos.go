package initialize

import (
	"api-http/user-web/global"
	"bytes"
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
)

func LoadNaocs() {
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: viper.GetString("nacos.host"), // Nacos的服务地址
			Port:   viper.GetUint64("nacos.port"), // Nacos的服务端口
			// ContextPath: "", // Nacos的ContextPath，默认/nacos，在2.0中不需要设置
			// Scheme:      "", // Nacos的服务地址前缀，默认http，在2.0中不需要设置
			// GrpcPort:    "", // Nacos的 grpc 服务端口, 默认为 服务端口+1000, 不是必填
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         viper.GetString("nacos.NamespaceId"), // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
		TimeoutMs:           viper.GetUint64("nacos.TimeoutMs"),
		NotLoadCacheAtStart: viper.GetBool("nacos.NotLoadCacheAtStart"),
		LogDir:              viper.GetString("nacos.LogDir"),
		CacheDir:            viper.GetString("nacos.CacheDir"),
		LogLevel:            viper.GetString("nacos.LogLevel"),
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
		DataId: viper.GetString("nacos.DataId"),
		Group:  viper.GetString("nacos.Group"),
	})

	if err != nil {
		panic(fmt.Sprint("GetConfig error:", err.Error()))
	}

	// 设置配置类型（"json", "yaml", "toml"等）
	bindingOnConfigure(content)

	configClient.ListenConfig(vo.ConfigParam{
		DataId: viper.GetString("nacos.DataId"),
		Group:  viper.GetString("nacos.Group"),
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("nacos listening on config change")
			// fmt.Println("namespace, group, dataId, data string:", namespace, group, dataId, data)
			fmt.Println(data)
			bindingOnConfigure(data)
		},
	})
}

func bindingOnConfigure(content string) {
	fmt.Println("content:", content)
	v := viper.New()

	v.SetConfigType("yaml")
	v.ReadConfig(bytes.NewBuffer([]byte(content)))
	err := v.Unmarshal(&global.Conf)
	if err != nil {
		panic(fmt.Sprint("nacos config binding failed on struct Configure: ", err.Error()))
	}
	fmt.Println("Configure conf:", global.Conf)
}
