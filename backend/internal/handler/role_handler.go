package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	service *service.RoleService
}

func NewRoleHandler(service *service.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	roles, err := h.service.GetAll()
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, roles)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	role, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Role not found")
		return
	}
	Success(c, role)
}

func (h *RoleHandler) GetAllPermissions(c *gin.Context) {
	perms, err := h.service.GetAllPermissions()
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, perms)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		BadRequest(c, "Invalid role payload", err.Error())
		return
	}

	created, err := h.service.Create(&role)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Role created successfully")
}

func (h *RoleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		BadRequest(c, "Invalid role payload", err.Error())
		return
	}

	updated, err := h.service.Update(id, &role)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Role updated successfully")
}

type UpdatePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

func (h *RoleHandler) UpdatePermissions(c *gin.Context) {
	id := c.Param("id")
	var req UpdatePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Permissions list is required", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.service.UpdatePermissions(id, req.Permissions, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id, "permissions": req.Permissions}, "Permissions updated successfully")
}
