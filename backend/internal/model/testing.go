package model

import (
	"time"

	"gorm.io/gorm"
)

// Test Plan Statuses
const (
	TestPlanStatusActive   = "ACTIVE"
	TestPlanStatusDraft    = "DRAFT"
	TestPlanStatusArchived = "ARCHIVED"
)

// Test Plan Version Statuses
const (
	PlanVersionStatusDraft    = "DRAFT"
	PlanVersionStatusReleased = "RELEASED"
	PlanVersionStatusArchived = "ARCHIVED"
)

// Test Evaluation Types
const (
	TestTypeNumericRange        = "NUMERIC_RANGE"
	TestTypeNumericMin          = "NUMERIC_MIN"
	TestTypeNumericMax          = "NUMERIC_MAX"
	TestTypeBooleanPassFail     = "BOOLEAN_PASS_FAIL"
	TestTypeVisualInspection    = "VISUAL_INSPECTION"
	TestTypeTextLogVerification = "TEXT_LOG_VERIFICATION"
)

// Prototype Test Session Statuses
const (
	SessionStatusDraft      = "DRAFT"
	SessionStatusReady      = "READY"
	SessionStatusInProgress = "IN_PROGRESS"
	SessionStatusCompleted  = "COMPLETED"
)

// Test Execution Statuses
const (
	ExecStatusPending    = "PENDING"
	ExecStatusInProgress = "IN_PROGRESS"
	ExecStatusPassed     = "PASSED"
	ExecStatusFailed     = "FAILED"
	ExecStatusSkipped    = "SKIPPED"
	ExecStatusBlocked    = "BLOCKED"
)

// Engineering Finding Severities
const (
	FindingSeverityCritical = "CRITICAL"
	FindingSeverityHigh     = "HIGH"
	FindingSeverityMedium   = "MEDIUM"
	FindingSeverityLow      = "LOW"
	FindingSeverityInfo     = "INFO"
)

// Finding Dispositions
const (
	FindingDispOpen         = "OPEN"
	FindingDispBenchRework  = "BENCH_REWORK"
	FindingDispDesignChange = "DESIGN_CHANGE_RECOMMENDED"
	FindingDispAcceptedAsIs = "ACCEPTED_AS_IS"
	FindingDispClosed       = "CLOSED"
)

// Change Candidate Statuses
const (
	ChangeCandidateNone         = "NONE"
	ChangeCandidateCandidate    = "CANDIDATE"
	ChangeCandidateDeferred     = "DEFERRED"
	ChangeCandidateAcceptedRisk = "ACCEPTED_RISK"
)

// Recommended Change Scopes
const (
	ChangeScopeNone       = "NONE"
	ChangeScopeSchematic  = "SCHEMATIC_CIRCUIT"
	ChangeScopePCBLayout  = "PCB_LAYOUT"
	ChangeScopeBOM        = "BOM_COMPONENT"
	ChangeScopeMechanical = "MECHANICAL_ENCLOSURE"
	ChangeScopeFirmware   = "FIRMWARE_LOGIC"
)

// Validation Decisions
const (
	ValidationDecisionValidated              = "VALIDATED"
	ValidationDecisionConditionallyValidated = "CONDITIONALLY_VALIDATED"
	ValidationDecisionRejected               = "REJECTED"
	ValidationDecisionPendingReview          = "PENDING_REVIEW"
)

// TestPlan represents a permanent protocol test suite identity.
type TestPlan struct {
	ID               string            `gorm:"primaryKey;size:64" json:"id"`
	Code             string            `gorm:"size:64;not null;index" json:"code"`
	Name             string            `gorm:"size:255;not null" json:"name"`
	Description      string            `gorm:"type:text" json:"description"`
	Category         string            `gorm:"size:64;index" json:"category"`
	ProductID        *string           `gorm:"size:64;index" json:"productId,omitempty"`
	TargetRevisionID *string           `gorm:"size:64;index" json:"targetRevisionId,omitempty"`
	Status           string            `gorm:"size:32;default:'ACTIVE'" json:"status"`
	CreatedBy        string            `gorm:"size:255" json:"createdBy"`
	Versions         []TestPlanVersion `gorm:"foreignKey:TestPlanID;constraint:OnDelete:CASCADE" json:"versions,omitempty"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt    `gorm:"index" json:"-"`
}

// TestPlanVersion represents an immutable or draft version of a Test Plan.
type TestPlanVersion struct {
	ID            string         `gorm:"primaryKey;size:64" json:"id"`
	TestPlanID    string         `gorm:"size:64;not null;index" json:"testPlanId"`
	VersionNumber string         `gorm:"size:32;not null" json:"versionNumber"`
	Status        string         `gorm:"size:32;default:'DRAFT'" json:"status"` // DRAFT, RELEASED, ARCHIVED
	Description   string         `gorm:"type:text" json:"description"`
	ReleaseNotes  string         `gorm:"type:text" json:"releaseNotes"`
	CreatedBy     string         `gorm:"size:255" json:"createdBy"`
	ReleasedBy    string         `gorm:"size:255" json:"releasedBy"`
	ReleasedAt    *time.Time     `json:"releasedAt,omitempty"`
	TestCases     []TestCase     `gorm:"foreignKey:TestPlanVersionID;constraint:OnDelete:CASCADE" json:"testCases,omitempty"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TestCase represents a single parametric or procedure specification in a TestPlanVersion.
type TestCase struct {
	ID                string         `gorm:"primaryKey;size:64" json:"id"`
	TestPlanVersionID string         `gorm:"size:64;not null;index" json:"testPlanVersionId"`
	Sequence          int            `gorm:"default:1;index" json:"sequence"`
	Code              string         `gorm:"size:64;not null" json:"code"`
	Name              string         `gorm:"size:255;not null" json:"name"`
	Category          string         `gorm:"size:64;index" json:"category"`
	Description       string         `gorm:"type:text" json:"description"`
	TestType          string         `gorm:"size:64;not null" json:"testType"` // NUMERIC_RANGE, etc.
	Unit              string         `gorm:"size:32" json:"unit"`
	MinimumValue      *float64       `json:"minimumValue,omitempty"`
	MaximumValue      *float64       `json:"maximumValue,omitempty"`
	TargetValue       *float64       `json:"targetValue,omitempty"`
	ExpectedBoolean   *bool          `json:"expectedBoolean,omitempty"`
	ExpectedText      string         `gorm:"size:255" json:"expectedText,omitempty"`
	Instructions      string         `gorm:"type:text" json:"instructions"`
	AcceptanceNotes   string         `gorm:"type:text" json:"acceptanceNotes,omitempty"`
	Status            string         `gorm:"size:32;default:'ACTIVE'" json:"status"`
	CreatedAt         time.Time      `json:"createdAt"`
	UpdatedAt         time.Time      `json:"updatedAt"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

// PrototypeTestSession represents an instantiated physical testing session on a PrototypeBuild.
type PrototypeTestSession struct {
	ID                   string         `gorm:"primaryKey;size:64" json:"id"`
	PrototypeBuildID     string         `gorm:"size:64;not null;index" json:"prototypeBuildId"`
	TestPlanID           string         `gorm:"size:64;not null;index" json:"testPlanId"`
	TestPlanVersionID    string         `gorm:"size:64;not null;index" json:"testPlanVersionId"`
	SessionCode          string         `gorm:"size:64;not null;index" json:"sessionCode"`
	Name                 string         `gorm:"size:255;not null" json:"name"`
	Description          string         `gorm:"type:text" json:"description"`
	Status               string         `gorm:"size:32;default:'IN_PROGRESS'" json:"status"` // DRAFT, READY, IN_PROGRESS, COMPLETED
	TestPlanSnapshotJSON string         `gorm:"type:longtext;not null" json:"testPlanSnapshotJson"`
	TotalCases           int            `gorm:"default:0" json:"totalCases"`
	PassedCases          int            `gorm:"default:0" json:"passedCases"`
	FailedCases          int            `gorm:"default:0" json:"failedCases"`
	BlockedCases         int            `gorm:"default:0" json:"blockedCases"`
	SkippedCases         int            `gorm:"default:0" json:"skippedCases"`
	PassRatePercentage   float64        `gorm:"default:0.0" json:"passRatePercentage"`
	StartedAt            *time.Time     `json:"startedAt,omitempty"`
	CompletedAt          *time.Time     `json:"completedAt,omitempty"`
	CreatedBy            string         `gorm:"size:255" json:"createdBy"`
	Executions           []TestExecution `gorm:"foreignKey:PrototypeTestSessionID;constraint:OnDelete:CASCADE" json:"executions,omitempty"`
	Decisions            []PrototypeValidationDecision `gorm:"foreignKey:PrototypeTestSessionID;constraint:OnDelete:CASCADE" json:"decisions,omitempty"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

// TestExecution represents an execution iteration of a test case within a session.
type TestExecution struct {
	ID                     string            `gorm:"primaryKey;size:64" json:"id"`
	PrototypeTestSessionID string            `gorm:"size:64;not null;index" json:"prototypeTestSessionId"`
	TestCaseSnapshotID     string            `gorm:"size:64;not null;index" json:"testCaseSnapshotId"`
	TestCaseCode           string            `gorm:"size:64;not null" json:"testCaseCode"`
	TestCaseName           string            `gorm:"size:255;not null" json:"testCaseName"`
	Category               string            `gorm:"size:64;index" json:"category"`
	TestType               string            `gorm:"size:64;not null" json:"testType"`
	RunNumber              int               `gorm:"default:1" json:"runNumber"`
	Status                 string            `gorm:"size:32;default:'PENDING'" json:"status"` // PENDING, IN_PROGRESS, PASSED, FAILED, SKIPPED, BLOCKED
	MeasuredValue          *float64          `json:"measuredValue,omitempty"`
	MeasuredBoolean        *bool             `json:"measuredBoolean,omitempty"`
	ObservedText           string            `gorm:"type:text" json:"observedText,omitempty"`
	Unit                   string            `gorm:"size:32" json:"unit,omitempty"`
	CalculatedResult       string            `gorm:"size:32" json:"calculatedResult,omitempty"` // PASSED, FAILED
	Notes                  string            `gorm:"type:text" json:"notes,omitempty"`
	ExecutedBy             string            `gorm:"size:255" json:"executedBy"`
	StartedAt              *time.Time        `json:"startedAt,omitempty"`
	CompletedAt            *time.Time        `json:"completedAt,omitempty"`
	Measurements           []TestMeasurement `gorm:"foreignKey:TestExecutionID;constraint:OnDelete:CASCADE" json:"measurements,omitempty"`
	Evidences              []TestEvidence    `gorm:"foreignKey:TestExecutionID;constraint:OnDelete:CASCADE" json:"evidences,omitempty"`
	CreatedAt              time.Time         `json:"createdAt"`
	UpdatedAt              time.Time         `json:"updatedAt"`
}

// TestMeasurement represents a discrete quantitative measurement sample.
type TestMeasurement struct {
	ID               string    `gorm:"primaryKey;size:64" json:"id"`
	TestExecutionID  string    `gorm:"size:64;not null;index" json:"testExecutionId"`
	SampleIndex      int       `gorm:"default:1" json:"sampleIndex"`
	MeasuredValue    *float64  `json:"measuredValue,omitempty"`
	MeasuredBoolean  *bool     `json:"measuredBoolean,omitempty"`
	ObservedText     string    `gorm:"type:text" json:"observedText,omitempty"`
	Unit             string    `gorm:"size:32" json:"unit"`
	CalculatedResult string    `gorm:"size:32;not null" json:"calculatedResult"` // PASSED, FAILED
	RecordedBy       string    `gorm:"size:255" json:"recordedBy"`
	RecordedAt       time.Time `json:"recordedAt"`
	CreatedAt        time.Time `json:"createdAt"`
}

// TestEvidence represents an attached waveform, thermal image, or datalog file.
type TestEvidence struct {
	ID              string    `gorm:"primaryKey;size:64" json:"id"`
	TestExecutionID string    `gorm:"size:64;not null;index" json:"testExecutionId"`
	FileName        string    `gorm:"size:255;not null" json:"fileName"`
	FilePath        string    `gorm:"size:512;not null" json:"filePath"`
	FileType        string    `gorm:"size:64" json:"fileType"`
	FileSize        int64     `json:"fileSize"`
	Description     string    `gorm:"size:512" json:"description"`
	UploadedBy      string    `gorm:"size:255" json:"uploadedBy"`
	UploadedAt      time.Time `json:"uploadedAt"`
}

// EngineeringFinding represents a defect, anomaly, or engineering observation.
type EngineeringFinding struct {
	ID                     string         `gorm:"primaryKey;size:64" json:"id"`
	PrototypeBuildID       string         `gorm:"size:64;not null;index" json:"prototypeBuildId"`
	PrototypeTestSessionID *string        `gorm:"size:64;index" json:"prototypeTestSessionId,omitempty"`
	TestExecutionID        *string        `gorm:"size:64;index" json:"testExecutionId,omitempty"`
	Code                   string         `gorm:"size:64;not null;index" json:"code"`
	Title                  string         `gorm:"size:255;not null" json:"title"`
	Description            string         `gorm:"type:text;not null" json:"description"`
	Category               string         `gorm:"size:64;index" json:"category"`
	Severity               string         `gorm:"size:32;not null;index" json:"severity"` // CRITICAL, HIGH, MEDIUM, LOW, INFO
	FindingDisposition     string         `gorm:"size:64;default:'OPEN';index" json:"findingDisposition"` // OPEN, BENCH_REWORK, DESIGN_CHANGE_RECOMMENDED, ACCEPTED_AS_IS, CLOSED
	ChangeCandidateStatus  string         `gorm:"size:64;default:'NONE';index" json:"changeCandidateStatus"` // NONE, CANDIDATE, DEFERRED, ACCEPTED_RISK
	RecommendedChangeScope string         `gorm:"size:64;default:'NONE'" json:"recommendedChangeScope"` // NONE, SCHEMATIC_CIRCUIT, PCB_LAYOUT, BOM_COMPONENT, MECHANICAL_ENCLOSURE, FIRMWARE_LOGIC
	RootCause              string         `gorm:"type:text" json:"rootCause"`
	ContainmentAction      string         `gorm:"type:text" json:"containmentAction"`
	ResolutionNotes        string         `gorm:"type:text" json:"resolutionNotes"`
	Status                 string         `gorm:"size:32;default:'OPEN'" json:"status"`
	ReportedBy             string         `gorm:"size:255" json:"reportedBy"`
	AssignedTo             string         `gorm:"size:255" json:"assignedTo"`
	ClosedBy               string         `gorm:"size:255" json:"closedBy,omitempty"`
	ClosedAt               *time.Time     `json:"closedAt,omitempty"`
	CreatedAt              time.Time      `json:"createdAt"`
	UpdatedAt              time.Time      `json:"updatedAt"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"-"`
}

// PrototypeValidationDecision represents an append-only formal engineering validation sign-off.
type PrototypeValidationDecision struct {
	ID                     string    `gorm:"primaryKey;size:64" json:"id"`
	PrototypeTestSessionID string    `gorm:"size:64;not null;index" json:"prototypeTestSessionId"`
	Decision               string    `gorm:"size:64;not null;index" json:"decision"` // VALIDATED, CONDITIONALLY_VALIDATED, REJECTED, PENDING_REVIEW
	Justification          string    `gorm:"type:text" json:"justification"`
	DecisionSummary        string    `gorm:"type:text" json:"decisionSummary"`
	DecidedBy              string    `gorm:"size:255;not null" json:"decidedBy"`
	IsCurrent              bool      `gorm:"default:true;index" json:"isCurrent"`
	CreatedAt              time.Time `json:"createdAt"`
}

// FrozenTestPlanSnapshot represents the immutable JSON payload frozen when a test session starts.
type FrozenTestPlanSnapshot struct {
	TestPlanID      string                 `json:"testPlanId"`
	TestPlanCode    string                 `json:"testPlanCode"`
	TestPlanName    string                 `json:"testPlanName"`
	Category        string                 `json:"category"`
	VersionID       string                 `json:"versionId"`
	VersionNumber   string                 `json:"versionNumber"`
	TotalTestCases  int                    `json:"totalTestCases"`
	FrozenAt        time.Time              `json:"frozenAt"`
	FrozenBy        string                 `json:"frozenBy"`
	TestCases       []FrozenTestCaseItem   `json:"testCases"`
}

// FrozenTestCaseItem represents an individual test case inside the frozen JSON snapshot.
type FrozenTestCaseItem struct {
	ID              string   `json:"id"`
	Sequence        int      `json:"sequence"`
	Code            string   `json:"code"`
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	TestType        string   `json:"testType"`
	Unit            string   `json:"unit,omitempty"`
	MinimumValue    *float64 `json:"minimumValue,omitempty"`
	MaximumValue    *float64 `json:"maximumValue,omitempty"`
	TargetValue     *float64 `json:"targetValue,omitempty"`
	ExpectedBoolean *bool    `json:"expectedBoolean,omitempty"`
	ExpectedText    string   `json:"expectedText,omitempty"`
	Instructions    string   `json:"instructions"`
}

// TestingDashboardStats holds aggregated KPIs for the Testing Command Center.
type TestingDashboardStats struct {
	TotalTestPlans       int64   `json:"totalTestPlans"`
	ActiveTestSessions   int64   `json:"activeTestSessions"`
	TotalExecutions      int64   `json:"totalExecutions"`
	TotalPassed          int64   `json:"totalPassed"`
	TotalFailed          int64   `json:"totalFailed"`
	TotalBlocked         int64   `json:"totalBlocked"`
	OverallPassRate      float64 `json:"overallPassRate"`
	OpenFindingsCount    int64   `json:"openFindingsCount"`
	CriticalFindingsCount int64   `json:"criticalFindingsCount"`
	HighFindingsCount    int64   `json:"highFindingsCount"`
	ChangeCandidatesCount int64  `json:"changeCandidatesCount"`
	PendingValidations   int64   `json:"pendingValidations"`
	TotalValidatedSessions int64 `json:"totalValidatedSessions"`
	TotalRejectedSessions int64  `json:"totalRejectedSessions"`
}
