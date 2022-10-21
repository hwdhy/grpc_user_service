package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc_demo/db"
	"grpc_demo/pb"
	"grpc_demo/service"
	"log"
	"math/rand"
	"net"
	"net/http"
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
	db.InitConnectionPgsql() // 初始化数据库连接

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if *serverType == "grpc" {
		logrus.Fatal(runGRPCServer(listen))
	} else {
		logrus.Fatal(runRESTServer(listen))
	}
}

func runGRPCServer(listen net.Listener) error {
	server := grpc.NewServer()
	pb.RegisterUserServer(server, &service.User{})
	logrus.Printf("server listening at %v", listen.Addr())

	return server.Serve(listen)
}

func runRESTServer(listen net.Listener) error {
	mux := runtime.NewServeMux()
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	err := pb.RegisterUserHandlerFromEndpoint(ctx, mux, *endpoint, dialOptions)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Printf("start REST server at %s", listen.Addr())
	return http.Serve(listen, mux)
}
