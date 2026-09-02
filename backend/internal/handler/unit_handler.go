package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UnitHandler struct {
	service *service.UnitService
}

func NewUnitHandler(service *service.UnitService) *UnitHandler {
	return &UnitHandler{service: service}
}

func (h *UnitHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	units, meta, err := h.service.GetAll(p)
	if err != nil {
		ServerError(c, err)
		return
	}
	SuccessWithMeta(c, units, meta)
}

func (h *UnitHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	unit, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Unit not found")
		return
	}
	Success(c, unit)
}

func (h *UnitHandler) Create(c *gin.Context) {
	var unit model.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		BadRequest(c, "Invalid unit payload", err.Error())
		return
	}

	created, err := h.service.Create(&unit)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, created, "Unit of measure created successfully")
}

func (h *UnitHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var unit model.Unit
	if err := c.ShouldBindJSON(&unit); err != nil {
		BadRequest(c, "Invalid unit payload", err.Error())
		return
	}

	updated, err := h.service.Update(id, &unit)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, updated, "Unit of measure updated successfully")
}

func (h *UnitHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"id": id}, "Unit of measure deleted successfully")
}
