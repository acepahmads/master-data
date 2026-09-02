package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"
)

type ExchangeRateHandler struct {
	rateService *service.ExchangeRateService
}

func NewExchangeRateHandler(rateService *service.ExchangeRateService) *ExchangeRateHandler {
	return &ExchangeRateHandler{rateService: rateService}
}

// GetLatestRate returns the active exchange rate between base and quote.
func (h *ExchangeRateHandler) GetLatestRate(c *gin.Context) {
	base := c.DefaultQuery("base", "USD")
	quote := c.DefaultQuery("quote", "IDR")

	rate, err := h.rateService.GetRate(base, quote, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve exchange rate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Exchange rate retrieved successfully",
		"data":    rate,
	})
}

// RefreshRate forces a live update of the exchange rate from the external provider.
func (h *ExchangeRateHandler) RefreshRate(c *gin.Context) {
	base := c.DefaultQuery("base", "USD")
	quote := c.DefaultQuery("quote", "IDR")

	rate, err := h.rateService.GetRate(base, quote, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh exchange rate"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Exchange rate refreshed successfully from provider",
		"data":    rate,
	})
}

// CalculatePricing calculates interactive multi-currency pricing given inputs.
func (h *ExchangeRateHandler) CalculatePricing(c *gin.Context) {
	var req model.PricingCalculationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.rateService.CalculatePricing(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Pricing calculated successfully",
		"data":    result,
	})
}
