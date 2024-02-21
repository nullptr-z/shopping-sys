package middleware

import (
	"api-http/user-web/global"
	"api-http/user-web/utils"
	. "api-http/user-web/utils"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Authorization(g *gin.Context) {
	author := g.GetHeader("Authorization")
	bearerSchema := "Bearer "
	// 从请求头提取 TOken
	zap.S().Infof("author:", author)
	if author == "" {
		ResponseError(g, CodeInvalidateAuth)
		g.Abort()
		return
	}
	// 分割取出 Token
	tokenString := strings.TrimPrefix(author, bearerSchema)
	zap.S().Infof("token:", tokenString)
	if tokenString == author {
		ResponseError(g, CodeInvalidateAuth)
		g.Abort()
		return
	}

	// 解析验证 Token
	if claims, err := utils.ParseTOken(tokenString); err != nil {
		ResponseError(g, CodeInvalidateToken)
		g.Abort()
		return
	} else {
		// Token 是有效的
		zap.S().Infof("解析token, User ID: %d, Username: %s\n", claims.ID, claims.NickName)
		// 将用户信息添加到请求的上下文中
		g.Set(global.ContextUserIDKey, claims)
		g.Next() // 处理下一个请求
	}
}

func AuthAdmin(g *gin.Context) {
	claim := global.GetUserClaim(g)
	fmt.Println("claim:", claim)
	if claim.AuthorityId != global.ADMIN {
		ResponseError(g, CodeInsufficientPerms, "需要管理员访问权限")
		g.Abort()
		return
	}
	g.Next()
}
