package handler

import (
	"context"
	"fmt"
	"time"

	. "shopping-sys/user_server/global"
	"shopping-sys/user_server/model"
	"shopping-sys/user_server/proto"
	. "shopping-sys/user_server/proto"
	"shopping-sys/user_server/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type UserService struct {
	proto.UnimplementedUserServer // 嵌入这个结构体
}

// 分页逻辑
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}

		// switch {
		// case pageSize > 100:
		// 	pageSize = 100
		// case pageSize <= 0:
		// 	pageSize = 10
		// }

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

// 分页查询 users
func (*UserService) GetUserList(ctx context.Context, req *PageInfo) (*UserListResponse, error) {
	var users []model.User
	var rsp UserListResponse
	var total int64

	// 获取总条数
	result := DB.Model(&model.User{}).Count(&total)
	if result.Error != nil {
		return nil, result.Error
	}
	rsp.Total = uint64(total) // 总条数

	// 获取分页数据
	// page := req.Page
	// pageSize := req.PageSize
	// if page <= 0 {
	// 	page = 1
	// }
	// offset := (page - 1) * pageSize
	// fmt.Println("pageSize:", pageSize)
	// fmt.Println("offset:", offset)
	// DB.Offset(int(offset)).Limit(int(pageSize)).Find(&users)
	result = DB.Scopes(Paginate(int(req.Page), int(req.PageSize))).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

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
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	userInfo := IntoDbUseUserInfoResponse(user)

	return userInfo, nil
}

// 注册用户
func (slf *UserService) CreateUser(ctx context.Context, req *CreateUserInfo) (*UserInfoResponse, error) {
	_, err := slf.GetUserByMobile(ctx, &MobileRequest{Mobile: req.Mobile})
	if err == nil {
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
		BaseModel: model.BaseModel{
			CratedAt:  time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	result := DB.Create(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprint("用户创建失败", result.Error.Error()))
	}

	userInfo := IntoDbUseUserInfoResponse(user)

	return userInfo, nil
}

// 修改用户信息
func (*UserService) UpdateUser(ctx context.Context, req *UpdateUserInfo) (*emptypb.Empty, error) {
	var user model.User

	result := DB.First(&user, req.Id)
	if result.Error != nil {
		return nil, status.Errorf(codes.NotFound, "用户不存在")
	}

	user.Mobile = req.Mobile
	user.NickName = req.NickName
	user.Gender = req.Gender
	if req.Birthday > 0 {
		birthday := time.Unix(int64(req.Birthday), 0)
		user.Birthday = &birthday
	}

	result = DB.Save(&user)
	if result.Error != nil {
		return nil, status.Errorf(codes.Internal, result.Error.Error())
	}
	return &emptypb.Empty{}, nil
}

// 验证密码
func (*UserService) CheckPassword(ctx context.Context, req *CheckPasswordInfo) (*CheckedResponse, error) {
	pwd, salt := utils.SplitPasswordSalt(req.EncryptedPassword)
	is_validate := utils.ValidPassword(pwd, salt, req.Password)
	// if !is_validate {
	// 	return &CheckedResponse{Success: false}, status.Errorf(codes.PermissionDenied, "密码校验不通过")
	// }

	return &CheckedResponse{Success: is_validate}, nil
}
