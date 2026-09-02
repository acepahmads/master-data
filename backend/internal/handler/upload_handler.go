package handler

import (
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	service *service.UploadService
}

func NewUploadHandler(service *service.UploadService) *UploadHandler {
	return &UploadHandler{service: service}
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		BadRequest(c, "File is required in form-data", err.Error())
		return
	}

	result, err := h.service.SaveUploadedFile(file)
	if err != nil {
		BadRequest(c, "Failed to upload file", err.Error())
		return
	}

	Created(c, result, "File uploaded successfully")
}
