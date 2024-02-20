package initialize

import (
	"api-http/user-web/router"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Routers() *gin.Engine {
	if viper.Get("mode") == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 发布模式
	}
	g := gin.New()
	g.Use(gin.Logger(), gin.Recovery(), LoggerFormateOutput)

	// 跨域访问配置
	// g.Use(func(c *gin.Context) {
	// 	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	// 	// c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// 	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
	// 	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

	// 	if c.Request.Method == "OPTIONS" {
	// 		// 处理探测性请求
	// 		c.AbortWithStatus(http.StatusNoContent)
	// 		return
	// 	}

	// 	c.Next()
	// })

	routeRoot := g.Group("/v1")
	router.InitUserRouter(routeRoot)

	fmt.Println("Register routers .......")
	return g
}
