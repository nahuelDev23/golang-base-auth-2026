package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Service) AuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		header := c.GetHeader("Authorization")

		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "missing authorization header",
			})

			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header",
			})

			c.Abort()
			return
		}

		token := parts[1]

		claims, err := s.ValidateToken(
			c.Request.Context(),
			token,
		)

		if err != nil {

			if errors.Is(err, ErrInvalidSession) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid session",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "invalid token",
				})
			}

			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextSessionID, claims.SessionID)

		c.Next()
	}
}
