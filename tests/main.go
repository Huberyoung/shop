package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"

	pUser "shop_srvs/proto/user"
)

var userClient pUser.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = pUser.NewUserClient(conn)
}

func testGetUserById() {
	rsp, err := userClient.GetUserById(context.Background(), &pUser.IdRequest{Id: 13})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nGetUserById: %+v\n", rsp)
}

func testGetUserByMobile() {
	rsp, err := userClient.GetUserByMobile(context.Background(), &pUser.MobileRequest{Mobile: "18159567370"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n\nGetUserByMobile: %+v\n", rsp)
}

func TestGetUserList() {
	rsp, err := userClient.GetUserList(context.Background(), &pUser.PageInfo{PageNum: 1, PageSize: 2})
	if err != nil {
		panic(err)
	}

	for _, user := range rsp.Data {
		fmt.Printf("id:%d  gender:%d  mobile:%s  password:%s\n\n", user.Id, user.Gender, user.Mobile, user.Password)

		rsp1, err1 := userClient.CheckPassword(context.Background(), &pUser.PasswordCheckRequest{
			Password:          "123456",
			EncryptedPassword: user.Password,
		})
		if err1 != nil {
			panic(err1)
		}
		fmt.Println(rsp1.Success)

	}
}

func main() {
	Init()
	defer conn.Close()
	//TestGetUserList()
	//testGetUserById()
	//testGetUserByMobile()

	for i := 1; i < 10; i++ {
		user, err := userClient.CreateUser(context.Background(), &pUser.CreateUserRequest{
			Nickname: fmt.Sprintf("喵喵%d", i),
			Password: "123456",
			Mobile:   fmt.Sprintf("1815956736%d", i),
			Gender:   pUser.CreateUserRequest_MALE.Enum(),
			BirthDay: uint64(time.Now().Unix()),
		})
		if err != nil {
			fmt.Printf("userClient.CreateUser err:%s\n\n", err)
		} else {
			fmt.Printf("user:%+v\n", user)
		}
		time.Sleep(1 * time.Second)
	}
}
