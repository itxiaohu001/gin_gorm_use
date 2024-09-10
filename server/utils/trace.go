package utils

import (
	"context"
	"github.com/google/uuid"
)

type contextKey string

const (
	TraceIDKey contextKey = "traceID"
	UserIDKey  contextKey = "userID"
)

// NewContextWithTraceID 创建一个包含 traceID 的新 context
func NewContextWithTraceID(ctx context.Context) context.Context {
	traceID := uuid.New().String()
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// GetTraceID 从 context 中获取 traceID
func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	if !ok {
		return "unknown"
	}
	return traceID
}

// SetUserID 将用户ID设置到 context 中
func SetUserID(ctx context.Context, userID uint) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// GetUserID 从 context 中获取用户ID
func GetUserID(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value(UserIDKey).(uint)
	return userID, ok
}
