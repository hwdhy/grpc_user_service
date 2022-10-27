package main

import (
	"context"
	"flag"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dial server err: %v", err)
	}
	defer conn.Close()
	c := user_pb.NewUserClient(conn)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second)
	defer cancelFunc()

	//response, err := c.Register(ctx, &pb.UserRegisterRequest{
	//	Username: "1111",
	//	Password: "2222",
	//})

	response, err := c.Login(ctx, &user_pb.UserLoginRequest{
		Username: "1111",
		Password: "2222",
	})
	if err != nil {
		logrus.Fatalf("could not greet: %v", err)
	}
	logrus.Printf("create user %v", response.GetToken())
}
