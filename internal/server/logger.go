package server

import (
	"time"

	"github.com/gin-gonic/gin"

	"test.nahueldev23.com/internal/auth"
	"test.nahueldev23.com/internal/logging"
)

func LoggerMiddleware(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		ctx := c.Request.Context()
		duration := time.Since(start)
		status := c.Writer.Status()

		args := []any{
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		}

		if userID, exists := c.Get(auth.ContextUserID); exists {
			args = append(args, "user_id", userID)
		}

		if sessionID, exists := c.Get(auth.ContextSessionID); exists {
			args = append(args, "session_id", sessionID)
		}

		switch {
		case status >= 500:
			logger.Error(ctx, "http request", args...)

		case status >= 400:
			logger.Warn(ctx, "http request", args...)

		default:
			logger.Info(ctx, "http request", args...)
		}
	}
}
