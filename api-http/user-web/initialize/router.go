package initialize

import (
	"api-http/user-web/router"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Routers() *gin.Engine {
	if viper.Get("mode") == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 发布模式
	}
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery(), LoggerFormateOutput)

	g.GET("/health", func(ctx *gin.Context) {
		// 用于 consul 健康检查
		ctx.AbortWithStatus(http.StatusOK)
	})

	routeRoot := g.Group("/u/v1")
	router.InitUserRouter(routeRoot)

	fmt.Println("Register routers .......")
	return g
}
