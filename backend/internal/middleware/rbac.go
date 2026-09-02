package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if exists && roleVal.(string) == "Super Admin" {
			c.Next()
			return
		}

		permsVal, exists := c.Get("permissions")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access forbidden: insufficient permissions",
			})
			c.Abort()
			return
		}

		permissions, ok := permsVal.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access forbidden: invalid permission state",
			})
			c.Abort()
			return
		}

		hasPermission := false
		for _, p := range permissions {
			if p == permissionCode || p == "*" {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Access forbidden: missing required permission '" + permissionCode + "'",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
