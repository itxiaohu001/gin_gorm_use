package dto

type UserRegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UpdateRoleInput struct {
	Role string `json:"role" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

// RegisterResponse 用户注册响应
type RegisterResponse struct {
	Message string `json:"message"`
}

// DeleteUserResponse 删除用户响应
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// UpdateRoleResponse 更新用户角色响应
type UpdateRoleResponse struct {
	Message string `json:"message"`
}

// RegisterUserResponse 用户注册响应
type RegisterUserResponse struct {
	Message string `json:"message"`
}

// LoginUserResponse 用户登录响应
type LoginUserResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

// UpdateUserRoleResponse 更新用户角色响应
type UpdateUserRoleResponse struct {
	Message string `json:"message"`
}
