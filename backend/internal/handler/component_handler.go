package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ComponentHandler struct {
	service *service.ComponentService
}

func NewComponentHandler(service *service.ComponentService) *ComponentHandler {
	return &ComponentHandler{service: service}
}

func (h *ComponentHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	filters := repository.ComponentFilters{
		CategoryID:     c.Query("categoryId"),
		ManufacturerID: c.Query("manufacturerId"),
		SupplierID:     c.Query("supplierId"),
		Status:         c.Query("status"),
		MountingType:   c.Query("mountingType"),
		RoHS:           c.Query("rohs"),
		UnitID:         c.Query("unitId"),
	}

	components, meta, err := h.service.GetAll(p, filters)
	if err != nil {
		ServerError(c, err)
		return
	}

	SuccessWithMeta(c, components, meta)
}

func (h *ComponentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	component, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Component not found")
		return
	}
	Success(c, component)
}

func (h *ComponentHandler) Create(c *gin.Context) {
	var component model.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		BadRequest(c, "Invalid component payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.service.Create(&component, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Component registered successfully")
}

func (h *ComponentHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var component model.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		BadRequest(c, "Invalid component payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.Update(id, &component, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Component updated successfully")
}

func (h *ComponentHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.service.Delete(id, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id}, "Component removed successfully")
}

func (h *ComponentHandler) AddDocument(c *gin.Context) {
	id := c.Param("id")
	var doc model.ComponentDocument
	if err := c.ShouldBindJSON(&doc); err != nil {
		BadRequest(c, "Invalid document payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	createdDoc, err := h.service.AddDocument(id, &doc, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, createdDoc, "Document attached successfully")
}

func (h *ComponentHandler) DeleteDocument(c *gin.Context) {
	docID := c.Param("docId")
	if err := h.service.DeleteDocument(docID); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"id": docID}, "Document removed successfully")
}
