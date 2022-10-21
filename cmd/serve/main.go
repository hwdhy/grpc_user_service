package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc_demo/db"
	"grpc_demo/pb"
	"grpc_demo/service"
	"log"
	"math/rand"
	"net"
	"time"
)

var (
	port = flag.Int("port", 50051, "the server port")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	db.InitConnectionPgsql() // 初始化数据库连接

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServer(s, &service.User{})
	logrus.Printf("server listening at %v", listen.Addr())
	// 启动服务
	logrus.Fatal(s.Serve(listen))
}
