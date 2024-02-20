package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandlerGetList(g *gin.Context) {
	// 获取参数
	param := g.PostForm("")
	fmt.Println("param:", param)
	// 参数校验
	// 响应结果
	g.JSON(http.StatusOK, gin.H{"message": ""})
}
