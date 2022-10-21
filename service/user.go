package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/sirupsen/logrus"
	"grpc_demo/db"
	"grpc_demo/models"
	"grpc_demo/pb"
	"grpc_demo/tools"
)

type User struct {
	pb.UnimplementedUserServer
}

// Register 用户注册
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
		logrus.Printf("register user err: %v", err)
		return &pb.UserRegisterResponse{Status: "error"}, err
	}
	logrus.Printf("create user(%+v) success", userData)

	return &pb.UserRegisterResponse{Status: "success"}, nil
}

// Login 用户登录
func (u *User) Login(ctx context.Context, input *pb.UserLoginRequest) (*pb.UserLoginResponse, error) {
	// 1. 判断用户是否存在
	var user models.User
	if err := db.PgsqlDB.Model(models.User{}).Where("username = ?", input.Username).First(&user).Error; err != nil {
		logrus.Errorf("find user(%s) err: %v", input.Username, err)
	}
	if user.ID == 0 {
		return nil, errors.New("user not exist")
	}
	// 2. 密码校验
	hash := md5.New()
	hash.Write([]byte(input.Password))
	hashPassword := hex.EncodeToString(hash.Sum([]byte(user.Salt)))
	if hashPassword != user.Password {
		return nil, errors.New("password does not match")
	}
	logrus.Printf("user(%s) login success", input.Username)
	return &pb.UserLoginResponse{Status: "success"}, nil
}

func (u *User) List(ctx context.Context, input *pb.UserListRequest) (*pb.UserListResponse, error) {
	offset := (input.Page - 1) * input.PageSize

	var users []models.User
	if err := db.PgsqlDB.Model(models.User{}).Offset(int(offset)).Limit(int(input.PageSize)).Find(&users).Error; err != nil {
		logrus.Errorf("select user err:%v", err)
	}

	res := make([]*pb.UserList, len(users))
	for index, user := range users {

		res[index] = &pb.UserList{
			Id:         uint32(user.ID),
			Username:   user.Username,
			Password:   user.Password,
			Type:       user.Type,
			Ip:         user.IP,
			CreateTime: user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &pb.UserListResponse{
		Data: res,
	}, nil
}
