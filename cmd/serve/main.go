package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc_demo/db"
	"grpc_demo/service"
	"grpc_tools/common"
	"grpc_tools/etcd"
	"grpc_tools/pb/user_pb"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

var (
	port       = flag.Int("port", 50051, "the server port")
	serverType = flag.String("type", "grpc", "type of server(grpc/rest)")
	endpoint   = flag.String("endpoint", "0.0.0.0:8080", "grpc endpoint")
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	if *serverType == "grpc" {
		logrus.Fatal(runGRPCServer(listen))
	} else {
		logrus.Fatal(runRESTServer(listen))
	}
}

// 启动grpc服务
func runGRPCServer(listen net.Listener) error {
	db.InitConnectionPgsql() // 初始化数据库连接

	//更新接口权限
	e := common.InitAdapter([]map[string]int{
		service.UserPermission,
	})

	interceptor := service.NewAuthInterceptor()
	serverOptions := []grpc.ServerOption{
		grpc.UnaryInterceptor(interceptor.Unary(e)),
	}

	server := grpc.NewServer(serverOptions...)
	user_pb.RegisterUserServer(server, &service.User{})

	etcdRegister, err := etcd.NewEtcdRegister()
	if err != nil {
		logrus.Fatal(err)
	}
	defer etcdRegister.Close()
	serviceName := "user_service"
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(*port))

	err = etcdRegister.RegisterServer("/etcd/"+serviceName, addr, 5)
	if err != nil {
		logrus.Fatalf("register error %v ", err)
	}

	logrus.Printf("server listening at %v", listen.Addr())
	return server.Serve(listen)
}

// 启动rest服务
func runRESTServer(listen net.Listener) error {
	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	err := user_pb.RegisterUserHandlerFromEndpoint(ctx, mux, *endpoint, dialOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Printf("start REST server at %s", listen.Addr())
	return http.Serve(listen, mux)
}
