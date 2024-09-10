package main

import (
	"fmt"
	"test/config"
	"test/database"
	_ "test/docs" // 导入 Swagger 文档
	"test/logger"
	"test/middleware"
	"test/router"
	"test/storage"
)

// @title 用户认证 API
// @version 1.0
// @description 这是一个用户认证系统的 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	if err := config.LoadConfig(); err != nil {
		panic(fmt.Sprintf("加载配置失败: %v", err))
	}

	if err := logger.Init(); err != nil {
		panic(fmt.Sprintf("初始化日志失败: %v", err))
	}
	defer logger.Log.Sync()

	// 初始化数据库连接d
	if err := database.ConnectDatabase(); err != nil {
		logger.Sugar.Fatalf("连接数据库失败: %v", err)
	}

	// 创建超级管理员账号
	if err := database.SeedSuperAdmin(database.DB); err != nil {
		logger.Sugar.Fatalf("创建超级管理员账号失败: %v", err)
	}

	// 初始化存储
	if err := storage.InitStorage(); err != nil {
		logger.Sugar.Fatalf("初始化存储失败: %v", err)
	}

	// 设置路由
	r := router.SetupRouter()

	// 使用 TraceMiddleware
	r.Use(middleware.TraceMiddleware())

	// 启动服务器
	serverAddr := config.AppConfig.Server.Port
	logger.Sugar.Infof("服务器正在启动，监听端口 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		logger.Sugar.Fatalf("启动服务器失败: %v", err)
	}
}
