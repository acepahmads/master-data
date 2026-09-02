package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	users, meta, err := h.service.GetAll(p)
	if err != nil {
		ServerError(c, err)
		return
	}
	SuccessWithMeta(c, users, meta)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "User not found")
		return
	}
	Success(c, user)
}

type CreateUserRequest struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	RoleID     string `json:"roleId" binding:"required"`
	Department string `json:"department"`
	Password   string `json:"password"`
	Avatar     string `json:"avatar"`
	TwoFactor  bool   `json:"twoFactor"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid user payload", err.Error())
		return
	}

	user := model.User{
		Name:       req.Name,
		Email:      req.Email,
		RoleID:     req.RoleID,
		Department: req.Department,
		Avatar:     req.Avatar,
		TwoFactor:  req.TwoFactor,
	}

	created, err := h.service.Create(&user, req.Password)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "User created successfully")
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		BadRequest(c, "Invalid user payload", err.Error())
		return
	}

	updated, err := h.service.Update(id, &user)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "User updated successfully")
}

func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id := c.Param("id")
	updated, err := h.service.ToggleStatus(id)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, updated, "User status updated successfully")
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id := c.Param("id")
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "New password is required", err.Error())
		return
	}

	if err := h.service.ResetPassword(id, req.Password); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id}, "Password reset successfully")
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"id": id}, "User deleted successfully")
}
