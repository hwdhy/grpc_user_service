package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
	"user_service"
	"user_service/models"
)

var PgsqlDB *gorm.DB

// InitConnectionPgsql 数据库连接
func InitConnectionPgsql() {
	pDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
			user_service.PgsqlUsername, user_service.PgsqlPassword, user_service.PgsqlDbname, user_service.PgsqlHost, user_service.PgsqlPort),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	DB, err := pDB.DB()
	if err != nil {
		logrus.Fatalf("connect pgsql db err: %v", err)
	}
	DB.SetMaxOpenConns(100)
	DB.SetConnMaxIdleTime(10)
	DB.SetConnMaxLifetime(time.Minute)

	PgsqlDB = pDB
	_ = AutoMigrate()
}

// AutoMigrate 自动建表
func AutoMigrate() error {
	return PgsqlDB.AutoMigrate(
		&models.User{},
	)
}
