package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"test/handlers"
	"test/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Swagger 文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 添加 TraceMiddleware
	r.Use(middleware.TraceMiddleware())

	// 登录
	r.POST("/login", handlers.Login)

	// 管理员路由
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminAuthMiddleware())
	{
		admin.POST("/register", handlers.Register)
		admin.DELETE("/users/:id", handlers.DeleteUser)
		admin.PUT("/users/:id/role", handlers.UpdateUserRole)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// 文件相关路由
		files := protected.Group("/files")
		{
			files.GET("", handlers.ListFiles) // 新添加的路由
			files.POST("/upload", handlers.UploadFile)
			files.GET("/download/:fileID", handlers.DownloadFile)
			files.DELETE("/:fileID", handlers.DeleteFile)
		}
		// 添加其他需要认证的路由
	}

	return r
}

// 展示登录页面
func showLoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}
