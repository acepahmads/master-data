package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PartnerHandler struct {
	mfgService      *service.ManufacturerService
	supplierService *service.SupplierService
}

func NewPartnerHandler(mfgService *service.ManufacturerService, supplierService *service.SupplierService) *PartnerHandler {
	return &PartnerHandler{
		mfgService:      mfgService,
		supplierService: supplierService,
	}
}

// MANUFACTURERS

func (h *PartnerHandler) GetAllManufacturers(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	status := c.Query("status")
	mfgs, meta, err := h.mfgService.GetAll(p, status)
	if err != nil {
		ServerError(c, err)
		return
	}
	SuccessWithMeta(c, mfgs, meta)
}

func (h *PartnerHandler) GetManufacturerByID(c *gin.Context) {
	id := c.Param("id")
	mfg, err := h.mfgService.GetByID(id)
	if err != nil {
		NotFound(c, "Manufacturer not found")
		return
	}
	Success(c, mfg)
}

func (h *PartnerHandler) CreateManufacturer(c *gin.Context) {
	var mfg model.Manufacturer
	if err := c.ShouldBindJSON(&mfg); err != nil {
		BadRequest(c, "Invalid manufacturer payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.mfgService.Create(&mfg, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Manufacturer created successfully")
}

func (h *PartnerHandler) UpdateManufacturer(c *gin.Context) {
	id := c.Param("id")
	var mfg model.Manufacturer
	if err := c.ShouldBindJSON(&mfg); err != nil {
		BadRequest(c, "Invalid manufacturer payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.mfgService.Update(id, &mfg, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Manufacturer updated successfully")
}

func (h *PartnerHandler) DeleteManufacturer(c *gin.Context) {
	id := c.Param("id")
	if err := h.mfgService.Delete(id); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"id": id}, "Manufacturer deleted successfully")
}

// SUPPLIERS

func (h *PartnerHandler) GetAllSuppliers(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	status := c.Query("status")
	suppliers, meta, err := h.supplierService.GetAll(p, status)
	if err != nil {
		ServerError(c, err)
		return
	}
	SuccessWithMeta(c, suppliers, meta)
}

func (h *PartnerHandler) GetSupplierByID(c *gin.Context) {
	id := c.Param("id")
	supplier, err := h.supplierService.GetByID(id)
	if err != nil {
		NotFound(c, "Supplier not found")
		return
	}
	Success(c, supplier)
}

func (h *PartnerHandler) CreateSupplier(c *gin.Context) {
	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		BadRequest(c, "Invalid supplier payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.supplierService.Create(&supplier, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Supplier created successfully")
}

func (h *PartnerHandler) UpdateSupplier(c *gin.Context) {
	id := c.Param("id")
	var supplier model.Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		BadRequest(c, "Invalid supplier payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.supplierService.Update(id, &supplier, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Supplier updated successfully")
}

func (h *PartnerHandler) DeleteSupplier(c *gin.Context) {
	id := c.Param("id")
	if err := h.supplierService.Delete(id); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"id": id}, "Supplier deleted successfully")
}
