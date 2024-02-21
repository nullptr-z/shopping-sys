package router

import (
	"api-http/user-web/api"
	"api-http/user-web/middleware"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(route *gin.RouterGroup) {
	userRoute := route.Group("user")
	{
		userRoute.GET("/list", middleware.Authorization, middleware.AuthAdmin, api.HandlerGetList)
		userRoute.POST("/register", api.HandlerRegister)
		userRoute.POST("/login", api.HandlerLogin)
	}

}
