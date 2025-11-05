package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware 使用 slog 记录结构化日志
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 记录日志
		slog.Info("request",
			"method", method,
			"path", path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
