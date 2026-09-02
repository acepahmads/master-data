package model

import (
	"time"

	"gorm.io/gorm"
)

// ExchangeRate holds cached and persisted currency exchange rates.
type ExchangeRate struct {
	ID            string         `gorm:"primaryKey;size:64" json:"id"`
	BaseCurrency  string         `gorm:"size:8;index;not null" json:"baseCurrency"`  // e.g. USD
	QuoteCurrency string         `gorm:"size:8;index;not null" json:"quoteCurrency"` // e.g. IDR
	Rate          float64        `gorm:"type:decimal(14,4);not null" json:"rate"`
	Provider      string         `gorm:"size:64;not null" json:"provider"` // e.g. Frankfurter API, Manual Override
	Source        string         `gorm:"size:255" json:"source"`
	FetchedAt     time.Time      `json:"fetchedAt"`
	EffectiveDate string         `gorm:"size:32" json:"effectiveDate"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// PricingCalculationRequest represents an interactive request from the UI to calculate pricing.
type PricingCalculationRequest struct {
	PurchasePrice          float64 `json:"purchasePrice"`
	PurchaseCurrency       string  `json:"purchaseCurrency"`
	ProfitMarginPercentage float64 `json:"profitMarginPercentage"`
	SellingCurrency        string  `json:"sellingCurrency"`
	ExchangeRateMode       string  `json:"exchangeRateMode"`   // AUTO, MANUAL
	ManualExchangeRate     float64 `json:"manualExchangeRate"` // user override if MANUAL
}

// PricingCalculationResult returns the calculated multi-currency pricing breakdown.
type PricingCalculationResult struct {
	PurchasePrice          float64    `json:"purchasePrice"`
	PurchaseCurrency       string     `json:"purchaseCurrency"`
	ProfitMarginPercentage float64    `json:"profitMarginPercentage"`
	BaseSellingPrice       float64    `json:"baseSellingPrice"` // PurchasePrice * (1 + Margin/100) in PurchaseCurrency
	SellingPrice           float64    `json:"sellingPrice"`     // In SellingCurrency
	SellingCurrency        string     `json:"sellingCurrency"`
	ExchangeRateMode       string     `json:"exchangeRateMode"`
	LiveExchangeRate       float64    `json:"liveExchangeRate"`
	ManualExchangeRate     float64    `json:"manualExchangeRate"`
	EffectiveExchangeRate  float64    `json:"effectiveExchangeRate"`
	ExchangeRateProvider   string     `json:"exchangeRateProvider"`
	ExchangeRateUpdatedAt  *time.Time `json:"exchangeRateUpdatedAt"`
	PurchasePriceInIDR     float64    `json:"purchasePriceInIDR"`
	SellingPriceInIDR      float64    `json:"sellingPriceInIDR"`
}
