package model

import (
	"time"

	"gorm.io/gorm"
)

// Product types
const (
	ProductTypeTrading     = "TRADING"
	ProductTypeProject     = "PROJECT"
	ProductTypeManufacture = "MANUFACTURE"
	ProductTypeRND         = "MANUFACTURE" // backward-compatibility alias
)

// Project item types
const (
	ProjectItemTypeProduct     = "PRODUCT"
	ProjectItemTypeComponent   = "COMPONENT"
	ProjectItemTypeSubAssembly = "SUB_ASSEMBLY"
)

type Product struct {
	ID                string                `gorm:"primaryKey;size:64" json:"id"`
	Code              string                `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name              string                `gorm:"size:255;index;not null" json:"name"`
	ProductType       string                `gorm:"size:32;default:'MANUFACTURE';index" json:"productType"` // TRADING, PROJECT, MANUFACTURE
	Category          string                `gorm:"size:100;index" json:"category"`
	CategoryID        *string               `gorm:"size:64;index" json:"categoryId"`
	Status            string                `gorm:"size:32;default:'Active';index" json:"status"`
	Description       string                `gorm:"type:text" json:"description"`
	LongDescription   string                `gorm:"type:text" json:"longDescription"`
	Image             string                `gorm:"size:512" json:"image"`
	Images            string                `gorm:"type:text" json:"images"` // JSON array string of up to 5 image URLs
	ProjectLead       string                `gorm:"size:255" json:"projectLead"`
	TargetMarket      string                `gorm:"size:255" json:"targetMarket"`
	PowerRating       string                `gorm:"size:128" json:"powerRating"`
	IngressProtection string                `gorm:"size:128" json:"ingressProtection"`
	OperatingTemp     string                `gorm:"size:128" json:"operatingTemp"`
	Certifications    string                `gorm:"type:text" json:"certifications"` // JSON array string e.g. ["CE","FCC"]
	ComponentsCount     int                   `gorm:"default:0" json:"componentsCount"`
	TradingDetail       *TradingProductDetail `gorm:"foreignKey:ProductID" json:"tradingDetail,omitempty"`
	ProjectItems        []ProjectItem         `gorm:"foreignKey:ProjectProductID" json:"projectItems,omitempty"`
	CostSummary         *PriceSummary         `gorm:"-" json:"costSummary,omitempty"`
	SellingPriceSummary *PriceSummary         `gorm:"-" json:"sellingPriceSummary,omitempty"`
	CreatedAt           time.Time             `json:"createdAt"`
	UpdatedAt           time.Time             `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt        `gorm:"index" json:"-"`
}

// PriceSummary holds normalized display-ready cost and pricing information for the Product List.
type PriceSummary struct {
	Amount         float64  `json:"amount"`
	Currency       string   `json:"currency"`
	AmountInIDR    *float64 `json:"amountInIDR,omitempty"`
	Label          string   `json:"label"`
	Available      bool     `json:"available"`
	Margin         *float64 `json:"margin,omitempty"`
	RateMode       string   `json:"rateMode,omitempty"`
	ExchangeRate   *float64 `json:"exchangeRate,omitempty"`
}

// TradingProductDetail holds commercial and sourcing information for TRADING products.
type TradingProductDetail struct {
	ID                     string         `gorm:"primaryKey;size:64" json:"id"`
	ProductID              string         `gorm:"size:64;uniqueIndex;not null" json:"productId"`
	ManufacturerID         *string        `gorm:"size:64;index" json:"manufacturerId"`
	ManufacturerName       string         `gorm:"size:255" json:"manufacturerName"`
	SupplierID             *string        `gorm:"size:64;index" json:"supplierId"`
	SupplierName           string         `gorm:"size:255" json:"supplierName"`
	ManufacturerPartNumber string         `gorm:"size:128;index" json:"manufacturerPartNumber"`
	SupplierPartNumber     string         `gorm:"size:128" json:"supplierPartNumber"`

	// Multi-Currency & Profit Engine
	PurchasePrice          float64        `gorm:"type:decimal(14,4);default:0.0000" json:"purchasePrice"`
	PurchaseCurrency       string         `gorm:"size:8;default:'USD'" json:"purchaseCurrency"`
	ProfitMarginPercentage float64        `gorm:"type:decimal(6,2);default:25.00" json:"profitMarginPercentage"`
	SellingPrice           float64        `gorm:"type:decimal(14,4);default:0.0000" json:"sellingPrice"`
	SellingCurrency        string         `gorm:"size:8;default:'USD'" json:"sellingCurrency"`

	// Exchange Rate Governance
	ExchangeRateMode       string         `gorm:"size:16;default:'AUTO'" json:"exchangeRateMode"` // AUTO, MANUAL
	LiveExchangeRate       float64        `gorm:"type:decimal(14,4);default:16500.0000" json:"liveExchangeRate"`
	ManualExchangeRate     float64        `gorm:"type:decimal(14,4);default:16500.0000" json:"manualExchangeRate"`
	EffectiveExchangeRate  float64        `gorm:"type:decimal(14,4);default:16500.0000" json:"effectiveExchangeRate"`
	ExchangeRateProvider   string         `gorm:"size:64;default:'Frankfurter API'" json:"exchangeRateProvider"`
	ExchangeRateUpdatedAt  *time.Time     `json:"exchangeRateUpdatedAt"`

	// Converted IDR References
	PurchasePriceInIDR     float64        `gorm:"type:decimal(16,2);default:0.00" json:"purchasePriceInIDR"`
	SellingPriceInIDR      float64        `gorm:"type:decimal(16,2);default:0.00" json:"sellingPriceInIDR"`

	Currency               string         `gorm:"size:8;default:'USD'" json:"currency"` // backwards compatibility
	LeadTimeDays           int            `gorm:"default:0" json:"leadTimeDays"`
	WarrantyInformation    string         `gorm:"size:255" json:"warrantyInformation"`
	Notes                  string         `gorm:"type:text" json:"notes"`
	CreatedAt              time.Time      `json:"createdAt"`
	UpdatedAt              time.Time      `json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// ProjectItem represents a hierarchical composition element within a PROJECT product.
type ProjectItem struct {
	ID               string         `gorm:"primaryKey;size:64" json:"id"`
	ProjectProductID string         `gorm:"size:64;index;not null" json:"projectProductId"`
	ParentItemID     *string        `gorm:"size:64;index" json:"parentItemId"`
	ItemType         string         `gorm:"size:32;index;not null" json:"itemType"` // PRODUCT, COMPONENT, SUB_ASSEMBLY
	Name             string         `gorm:"size:255;not null" json:"name"`
	ProductID        *string        `gorm:"size:64;index" json:"productId"`
	ProductCode      string         `gorm:"size:64" json:"productCode"`
	ProductType      string         `gorm:"size:32" json:"productType"` // referenced product's type (TRADING or RND)
	ComponentID      *string        `gorm:"size:64;index" json:"componentId"`
	ComponentMPN     string         `gorm:"size:128" json:"componentMpn"`
	Quantity         float64        `gorm:"type:decimal(10,2);default:1.00" json:"quantity"`
	UnitID           *string        `gorm:"size:64" json:"unitId"`
	UnitName         string         `gorm:"size:64;default:'pcs'" json:"unitName"`
	SortOrder        int            `gorm:"default:0" json:"sortOrder"`
	Notes            string         `gorm:"type:text" json:"notes"`
	Status           string         `gorm:"size:32;default:'Active'" json:"status"`
	CostPrice        float64        `gorm:"type:decimal(14,2);default:0.00" json:"costPrice"`
	SellingPrice     float64        `gorm:"type:decimal(14,2);default:0.00" json:"sellingPrice"`
	Currency         string         `gorm:"size:8;default:'IDR'" json:"currency"`
	SubItems         []ProjectItem  `gorm:"foreignKey:ParentItemID" json:"subItems,omitempty"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}
