package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pUser "shop_srvs/proto/user"
	"time"
)

var userClient pUser.UserClient
var conn *grpc.ClientConn

func init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	userClient = pUser.NewUserClient(conn)
}
func main() {
	defer conn.Close()
	for i := 5; i < 7; i++ {
		user, err := userClient.CreateUser(context.Background(), &pUser.CreateUserRequest{
			Nickname: fmt.Sprintf("喵喵%d", i),
			Password: "123456",
			Mobile:   fmt.Sprintf("1815956737%d", i),
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
