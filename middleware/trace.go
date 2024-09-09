package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const TraceIDKey = "TraceID"

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set(TraceIDKey, traceID)
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}
