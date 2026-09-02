package handler

import (
	"net/http"

	"iot-rd-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Success(c *gin.Context, data interface{}, message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func SuccessWithMeta(c *gin.Context, data interface{}, meta repository.PaginationMeta, message ...string) {
	msg := "Success"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
		Meta:    meta,
	})
}

func Created(c *gin.Context, data interface{}, message ...string) {
	msg := "Created successfully"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func BadRequest(c *gin.Context, message string, errors ...interface{}) {
	var errDetails interface{}
	if len(errors) > 0 {
		errDetails = errors[0]
	}
	c.JSON(http.StatusBadRequest, APIResponse{
		Success: false,
		Message: message,
		Errors:  errDetails,
	})
}

func NotFound(c *gin.Context, message ...string) {
	msg := "Resource not found"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	c.JSON(http.StatusNotFound, APIResponse{
		Success: false,
		Message: msg,
	})
}

func ServerError(c *gin.Context, err error) {
	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Message: "An internal server error occurred",
		Errors:  err.Error(),
	})
}
