package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"shop_srvs/global"
	"shop_srvs/handler"
	"shop_srvs/model/starup"
	"shop_srvs/model/user"
	pUser "shop_srvs/proto/user"
)

func init() {
	err := starup.SetUpSetting()
	if err != nil {
		log.Fatalf("Init.setUpSetting err:%v", err)
	}

	err = starup.SetUpDBEngine()
	if err != nil {
		log.Fatalf("Init.SetUpDBEngine err:%v", err)
	}
}

func main() {
	var port int
	var ip string
	flag.StringVar(&ip, "ip", "127.0.01", "ip地址")
	flag.IntVar(&port, "port", 50051, "端口号")
	address := fmt.Sprintf("%s:%d", ip, port)
	fmt.Println(address)

	server := grpc.NewServer()
	pUser.RegisterUserServer(server, &handler.UserServer{})
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("net.Listen err:%s\n", err)
	}

	err = global.DBEngine.AutoMigrate(&user.User{})
	if err != nil {
		log.Fatalf("global.DBEngine.AutoMigrate err:%s\n", err)
	}

	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("server.Serve err:%s\n", err)
	}
}
