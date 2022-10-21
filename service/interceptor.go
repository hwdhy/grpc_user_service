package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type AuthInterceptor struct {
}

func NewAuthInterceptor() *AuthInterceptor {
	return &AuthInterceptor{}
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		logrus.Printf("--- interceptor: %s", info.FullMethod)
		return handler(ctx, req)
	}
}
