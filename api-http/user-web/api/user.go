package api

import (
	"api-http/user-web/proto"
	. "api-http/user-web/utils"
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func HandlerGetList(g *gin.Context) {
	// zap.S().Debug("get user list")
	// 链接 grpc 服务
	ip, port := viper.Get("rpc.ip"), viper.Get("rpc.port")
	rpcAddress := fmt.Sprintf("%s:%d", ip, port)
	conn, err := grpc.Dial(rpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		zap.S().Errorw("[HandlerGetList] connect rpc server of user failed", "msg", err.Error())
	}
	// 获取请求参数
	page, _ := strconv.Atoi(g.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(g.DefaultQuery("pageSize", "5"))
	// 创建 grpc client
	client := proto.NewUserClient(conn)
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
