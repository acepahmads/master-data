package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid login credentials format", err.Error())
		return
	}

	res, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, res, "Login successful")
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userVal, exists := c.Get("user")
	if !exists {
		BadRequest(c, "User context not found")
		return
	}

	user := userVal.(*model.User)
	permsVal, _ := c.Get("permissions")
	permissions := permsVal.([]string)

	Success(c, gin.H{
		"user":        user,
		"permissions": permissions,
	})
}
