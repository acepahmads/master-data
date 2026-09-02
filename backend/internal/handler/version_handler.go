package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type VersionHandler struct {
	service *service.VersionService
}

func NewVersionHandler(service *service.VersionService) *VersionHandler {
	return &VersionHandler{service: service}
}

func (h *VersionHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	productID := c.Query("product_id")
	if productID == "" {
		productID = c.Query("productId")
	}
	status := c.Query("status")
	owner := c.Query("owner")

	versions, meta, err := h.service.GetAll(p, productID, status, owner)
	if err != nil {
		ServerError(c, err)
		return
	}

	SuccessWithMeta(c, versions, meta)
}

func (h *VersionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	version, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Product version not found")
		return
	}
	Success(c, version)
}

func (h *VersionHandler) GetLifecycle(c *gin.Context) {
	productID := c.Param("id")
	lifecycle, err := h.service.GetLifecycle(productID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, lifecycle)
}

func (h *VersionHandler) Compare(c *gin.Context) {
	baseID := c.Query("base_id")
	if baseID == "" {
		baseID = c.Query("baseId")
	}
	targetID := c.Query("target_id")
	if targetID == "" {
		targetID = c.Query("targetId")
	}

	if baseID == "" || targetID == "" {
		BadRequest(c, "Both base_id and target_id are required for version comparison")
		return
	}

	res, err := h.service.Compare(baseID, targetID)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, res)
}

func (h *VersionHandler) Create(c *gin.Context) {
	var version model.ProductVersion
	if err := c.ShouldBindJSON(&version); err != nil {
		BadRequest(c, "Invalid version payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.service.Create(&version, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Product version created successfully")
}

func (h *VersionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var version model.ProductVersion
	if err := c.ShouldBindJSON(&version); err != nil {
		BadRequest(c, "Invalid version payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.Update(id, &version, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Product version updated successfully")
}

func (h *VersionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.service.Delete(id, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id}, "Product version deleted successfully")
}
