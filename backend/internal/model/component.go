package model

import (
	"time"

	"gorm.io/gorm"
)

type Component struct {
	ID                  string                   `gorm:"primaryKey;size:64" json:"id"`
	PartNumber          string                   `gorm:"size:128;uniqueIndex;not null" json:"partNumber"`
	InternalCode        string                   `gorm:"size:128;uniqueIndex;not null" json:"internalCode"`
	Name                string                   `gorm:"size:255;index;not null" json:"name"`
	CategoryID          string                   `gorm:"size:64;index;not null" json:"categoryId"`
	Category            *ComponentCategory       `gorm:"foreignKey:CategoryID" json:"categoryObj,omitempty"`
	CategoryName        string                   `gorm:"-" json:"category"`
	ManufacturerID      string                   `gorm:"size:64;index;not null" json:"manufacturerId"`
	Manufacturer        *Manufacturer            `gorm:"foreignKey:ManufacturerID" json:"manufacturerObj,omitempty"`
	ManufacturerName    string                   `gorm:"-" json:"manufacturer"`
	SupplierID          string                   `gorm:"size:64;index;not null" json:"supplierId"`
	Supplier            *Supplier                `gorm:"foreignKey:SupplierID" json:"supplierObj,omitempty"`
	SupplierName        string                   `gorm:"-" json:"supplier"`
	UnitID              string                   `gorm:"size:64;index" json:"unitId"`
	Unit                *Unit                    `gorm:"foreignKey:UnitID" json:"unitObj,omitempty"`
	UnitCode            string                   `gorm:"-" json:"unit"`
	PackageType         string                   `gorm:"size:128" json:"packageType"`
	MountingType        string                   `gorm:"size:128" json:"mountingType"`
	LeadTime            string                   `gorm:"size:64" json:"leadTime"`
	RoHSCompliant       bool                     `gorm:"default:true" json:"rohsCompliant"`
	REACHCompliant      bool                     `gorm:"default:true" json:"reachCompliant"`
	EstimatedUnitCost   float64                  `gorm:"type:decimal(12,4);default:0.0000" json:"estimatedUnitCost"`
	Currency            string                   `gorm:"size:8;default:'USD'" json:"currency"`
	CostSource          string                   `gorm:"size:64;default:'MANUAL_ESTIMATE'" json:"costSource"`
	ShortDescription    string                   `gorm:"type:text" json:"shortDescription"`
	DetailedDescription string                   `gorm:"type:text" json:"detailedDescription"`
	Status              string                   `gorm:"size:32;default:'Active';index" json:"status"`
	Image               string                   `gorm:"size:512" json:"image"`
	Specs               []ComponentSpecification `gorm:"foreignKey:ComponentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"specs"`
	Documents           []ComponentDocument      `gorm:"foreignKey:ComponentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"documents"`
	CreatedAt           time.Time                `json:"createdAt"`
	UpdatedAt           time.Time                `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt           `gorm:"index" json:"-"`
}

type ComponentSpecification struct {
	ID                 string    `gorm:"primaryKey;size:64" json:"id"`
	ComponentID        string    `gorm:"size:64;index;not null" json:"componentId"`
	SpecificationName  string    `gorm:"size:128;not null" json:"name"`
	SpecificationValue string    `gorm:"size:255;not null" json:"value"`
	SpecificationUnit  string    `gorm:"size:64" json:"unit"`
	DisplayOrder       int       `gorm:"default:0" json:"displayOrder"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

type ComponentDocument struct {
	ID           string    `gorm:"primaryKey;size:64" json:"id"`
	ComponentID  string    `gorm:"size:64;index;not null" json:"componentId"`
	DocumentType string    `gorm:"size:64;not null" json:"type"`
	FileName     string    `gorm:"size:255;not null" json:"name"`
	FilePath     string    `gorm:"size:512;not null" json:"path"`
	FileSize     string    `gorm:"size:64" json:"size"`
	MimeType     string    `gorm:"size:128" json:"mimeType"`
	UploadedBy   string    `gorm:"size:255" json:"uploadedBy"`
	CreatedAt    time.Time `json:"date"`
}
