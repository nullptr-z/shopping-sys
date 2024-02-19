package handler

import (
	"context"
	"fmt"

	. "shopping-sys/user_server/global"
	"shopping-sys/user_server/model"
	. "shopping-sys/user_server/proto"
	"shopping-sys/user_server/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type UserService struct{}

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

// 分页查询 users
func (*UserService) GetUserList(ctx context.Context, req *PageInfo) (*UserListResponse, error) {
	var users []model.User
	var rsp UserListResponse

	// 获取总条数
	result := DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp.Total = uint64(result.RowsAffected) // 总条数

	// 获取分页数据
	DB.Scopes(Paginate(int(req.Page), int(req.PageSize))).Find(&users)

	for _, u := range users {
		userRsp := IntoDbUseUserInfoResponse(u)
		rsp.Data = append(rsp.Data, userRsp)
	}

	return &rsp, nil
}

// 使用手机查询用户
func (*UserService) GetUserByMobile(ctx context.Context, req *MobileRequest) (*UserInfoResponse, error) {
	var user model.User

	result := DB.Where("mobile", req.Mobile).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.Error != nil {
		return nil, status.Error(codes.NotFound, "用户不存在")
	}

	userInfo := IntoDbUseUserInfoResponse(user)

	return userInfo, nil
}

// 使用ID查询用户
func (*UserService) GetUserById(ctx context.Context, req *IdRequest) (*UserInfoResponse, error) {
	var user model.User

	result := DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfo := IntoDbUseUserInfoResponse(user)

	return userInfo, nil
}

// 注册用户
func (slf *UserService) CreateUser(ctx context.Context, req *CreateUserInfo) (*UserInfoResponse, error) {
	_, err := slf.GetUserByMobile(ctx, &MobileRequest{Mobile: req.Mobile})
	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "用户已存在")
	}

	//  密码加密
	enPwd, salt := utils.CryptoPasswordWithSalt(req.Password)
	saltPwd := utils.MergePasswordSalt(salt, enPwd)
	// 创建用户
	user := model.User{
		Mobile:   req.Mobile,
		NickName: req.NickName,
		Password: saltPwd,
	}

	result := DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprint("用户创建失败", result.Error.Error()))
	}

	userInfo := IntoDbUseUserInfoResponse(user)

	return userInfo, nil
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
