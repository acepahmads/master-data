package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *service.CategoryService
}

func NewCategoryHandler(service *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetTree(c *gin.Context) {
	tree, err := h.service.GetTree()
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, tree)
}

func (h *CategoryHandler) GetFlat(c *gin.Context) {
	flat, err := h.service.GetFlat()
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, flat)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	cat, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Category not found")
		return
	}
	Success(c, cat)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var cat model.ComponentCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		BadRequest(c, "Invalid category payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.service.Create(&cat, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Category created successfully")
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var cat model.ComponentCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		BadRequest(c, "Invalid category payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.Update(id, &cat, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Category updated successfully")
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"id": id}, "Category deleted successfully")
}
