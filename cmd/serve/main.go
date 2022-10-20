package main

import (
	"flag"
	"fmt"
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
	log.Printf("server listening at %v", listen.Addr())
	log.Fatal(s.Serve(listen))
}
