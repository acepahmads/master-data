package model

import (
	"time"

	"gorm.io/gorm"
)

// EngineeringChangeRequest represents a formal engineering change request (ECR)
type EngineeringChangeRequest struct {
	ID                       string          `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Code                     string          `gorm:"type:varchar(64);not null;index:idx_ecr_code_del,unique" json:"code"`
	Title                    string          `gorm:"type:varchar(255);not null" json:"title"`
	Origin                   string          `gorm:"type:varchar(64);not null;default:'TEST_FINDING'" json:"origin"` // TEST_FINDING, CUSTOMER_REQUIREMENT, ENGINEERING_REVIEW, REGULATORY, SUPPLIER_CHANGE, OBSOLESCENCE, COST_REDUCTION, RELIABILITY, SECURITY, OTHER
	SourceFindingID          *string         `gorm:"type:varchar(64);index" json:"sourceFindingId"`
	ProductID                string          `gorm:"type:varchar(64);not null;index" json:"productId"`
	ProductVersionID         string          `gorm:"type:varchar(64);not null;index" json:"productVersionId"`
	SourceHardwareRevisionID string          `gorm:"type:varchar(64);not null;index" json:"sourceHardwareRevisionId"`
	ChangeType               string          `gorm:"type:varchar(64);not null;default:'SCHEMATIC_CIRCUIT'" json:"changeType"` // SCHEMATIC_CIRCUIT, PCB_LAYOUT, BOM_COMPONENT, COMPONENT_VALUE, MECHANICAL_ENCLOSURE, FIRMWARE_LOGIC, DOCUMENTATION, OTHER
	Priority                 string          `gorm:"type:varchar(32);not null;default:'MEDIUM'" json:"priority"`               // LOW, MEDIUM, HIGH, CRITICAL
	Description              string          `gorm:"type:text;not null" json:"description"`
	BusinessReason           string          `gorm:"type:text" json:"businessReason"`
	TechnicalJustification   string          `gorm:"type:text" json:"technicalJustification"`
	Status                   string          `gorm:"type:varchar(32);not null;default:'DRAFT'" json:"status"` // DRAFT, SUBMITTED, UNDER_REVIEW, IMPACT_ANALYSIS, PENDING_APPROVAL, APPROVED, REJECTED, CANCELLED, CONVERTED_TO_ECO
	RequestedBy              string          `gorm:"type:varchar(128);not null" json:"requestedBy"`
	AssignedLead             string          `gorm:"type:varchar(128)" json:"assignedLead"`
	EstimatedCost            float64         `gorm:"type:decimal(12,4);not null;default:0.0000" json:"estimatedCost"`
	TargetReleaseDate        *time.Time      `json:"targetReleaseDate"`
	CreatedAt                time.Time       `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt                time.Time       `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt                gorm.DeletedAt  `gorm:"index" json:"deletedAt,omitempty"`

	// Associations
	Product                *Product                `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductVersion         *ProductVersion         `gorm:"foreignKey:ProductVersionID" json:"productVersion,omitempty"`
	SourceHardwareRevision *HardwareRevision       `gorm:"foreignKey:SourceHardwareRevisionID" json:"sourceHardwareRevision,omitempty"`
	SourceFinding          *EngineeringFinding     `gorm:"foreignKey:SourceFindingID" json:"sourceFinding,omitempty"`
	ChangeItems            []ECRChangeItem         `gorm:"foreignKey:ECRID;constraint:OnDelete:CASCADE" json:"changeItems,omitempty"`
	ImpactAnalyses         []ECRImpactAnalysis     `gorm:"foreignKey:ECRID;constraint:OnDelete:CASCADE" json:"impactAnalyses,omitempty"`
	Reviews                []ECRReview             `gorm:"foreignKey:ECRID;constraint:OnDelete:CASCADE" json:"reviews,omitempty"`
	Approvals              []ECRApproval           `gorm:"foreignKey:ECRID;constraint:OnDelete:CASCADE" json:"approvals,omitempty"`
	ECO                    *EngineeringChangeOrder `gorm:"foreignKey:ECRID" json:"eco,omitempty"`
}

func (EngineeringChangeRequest) TableName() string {
	return "engineering_change_requests"
}

// ECRChangeItem represents a discrete proposed engineering delta within an ECR
type ECRChangeItem struct {
	ID                     string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	ECRID                  string         `gorm:"type:varchar(64);not null;index:idx_item_ecr_seq" json:"ecrId"`
	Sequence               int            `gorm:"type:int;not null;default:1;index:idx_item_ecr_seq" json:"sequence"`
	ChangeType             string         `gorm:"type:varchar(64);not null" json:"changeType"` // SCHEMATIC_CIRCUIT, PCB_LAYOUT, BOM_COMPONENT, COMPONENT_VALUE, MECHANICAL_ENCLOSURE, FIRMWARE_LOGIC, DOCUMENTATION, OTHER
	Title                  string         `gorm:"type:varchar(255);not null" json:"title"`
	AffectedReference      string         `gorm:"type:varchar(128)" json:"affectedReference"` // e.g. U3, C18, R12
	AffectedComponentID    *string        `gorm:"type:varchar(64);index" json:"affectedComponentId"`
	CurrentState           string         `gorm:"type:text;not null" json:"currentState"`
	ProposedState          string         `gorm:"type:text;not null" json:"proposedState"`
	TechnicalJustification string         `gorm:"type:text" json:"technicalJustification"`
	RiskLevel              string         `gorm:"type:varchar(32);not null;default:'MEDIUM'" json:"riskLevel"` // LOW, MEDIUM, HIGH, CRITICAL
	Notes                  string         `gorm:"type:text" json:"notes"`
	CreatedAt              time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt              time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Associations
	AffectedComponent *Component `gorm:"foreignKey:AffectedComponentID" json:"affectedComponent,omitempty"`
}

func (ECRChangeItem) TableName() string {
	return "ecr_change_items"
}

// ECRImpactAnalysis represents domain risk and impact assessment for an ECR
type ECRImpactAnalysis struct {
	ID                  string    `gorm:"primaryKey;type:varchar(64)" json:"id"`
	ECRID               string    `gorm:"type:varchar(64);not null;index:idx_impact_ecr_domain" json:"ecrId"`
	Domain              string    `gorm:"type:varchar(64);not null;index:idx_impact_ecr_domain" json:"domain"` // ELECTRICAL, THERMAL, MECHANICAL, FIRMWARE, SOFTWARE, SAFETY_REGULATORY, BOM_COST, SCHEDULE, SCRAP_REWORK
	ImpactLevel         string    `gorm:"type:varchar(32);not null;default:'NOT_AFFECTED'" json:"impactLevel"` // NOT_AFFECTED, LOW, MEDIUM, HIGH, CRITICAL, UNKNOWN
	Description         string    `gorm:"type:text;not null" json:"description"`
	Mitigation          string    `gorm:"type:text" json:"mitigation"`
	CostDelta           float64   `gorm:"type:decimal(12,4);not null;default:0.0000" json:"costDelta"`
	ScheduleImpactWeeks float64   `gorm:"type:decimal(6,2);not null;default:0.00" json:"scheduleImpactWeeks"`
	ScrapDisposition    string    `gorm:"type:varchar(64);not null;default:'RUN_OUT'" json:"scrapDisposition"` // RUN_OUT, SCRAP_ALL, REWORK_EXISTING, USE_AS_IS
	ReviewerID          string    `gorm:"type:varchar(64)" json:"reviewerId"`
	ReviewerName        string    `gorm:"type:varchar(128)" json:"reviewerName"`
	ReviewStatus        string    `gorm:"type:varchar(32);not null;default:'PENDING'" json:"reviewStatus"` // PENDING, REVIEWED
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ECRImpactAnalysis) TableName() string {
	return "ecr_impact_analyses"
}

// ECRReview represents peer technical review for an ECR
type ECRReview struct {
	ID           string     `gorm:"primaryKey;type:varchar(64)" json:"id"`
	ECRID        string     `gorm:"type:varchar(64);not null;index:idx_rev_ecr_disc" json:"ecrId"`
	Discipline   string     `gorm:"type:varchar(64);not null;index:idx_rev_ecr_disc" json:"discipline"` // HARDWARE, FIRMWARE, MECHANICAL, SOFTWARE, QA, PROJECT_MANAGEMENT
	ReviewerID   string     `gorm:"type:varchar(64);not null" json:"reviewerId"`
	ReviewerName string     `gorm:"type:varchar(128);not null" json:"reviewerName"`
	Decision     string     `gorm:"type:varchar(32);not null;default:'PENDING'" json:"decision"` // PENDING, APPROVE, REJECT, REQUEST_CHANGES
	Comments     string     `gorm:"type:text" json:"comments"`
	ReviewedAt   *time.Time `json:"reviewedAt"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ECRReview) TableName() string {
	return "ecr_reviews"
}

// ECRApproval represents formal sign-off record in the CCB approval chain
type ECRApproval struct {
	ID            string     `gorm:"primaryKey;type:varchar(64)" json:"id"`
	ECRID         string     `gorm:"type:varchar(64);not null;index:idx_appr_ecr_stage" json:"ecrId"`
	Stage         string     `gorm:"type:varchar(64);not null;index:idx_appr_ecr_stage" json:"stage"` // TIER_1_TECH_LEAD, TIER_2_ENG_MANAGER, TIER_3_QUALITY_DIRECTOR
	ApproverID    string     `gorm:"type:varchar(64);not null" json:"approverId"`
	ApproverName  string     `gorm:"type:varchar(128);not null" json:"approverName"`
	ApproverRole  string     `gorm:"type:varchar(64);not null" json:"approverRole"`
	Decision      string     `gorm:"type:varchar(32);not null;default:'PENDING'" json:"decision"` // PENDING, APPROVED, REJECTED, REQUEST_CHANGES
	Justification string     `gorm:"type:text" json:"justification"`
	DecidedAt     *time.Time `json:"decidedAt"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (ECRApproval) TableName() string {
	return "ecr_approvals"
}

// EngineeringChangeOrder represents an authorized engineering change implementation order (ECO)
type EngineeringChangeOrder struct {
	ID                         string         `gorm:"primaryKey;type:varchar(64)" json:"id"`
	Code                       string         `gorm:"type:varchar(64);not null;index:idx_eco_code_del,unique" json:"code"`
	ECRID                      string         `gorm:"type:varchar(64);not null;uniqueIndex" json:"ecrId"`
	Title                      string         `gorm:"type:varchar(255);not null" json:"title"`
	SourceHardwareRevisionID   string         `gorm:"type:varchar(64);not null;index" json:"sourceHardwareRevisionId"`
	TargetHardwareRevisionID   *string        `gorm:"type:varchar(64);index" json:"targetHardwareRevisionId"`
	SourceBOMID                *string        `gorm:"type:varchar(64);index" json:"sourceBomId"`
	TargetBOMID                *string        `gorm:"type:varchar(64);index" json:"targetBomId"`
	ImplementationOwner        string         `gorm:"type:varchar(128);not null" json:"implementationOwner"`
	Status                     string         `gorm:"type:varchar(32);not null;default:'DRAFT'" json:"status"` // DRAFT, APPROVED, IN_IMPLEMENTATION, IMPLEMENTED, VERIFICATION_REQUIRED, VERIFIED, CLOSED, CANCELLED
	ImplementationNotes        string         `gorm:"type:text" json:"implementationNotes"`
	ScrapDisposition           string         `gorm:"type:varchar(64);not null;default:'RUN_OUT'" json:"scrapDisposition"` // RUN_OUT, SCRAP_EXISTING, BENCH_REWORK, USE_AS_IS
	RequiresRegressionTesting  bool           `gorm:"not null;default:true" json:"requiresRegressionTesting"`
	ApprovedAt                 *time.Time     `json:"approvedAt"`
	ImplementedAt              *time.Time     `json:"implementedAt"`
	VerifiedAt                 *time.Time     `json:"verifiedAt"`
	ClosedAt                   *time.Time     `json:"closedAt"`
	CreatedBy                  string         `gorm:"type:varchar(128);not null" json:"createdBy"`
	CreatedAt                  time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt                  time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	DeletedAt                  gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Associations
	ECR                    *EngineeringChangeRequest `gorm:"foreignKey:ECRID" json:"ecr,omitempty"`
	SourceHardwareRevision *HardwareRevision         `gorm:"foreignKey:SourceHardwareRevisionID" json:"sourceHardwareRevision,omitempty"`
	TargetHardwareRevision *HardwareRevision         `gorm:"foreignKey:TargetHardwareRevisionID" json:"targetHardwareRevision,omitempty"`
	SourceBOM              *EngineeringBOM           `gorm:"foreignKey:SourceBOMID" json:"sourceBom,omitempty"`
	TargetBOM              *EngineeringBOM           `gorm:"foreignKey:TargetBOMID" json:"targetBom,omitempty"`
}

func (EngineeringChangeOrder) TableName() string {
	return "engineering_change_orders"
}

// BOMChangePreviewItem represents a side-by-side BOM line comparison before and after an ECO
type BOMChangePreviewItem struct {
	ReferenceDesignator string   `json:"referenceDesignator"`
	ComponentID         string   `json:"componentId"`
	ComponentName       string   `json:"componentName"`
	MPN                 string   `json:"mpn"`
	Package             string   `json:"package"`
	SourceQuantity      int      `json:"sourceQuantity"`
	TargetQuantity      int      `json:"targetQuantity"`
	SourceUnitCost      float64  `json:"sourceUnitCost"`
	TargetUnitCost      float64  `json:"targetUnitCost"`
	SourceTotalCost     float64  `json:"sourceTotalCost"`
	TargetTotalCost     float64  `json:"targetTotalCost"`
	CostDelta           float64  `json:"costDelta"`
	ChangeStatus        string   `json:"changeStatus"` // ADDED, REMOVED, MODIFIED, UNCHANGED
	ChangeReason        string   `json:"changeReason"`
}

// BOMChangePreviewSummary provides cost rollups and counts for ECO preview
type BOMChangePreviewSummary struct {
	SourceRevisionCode string                 `json:"sourceRevisionCode"`
	TargetRevisionCode string                 `json:"targetRevisionCode"`
	SourceBOMCost      float64                `json:"sourceBomCost"`
	TargetBOMCost      float64                `json:"targetBomCost"`
	NetCostDelta       float64                `json:"netCostDelta"`
	PercentCostDelta   float64                `json:"percentCostDelta"`
	TotalAdded         int                    `json:"totalAdded"`
	TotalRemoved       int                    `json:"totalRemoved"`
	TotalModified      int                    `json:"totalModified"`
	TotalUnchanged     int                    `json:"totalUnchanged"`
	Items              []BOMChangePreviewItem `json:"items"`
}

// ChangeTraceabilityChain represents the closed-loop lineage from Finding to Validation
type ChangeTraceabilityChain struct {
	FindingID               *string `json:"findingId,omitempty"`
	FindingCode             *string `json:"findingCode,omitempty"`
	FindingTitle            *string `json:"findingTitle,omitempty"`
	ECRID                   string  `json:"ecrId"`
	ECRCode                 string  `json:"ecrCode"`
	ECRTitle                string  `json:"ecrTitle"`
	ECRStatus               string  `json:"ecrStatus"`
	ECOID                   *string `json:"ecoId,omitempty"`
	ECOCode                 *string `json:"ecoCode,omitempty"`
	ECOStatus               *string `json:"ecoStatus,omitempty"`
	SourceRevisionID        string  `json:"sourceRevisionId"`
	SourceRevisionCode      string  `json:"sourceRevisionCode"`
	TargetRevisionID        *string `json:"targetRevisionId,omitempty"`
	TargetRevisionCode      *string `json:"targetRevisionCode,omitempty"`
	TargetBOMID             *string `json:"targetBomId,omitempty"`
	TargetBOMCode           *string `json:"targetBomCode,omitempty"`
	PrototypeID             *string `json:"prototypeId,omitempty"`
	PrototypeCode           *string `json:"prototypeCode,omitempty"`
	TestSessionID           *string `json:"testSessionId,omitempty"`
	TestSessionCode         *string `json:"testSessionCode,omitempty"`
	ValidationDecisionID    *string `json:"validationDecisionId,omitempty"`
	ValidationDecisionState *string `json:"validationDecisionState,omitempty"`
	RegressionBuildRequired bool    `json:"regressionBuildRequired"`
	ClosedLoopVerified      bool    `json:"closedLoopVerified"`
	CCBApprovalCount        int     `json:"ccbApprovalCount"`
}
