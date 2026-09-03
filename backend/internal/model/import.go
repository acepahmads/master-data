package model

import (
	"time"
)

// ImportBatchStatus constants
const (
	BatchStatusUploaded          = "UPLOADED"
	BatchStatusParsing           = "PARSING"
	BatchStatusParsed            = "PARSED"
	BatchStatusValidating         = "VALIDATING"
	BatchStatusReadyForReview     = "READY_FOR_REVIEW"
	BatchStatusPartiallyApproved  = "PARTIALLY_APPROVED"
	BatchStatusApproved           = "APPROVED"
	BatchStatusImporting          = "IMPORTING"
	BatchStatusImported           = "IMPORTED"
	BatchStatusFailed             = "FAILED"
	BatchStatusCancelled          = "CANCELLED"
)

// StagedRowValidationStatus constants
const (
	RowValidationPending = "PENDING"
	RowValidationValid   = "VALID"
	RowValidationWarning = "WARNING"
	RowValidationError   = "ERROR"
)

// StagedRowDuplicateStatus constants
const (
	RowDuplicateNewRecord        = "NEW_RECORD"
	RowDuplicatePossibleMatch    = "POSSIBLE_DUPLICATE"
	RowDuplicateMatchedExisting  = "MATCHED_EXISTING"
	RowDuplicateBatchCollision   = "DUPLICATE_IN_BATCH"
)

// StagedRowReviewDecision constants
const (
	RowDecisionPending  = "PENDING"
	RowDecisionApproved = "APPROVED"
	RowDecisionRejected = "REJECTED"
	RowDecisionEdited   = "EDITED"
	RowDecisionSkipped  = "SKIPPED"
)

// StagedRowItemClassification constants
const (
	ItemClassProduct    = "PRODUCT"
	ItemClassComponent  = "COMPONENT"
	ItemClassSparePart  = "SPARE_PART"
	ItemClassAccessory  = "ACCESSORY"
	ItemClassService    = "SERVICE"
	ItemClassAssembly   = "ASSEMBLY"
	ItemClassConsumable = "CONSUMABLE"
	ItemClassOther      = "OTHER"
)

// ImportBatch represents a single multi-file Excel data migration session.
type ImportBatch struct {
	ID              string     `gorm:"primaryKey;size:36" json:"id"`
	BatchCode       string     `gorm:"size:32;uniqueIndex;not null" json:"batchCode"` // e.g. IMP-2026-0001
	SourceYear      string     `gorm:"size:4;index;default:'2026'" json:"sourceYear"`  // 2026 (Active) or 2016-2025 (Archive)
	SourceType      string     `gorm:"size:64;default:'EXCEL_WORKBOOK'" json:"sourceType"`
	SourceReference string     `gorm:"size:255" json:"sourceReference"`
	SourceNotes     string     `gorm:"type:text" json:"sourceNotes"`
	Status          string     `gorm:"size:32;index;default:'UPLOADED'" json:"status"`
	
	FilesCount      int        `gorm:"default:0" json:"filesCount"`
	RowsCount       int        `gorm:"default:0" json:"rowsCount"`
	ValidRows       int        `gorm:"default:0" json:"validRows"`
	WarningRows     int        `gorm:"default:0" json:"warningRows"`
	ErrorRows       int        `gorm:"default:0" json:"errorRows"`
	ApprovedRows    int        `gorm:"default:0" json:"approvedRows"`
	RejectedRows    int        `gorm:"default:0" json:"rejectedRows"`
	ImportedRows    int        `gorm:"default:0" json:"importedRows"`
	
	ColumnMappingJSON string   `gorm:"type:text" json:"columnMappingJson"` // Serialized Header Mapping definition
	
	UploadedBy      string     `gorm:"size:128;not null" json:"uploadedBy"`
	ApprovedBy      string     `gorm:"size:128" json:"approvedBy"`
	ApprovedAt      *time.Time `json:"approvedAt"`
	ImportedAt      *time.Time `json:"importedAt"`
	
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt       time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`

	// Relationships
	Files           []ImportFile      `gorm:"foreignKey:BatchID;constraint:OnDelete:CASCADE" json:"files,omitempty"`
	StagedRows      []ImportStagedRow `gorm:"foreignKey:BatchID;constraint:OnDelete:CASCADE" json:"stagedRows,omitempty"`
}

func (ImportBatch) TableName() string {
	return "imp_batches"
}

// ImportFile tracks every uploaded source workbook within an import batch.
type ImportFile struct {
	ID               string    `gorm:"primaryKey;size:36" json:"id"`
	BatchID          string    `gorm:"size:36;index;not null" json:"batchId"`
	OriginalFileName string    `gorm:"size:255;not null" json:"originalFileName"`
	StoredFileName   string    `gorm:"size:255;not null" json:"storedFileName"`
	StoragePath      string    `gorm:"size:512;not null" json:"storagePath"`
	FileSHA256       string    `gorm:"size:64;not null" json:"fileSha256"`
	FileSizeBytes    int64     `gorm:"not null" json:"fileSizeBytes"`
	MimeType         string    `gorm:"size:128;default:'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'" json:"mimeType"`
	SheetCount       int       `gorm:"default:1" json:"sheetCount"`
	RowCount         int       `gorm:"default:0" json:"rowCount"`
	SheetManifestJSON string   `gorm:"type:text" json:"sheetManifestJson"` // JSON list of sheets + row counts + headers
	Status           string    `gorm:"size:32;default:'UPLOADED'" json:"status"`
	SourceYear       string    `gorm:"size:4;default:'2026'" json:"sourceYear"`
	UploadedBy       string    `gorm:"size:128;not null" json:"uploadedBy"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ImportFile) TableName() string {
	return "imp_files"
}

// ImportStagedRow represents a single staged, normalized row with complete provenance.
type ImportStagedRow struct {
	ID                 string     `gorm:"primaryKey;size:36" json:"id"`
	BatchID            string     `gorm:"size:36;index;not null" json:"batchId"`
	FileID             string     `gorm:"size:36;index;not null" json:"fileId"`
	FileName           string     `gorm:"size:255" json:"fileName"`
	SheetName          string     `gorm:"size:128;not null" json:"sheetName"`
	RowNumber          int        `gorm:"not null" json:"rowNumber"`
	
	// Original Data Snapshot (IMMUTABLE)
	RawDataJSON        string     `gorm:"type:text;not null" json:"rawDataJson"`
	
	// Normalized Data Values
	NormalizedDataJSON string     `gorm:"type:text" json:"normalizedDataJson"`
	NormalizedName     string     `gorm:"size:255;index" json:"normalizedName"`
	NormalizedCode     string     `gorm:"size:64;index" json:"normalizedCode"`
	NormalizedCategory string     `gorm:"size:128" json:"normalizedCategory"`
	NormalizedMfg      string     `gorm:"size:128" json:"normalizedMfg"`
	NormalizedSupplier string     `gorm:"size:128" json:"normalizedSupplier"`
	NormalizedUnit     string     `gorm:"size:32" json:"normalizedUnit"`
	UnitCost           float64    `gorm:"type:decimal(20,4);default:0" json:"unitCost"`
	SellingPrice       float64    `gorm:"type:decimal(20,4);default:0" json:"sellingPrice"`
	Currency           string     `gorm:"size:8;default:'IDR'" json:"currency"`
	TargetMarket       string     `gorm:"size:255" json:"targetMarket"`
	Description        string     `gorm:"type:text" json:"description"`
	
	// Item Classification & Product Archetype
	ItemClassification string     `gorm:"size:32;default:'PRODUCT'" json:"itemClassification"` // PRODUCT, COMPONENT, SERVICE, etc.
	ProductType        string     `gorm:"size:32;default:'TRADING'" json:"productType"`        // TRADING, PROJECT, RND (when PRODUCT)
	TargetEntity       string     `gorm:"size:64;default:'products'" json:"targetEntity"`      // products, components, service_items
	
	// Validation Diagnostics
	ValidationStatus   string     `gorm:"size:32;index;default:'PENDING'" json:"validationStatus"` // VALID, WARNING, ERROR
	ValidationErrorsJSON   string `gorm:"type:text" json:"validationErrorsJson"`
	ValidationWarningsJSON string `gorm:"type:text" json:"validationWarningsJson"`
	
	// Duplicate Diagnostics
	DuplicateStatus         string `gorm:"size:32;default:'NEW_RECORD'" json:"duplicateStatus"` // NEW_RECORD, POSSIBLE_DUPLICATE, MATCHED_EXISTING
	DuplicateCandidatesJSON string `gorm:"type:text" json:"duplicateCandidatesJson"`
	MatchedMasterID         string `gorm:"size:36" json:"matchedMasterId"`
	MatchedMasterType       string `gorm:"size:32" json:"matchedMasterType"`
	
	// AI Co-Pilot Classification Suggestion
	AISuggestionsJSON  string     `gorm:"type:text" json:"aiSuggestionsJson"`
	AISuggestedType    string     `gorm:"size:32" json:"aiSuggestedType"`
	AISuggestedProdType string    `gorm:"size:32" json:"aiSuggestedProdType"`
	AISuggestedCategory string    `gorm:"size:128" json:"aiSuggestedCategory"`
	AIConfidence       float64    `gorm:"type:decimal(5,2);default:0" json:"aiConfidence"`
	
	// Human Review & Decision
	ReviewDecision     string     `gorm:"size:32;index;default:'PENDING'" json:"reviewDecision"` // PENDING, APPROVED, REJECTED, EDITED, SKIPPED
	ReviewNotes        string     `gorm:"type:text" json:"reviewNotes"`
	ReviewedBy         string     `gorm:"size:128" json:"reviewedBy"`
	ReviewedAt         *time.Time `json:"reviewedAt"`
	
	// Commit / Dispatch Target
	FinalStatus        string     `gorm:"size:32;default:'STAGED'" json:"finalStatus"` // STAGED, APPROVED, IMPORTED, FAILED
	CreatedMasterID    string     `gorm:"size:36;index" json:"createdMasterId"`
	TargetMasterTable  string     `gorm:"size:64" json:"targetMasterTable"`
	
	CreatedAt          time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt          time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ImportStagedRow) TableName() string {
	return "imp_staged_rows"
}
