package main

import (
	"context"
	"fmt"
	"testing"
	"user-server/proto"
	. "user-server/proto"

	"google.golang.org/grpc"
)

var client proto.UserClient

func Init() *grpc.ClientConn {
	conn, err := grpc.Dial("127.0.0.1:10001", grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprint("failed grpc dial", err.Error()))
	}
	// defer conn.Close()

	client = proto.NewUserClient(conn)

	return conn
}

func TestGetUserList(t *testing.T) {
	usersRsp, err := client.GetUserList(context.Background(), &PageInfo{Page: 1, PageSize: 3})
	if err != nil {
		panic(fmt.Sprint(" error:", err))
	}
	for i, u := range usersRsp.Data {
		fmt.Println(i, u.Id, u.Mobile, u.NickName, u.Role, u.Birthday)
		checkPwd, _ := client.CheckPassword(
			context.Background(),
			&CheckPasswordInfo{Password: "password123", EncryptedPassword: u.Password},
		)
		fmt.Println("checked password:", checkPwd.Success)
	}
}

func TestCreateUser(t *testing.T) {
	user, err := client.CreateUser(
		context.Background(),
		&CreateUserInfo{
			NickName: "zheng",
			Password: "zhengmr123",
			Mobile:   "13340717428",
		})
	if err != nil {
		panic(fmt.Sprint("error:", err))
	}
	fmt.Println("create user:", user)
}

func TestGetUserById(t *testing.T) {
	user, err := client.GetUserById(context.Background(), &IdRequest{Id: 3})
	if err != nil {
		panic(fmt.Sprint("error:", err))
	}
	fmt.Println("get user by id:", user.Mobile)
}

func TestGetUserByMobile(t *testing.T) {
	user, err := client.GetUserByMobile(context.Background(), &MobileRequest{Mobile: "13340717428"})
	if err != nil {
		panic(fmt.Sprint("error:", err))
	}
	fmt.Println("get user by Mobile:", user.Mobile)
}

func TestUserUpdate(t *testing.T) {
	_, err := client.UpdateUser(context.Background(), &UpdateUserInfo{
		Id:       10,
		NickName: "zhouzheng",
		Mobile:   "13340717428",
	})
	if err != nil {
		panic(fmt.Sprint("error:", err))
	}
}

func main() {
	conn := Init()
	defer conn.Close()

	// TestCreateUser(nil) // 创建
	// TestGetUserList(nil) // 查列表
	// TestGetUserById(nil)// 主键查询
	// TestGetUserByMobile(nil)// 条件查询
	TestUserUpdate(nil) // 更新用户信息

}
