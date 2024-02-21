package api

import (
	"api-http/user-web/forms"
	"api-http/user-web/proto"
	"api-http/user-web/utils"
	. "api-http/user-web/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connGrpc() proto.UserClient {
	// 链接 grpc 服务
	ip, port := viper.Get("rpc.ip"), viper.Get("rpc.port")
	rpcAddress := fmt.Sprintf("%s:%d", ip, port)
	conn, err := grpc.Dial(rpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[HandlerGetList] connect rpc server of user failed", "msg", err.Error())
	}
	// 创建 grpc client
	client := proto.NewUserClient(conn)

	return client
}

func HandlerGetList(g *gin.Context) {
	client := connGrpc()
	// 获取请求参数
	page, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(g.DefaultQuery("pageSize", "5"))
	// 获取用户列表
	rsp, err := client.GetUserList(
		context.Background(),
		&proto.PageInfo{Page: uint32(page), PageSize: uint32(pageSize)},
	)
	if err != nil {
		zap.S().Errorw("[GetUserList] error", err.Error())
		ResponseError(g, CodeInvalidParams, "code", GRPCStatusToHTTP(err))
		return
	}
	// 响应结果
	g.JSON(http.StatusOK, gin.H{"total": rsp.Total, "data": rsp.Data})
}

// 密码登录
func HandlerLogin(g *gin.Context) {
	var logForm forms.PasswordLoginForm
	if err := g.ShouldBind(&logForm); err != nil {
		zap.S().Errorw("[ShouldBind] 参数绑定错误", err.Error())
		ResponseError(g, CodeInvalidParams)
		return
	}
	if _, err := govalidator.ValidateStruct(logForm); err != nil {
		zap.S().Errorw("[Validate form data] 参数验证不通过", err.Error())
		ResponseError(g, CodeInvalidParams)
		return
	}

	// 调用grpc登录服务
	client := connGrpc()
	user, err := client.GetUserByMobile(
		context.Background(),
		&proto.MobileRequest{Mobile: logForm.Mobile},
	)
	if err != nil {
		code := GRPCStatusToHTTP(err)
		zap.S().Errorw("[login GetUserByMobile] 用户没找到", "grpc code", code, err.Error())
		ResponseError(g, CodeUserNotExists)
		return
	}

	checked, _ := client.CheckPassword(
		context.Background(),
		&proto.CheckPasswordInfo{
			Password:          logForm.Password,
			EncryptedPassword: user.Password,
		},
	)

	if !checked.Success {
		zap.S().Errorw("[login CheckPassword] 密码校验失败", "grpc code", GRPCStatusToHTTP(err), err.Error())
		ResponseError(g, CodeInvalidPassword)
		return
	}

	token, err := utils.GenToken(&CustomClaims{
		ID:          uint(user.Id),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
	})
	if err != nil {
		zap.S().Errorw("[GenToken] 生成 token 出错")
		ResponseError(g, CodeServerInternal, "Token 生成失败")
		return
	}

	g.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "登陆成功",
		"token":   token,
		"expired": utils.TokenExpireDuration,
	})
}
