package handler

import (
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	service *service.ActivityService
}

func NewActivityHandler(service *service.ActivityService) *ActivityHandler {
	return &ActivityHandler{service: service}
}

func (h *ActivityHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	module := c.Query("module")
	entityType := c.Query("entityType")
	entityID := c.Query("entityId")

	activities, meta, err := h.service.GetAll(p, module, entityType, entityID)
	if err != nil {
		ServerError(c, err)
		return
	}

	SuccessWithMeta(c, activities, meta)
}
