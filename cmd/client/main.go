package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"grpc_demo/db"
	"grpc_tools/pb/user_pb"
	"log"
	"time"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()

	etcdResolverBuilder := db.NewEtcdResolverBuilder()
	resolver.Register(etcdResolverBuilder)

	conn, err := grpc.Dial("etcd:///", grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial server err: %v", err)
	}
	defer conn.Close()

	// 获取客户端对象
	c := user_pb.NewUserClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	//response, err := c.Register(ctx, &pb.UserRegisterRequest{
	//	Username: "1111",
	//	Password: "2222",
	//})

	response, err := c.Login(ctx, &user_pb.UserLoginRequest{
		Username: "admin",
		Password: "123456",
	})
	if err != nil {
		logrus.Fatalf("could not user: %v", err)
	}
	logrus.Printf("create user %v", response.GetToken())
}
