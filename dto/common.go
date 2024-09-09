package dto

// ErrorResponse 通用错误响应结构
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}