package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService *service.CustomerService
}

func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (h *CustomerHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	status := c.Query("status")
	source := c.Query("source")
	customers, meta, err := h.customerService.GetAll(p, status, source)
	if err != nil {
		ServerError(c, err)
		return
	}
	SuccessWithMeta(c, customers, meta)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	customer, err := h.customerService.GetByID(id)
	if err != nil {
		NotFound(c, "Customer not found")
		return
	}
	Success(c, customer)
}

func (h *CustomerHandler) Create(c *gin.Context) {
	var payload model.Customer
	if err := c.ShouldBindJSON(&payload); err != nil {
		BadRequest(c, "Invalid customer payload", err.Error())
		return
	}

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	created, err := h.customerService.Create(&payload, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, created)
}

func (h *CustomerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var payload model.Customer
	if err := c.ShouldBindJSON(&payload); err != nil {
		BadRequest(c, "Invalid customer payload", err.Error())
		return
	}
	payload.ID = id

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	updated, err := h.customerService.Update(&payload, currentUser)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, updated)
}

func (h *CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	var currentUser *model.User
	if u, exists := c.Get("user"); exists {
		currentUser = u.(*model.User)
	}

	if err := h.customerService.Delete(id, currentUser); err != nil {
		BadRequest(c, err.Error())
		return
	}
	Success(c, gin.H{"deleted": true})
}
