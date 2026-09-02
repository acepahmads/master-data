package model

import (
	"time"

	"gorm.io/gorm"
)

type ProductVersion struct {
	ID                 string             `gorm:"primaryKey;size:64" json:"id"`
	ProductID          string             `gorm:"size:64;index;not null" json:"productId"`
	ProductCode        string             `gorm:"size:64;index" json:"productCode"`
	ProductName        string             `gorm:"size:255" json:"productName"`
	VersionNumber      string             `gorm:"size:64;not null" json:"versionNumber"`
	VersionName        string             `gorm:"size:255;not null" json:"versionName"`
	Status             string             `gorm:"size:32;default:'Draft';index" json:"status"` // Draft, In Review, Active, Deprecated, Archived
	CurrentRevision    string             `gorm:"size:32;default:'REV-A'" json:"currentRevision"`
	RevisionsCount     int                `gorm:"default:0" json:"revisionsCount"`
	Owner              string             `gorm:"size:255" json:"owner"`
	OwnerAvatar        string             `gorm:"size:512" json:"ownerAvatar"`
	EngineeringOwnerID *string            `gorm:"size:64;index" json:"engineeringOwnerId,omitempty"`
	ReleaseDate        *string            `gorm:"size:32" json:"releaseDate"`
	Description        string             `gorm:"type:text" json:"description"`
	ReleaseNotes       string             `gorm:"type:text" json:"releaseNotes"`
	BaseVersionID      *string            `gorm:"size:64;index" json:"baseVersionId"`
	Documents          string             `gorm:"type:text" json:"documents"` // JSON array string
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt     `gorm:"index" json:"-"`

	// Associations
	HardwareRevisions  []HardwareRevision `gorm:"foreignKey:ProductVersionID" json:"hardwareRevisions,omitempty"`
}

type HardwareRevision struct {
	ID                 string         `gorm:"primaryKey;size:64" json:"id"`
	ProductVersionID   string         `gorm:"size:64;index;not null" json:"productVersionId"`
	VersionNumber      string         `gorm:"size:64" json:"versionNumber"`
	ProductCode        string         `gorm:"size:64;index" json:"productCode"`
	ProductName        string         `gorm:"size:255" json:"productName"`
	Code               string         `gorm:"size:32;not null" json:"code"` // REV-A, REV-B, REV-C
	Name               string         `gorm:"size:255;not null" json:"name"`
	Status             string         `gorm:"size:32;default:'Draft';index" json:"status"` // Draft, In Review, Approved, Released, Superseded, Obsolete
	PCBRevision        string         `gorm:"size:64" json:"pcbRevision"`
	SchematicRevision  string         `gorm:"size:64" json:"schematicRevision"`
	Engineer           string         `gorm:"size:255" json:"engineer"`
	ApprovedBy         *string        `gorm:"size:255" json:"approvedBy"`
	EngineeringOwnerID *string        `gorm:"size:64;index" json:"engineeringOwnerId,omitempty"`
	ReleaseDate        *string        `gorm:"size:32" json:"releaseDate"`
	ChangeSummary      string         `gorm:"type:text" json:"changeSummary"`
	ChangesList        string         `gorm:"type:text" json:"changesList"` // JSON array string
	Documents          string          `gorm:"type:text" json:"documents"`   // JSON array string
	EngineeringBOM     *EngineeringBOM `gorm:"foreignKey:HardwareRevisionID" json:"engineeringBom,omitempty"`
	CreatedAt          time.Time       `json:"createdAt"`
	UpdatedAt          time.Time       `json:"updatedAt"`
	DeletedAt          gorm.DeletedAt  `gorm:"index" json:"-"`
}

type ProductLifecycleResponse struct {
	ProductCode string                  `json:"productCode"`
	ProductName string                  `json:"productName"`
	Versions    []ProductVersionSummary `json:"versions"`
}

type ProductVersionSummary struct {
	ID              string                    `json:"id"`
	VersionNumber   string                    `json:"versionNumber"`
	VersionName     string                    `json:"versionName"`
	Status          string                    `json:"status"`
	CurrentRevision string                    `json:"currentRevision"`
	Owner           string                    `json:"owner"`
	ReleaseDate     *string                   `json:"releaseDate"`
	Revisions       []HardwareRevisionSummary `json:"revisions"`
}

type HardwareRevisionSummary struct {
	ID          string  `json:"id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Status      string  `json:"status"`
	PCBRevision string  `json:"pcbRevision"`
	ReleaseDate *string `json:"releaseDate"`
}
