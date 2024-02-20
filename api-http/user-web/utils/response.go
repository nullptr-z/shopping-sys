package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
```json
response struct
{
	statusCode: 400,
	msg: interface{ }
	data: interface{ }
}
```
*/

type ResponseData struct {
	Code RespCodeType `json:"code"`
	Msg  interface{}  `json:"msg"`
	Data interface{}  `json:"data,omitempty"`
}

func ResponseError(g *gin.Context, code RespCodeType, msg ...any) {
	var msgarr = []string{code.getMsg()}
	for _, arg := range msg {
		switch v := arg.(type) {

		case string:
			msgarr = append(msgarr, v)
		}
	}
	rsp := &ResponseData{
		Code: code,
		Msg:  msgarr,
		Data: nil,
	}

	g.JSON(http.StatusOK, rsp)
}

func ResponseErrorWithMsg(g *gin.Context, code RespCodeType, msg interface{}) {
	rsp := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	g.JSON(http.StatusOK, rsp)
}

func ResponseSuccess(g *gin.Context, data interface{}, msg ...any) {
	var msgarr = []string{CodeSuccess.getMsg()}
	for _, arg := range msg {
		switch v := arg.(type) {
		case string:
			msgarr = append(msgarr, v)
		}
	}
	rsp := &ResponseData{
		Code: CodeSuccess,
		Msg:  msgarr,
		Data: data,
	}

	g.JSON(http.StatusOK, rsp)
}
