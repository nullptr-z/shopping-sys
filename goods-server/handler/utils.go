package handler

import (
	"goods-server/model"
	"goods-server/proto"
)

// 数据库查询的 user 数据转换到 proto.UserInfoResponse
func IntoDbUseUserInfoResponse(u model.User) *proto.UserInfoResponse {
	var user = proto.UserInfoResponse{
		Id:       u.ID,
		Mobile:   u.Mobile,
		Password: u.Password,
		NickName: u.NickName,
		Gender:   u.Gender,
		Role:     int32(u.Role),
	}
	if u.Birthday != nil {
		user.Birthday = uint64(u.Birthday.Unix())
	}
	return &user
}
