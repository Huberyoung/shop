package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pUser "shop_srvs/proto/user"
)

var userClient pUser.UserClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = pUser.NewUserClient(conn)
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
	TestGetUserList()

	//
	//for i := 5; i < 7; i++ {
	//	user, err := userClient.CreateUser(context.Background(), &pUser.CreateUserRequest{
	//		Nickname: fmt.Sprintf("喵喵%d", i),
	//		Password: "123456",
	//		Mobile:   fmt.Sprintf("1815956737%d", i),
	//		Gender:   pUser.CreateUserRequest_MALE.Enum(),
	//		BirthDay: uint64(time.Now().Unix()),
	//	})
	//	if err != nil {
	//		fmt.Printf("userClient.CreateUser err:%s\n\n", err)
	//	} else {
	//		fmt.Printf("user:%+v\n", user)
	//	}
	//	time.Sleep(1 * time.Second)
	//}
}
