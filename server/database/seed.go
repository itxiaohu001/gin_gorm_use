package database

import (
	"errors"
	"gorm.io/gorm"
	"test/logger"
	"test/models"
	"test/utils"
)

func SeedSuperAdmin(db *gorm.DB) error {
	var existingAdmin models.User
	result := db.Where("role = ?", models.RoleAdmin).First(&existingAdmin)

	if result.Error == nil {
		// 超级管理员已存在，记录日志并返回
		logger.Sugar.Info("超级管理员已存在，用户名：", existingAdmin.Username)
		return nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 发生了其他错误
		return result.Error
	}

	// 超级管理员不存在，创建新的超级管理员
	hashedPassword, err := utils.HashPassword("superadmin123")
	if err != nil {
		return err
	}

	superAdmin := models.User{
		Username: "superadmin",
		Password: hashedPassword,
		Role:     models.RoleAdmin,
	}

	result = db.Create(&superAdmin)
	if result.Error != nil {
		return result.Error
	}

	logger.Sugar.Info("成功创建超级管理员，用户名：", superAdmin.Username)
	return nil
}
