package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户表
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique"` // 用户名称
	Password  string //用户密码
	Salt      string // 密码盐值
	Type      string `gorm:"default:0"`        // 用户类型 0:普通用户
	Role      string `gorm:"default:'member'"` // 角色
	IP        string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (*User) TableName() string {
	return "user"
}
