package handler

import (
	"context"
	"shopping-sys/user_server/global"
	"shopping-sys/user_server/model"
	"shopping-sys/user_server/proto"

	"gorm.io/gorm"
)

type UserServer struct{}

// 分页查询 users
func (*UserServer) GetUserList(ctx context.Context, req *proto.PageInfo) (*proto.UserListResponse, error) {
	var users []model.User
	var rsp proto.UserListResponse

	// 获取总条数
	result := global.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp.Total = uint64(result.RowsAffected) // 总条数

	// 获取分页数据
	global.DB.Scopes(Paginate(int(req.Page), int(req.PageSize))).Find(&users)

	for _, u := range users {
		userRsp := IntoDbUseUserInfoResponse(u)
		rsp.Data = append(rsp.Data, &userRsp)
	}

	return &rsp, nil
}

// 分页逻辑
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func IntoDbUseUserInfoResponse(u model.User) proto.UserInfoResponse {
	var userRsp proto.UserInfoResponse

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
	return userRsp
}

// type UserServer interface {
// 	// 注册
// 	CreateUser(context.Context, *CreateUserInfo) (*UserInfoResponse, error)
// 	// // 登录
// 	// rpc Login(CreateUserInfo) returns(UserInfoResponse);
// 	// 验证密码接口
// 	CheckPassword(context.Context, *CheckPasswordInfo) (*CheckedResponse, error)
// 	// 获取用户列表
// 	GetUserList(context.Context, *PageInfo) (*UserListfoResponse, error)
// 	// 通过手机号码获取用户
// 	GetUserByMobile(context.Context, *MobileRequest) (*UserInfoResponse, error)
// 	// 通过用户 ID 获取用户
// 	GetUserById(context.Context, *IdRequest) (*UserInfoResponse, error)
// 	// 更新用户信息
// 	UpdateUser(context.Context, *UpdateUserInfo) (*emptypb.Empty, error)
// 	mustEmbedUnimplementedUserServer()
// }
