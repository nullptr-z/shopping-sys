package main

import (
	"api-http/user-web/initialize"
	"fmt"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置 viper
	if err := initialize.ViperConfig(); err != nil {
		fmt.Println("Failed initialized Init ViperConfig", err)
		return
	}
	// 2. 初始化日志
	if err := initialize.Logger(); err != nil {
		fmt.Println("Failed initialized Init Logger", err)
		return
	}
	// 注册路由
	g := initialize.Routers()

	// 启动服务
	port := viper.Get("port")
	zap.S().Infof("run service, port:", port)
	if err := g.Run(fmt.Sprintf(":%d", port)); err != nil {
		zap.S().Panic("run service failed", err.Error())
	}
}
