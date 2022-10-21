package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"grpc_demo"
	"grpc_demo/models"
	"time"
)

var PgsqlDB *gorm.DB

// InitConnectionPgsql 数据库连接
func InitConnectionPgsql() {
	pDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			grpc_demo.PgsqlUsername, grpc_demo.PgsqlPassword, grpc_demo.PgsqlDbname, grpc_demo.PgsqlPort),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	DB, err := pDB.DB()
	DB.SetMaxOpenConns(100)
	DB.SetConnMaxIdleTime(10)
	DB.SetConnMaxLifetime(time.Minute)
	if err != nil {
		logrus.Fatalf("connect pgsql db err: %v", err)
	}

	PgsqlDB = pDB
	_ = AutoMigrate()
}

// AutoMigrate 自动建表
func AutoMigrate() error {
	return PgsqlDB.AutoMigrate(
		&models.User{},
	)
}
