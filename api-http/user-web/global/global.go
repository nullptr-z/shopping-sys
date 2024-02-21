package global

import (
	"api-http/user-web/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	ADMIN uint = iota + 1
	USER
	VISITOR
)

// 全局上下文可访问字段, 后续上文件中可以 g.Get(key) 可以直接获取
// middle.authorization中设置的
var (
	ContextUserIDKey = "claim"
)

func GetUserId(g *gin.Context) string {
	id := GetUserClaim(g).ID
	return strconv.FormatInt(int64(id), 10)
}

func GetUserClaim(g *gin.Context) *utils.CustomClaims {
	claim := g.MustGet(ContextUserIDKey)
	claims := claim.(*utils.CustomClaims)
	return claims
}
