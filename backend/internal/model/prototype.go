package model

import (
	"time"

	"gorm.io/gorm"
)

// Prototype Build Statuses
const (
	PrototypeStatusPlanned    = "PLANNED"
	PrototypeStatusPreparing  = "PREPARING"
	PrototypeStatusInProgress = "IN_PROGRESS"
	PrototypeStatusAssembly   = "ASSEMBLY"
	PrototypeStatusPowerOn    = "POWER_ON"
	PrototypeStatusCompleted  = "COMPLETED"
	PrototypeStatusOnHold     = "ON_HOLD"
	PrototypeStatusCancelled  = "CANCELLED"
)

// Prototype Assembly Stages (1-8)
const (
	AssemblyStagePCBPreparation       = "PCB_PREPARATION"
	AssemblyStageComponentPreparation = "COMPONENT_PREPARATION"
	AssemblyStageSMTAssembly          = "SMT_ASSEMBLY"
	AssemblyStageThroughHoleAssembly  = "THROUGH_HOLE_ASSEMBLY"
	AssemblyStageMechanicalFitting    = "MECHANICAL_FITTING"
	AssemblyStageWiring               = "WIRING"
	AssemblyStageFirmwareFlashing     = "FIRMWARE_FLASHING"
	AssemblyStageInitialPowerOn       = "INITIAL_POWER_ON"
)

// Assembly Stage Statuses
const (
	StageStatusNotStarted = "NOT_STARTED"
	StageStatusInProgress = "IN_PROGRESS"
	StageStatusCompleted  = "COMPLETED"
	StageStatusBlocked    = "BLOCKED"
	StageStatusSkipped    = "SKIPPED"
)

// Component Preparation Statuses
const (
	CompPrepAvailable          = "AVAILABLE"
	CompPrepPartial            = "PARTIAL"
	CompPrepMissing            = "MISSING"
	CompPrepSubstituteRequired = "SUBSTITUTE_REQUIRED"
)

// Engineering Note Categories
const (
	NoteCategoryGeneral           = "GENERAL"
	NoteCategoryAssemblyIssue     = "ASSEMBLY_ISSUE"
	NoteCategoryDesignObservation = "DESIGN_OBSERVATION"
	NoteCategoryKnownIssue        = "KNOWN_ISSUE"
	NoteCategoryRecommendation    = "RECOMMENDATION"
)

// PrototypeBuild represents an instantiated physical prototype build unit run tied to a Hardware Revision.
type PrototypeBuild struct {
	ID                   string                          `gorm:"primaryKey;size:64" json:"id"`
	Code                 string                          `gorm:"size:64;not null" json:"code"`
	Name                 string                          `gorm:"size:255;not null" json:"name"`
	HardwareRevisionID   string                          `gorm:"size:64;index;not null" json:"hardwareRevisionId"`
	ProductID            string                          `gorm:"size:64;index;not null" json:"productId"`
	ProductCode          string                          `gorm:"size:64;index" json:"productCode"`
	ProductName          string                          `gorm:"size:255" json:"productName"`
	VersionID            string                          `gorm:"size:64;index;not null" json:"versionId"`
	VersionNumber        string                          `gorm:"size:64" json:"versionNumber"`
	VersionName          string                          `gorm:"size:255" json:"versionName"`
	HardwareRevisionCode string                          `gorm:"size:64" json:"hardwareRevisionCode"`
	HardwareRevisionName string                          `gorm:"size:255" json:"hardwareRevisionName"`
	BuildNumber          int                             `gorm:"not null" json:"buildNumber"`
	Quantity             int                             `gorm:"default:1;not null" json:"quantity"`
	Purpose              string                          `gorm:"size:128;not null" json:"purpose"`
	Objective            string                          `gorm:"type:text" json:"objective"`
	Status               string                          `gorm:"size:32;default:'PLANNED';index;not null" json:"status"`
	OwnerID              *string                         `gorm:"size:64;index" json:"ownerId,omitempty"`
	EngineeringOwner     string                          `gorm:"size:255" json:"engineeringOwner"`
	StartedAt            *time.Time                      `json:"startedAt,omitempty"`
	TargetCompletionAt   *string                         `gorm:"size:32" json:"targetCompletionAt,omitempty"`
	CompletedAt          *time.Time                      `json:"completedAt,omitempty"`
	Notes                string                          `gorm:"type:text" json:"notes"`
	BOMSnapshotJSON      string                          `gorm:"type:longtext;not null" json:"bomSnapshotJson"`
	BOMSnapshotCost      float64                         `gorm:"type:decimal(12,4);default:0.0000" json:"bomSnapshotCost"`
	BOMSnapshotID        string                          `gorm:"size:64" json:"bomSnapshotId"`
	BOMSnapshotCurrency  string                          `gorm:"size:8;default:'USD'" json:"bomSnapshotCurrency"`
	ReadinessPercentage  float64                         `gorm:"type:decimal(5,2);default:0.00" json:"readinessPercentage"`
	CreatedAt            time.Time                       `json:"createdAt"`
	UpdatedAt            time.Time                       `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt                  `gorm:"index" json:"-"`

	// Associations
	HardwareRevision      *HardwareRevision                `gorm:"foreignKey:HardwareRevisionID" json:"hardwareRevision,omitempty"`
	ComponentPreparations []PrototypeComponentPreparation `gorm:"foreignKey:PrototypeBuildID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"componentPreparations,omitempty"`
	AssemblyStages        []PrototypeAssemblyStage        `gorm:"foreignKey:PrototypeBuildID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"assemblyStages,omitempty"`
	EngineeringNotes      []PrototypeEngineeringNote      `gorm:"foreignKey:PrototypeBuildID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"engineeringNotes,omitempty"`
}

// PrototypeComponentPreparation tracks part availability and kitting readiness for prototype builds.
type PrototypeComponentPreparation struct {
	ID                  string         `gorm:"primaryKey;size:64" json:"id"`
	PrototypeBuildID    string         `gorm:"size:64;index;not null" json:"prototypeBuildId"`
	BOMItemID           *string        `gorm:"size:64;index" json:"bomItemId,omitempty"`
	ComponentID         *string        `gorm:"size:64;index" json:"componentId,omitempty"`
	ComponentName       string         `gorm:"size:255;not null" json:"componentName"`
	ComponentPartNumber string         `gorm:"size:128" json:"componentPartNumber"`
	ComponentCategory   string         `gorm:"size:100" json:"componentCategory"`
	RequiredQuantity    float64        `gorm:"type:decimal(10,4);default:1.0000;not null" json:"requiredQuantity"`
	AvailableQuantity   float64        `gorm:"type:decimal(10,4);default:0.0000;not null" json:"availableQuantity"`
	MissingQuantity     float64        `gorm:"type:decimal(10,4);default:0.0000" json:"missingQuantity"`
	UnitCode            string         `gorm:"size:32;default:'PCS'" json:"unitCode"`
	Status              string         `gorm:"size:32;default:'MISSING';index;not null" json:"status"` // AVAILABLE, PARTIAL, MISSING, SUBSTITUTE_REQUIRED
	Notes               string         `gorm:"type:text" json:"notes"`
	CreatedAt           time.Time      `json:"createdAt"`
	UpdatedAt           time.Time      `json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}

// PrototypeAssemblyStage tracks progress through the 8 hardware assembly stages.
type PrototypeAssemblyStage struct {
	ID               string         `gorm:"primaryKey;size:64" json:"id"`
	PrototypeBuildID string         `gorm:"size:64;index;not null" json:"prototypeBuildId"`
	Stage            string         `gorm:"size:64;index;not null" json:"stage"`
	Sequence         int            `gorm:"not null" json:"sequence"` // 1 through 8
	Name             string         `gorm:"size:255;not null" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	Status           string         `gorm:"size:32;default:'NOT_STARTED';index;not null" json:"status"` // NOT_STARTED, IN_PROGRESS, COMPLETED, BLOCKED, SKIPPED
	StartedAt        *time.Time     `json:"startedAt,omitempty"`
	CompletedAt      *time.Time     `json:"completedAt,omitempty"`
	PerformedBy      string         `gorm:"size:255" json:"performedBy"`
	Notes            string         `gorm:"type:text" json:"notes"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// PrototypeEngineeringNote represents a categorized technical observation logged during build.
type PrototypeEngineeringNote struct {
	ID               string         `gorm:"primaryKey;size:64" json:"id"`
	PrototypeBuildID string         `gorm:"size:64;index;not null" json:"prototypeBuildId"`
	Category         string         `gorm:"size:32;index;not null" json:"category"` // GENERAL, ASSEMBLY_ISSUE, DESIGN_OBSERVATION, KNOWN_ISSUE, RECOMMENDATION
	Title            string         `gorm:"size:255;not null" json:"title"`
	Content          string         `gorm:"type:text;not null" json:"content"`
	CreatedBy        string         `gorm:"size:255;not null" json:"createdBy"`
	UserID           *string        `gorm:"size:64;index" json:"userId,omitempty"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// FrozenBOMSnapshotData represents the unmarshaled immutable BOM structure saved in PrototypeBuild.
type FrozenBOMSnapshotData struct {
	SnapshotID         string             `json:"snapshotId"`
	FrozenAt           time.Time          `json:"frozenAt"`
	HardwareRevisionID string             `json:"hardwareRevisionId"`
	HardwareRevisionCode string           `json:"hardwareRevisionCode"`
	TotalItems         int                `json:"totalItems"`
	TotalQuantity      float64            `json:"totalQuantity"`
	TotalEstimatedCost float64            `json:"totalEstimatedCost"`
	Currency           string             `json:"currency"`
	Items              []FrozenBOMItemDTO `json:"items"`
}

// FrozenBOMItemDTO represents an individual frozen line item in the BOM snapshot.
type FrozenBOMItemDTO struct {
	ID                   string  `json:"id"`
	BOMItemID            string  `json:"bomItemId"`
	ComponentID          *string `json:"componentId,omitempty"`
	ComponentName        string  `json:"componentName"`
	MPN                  string  `json:"mpn"`
	Category             string  `json:"category"`
	Quantity             float64 `json:"quantity"`
	UnitCode             string  `json:"unitCode"`
	UnitCost             float64 `json:"unitCost"`
	LineCost             float64 `json:"lineCost"`
	ReferenceDesignators string  `json:"refDesignators"`
	Position             int     `json:"position"`
	Status               string  `json:"status"`
}

// PrototypeDashboardStats represents high-level metrics for the command center.
type PrototypeDashboardStats struct {
	TotalBuilds        int64 `json:"totalBuilds"`
	PlannedCount       int64 `json:"plannedCount"`
	ComponentPrepCount int64 `json:"componentPrepCount"`
	InAssemblyCount    int64 `json:"inAssemblyCount"`
	PowerOnCount       int64 `json:"powerOnCount"`
	ReadyForTestingCount int64 `json:"readyForTestingCount"`
	CompletedCount     int64 `json:"completedCount"`
	OnHoldCount        int64 `json:"onHoldCount"`
}
