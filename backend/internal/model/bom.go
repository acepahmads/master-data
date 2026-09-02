package model

import (
	"time"

	"gorm.io/gorm"
)

// BOM Item Types
const (
	BOMItemTypeComponent   = "COMPONENT"
	BOMItemTypeSubAssembly = "SUB_ASSEMBLY"
)

// BOM Statuses
const (
	BOMStatusDraft    = "Draft"
	BOMStatusActive   = "Active"
	BOMStatusObsolete = "Obsolete"
	BOMStatusArchived = "Archived"
)

// EngineeringBOM represents the engineering Bill of Materials for a specific Hardware Revision.
type EngineeringBOM struct {
	ID                     string         `gorm:"primaryKey;size:64" json:"id"`
	HardwareRevisionID     string         `gorm:"size:64;uniqueIndex;not null" json:"hardwareRevisionId"`
	BOMCode                string         `gorm:"size:64;uniqueIndex;not null" json:"bomCode"`
	Name                   string         `gorm:"size:255;not null" json:"name"`
	Description            string         `gorm:"type:text" json:"description"`
	Status                 string         `gorm:"size:32;default:'Active';index" json:"status"`
	TotalCost              float64        `gorm:"type:decimal(12,4);default:0.0000" json:"totalCost"`
	Currency               string         `gorm:"size:8;default:'USD'" json:"currency"`
	TargetMargin           float64        `gorm:"type:decimal(5,2);default:35.00" json:"targetMargin"` // Percentage (e.g. 35.00%)
	EstimatedSellingPrice  float64        `gorm:"type:decimal(12,4);default:0.0000" json:"estimatedSellingPrice"`
	Notes                  string         `gorm:"type:text" json:"notes"`
	CreatedBy              string         `gorm:"size:255" json:"createdBy"`
	Items                  []BOMItem      `gorm:"foreignKey:EngineeringBOMID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"items,omitempty"`
	CreatedAt              time.Time      `json:"createdAt"`
	UpdatedAt              time.Time      `json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// BOMItem represents an individual line item or sub-assembly within an Engineering BOM.
type BOMItem struct {
	ID                     string         `gorm:"primaryKey;size:64" json:"id"`
	EngineeringBOMID       string         `gorm:"size:64;index;not null" json:"engineeringBomId"`
	ParentItemID           *string        `gorm:"size:64;index" json:"parentItemId"`
	ItemType               string         `gorm:"size:32;index;not null" json:"itemType"` // COMPONENT, SUB_ASSEMBLY
	Name                   string         `gorm:"size:255;not null" json:"name"`
	ComponentID            *string        `gorm:"size:64;index" json:"componentId"`
	ComponentPartNumber    string         `gorm:"size:128" json:"componentPartNumber"`
	ComponentCategory      string         `gorm:"size:100" json:"componentCategory"`
	Quantity               float64        `gorm:"type:decimal(10,4);default:1.0000" json:"quantity"`
	UnitID                 *string        `gorm:"size:64" json:"unitId"`
	UnitCode               string         `gorm:"size:32;default:'PCS'" json:"unitCode"`
	UnitCost               float64        `gorm:"type:decimal(12,4);default:0.0000" json:"unitCost"`
	LineCost               float64        `gorm:"type:decimal(12,4);default:0.0000" json:"lineCost"`
	ReferenceDesignators   string         `gorm:"type:text" json:"referenceDesignators"` // e.g. "R1, R2, R3" or "C101-C104"
	Position               int            `gorm:"default:0" json:"position"`
	Notes                  string         `gorm:"type:text" json:"notes"`
	Component              *Component     `gorm:"foreignKey:ComponentID" json:"component,omitempty"`
	SubItems               []BOMItem      `gorm:"foreignKey:ParentItemID" json:"subItems,omitempty"`
	CreatedAt              time.Time      `json:"createdAt"`
	UpdatedAt              time.Time      `json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// BOMCostSummary provides a breakdown of cost drivers for an Engineering BOM.
type BOMCostSummary struct {
	TotalBOMCost           float64                `json:"totalBomCost"`
	Currency               string                 `json:"currency"`
	TargetMargin           float64                `json:"targetMargin"`
	EstimatedSellingPrice  float64                `json:"estimatedSellingPrice"`
	TotalComponentsCount   int                    `json:"totalComponentsCount"`
	TotalAssembliesCount   int                    `json:"totalAssembliesCount"`
	CostByCategory         map[string]float64     `json:"costByCategory"`
	TopCostDrivers         []TopCostItem          `json:"topCostDrivers"`
}

type TopCostItem struct {
	ItemName               string  `json:"itemName"`
	PartNumber             string  `json:"partNumber"`
	Quantity               float64 `json:"quantity"`
	UnitCost               float64 `json:"unitCost"`
	LineCost               float64 `json:"lineCost"`
	PercentageOfTotal      float64 `json:"percentageOfTotal"`
}

// ProductCostSummary represents the cost rollup for an R&D product.
type ProductCostSummary struct {
	ProductID              string  `json:"productId"`
	ProductCode            string  `json:"productCode"`
	ProductName            string  `json:"productName"`
	ProductType            string  `json:"productType"`
	CurrentVersion         string  `json:"currentVersion"`
	CurrentRevision        string  `json:"currentRevision"`
	RevisionID             string  `json:"revisionId"`
	BOMID                  string  `json:"bomId"`
	BOMCode                string  `json:"bomCode"`
	EstimatedBOMCost       float64 `json:"estimatedBomCost"`
	Currency               string  `json:"currency"`
	TargetMargin           float64 `json:"targetMargin"`
	EstimatedSellingPrice  float64 `json:"estimatedSellingPrice"`
}

// ProjectCostSummary represents the cost aggregation for a Project product.
type ProjectCostSummary struct {
	ProjectID              string  `json:"projectId"`
	ProjectCode            string  `json:"projectCode"`
	ProjectName            string  `json:"projectName"`
	ProductType            string  `json:"productType"`
	TradingProductsCost    float64 `json:"tradingProductsCost"`
	RNDProductsCost        float64 `json:"rndProductsCost"`
	DirectComponentsCost   float64 `json:"directComponentsCost"`
	SubAssembliesCost      float64 `json:"subAssembliesCost"`
	TotalEstimatedCost     float64 `json:"totalEstimatedCost"`
	Currency               string  `json:"currency"`
	TargetMargin           float64 `json:"targetMargin"`
	EstimatedSellingPrice  float64 `json:"estimatedSellingPrice"`
	ItemsCount             int     `json:"itemsCount"`
}
