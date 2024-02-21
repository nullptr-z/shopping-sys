package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Authorization, accept, origin, Cache-Control")
	c.Writer.Header().Set("Access-Control-Expose-Headers", "Access-Control-Allow-Origin, Content-Length")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // 允许跨域请求携带认证信息（如Cookies）

	if c.Request.Method == "OPTIONS" {
		// 处理探测性请求
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	c.Next()
}
