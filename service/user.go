package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"grpc_demo/db"
	"grpc_demo/models"
	"grpc_demo/pb"
	"grpc_demo/tools"
	"log"
)

type User struct {
	pb.UnimplementedUserServer
}

func (u *User) Register(ctx context.Context, input *pb.UserRegisterRequest) (*pb.UserRegisterResponse, error) {
	// 生成4为随机盐值
	salt := tools.RandomString(4)
	hash := md5.New()
	hash.Write([]byte(input.GetPassword()))
	hashPassword := hex.EncodeToString(hash.Sum([]byte(salt)))

	userData := models.User{
		Username: input.GetUsername(),
		Password: hashPassword,
		Salt:     salt,
	}
	if err := db.PgsqlDB.Model(models.User{}).Create(&userData).Error; err != nil {
		log.Println("register user err:", err)
		return &pb.UserRegisterResponse{Status: "error"}, err
	}

	return &pb.UserRegisterResponse{Status: "success"}, nil
}
