package server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"test.nahueldev23.com/internal/logging"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()

		ctx := logging.WithRequestID(
			c.Request.Context(),
			requestID,
		)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
