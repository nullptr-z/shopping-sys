package main

import (
	"api-http/user-web/initialize"
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	startService(g)
}

func startService(g *gin.Engine) {
	// binding socket
	host := viper.GetString("host")
	port := viper.GetInt("port")
	addr := fmt.Sprintf("%s:%d", host, port)
	srv := &http.Server{Addr: addr, Handler: g}

	go func() {
		fmt.Print("\n Listening on: http://", addr, "\n\n")
		if err := srv.ListenAndServe(); err != nil {
			zap.L().Error("Listening:", zap.Error(err))
		}
	}()

	// 等待系统中断信号来关闭服务，为关闭服务设置一个 5 秒的超时
	// kill default syscaLL.SIGTERM
	// kill -2  syscaLL.SIGINT 我们常用的ctrl+C就是触发系统 SIGINT 信号
	// kill -9  syscalL.SIGKILL 不能被捕获，所以不需要添加它
	quit := make(chan os.Signal, 1)
	// signal.Notify 会把收到的 syscalL.SIGINT 或 syscaLL.SIGTERM 信号转发给 quit
	<-quit // 阻塞等待关闭信号
	zap.L().Info("Shutdown Server...")
	// 定时 5 的Chan
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 延迟 5 秒，处理还未完成的请求扫尾，然后优雅停机
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown", zap.Error(err))
	}
	zap.L().Info("Service exiting")
}
