package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type BOMHandler struct {
	service *service.BOMService
}

func NewBOMHandler(service *service.BOMService) *BOMHandler {
	return &BOMHandler{service: service}
}

func (h *BOMHandler) GetByRevision(c *gin.Context) {
	revisionID := c.Param("id")
	bom, err := h.service.GetBOMByRevisionID(revisionID)
	if err != nil {
		NotFound(c, "Engineering BOM not found for this hardware revision")
		return
	}
	Success(c, bom)
}

func (h *BOMHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	bom, err := h.service.GetBOMByID(id)
	if err != nil {
		NotFound(c, "Engineering BOM not found")
		return
	}
	Success(c, bom)
}

func (h *BOMHandler) Create(c *gin.Context) {
	revisionID := c.Param("id")
	var bom model.EngineeringBOM
	if err := c.ShouldBindJSON(&bom); err != nil {
		BadRequest(c, "Invalid BOM payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.service.CreateBOM(revisionID, &bom, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Engineering BOM created successfully")
}

func (h *BOMHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var bom model.EngineeringBOM
	if err := c.ShouldBindJSON(&bom); err != nil {
		BadRequest(c, "Invalid BOM payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.UpdateBOM(id, &bom, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Engineering BOM updated successfully")
}

func (h *BOMHandler) AddItem(c *gin.Context) {
	bomID := c.Param("id")
	var item model.BOMItem
	if err := c.ShouldBindJSON(&item); err != nil {
		BadRequest(c, "Invalid BOM item payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.service.AddItem(bomID, &item, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "BOM item added successfully")
}

func (h *BOMHandler) UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item model.BOMItem
	if err := c.ShouldBindJSON(&item); err != nil {
		BadRequest(c, "Invalid BOM item payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.UpdateItem(id, &item, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "BOM item updated successfully")
}

func (h *BOMHandler) DeleteItem(c *gin.Context) {
	id := c.Param("id")
	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.service.DeleteItem(id, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id}, "BOM item removed successfully")
}

func (h *BOMHandler) ReorderItems(c *gin.Context) {
	bomID := c.Param("id")
	var req struct {
		ItemIDs []string `json:"itemIds" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid reorder payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.service.ReorderItems(bomID, req.ItemIDs, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"status": "reordered"}, "BOM items reordered successfully")
}

func (h *BOMHandler) GetCostSummary(c *gin.Context) {
	bomID := c.Param("id")
	summary, err := h.service.GetBOMCostSummary(bomID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, summary)
}

func (h *BOMHandler) GetProductCostSummary(c *gin.Context) {
	productID := c.Param("id")
	summary, err := h.service.GetProductCostSummary(productID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, summary)
}

func (h *BOMHandler) GetProjectCostSummary(c *gin.Context) {
	projectID := c.Param("id")
	summary, err := h.service.GetProjectCostSummary(projectID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, summary)
}

func (h *BOMHandler) GetProductBOM(c *gin.Context) {
	productID := c.Param("id")
	bom, err := h.service.GetOrCreateActiveBOMForProduct(productID)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, bom)
}

func (h *BOMHandler) SaveProductBOM(c *gin.Context) {
	productID := c.Param("id")
	var req service.ProductBOMInput
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid BOM payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	bom, err := h.service.SaveProductBOM(productID, &req, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, bom, "Product BOM updated successfully")
}
