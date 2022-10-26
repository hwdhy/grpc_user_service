package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/metadata"
	"grpc_demo"
	"grpc_demo/db"
	"grpc_demo/models"
	"hwdhy/Htools/common"
	"hwdhy/Htools/pb/userPB"
)

var UserPermission = map[string]int{
	"/grpc_hwdhy.User/List":     common.Admin,
	"/grpc_hwdhy.User/Register": common.NotLogged,
	"/grpc_hwdhy.User/Login":    common.NotLogged,
}

type User struct {
	userPB.UnimplementedUserServer
}

// Register 用户注册
func (u *User) Register(ctx context.Context, input *userPB.UserRegisterRequest) (*userPB.UserRegisterResponse, error) {
	// 生成4为随机盐值
	salt := common.RandomString(4)
	hashPassword := common.StringHash(input.GetPassword(), salt)

	// get request ip
	md, _ := metadata.FromIncomingContext(ctx)
	remoteIP := md["x-forwarded-for"][0]

	userData := models.User{
		Username: input.GetUsername(),
		Password: hashPassword,
		Salt:     salt,
		IP:       remoteIP,
	}
	if err := db.PgsqlDB.Model(models.User{}).Create(&userData).Error; err != nil {
		logrus.Printf("register user err: %v", err)
		return &userPB.UserRegisterResponse{Status: "error"}, err
	}
	logrus.Printf("create user(%+v) success", userData)

	return &userPB.UserRegisterResponse{Status: "success"}, nil
}

// Login 用户登录
func (u *User) Login(ctx context.Context, input *userPB.UserLoginRequest) (*userPB.UserLoginResponse, error) {
	// 1. 判断用户是否存在
	var user models.User
	if err := db.PgsqlDB.Model(models.User{}).Where("username = ?", input.Username).First(&user).Error; err != nil {
		logrus.Errorf("find user(%s) err: %v", input.Username, err)
	}
	if user.ID == 0 {
		return nil, errors.New("user not exist")
	}
	// 2. 密码校验
	hashPassword := common.StringHash(input.Password, user.Salt)
	if hashPassword != user.Password {
		return nil, errors.New("password does not match")
	}
	logrus.Printf("user(%s) login success", input.Username)

	token := common.GenerateToken(uint64(user.ID), user.Role, grpc_demo.TokenKey)
	return &userPB.UserLoginResponse{Token: token}, nil
}

func (u *User) List(ctx context.Context, input *userPB.UserListRequest) (*userPB.UserListResponse, error) {
	offset := (input.Page - 1) * input.PageSize

	var users []models.User
	if err := db.PgsqlDB.Model(models.User{}).Offset(int(offset)).Limit(int(input.PageSize)).Find(&users).Error; err != nil {
		logrus.Errorf("select user err:%v", err)
	}

	res := make([]*userPB.UserList, len(users))
	for index, user := range users {

		res[index] = &userPB.UserList{
			Id:         uint32(user.ID),
			Username:   user.Username,
			Password:   user.Password,
			Type:       user.Type,
			Ip:         user.IP,
			CreateTime: user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &userPB.UserListResponse{
		Data: res,
	}, nil
}
