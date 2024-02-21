package router

import (
	"api-http/user-web/api"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(route *gin.RouterGroup) {
	userRoute := route.Group("user")
	{
		userRoute.GET("/list", api.HandlerGetList)
		userRoute.POST("/login", api.HandlerLogin)
	}

}
