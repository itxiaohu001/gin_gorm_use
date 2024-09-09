package handlers

import (
	"net/http"
	"test/dto"
	"test/logger"
	"test/services"
	"test/utils"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary 注册新用户
// @Description 创建一个新用户账号
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body dto.UserRegisterInput true "用户注册信息"
// @Success 201 {object} dto.RegisterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /register [post]
func Register(c *gin.Context) {
	var input dto.UserRegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := services.RegisterUser(c.Request.Context(), input.Username, input.Password, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "注册用户失败"})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 根据用户ID删除指定用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} dto.DeleteUserResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	resp, err := services.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateUserRole godoc
// @Summary 更新用户角色
// @Description 根据用户ID更新指定用户的角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Param role body dto.UpdateRoleInput true "新的用户角色"
// @Success 200 {object} dto.UpdateRoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /users/{id}/role [put]
func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")
	var input dto.UpdateRoleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	resp, err := services.UpdateUserRole(c.Request.Context(), userID, input.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "更户角色失败"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Login godoc
// @Summary 用户登录
// @Description 验证用户凭证并重定向到仪表板或返回错误信息
// @Tags 用户管理
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 302 {string} string "重定向到仪表板"
// @Failure 400 {object} dto.ErrorResponse "请求格式错误"
// @Failure 401 {object} dto.ErrorResponse "用户名或密码错误"
// @Failure 500 {object} dto.ErrorResponse "内部服务器错误"
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "请求格式错误", Message: "用户名和密码不能为空"})
		return
	}

	logger.Sugar.Infow("尝试登录",
		"traceID", utils.GetTraceID(c),
		"username", username,
	)

	resp, err := services.LoginUser(c.Request.Context(), username, password)
	if err != nil {
		logger.Sugar.Errorw("登录失败",
			"traceID", utils.GetTraceID(c),
			"error", err,
		)
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{Error: "用户名或密码错误", Message: err.Error()})
		return
	}

	logger.Sugar.Infow("登录成功",
		"traceID", utils.GetTraceID(c),
		"username", username,
	)

	// 设置 cookie
	c.SetCookie("token", resp.Token, 3600, "/", "", false, true)

	// 重定向到仪表板
	c.Redirect(http.StatusFound, "/dashboard")
}
