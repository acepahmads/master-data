package middleware

import (
	"net/http"
	"strings"

	"iot-rd-backend/internal/config"
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(cfg *config.Config, authService *service.AuthService, tokenService *service.APITokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token format. Expected 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])

		// 1. Check if token is a Static API Token (e.g. iot_live_...)
		if strings.HasPrefix(tokenString, "iot_live_") && tokenService != nil {
			apiToken, err := tokenService.ValidateToken(tokenString)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": err.Error(),
				})
				c.Abort()
				return
			}

			// Grant standard API permissions
			dummyUser := &model.User{
				ID:       "API-" + apiToken.ID,
				Username: apiToken.Name,
				Name:     "API Client (" + apiToken.Name + ")",
				Email:    apiToken.Name + "@api.local",
				RoleID:   "ROLE-SUPERADMIN",
				Status:   "Active",
			}
			allPerms := []string{
				"products.view", "products.create", "products.edit",
				"components.view", "components.create", "components.edit",
				"bom.view", "bom.create", "bom.edit",
				"versions.view", "prototypes.view",
			}
			c.Set("user", dummyUser)
			c.Set("userId", dummyUser.ID)
			c.Set("role", "Super Admin")
			c.Set("permissions", allPerms)
			c.Next()
			return
		}

		// 2. Standard JWT User Authentication
		claims := &service.JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired authentication token",
			})
			c.Abort()
			return
		}

		user, permissions, err := authService.GetMe(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "User session not found or suspended",
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userId", user.ID)
		c.Set("role", claims.Role)
		c.Set("permissions", permissions)

		c.Next()
	}
}
