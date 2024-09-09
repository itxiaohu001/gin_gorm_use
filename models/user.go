package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(255);uniqueIndex"` // 修改这里
	Password string `gorm:"type:varchar(255)"`             // 修改这里
	Role     string `gorm:"type:varchar(50);default:user"` // 修改这里
}

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

// UpdateRoleInput 结构体用于更新用户角色
type UpdateRoleInput struct {
	Role string `json:"role" binding:"required"`
}

// LoginInput 结构体用于用户登录
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 结构体用于登录响应
type LoginResponse struct {
	Token string `json:"token"`
}
