package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc_demo"
	"hwdhy/utools/common"
)

type AuthInterceptor struct {
}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		md, _ := metadata.FromIncomingContext(ctx)

		token := md["grpcgateway-cookie"][0]
		userID := common.GetUserID(token, grpc_demo.TokenKey)
		if userID == 0 {
			return nil, fmt.Errorf("user not exist")
		}

		logrus.Printf("--- interceptor: %s", info.FullMethod)
		return handler(ctx, req)
	}
}
