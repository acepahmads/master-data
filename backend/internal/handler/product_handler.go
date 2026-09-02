package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service        *service.ProductService
	projectService *service.ProjectItemService
}

func NewProductHandler(service *service.ProductService, projectService *service.ProjectItemService) *ProductHandler {
	return &ProductHandler{
		service:        service,
		projectService: projectService,
	}
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	category := c.Query("category")
	status := c.Query("status")
	productType := c.Query("product_type")
	if productType == "" {
		productType = c.Query("productType")
	}

	products, meta, err := h.service.GetAll(p, category, status, productType)
	if err != nil {
		ServerError(c, err)
		return
	}

	SuccessWithMeta(c, products, meta)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Product master not found")
		return
	}
	Success(c, product)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		BadRequest(c, "Invalid product payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.service.Create(&product, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Product master created successfully")
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var product model.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		BadRequest(c, "Invalid product payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.Update(id, &product, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Product master updated successfully")
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.service.Delete(id, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id}, "Product master deleted successfully")
}

// -----------------------------------------------------------------------------
// TRADING PRODUCT DETAIL HANDLERS
// -----------------------------------------------------------------------------

func (h *ProductHandler) GetTradingDetail(c *gin.Context) {
	id := c.Param("id")
	detail, err := h.service.GetTradingDetail(id)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, detail)
}

func (h *ProductHandler) UpdateTradingDetail(c *gin.Context) {
	id := c.Param("id")
	var detail model.TradingProductDetail
	if err := c.ShouldBindJSON(&detail); err != nil {
		BadRequest(c, "Invalid trading detail payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.service.UpdateTradingDetail(id, &detail, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Trading product commercial details updated successfully")
}

// -----------------------------------------------------------------------------
// PROJECT PRODUCT STRUCTURE HANDLERS
// -----------------------------------------------------------------------------

func (h *ProductHandler) GetProjectItems(c *gin.Context) {
	id := c.Param("id")
	format := c.Query("format")
	if format == "tree" {
		tree, err := h.projectService.GetProjectTree(id)
		if err != nil {
			BadRequest(c, err.Error())
			return
		}
		Success(c, tree)
		return
	}

	items, err := h.projectService.GetProjectItems(id)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, items)
}

func (h *ProductHandler) CreateProjectItem(c *gin.Context) {
	projectID := c.Param("id")
	var item model.ProjectItem
	if err := c.ShouldBindJSON(&item); err != nil {
		BadRequest(c, "Invalid project item payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.projectService.CreateItem(projectID, &item, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Project item added successfully")
}

func (h *ProductHandler) UpdateProjectItem(c *gin.Context) {
	id := c.Param("id")
	var item model.ProjectItem
	if err := c.ShouldBindJSON(&item); err != nil {
		BadRequest(c, "Invalid project item payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.projectService.UpdateItem(id, &item, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Project item updated successfully")
}

func (h *ProductHandler) DeleteProjectItem(c *gin.Context) {
	id := c.Param("id")
	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.projectService.DeleteItem(id, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"id": id}, "Project item deleted successfully")
}

func (h *ProductHandler) ReorderProjectItems(c *gin.Context) {
	projectID := c.Param("id")
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

	if err := h.projectService.ReorderItems(projectID, req.ItemIDs, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"status": "reordered"}, "Project items reordered successfully")
}
