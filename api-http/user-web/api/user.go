package api

import (
	"api-http/user-web/forms"
	"api-http/user-web/global"
	"api-http/user-web/proto"
	"api-http/user-web/utils"
	. "api-http/user-web/utils"
	"context"
	"net/http"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandlerGetList(g *gin.Context) {
	// 获取请求参数
	page, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(g.DefaultQuery("pageSize", "5"))
	// 获取用户列表
	rsp, err := global.UserRpc.GetUserList(
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
	user, err := global.UserRpc.GetUserByMobile(
		context.Background(),
		&proto.MobileRequest{Mobile: logForm.Mobile},
	)
	if err != nil {
		code := GRPCStatusToHTTP(err)
		zap.S().Errorw("[login GetUserByMobile] 用户没找到", "grpc code", code, err.Error())
		ResponseError(g, CodeUserNotExists)
		return
	}

	checked, _ := global.UserRpc.CheckPassword(
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

func HandlerRegister(g *gin.Context) {
	var regisForm forms.RegisterUser
	if err := g.ShouldBind(&regisForm); err != nil {
		zap.S().Errorw("[Register ShouldBind] 参数绑定错误", err.Error())
		ResponseError(g, CodeInvalidParams)
		return
	}
	if _, err := govalidator.ValidateStruct(regisForm); err != nil {
		zap.S().Errorw("[Register Validate form data] 参数验证不通过", err.Error())
		ResponseError(g, CodeInvalidParams)
		return
	}
	if regisForm.ConfirmPwd != regisForm.Password {
		zap.S().Errorw("[Register Validate form data] 两次密码不一致")
		ResponseError(g, CodeInvalidParams, "两次密码不一致")
		return
	}

	// 调用grpc登录服务
	_, err := global.UserRpc.GetUserByMobile(context.Background(), &proto.MobileRequest{Mobile: regisForm.Mobile})
	if err == nil {
		zap.S().Errorw("[Register] 用户已存在")
		ResponseError(g, CodeUserExists, "用户已存在")
		return
	}

	_, err = global.UserRpc.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: regisForm.ConfirmPwd,
		Password: regisForm.Password,
		Mobile:   regisForm.Mobile,
	})
	if err != nil {
		zap.S().Errorw("[Register CreateUser] 创建用户出错")
		ResponseError(g, CodeServerInternal, "创建用户出错")
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"success": true,
		"msg":     "用户创建成功",
	})
}
