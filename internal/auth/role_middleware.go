package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"test.nahueldev23.com/internal/user"
)

func RequireRole(allowedRoles ...user.Role) gin.HandlerFunc {

	return func(c *gin.Context) {

		value, exists := c.Get(ContextUserRole)

		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "role not found",
			})
			c.Abort()
			return
		}

		userRole, ok := value.(user.Role)

		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "invalid role",
			})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "forbidden",
		})
		c.Abort()
	}
}
