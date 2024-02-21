package utils

type RespCodeType int64

const (
	CodeSuccess RespCodeType = 1000 + iota
	CodeInvalidParams
	CodeInvalidPassword
	CodeUserExists
	CodeUserNotExists
	CodeServerBusy
	CodeServerInternal
	// login
	CodeInvalidateAuth RespCodeType = 1100 + iota
	CodeNeedLogin
	CodeInvalidateToken
	CodeInsufficientPerms
)

var responseCodeMessage = map[RespCodeType]string{
	CodeSuccess:           "success",
	CodeInvalidParams:     "Invalid parameters",
	CodeInvalidPassword:   "Invalid password",
	CodeUserExists:        "User exists",
	CodeUserNotExists:     "User not exists or password error",
	CodeServerBusy:        "Service busy",
	CodeServerInternal:    "Server internal",
	CodeInvalidateAuth:    "Invalid Authorization",
	CodeNeedLogin:         "need login",
	CodeInvalidateToken:   "Invalid token",
	CodeInsufficientPerms: "Insufficient permissions",
}

func (code RespCodeType) getMsg() string {
	msg, ok := responseCodeMessage[code]
	if !ok {
		msg = "Service busy"
	}
	return msg

}
