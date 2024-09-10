package database

import (
	"fmt"
	"test/config"
	"test/logger"
	"test/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase 连接到数据库
func ConnectDatabase() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.Database.Username,
		config.AppConfig.Database.Password,
		config.AppConfig.Database.Host,
		config.AppConfig.Database.Port,
		config.AppConfig.Database.Name)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("无法连接到数据库: %v", err)
	}

	// 自动迁移数据库结构
	err = DB.AutoMigrate(&models.User{}, &models.File{})
	if err != nil {
		return fmt.Errorf("自动迁移失败: %v", err)
	}

	logger.Sugar.Info("数据库连接成功")
	return nil
}
