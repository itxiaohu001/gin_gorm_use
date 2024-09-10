package services

import (
	"context"
	"errors"
	"runtime"
	"test/database"
	"test/dto"
	"test/logger"
	"test/models"
	"test/utils"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(ctx context.Context, username, password, role string) (dto.RegisterUserResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "RegisterUser",
	)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Username: username, Password: string(hashedPassword), Role: role}

	result := database.DB.Create(&user)
	if result.Error != nil {
		logger.Sugar.Errorw("用户注册失败",
			"error", result.Error,
			"username", username,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.RegisterUserResponse{}, errors.New("用户名已存在")
	}
	return dto.RegisterUserResponse{Message: "用户注册成功"}, nil
}

func LoginUser(ctx context.Context, username, password string) (dto.LoginUserResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "LoginUser",
	)

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		logger.Sugar.Errorw("用户登录失败：用户不存在",
			"error", err,
			"username", username,
		)
		return dto.LoginUserResponse{}, errors.New("用户名或密码无效")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		logger.Sugar.Errorw("用户登录失败：密码错误",
			"error", err,
			"username", username,
		)
		return dto.LoginUserResponse{}, errors.New("用户名或密码无效")
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		logger.Sugar.Errorw("生成令牌失败",
			"error", err,
			"userID", user.ID,
		)
		return dto.LoginUserResponse{}, errors.New("生成令牌失败")
	}

	return dto.LoginUserResponse{Token: token, Role: user.Role}, nil
}

func DeleteUser(ctx context.Context, userID string) (dto.DeleteUserResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "DeleteUser",
	)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		logger.Sugar.Errorw("删除用户失败：用户不存在",
			"error", err,
			"userID", userID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.DeleteUserResponse{}, errors.New("用户不存在")
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		logger.Sugar.Errorw("删除用户失败",
			"error", err,
			"userID", userID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.DeleteUserResponse{}, errors.New("删除用户失败")
	}

	return dto.DeleteUserResponse{Message: "用户删除成功"}, nil
}

func UpdateUserRole(ctx context.Context, userID, newRole string) (dto.UpdateUserRoleResponse, error) {
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", "UpdateUserRole",
	)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		logger.Sugar.Errorw("更新用户角色失败：用户不存在",
			"error", err,
			"userID", userID,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.UpdateUserRoleResponse{}, errors.New("用户不存在")
	}

	user.Role = newRole
	if err := database.DB.Save(&user).Error; err != nil {
		logger.Sugar.Errorw("更新用户角色失败",
			"error", err,
			"userID", userID,
			"newRole", newRole,
			"traceID", utils.GetTraceID(ctx),
		)
		return dto.UpdateUserRoleResponse{}, errors.New("更新用户角色失败")
	}

	return dto.UpdateUserRoleResponse{Message: "用户角色更新成功"}, nil
}

// logFunctionName 打印当前函数名
func logFunctionName(ctx context.Context) {
	pc, _, _, _ := runtime.Caller(1)
	functionName := runtime.FuncForPC(pc).Name()
	logger.Sugar.Infow("执行函数",
		"traceID", utils.GetTraceID(ctx),
		"functionName", functionName,
	)
}
