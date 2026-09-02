package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type RevisionHandler struct {
	service *service.RevisionService
}

func NewRevisionHandler(service *service.RevisionService) *RevisionHandler {
	return &RevisionHandler{service: service}
}

func (h *RevisionHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	versionID := c.Query("version_id")
	if versionID == "" {
		versionID = c.Query("product_version_id")
	}
	productID := c.Query("product_id")
	if productID == "" {
		productID = c.Query("productCode")
	}
	status := c.Query("status")

	revisions, meta, err := h.service.GetAll(p, versionID, productID, status)
	if err != nil {
		ServerError(c, err)
		return
	}

	SuccessWithMeta(c, revisions, meta)
}

func (h *RevisionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	revision, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Hardware revision not found")
		return
	}
	Success(c, revision)
}

func (h *RevisionHandler) Create(c *gin.Context) {
	var revision model.HardwareRevision
	if err := c.ShouldBindJSON(&revision); err != nil {
		BadRequest(c, "Invalid revision payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.service.Create(&revision, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Hardware revision created successfully")
}

func (h *RevisionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var revision model.HardwareRevision
	if err := c.ShouldBindJSON(&revision); err != nil {
		BadRequest(c, "Invalid revision payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.Update(id, &revision, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Hardware revision updated successfully")
}

func (h *RevisionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.service.Delete(id, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id}, "Hardware revision deleted successfully")
}
