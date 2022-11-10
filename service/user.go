package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"grpc_tools/common"
	"grpc_tools/pb/user_pb"
	"user_service/db"
	"user_service/models"
)

var UserPermission = map[string]int{
	"/grpc_hwdhy.User/List":     common.Admin,
	"/grpc_hwdhy.User/Register": common.NotLogged,
	"/grpc_hwdhy.User/Login":    common.NotLogged,
}

type User struct {
	user_pb.UnimplementedUserServer
}

// Register 用户注册
func (u *User) Register(ctx context.Context, input *user_pb.UserRegisterRequest) (*user_pb.UserRegisterResponse, error) {
	// 判断用户是否存在
	var findUser models.User
	db.PgsqlDB.Model(models.User{}).Where("username = ?", input.GetUsername()).First(&findUser)
	if findUser.ID != 0 {
		return nil, status.Errorf(codes.AlreadyExists, "user(%s) is already exists", input.GetUsername())
	}

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
		return nil, status.Errorf(codes.Internal, "create user err: %v", err)
	}
	logrus.Printf("create user(%+v) success", userData)

	return &user_pb.UserRegisterResponse{Code: uint32(codes.OK), Msg: "success"}, nil
}

// Login 用户登录
func (u *User) Login(ctx context.Context, input *user_pb.UserLoginRequest) (*user_pb.UserLoginResponse, error) {
	// 1. 判断用户是否存在
	var user models.User
	if err := db.PgsqlDB.Model(models.User{}).Where("username = ?", input.Username).First(&user).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "can't find user, err: %v", err)
	}
	if user.ID == 0 {
		return nil, status.Errorf(codes.NotFound, "incorrect username/password")
	}
	// 2. 密码校验
	hashPassword := common.StringHash(input.Password, user.Salt)
	if hashPassword != user.Password {
		return nil, status.Error(codes.Internal, "password does not match")
	}
	logrus.Printf("user(%s) login success", input.Username)

	token := common.GenerateToken(uint64(user.ID), user.Role)
	return &user_pb.UserLoginResponse{Code: uint32(codes.OK), Token: token}, nil
}

func (u *User) List(ctx context.Context, input *user_pb.UserListRequest) (*user_pb.UserListResponse, error) {
	offset := (input.Page - 1) * input.PageSize

	var users []models.User
	if err := db.PgsqlDB.Model(models.User{}).Offset(int(offset)).Limit(int(input.PageSize)).Find(&users).Error; err != nil {
		return nil, status.Errorf(codes.Internal, "select user err: %v", err)
	}

	dataCount := len(users)
	res := make([]*user_pb.UserList, dataCount)
	for index, user := range users {
		res[index] = &user_pb.UserList{
			Id:         uint32(user.ID),
			Username:   user.Username,
			Password:   user.Password,
			Type:       user.Type,
			Ip:         user.IP,
			CreateTime: user.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}
	return &user_pb.UserListResponse{
		Code:  uint32(codes.OK),
		Count: uint32(dataCount),
		Data:  res,
	}, nil
}
