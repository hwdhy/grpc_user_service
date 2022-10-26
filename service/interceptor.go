package service

import (
	"context"
	"errors"
	"github.com/casbin/casbin/v2"
	"github.com/hwdhy/grpc_tools/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc_demo"
)

type AuthInterceptor struct {
}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (interceptor *AuthInterceptor) Unary(enforcer *casbin.Enforcer) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		md, _ := metadata.FromIncomingContext(ctx)
		token := md["grpcgateway-cookie"][0]

		_, role := common.GetUserID(token, grpc_demo.TokenKey)
		if role == "" {
			role = "tourists"
		}

		res, err := enforcer.Enforce(role, info.FullMethod, info.Server)
		if err != nil {
			return nil, errors.New("permission verification failure")
		}
		if res {
			return handler(ctx, req)
		} else {
			return nil, errors.New("unauthorized")
		}
	}
}
