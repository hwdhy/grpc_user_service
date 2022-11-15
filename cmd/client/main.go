package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"grpc_tools/etcd"
	"grpc_tools/pb/user_pb"
	"time"
)

func main() {
	conn := etcd.ClientConn("userService", 2, "hwdhy")
	if conn == nil {
		logrus.Fatalf("get grpc client err")
	}
	defer conn.Close()

	// 获取客户端对象
	c := user_pb.NewUserClient(conn)
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	for i := 0; i < 10; i++ {
		response, err := c.List(ctx, &user_pb.UserListRequest{
			Page:     1,
			PageSize: 10,
		})
		if err != nil {
			logrus.Fatalf("user not find: %v", err)
		}
		logrus.Printf("login success, count : %d, token: %v", i, response.GetCount())
	}
}
