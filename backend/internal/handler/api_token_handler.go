package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type APITokenHandler struct {
	service *service.APITokenService
}

func NewAPITokenHandler(service *service.APITokenService) *APITokenHandler {
	return &APITokenHandler{service: service}
}

type GenerateTokenRequest struct {
	Name          string `json:"name" binding:"required"`
	ExpiresInDays int    `json:"expiresInDays"` // 0 for Unlimited / Never Expire
	Scope         string `json:"scope"`         // ALL_READ, FULL_ACCESS
}

func (h *APITokenHandler) Generate(c *gin.Context) {
	var req GenerateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request payload", err.Error())
		return
	}

	userName := "system"
	if u, exists := c.Get("user"); exists {
		if userObj, ok := u.(*model.User); ok {
			userName = userObj.Username
		}
	}

	res, err := h.service.GenerateToken(req.Name, req.ExpiresInDays, req.Scope, userName)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, res, "API Token generated successfully. Save this token securely as it will not be shown again.")
}

func (h *APITokenHandler) GetAll(c *gin.Context) {
	tokens, err := h.service.GetAll()
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, tokens)
}

func (h *APITokenHandler) Revoke(c *gin.Context) {
	id := c.Param("id")
	userName := "system"
	if u, exists := c.Get("user"); exists {
		if userObj, ok := u.(*model.User); ok {
			userName = userObj.Username
		}
	}

	if err := h.service.Revoke(id, userName); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, nil, "API Token revoked successfully")
}
