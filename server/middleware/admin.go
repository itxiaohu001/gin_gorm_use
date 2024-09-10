package middleware

import (
	"net/http"
	"test/database"
	"test/models"
	"test/utils"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户ID
		userID, exists := utils.GetUserID(c.Request.Context())
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			c.Abort()
			return
		}

		// 检查用户是否为管理员
		var user models.User
		result := database.DB.First(&user, userID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "检查用户权限时出错"})
			c.Abort()
			return
		}

		if user.Role != models.RoleAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}

		c.Next()
	}
}
