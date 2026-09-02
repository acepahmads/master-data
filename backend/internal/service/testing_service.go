package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type TestingService interface {
	GetDashboardStats() (*model.TestingDashboardStats, error)

	// Test Plans & Versions
	GetTestPlans(category, status, query string, page, limit int) ([]model.TestPlan, int64, error)
	GetTestPlanByID(id string) (*model.TestPlan, error)
	CreateTestPlan(req *CreateTestPlanRequest, user *model.User) (*model.TestPlan, error)
	UpdateTestPlan(id string, req *UpdateTestPlanRequest, user *model.User) (*model.TestPlan, error)
	DeleteTestPlan(id string, user *model.User) error

	GetTestPlanVersionByID(id string) (*model.TestPlanVersion, error)
	CreateTestPlanVersion(testPlanID string, req *CreateTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error)
	UpdateTestPlanVersion(versionID string, req *UpdateTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error)
	ReleaseTestPlanVersion(versionID string, req *ReleaseTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error)
	CloneTestPlanVersion(versionID string, req *CloneTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error)

	// Test Cases
	CreateTestCase(versionID string, req *CreateTestCaseRequest, user *model.User) (*model.TestCase, error)
	UpdateTestCase(caseID string, req *UpdateTestCaseRequest, user *model.User) (*model.TestCase, error)
	DeleteTestCase(caseID string, user *model.User) error

	// Prototype Test Sessions
	GetTestSessionsByPrototype(prototypeID string) ([]model.PrototypeTestSession, error)
	GetTestSessionByID(id string) (*model.PrototypeTestSession, error)
	CreatePrototypeTestSession(prototypeID string, req *CreateTestSessionRequest, user *model.User) (*model.PrototypeTestSession, error)
	UpdateTestSession(sessionID string, req *UpdateTestSessionRequest, user *model.User) (*model.PrototypeTestSession, error)

	// Test Executions & Measurements
	GetExecutionsBySession(sessionID string) ([]model.TestExecution, error)
	RecordExecutionRun(sessionID string, req *RecordExecutionRequest, user *model.User) (*model.TestExecution, error)
	CreateNextExecutionRun(sessionID, caseSnapshotID string, user *model.User) (*model.TestExecution, error)
	AddEvidence(executionID string, req *AddEvidenceRequest, user *model.User) (*model.TestEvidence, error)

	// Engineering Findings
	GetFindings(prototypeID, sessionID, severity, disposition, changeStatus, query string, page, limit int) ([]model.EngineeringFinding, int64, error)
	GetFindingByID(id string) (*model.EngineeringFinding, error)
	CreateFinding(req *CreateFindingRequest, user *model.User) (*model.EngineeringFinding, error)
	UpdateFinding(id string, req *UpdateFindingRequest, user *model.User) (*model.EngineeringFinding, error)
	DeleteFinding(id string, user *model.User) error

	// Validation Decisions
	SubmitValidationDecision(sessionID string, req *SubmitValidationRequest, user *model.User) (*model.PrototypeValidationDecision, error)
	GetValidationDecisions(sessionID string) ([]model.PrototypeValidationDecision, error)
}

type testingService struct {
	testingRepo repository.TestingRepository
	protoRepo   *repository.PrototypeRepository
	actRepo     *repository.ActivityRepository
}

func NewTestingService(
	testingRepo repository.TestingRepository,
	protoRepo *repository.PrototypeRepository,
	actRepo *repository.ActivityRepository,
) TestingService {
	return &testingService{
		testingRepo: testingRepo,
		protoRepo:   protoRepo,
		actRepo:     actRepo,
	}
}

// -----------------------------------------------------------------------------
// Request DTOs
// -----------------------------------------------------------------------------
type CreateTestPlanRequest struct {
	Code             string  `json:"code"`
	Name             string  `json:"name" binding:"required"`
	Description      string  `json:"description"`
	Category         string  `json:"category"`
	ProductID        *string `json:"productId"`
	TargetRevisionID *string `json:"targetRevisionId"`
	InitialVersion   string  `json:"initialVersion"`
}

type UpdateTestPlanRequest struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Category         string  `json:"category"`
	ProductID        *string `json:"productId"`
	TargetRevisionID *string `json:"targetRevisionId"`
	Status           string  `json:"status"`
}

type CreateTestPlanVersionRequest struct {
	VersionNumber string `json:"versionNumber" binding:"required"`
	Description   string `json:"description"`
	ReleaseNotes  string `json:"releaseNotes"`
}

type UpdateTestPlanVersionRequest struct {
	Description  string `json:"description"`
	ReleaseNotes string `json:"releaseNotes"`
}

type ReleaseTestPlanVersionRequest struct {
	ReleaseNotes string `json:"releaseNotes"`
}

type CloneTestPlanVersionRequest struct {
	NewVersionNumber string `json:"newVersionNumber" binding:"required"`
	Description      string `json:"description"`
}

type CreateTestCaseRequest struct {
	Sequence        int      `json:"sequence"`
	Code            string   `json:"code" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Category        string   `json:"category"`
	Description     string   `json:"description"`
	TestType        string   `json:"testType" binding:"required"`
	Unit            string   `json:"unit"`
	MinimumValue    *float64 `json:"minimumValue"`
	MaximumValue    *float64 `json:"maximumValue"`
	TargetValue     *float64 `json:"targetValue"`
	ExpectedBoolean *bool    `json:"expectedBoolean"`
	ExpectedText    string   `json:"expectedText"`
	Instructions    string   `json:"instructions"`
	AcceptanceNotes string   `json:"acceptanceNotes"`
}

type UpdateTestCaseRequest struct {
	Sequence        int      `json:"sequence"`
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	Description     string   `json:"description"`
	TestType        string   `json:"testType"`
	Unit            string   `json:"unit"`
	MinimumValue    *float64 `json:"minimumValue"`
	MaximumValue    *float64 `json:"maximumValue"`
	TargetValue     *float64 `json:"targetValue"`
	ExpectedBoolean *bool    `json:"expectedBoolean"`
	ExpectedText    string   `json:"expectedText"`
	Instructions    string   `json:"instructions"`
	AcceptanceNotes string   `json:"acceptanceNotes"`
	Status          string   `json:"status"`
}

type CreateTestSessionRequest struct {
	TestPlanID        string `json:"testPlanId" binding:"required"`
	TestPlanVersionID string `json:"testPlanVersionId" binding:"required"`
	Name              string `json:"name" binding:"required"`
	Description       string `json:"description"`
}

type UpdateTestSessionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type RecordExecutionRequest struct {
	TestCaseSnapshotID string   `json:"testCaseSnapshotId" binding:"required"`
	Status             string   `json:"status"` // PASSED, FAILED, SKIPPED, BLOCKED
	MeasuredValue      *float64 `json:"measuredValue"`
	MeasuredBoolean    *bool    `json:"measuredBoolean"`
	ObservedText       string   `json:"observedText"`
	Unit               string   `json:"unit"`
	Notes              string   `json:"notes"`
	CreateFinding      bool     `json:"createFinding"`
	FindingTitle       string   `json:"findingTitle"`
	FindingSeverity    string   `json:"findingSeverity"`
	FindingDescription string   `json:"findingDescription"`
}

type AddEvidenceRequest struct {
	FileName    string `json:"fileName" binding:"required"`
	FilePath    string `json:"filePath" binding:"required"`
	FileType    string `json:"fileType"`
	FileSize    int64  `json:"fileSize"`
	Description string `json:"description"`
}

type CreateFindingRequest struct {
	PrototypeBuildID       string  `json:"prototypeBuildId" binding:"required"`
	PrototypeTestSessionID *string `json:"prototypeTestSessionId"`
	TestExecutionID        *string `json:"testExecutionId"`
	Title                  string  `json:"title" binding:"required"`
	Description            string  `json:"description" binding:"required"`
	Category               string  `json:"category"`
	Severity               string  `json:"severity" binding:"required"` // CRITICAL, HIGH, MEDIUM, LOW, INFO
	FindingDisposition     string  `json:"findingDisposition"`          // OPEN, BENCH_REWORK, DESIGN_CHANGE_RECOMMENDED, ACCEPTED_AS_IS, CLOSED
	ChangeCandidateStatus  string  `json:"changeCandidateStatus"`       // NONE, CANDIDATE, DEFERRED, ACCEPTED_RISK
	RecommendedChangeScope string  `json:"recommendedChangeScope"`      // NONE, SCHEMATIC_CIRCUIT, PCB_LAYOUT, BOM_COMPONENT, MECHANICAL_ENCLOSURE, FIRMWARE_LOGIC
	RootCause              string  `json:"rootCause"`
	ContainmentAction      string  `json:"containmentAction"`
	ResolutionNotes        string  `json:"resolutionNotes"`
	AssignedTo             string  `json:"assignedTo"`
}

type UpdateFindingRequest struct {
	Title                  string  `json:"title"`
	Description            string  `json:"description"`
	Category               string  `json:"category"`
	Severity               string  `json:"severity"`
	FindingDisposition     string  `json:"findingDisposition"`
	ChangeCandidateStatus  string  `json:"changeCandidateStatus"`
	RecommendedChangeScope string  `json:"recommendedChangeScope"`
	RootCause              string  `json:"rootCause"`
	ContainmentAction      string  `json:"containmentAction"`
	ResolutionNotes        string  `json:"resolutionNotes"`
	Status                 string  `json:"status"`
	AssignedTo             string  `json:"assignedTo"`
	ClosedBy               *string `json:"closedBy"`
}

type SubmitValidationRequest struct {
	Decision        string `json:"decision" binding:"required"` // VALIDATED, CONDITIONALLY_VALIDATED, REJECTED, PENDING_REVIEW
	Justification   string `json:"justification"`
	DecisionSummary string `json:"decisionSummary"`
}

func getUserName(u *model.User) string {
	if u != nil && u.Name != "" {
		return u.Name
	}
	return "Alex Chen"
}

// -----------------------------------------------------------------------------
// Testing Dashboard Stats
// -----------------------------------------------------------------------------
func (s *testingService) GetDashboardStats() (*model.TestingDashboardStats, error) {
	return s.testingRepo.GetDashboardStats()
}

// -----------------------------------------------------------------------------
// Master Test Plans
// -----------------------------------------------------------------------------
func (s *testingService) GetTestPlans(category, status, query string, page, limit int) ([]model.TestPlan, int64, error) {
	return s.testingRepo.FindTestPlans(category, status, query, page, limit)
}

func (s *testingService) GetTestPlanByID(id string) (*model.TestPlan, error) {
	return s.testingRepo.GetTestPlanByID(id)
}

func (s *testingService) CreateTestPlan(req *CreateTestPlanRequest, user *model.User) (*model.TestPlan, error) {
	planCode := req.Code
	if planCode == "" {
		planCode = fmt.Sprintf("TP-%s", uuid.New().String()[:8])
	}

	author := getUserName(user)

	tp := &model.TestPlan{
		ID:               fmt.Sprintf("TP-%s", uuid.New().String()[:8]),
		Code:             planCode,
		Name:             req.Name,
		Description:      req.Description,
		Category:         req.Category,
		ProductID:        req.ProductID,
		TargetRevisionID: req.TargetRevisionID,
		Status:           model.TestPlanStatusActive,
		CreatedBy:        author,
	}

	if err := s.testingRepo.CreateTestPlan(tp); err != nil {
		return nil, fmt.Errorf("failed to create test plan: %w", err)
	}

	// Create initial draft version
	verNum := req.InitialVersion
	if verNum == "" {
		verNum = "v1.0"
	}
	ver := &model.TestPlanVersion{
		ID:            fmt.Sprintf("TPV-%s", uuid.New().String()[:8]),
		TestPlanID:    tp.ID,
		VersionNumber: verNum,
		Status:        model.PlanVersionStatusDraft,
		Description:   "Initial draft version",
		CreatedBy:     author,
	}
	_ = s.testingRepo.CreateTestPlanVersion(ver)

	// Log activity
	s.logActivity(user, "test_plans", "TestPlan", tp.ID, tp.Name, "Created Test Plan",
		fmt.Sprintf("Created master test plan '%s' (%s) with initial version %s", tp.Name, tp.Code, verNum), "blue")

	return s.testingRepo.GetTestPlanByID(tp.ID)
}

func (s *testingService) UpdateTestPlan(id string, req *UpdateTestPlanRequest, user *model.User) (*model.TestPlan, error) {
	tp, err := s.testingRepo.GetTestPlanByID(id)
	if err != nil {
		return nil, fmt.Errorf("test plan not found: %w", err)
	}

	if req.Name != "" {
		tp.Name = req.Name
	}
	if req.Description != "" {
		tp.Description = req.Description
	}
	if req.Category != "" {
		tp.Category = req.Category
	}
	if req.ProductID != nil {
		tp.ProductID = req.ProductID
	}
	if req.TargetRevisionID != nil {
		tp.TargetRevisionID = req.TargetRevisionID
	}
	if req.Status != "" {
		tp.Status = req.Status
	}

	if err := s.testingRepo.UpdateTestPlan(tp); err != nil {
		return nil, fmt.Errorf("failed to update test plan: %w", err)
	}

	s.logActivity(user, "test_plans", "TestPlan", tp.ID, tp.Name, "Updated Test Plan",
		fmt.Sprintf("Updated test plan '%s' metadata", tp.Name), "blue")

	return tp, nil
}

func (s *testingService) DeleteTestPlan(id string, user *model.User) error {
	tp, err := s.testingRepo.GetTestPlanByID(id)
	if err != nil {
		return fmt.Errorf("test plan not found: %w", err)
	}

	if err := s.testingRepo.DeleteTestPlan(tp.ID); err != nil {
		return fmt.Errorf("failed to delete test plan: %w", err)
	}

	s.logActivity(user, "test_plans", "TestPlan", tp.ID, tp.Name, "Deleted Test Plan",
		fmt.Sprintf("Deleted test plan '%s' (%s)", tp.Name, tp.Code), "red")

	return nil
}

// -----------------------------------------------------------------------------
// Test Plan Versions & Immutability
// -----------------------------------------------------------------------------
func (s *testingService) GetTestPlanVersionByID(id string) (*model.TestPlanVersion, error) {
	return s.testingRepo.GetTestPlanVersionByID(id)
}

func (s *testingService) CreateTestPlanVersion(testPlanID string, req *CreateTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error) {
	tp, err := s.testingRepo.GetTestPlanByID(testPlanID)
	if err != nil {
		return nil, fmt.Errorf("test plan not found: %w", err)
	}

	author := getUserName(user)

	ver := &model.TestPlanVersion{
		ID:            fmt.Sprintf("TPV-%s", uuid.New().String()[:8]),
		TestPlanID:    tp.ID,
		VersionNumber: req.VersionNumber,
		Status:        model.PlanVersionStatusDraft,
		Description:   req.Description,
		ReleaseNotes:  req.ReleaseNotes,
		CreatedBy:     author,
	}

	if err := s.testingRepo.CreateTestPlanVersion(ver); err != nil {
		return nil, fmt.Errorf("failed to create test plan version: %w", err)
	}

	s.logActivity(user, "test_plans", "TestPlanVersion", ver.ID, ver.VersionNumber, "Created Version",
		fmt.Sprintf("Created draft version %s for test plan '%s'", ver.VersionNumber, tp.Name), "blue")

	return ver, nil
}

func (s *testingService) UpdateTestPlanVersion(versionID string, req *UpdateTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error) {
	ver, err := s.testingRepo.GetTestPlanVersionByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	if ver.Status == model.PlanVersionStatusReleased {
		return nil, errors.New("cannot edit a RELEASED test plan version; clone to draft to make changes")
	}

	if req.Description != "" {
		ver.Description = req.Description
	}
	if req.ReleaseNotes != "" {
		ver.ReleaseNotes = req.ReleaseNotes
	}

	if err := s.testingRepo.UpdateTestPlanVersion(ver); err != nil {
		return nil, fmt.Errorf("failed to update version: %w", err)
	}
	return ver, nil
}

func (s *testingService) ReleaseTestPlanVersion(versionID string, req *ReleaseTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error) {
	ver, err := s.testingRepo.GetTestPlanVersionByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	if ver.Status == model.PlanVersionStatusReleased {
		return ver, nil // Already released
	}

	if len(ver.TestCases) == 0 {
		return nil, errors.New("cannot release a test plan version with 0 test cases; add at least one test case before releasing")
	}

	now := time.Now()
	author := getUserName(user)
	ver.Status = model.PlanVersionStatusReleased
	ver.ReleasedBy = author
	ver.ReleasedAt = &now
	if req.ReleaseNotes != "" {
		ver.ReleaseNotes = req.ReleaseNotes
	}

	if err := s.testingRepo.UpdateTestPlanVersion(ver); err != nil {
		return nil, fmt.Errorf("failed to release version: %w", err)
	}

	s.logActivity(user, "test_plans", "TestPlanVersion", ver.ID, ver.VersionNumber, "Released Version",
		fmt.Sprintf("Formally released test plan version %s (%d test cases)", ver.VersionNumber, len(ver.TestCases)), "emerald")

	return ver, nil
}

func (s *testingService) CloneTestPlanVersion(versionID string, req *CloneTestPlanVersionRequest, user *model.User) (*model.TestPlanVersion, error) {
	srcVer, err := s.testingRepo.GetTestPlanVersionByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("source version not found: %w", err)
	}

	author := getUserName(user)

	newVer := &model.TestPlanVersion{
		ID:            fmt.Sprintf("TPV-%s", uuid.New().String()[:8]),
		TestPlanID:    srcVer.TestPlanID,
		VersionNumber: req.NewVersionNumber,
		Status:        model.PlanVersionStatusDraft,
		Description:   req.Description,
		ReleaseNotes:  fmt.Sprintf("Cloned from %s", srcVer.VersionNumber),
		CreatedBy:     author,
	}

	if err := s.testingRepo.CreateTestPlanVersion(newVer); err != nil {
		return nil, fmt.Errorf("failed to clone version: %w", err)
	}

	// Copy all test cases
	for _, tc := range srcVer.TestCases {
		newCase := model.TestCase{
			ID:                fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
			TestPlanVersionID: newVer.ID,
			Sequence:          tc.Sequence,
			Code:              tc.Code,
			Name:              tc.Name,
			Category:          tc.Category,
			Description:       tc.Description,
			TestType:          tc.TestType,
			Unit:              tc.Unit,
			MinimumValue:      tc.MinimumValue,
			MaximumValue:      tc.MaximumValue,
			TargetValue:       tc.TargetValue,
			ExpectedBoolean:   tc.ExpectedBoolean,
			ExpectedText:      tc.ExpectedText,
			Instructions:      tc.Instructions,
			AcceptanceNotes:   tc.AcceptanceNotes,
			Status:            "ACTIVE",
		}
		_ = s.testingRepo.CreateTestCase(&newCase)
	}

	s.logActivity(user, "test_plans", "TestPlanVersion", newVer.ID, newVer.VersionNumber, "Cloned Version",
		fmt.Sprintf("Cloned version %s from %s (%d test cases copied)", newVer.VersionNumber, srcVer.VersionNumber, len(srcVer.TestCases)), "indigo")

	return s.testingRepo.GetTestPlanVersionByID(newVer.ID)
}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------
func (s *testingService) CreateTestCase(versionID string, req *CreateTestCaseRequest, user *model.User) (*model.TestCase, error) {
	ver, err := s.testingRepo.GetTestPlanVersionByID(versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	if ver.Status == model.PlanVersionStatusReleased {
		return nil, errors.New("cannot add test cases to an immutable RELEASED version; clone to draft first")
	}

	seq := req.Sequence
	if seq <= 0 {
		seq = len(ver.TestCases) + 1
	}

	tc := &model.TestCase{
		ID:                fmt.Sprintf("TC-%s", uuid.New().String()[:8]),
		TestPlanVersionID: ver.ID,
		Sequence:          seq,
		Code:              req.Code,
		Name:              req.Name,
		Category:          req.Category,
		Description:       req.Description,
		TestType:          req.TestType,
		Unit:              req.Unit,
		MinimumValue:      req.MinimumValue,
		MaximumValue:      req.MaximumValue,
		TargetValue:       req.TargetValue,
		ExpectedBoolean:   req.ExpectedBoolean,
		ExpectedText:      req.ExpectedText,
		Instructions:      req.Instructions,
		AcceptanceNotes:   req.AcceptanceNotes,
		Status:            "ACTIVE",
	}

	if err := s.testingRepo.CreateTestCase(tc); err != nil {
		return nil, fmt.Errorf("failed to create test case: %w", err)
	}

	s.logActivity(user, "test_plans", "TestCase", tc.ID, tc.Name, "Added Test Case",
		fmt.Sprintf("Added test case '%s' (%s) to version %s", tc.Name, tc.Code, ver.VersionNumber), "blue")

	return tc, nil
}

func (s *testingService) UpdateTestCase(caseID string, req *UpdateTestCaseRequest, user *model.User) (*model.TestCase, error) {
	tc, err := s.testingRepo.GetTestCaseByID(caseID)
	if err != nil {
		return nil, fmt.Errorf("test case not found: %w", err)
	}

	ver, err := s.testingRepo.GetTestPlanVersionByID(tc.TestPlanVersionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	if ver.Status == model.PlanVersionStatusReleased {
		return nil, errors.New("cannot modify test cases on an immutable RELEASED version; clone to draft first")
	}

	if req.Sequence > 0 {
		tc.Sequence = req.Sequence
	}
	if req.Name != "" {
		tc.Name = req.Name
	}
	if req.Category != "" {
		tc.Category = req.Category
	}
	if req.Description != "" {
		tc.Description = req.Description
	}
	if req.TestType != "" {
		tc.TestType = req.TestType
	}
	if req.Unit != "" {
		tc.Unit = req.Unit
	}
	if req.MinimumValue != nil {
		tc.MinimumValue = req.MinimumValue
	}
	if req.MaximumValue != nil {
		tc.MaximumValue = req.MaximumValue
	}
	if req.TargetValue != nil {
		tc.TargetValue = req.TargetValue
	}
	if req.ExpectedBoolean != nil {
		tc.ExpectedBoolean = req.ExpectedBoolean
	}
	if req.ExpectedText != "" {
		tc.ExpectedText = req.ExpectedText
	}
	if req.Instructions != "" {
		tc.Instructions = req.Instructions
	}
	if req.AcceptanceNotes != "" {
		tc.AcceptanceNotes = req.AcceptanceNotes
	}
	if req.Status != "" {
		tc.Status = req.Status
	}

	if err := s.testingRepo.UpdateTestCase(tc); err != nil {
		return nil, fmt.Errorf("failed to update test case: %w", err)
	}
	return tc, nil
}

func (s *testingService) DeleteTestCase(caseID string, user *model.User) error {
	tc, err := s.testingRepo.GetTestCaseByID(caseID)
	if err != nil {
		return fmt.Errorf("test case not found: %w", err)
	}

	ver, err := s.testingRepo.GetTestPlanVersionByID(tc.TestPlanVersionID)
	if err != nil {
		return fmt.Errorf("version not found: %w", err)
	}

	if ver.Status == model.PlanVersionStatusReleased {
		return errors.New("cannot delete test cases on an immutable RELEASED version; clone to draft first")
	}

	return s.testingRepo.DeleteTestCase(tc.ID)
}

// -----------------------------------------------------------------------------
// Prototype Test Sessions & Frozen Snapshots
// -----------------------------------------------------------------------------
func (s *testingService) GetTestSessionsByPrototype(prototypeID string) ([]model.PrototypeTestSession, error) {
	return s.testingRepo.FindTestSessionsByPrototype(prototypeID)
}

func (s *testingService) GetTestSessionByID(id string) (*model.PrototypeTestSession, error) {
	return s.testingRepo.GetTestSessionByID(id)
}

func (s *testingService) CreatePrototypeTestSession(prototypeID string, req *CreateTestSessionRequest, user *model.User) (*model.PrototypeTestSession, error) {
	// 1. Verify prototype build exists
	proto, err := s.protoRepo.FindByID(prototypeID)
	if err != nil {
		return nil, fmt.Errorf("prototype build not found: %w", err)
	}

	// 2. Fetch released test plan version
	ver, err := s.testingRepo.GetTestPlanVersionByID(req.TestPlanVersionID)
	if err != nil {
		return nil, fmt.Errorf("test plan version not found: %w", err)
	}

	tp, err := s.testingRepo.GetTestPlanByID(req.TestPlanID)
	if err != nil {
		return nil, fmt.Errorf("test plan not found: %w", err)
	}

	author := getUserName(user)

	// 3. Build Immutable Frozen Test Plan Snapshot JSON
	snapshotItems := make([]model.FrozenTestCaseItem, 0, len(ver.TestCases))
	for _, c := range ver.TestCases {
		snapshotItems = append(snapshotItems, model.FrozenTestCaseItem{
			ID:              c.ID,
			Sequence:        c.Sequence,
			Code:            c.Code,
			Name:            c.Name,
			Category:        c.Category,
			TestType:        c.TestType,
			Unit:            c.Unit,
			MinimumValue:    c.MinimumValue,
			MaximumValue:    c.MaximumValue,
			TargetValue:     c.TargetValue,
			ExpectedBoolean: c.ExpectedBoolean,
			ExpectedText:    c.ExpectedText,
			Instructions:    c.Instructions,
		})
	}

	now := time.Now()
	snapObj := model.FrozenTestPlanSnapshot{
		TestPlanID:     tp.ID,
		TestPlanCode:   tp.Code,
		TestPlanName:   tp.Name,
		Category:       tp.Category,
		VersionID:      ver.ID,
		VersionNumber:  ver.VersionNumber,
		TotalTestCases: len(snapshotItems),
		FrozenAt:       now,
		FrozenBy:       author,
		TestCases:      snapshotItems,
	}
	snapBytes, _ := json.Marshal(snapObj)

	sessionCode := fmt.Sprintf("SESS-%s-%s", proto.Code, uuid.New().String()[:6])

	session := &model.PrototypeTestSession{
		ID:                   fmt.Sprintf("SESS-%s", uuid.New().String()[:8]),
		PrototypeBuildID:     proto.ID,
		TestPlanID:           tp.ID,
		TestPlanVersionID:    ver.ID,
		SessionCode:          sessionCode,
		Name:                 req.Name,
		Description:          req.Description,
		Status:               model.SessionStatusInProgress,
		TestPlanSnapshotJSON: string(snapBytes),
		TotalCases:           len(snapshotItems),
		PassedCases:          0,
		FailedCases:          0,
		BlockedCases:         0,
		SkippedCases:         0,
		PassRatePercentage:   0.0,
		StartedAt:            &now,
		CreatedBy:            author,
	}

	if err := s.testingRepo.CreateTestSession(session); err != nil {
		return nil, fmt.Errorf("failed to create test session: %w", err)
	}

	// 4. Initialize Test Executions for each case in the snapshot (Run #1, Status PENDING)
	for _, tc := range snapshotItems {
		exec := model.TestExecution{
			ID:                     fmt.Sprintf("EXEC-%s", uuid.New().String()[:8]),
			PrototypeTestSessionID: session.ID,
			TestCaseSnapshotID:     tc.ID,
			TestCaseCode:           tc.Code,
			TestCaseName:           tc.Name,
			Category:               tc.Category,
			TestType:               tc.TestType,
			RunNumber:              1,
			Status:                 model.ExecStatusPending,
			Unit:                   tc.Unit,
			ExecutedBy:             author,
		}
		_ = s.testingRepo.CreateTestExecution(&exec)
	}

	s.logActivity(user, "testing", "PrototypeTestSession", session.ID, session.Name, "Initiated Test Session",
		fmt.Sprintf("Initiated test session '%s' on %s with frozen test snapshot (%d test cases)", session.Name, proto.Code, len(snapshotItems)), "indigo")

	return s.testingRepo.GetTestSessionByID(session.ID)
}

func (s *testingService) UpdateTestSession(sessionID string, req *UpdateTestSessionRequest, user *model.User) (*model.PrototypeTestSession, error) {
	session, err := s.testingRepo.GetTestSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	if req.Name != "" {
		session.Name = req.Name
	}
	if req.Description != "" {
		session.Description = req.Description
	}
	if req.Status != "" {
		session.Status = req.Status
		if req.Status == model.SessionStatusCompleted && session.CompletedAt == nil {
			now := time.Now()
			session.CompletedAt = &now
		}
	}

	if err := s.testingRepo.UpdateTestSession(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}
	return session, nil
}

// -----------------------------------------------------------------------------
// Test Executions & Server-Side Evaluation Engine
// -----------------------------------------------------------------------------
func (s *testingService) GetExecutionsBySession(sessionID string) ([]model.TestExecution, error) {
	return s.testingRepo.FindExecutionsBySession(sessionID)
}

func (s *testingService) RecordExecutionRun(sessionID string, req *RecordExecutionRequest, user *model.User) (*model.TestExecution, error) {
	session, err := s.testingRepo.GetTestSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	author := getUserName(user)

	// Parse snapshot to get test case criteria
	var snapshot model.FrozenTestPlanSnapshot
	_ = json.Unmarshal([]byte(session.TestPlanSnapshotJSON), &snapshot)

	var matchedCase *model.FrozenTestCaseItem
	for i := range snapshot.TestCases {
		if snapshot.TestCases[i].ID == req.TestCaseSnapshotID {
			matchedCase = &snapshot.TestCases[i]
			break
		}
	}
	if matchedCase == nil {
		return nil, errors.New("test case snapshot item not found in session snapshot")
	}

	// Server-Side Tolerance & Evaluation Calculation
	calculatedResult := s.evaluateTestCase(matchedCase, req.MeasuredValue, req.MeasuredBoolean, req.ObservedText)

	// If client provided explicit SKIPPED or BLOCKED, preserve that status
	execStatus := req.Status
	if execStatus == "" || execStatus == model.ExecStatusInProgress {
		if calculatedResult == "PASSED" {
			execStatus = model.ExecStatusPassed
		} else if calculatedResult == "FAILED" {
			execStatus = model.ExecStatusFailed
		} else {
			execStatus = model.ExecStatusInProgress
		}
	}

	// Find the latest execution for this test case in the session
	var targetExec *model.TestExecution
	for i := range session.Executions {
		e := &session.Executions[i]
		if e.TestCaseSnapshotID == req.TestCaseSnapshotID {
			if targetExec == nil || e.RunNumber > targetExec.RunNumber {
				targetExec = e
			}
		}
	}

	now := time.Now()

	if targetExec == nil || targetExec.Status != model.ExecStatusPending {
		// Create new execution run (Append-only)
		nextRun := 1
		if targetExec != nil {
			nextRun = targetExec.RunNumber + 1
		}
		newExec := model.TestExecution{
			ID:                     fmt.Sprintf("EXEC-%s", uuid.New().String()[:8]),
			PrototypeTestSessionID: session.ID,
			TestCaseSnapshotID:     matchedCase.ID,
			TestCaseCode:           matchedCase.Code,
			TestCaseName:           matchedCase.Name,
			Category:               matchedCase.Category,
			TestType:               matchedCase.TestType,
			RunNumber:              nextRun,
			Status:                 execStatus,
			MeasuredValue:          req.MeasuredValue,
			MeasuredBoolean:        req.MeasuredBoolean,
			ObservedText:           req.ObservedText,
			Unit:                   matchedCase.Unit,
			CalculatedResult:       calculatedResult,
			Notes:                  req.Notes,
			ExecutedBy:             author,
			StartedAt:              &now,
			CompletedAt:            &now,
		}
		if err := s.testingRepo.CreateTestExecution(&newExec); err != nil {
			return nil, fmt.Errorf("failed to create execution run: %w", err)
		}
		targetExec = &newExec
	} else {
		// Update existing pending run
		targetExec.Status = execStatus
		targetExec.MeasuredValue = req.MeasuredValue
		targetExec.MeasuredBoolean = req.MeasuredBoolean
		targetExec.ObservedText = req.ObservedText
		targetExec.CalculatedResult = calculatedResult
		targetExec.Notes = req.Notes
		targetExec.ExecutedBy = author
		targetExec.CompletedAt = &now

		if err := s.testingRepo.UpdateTestExecution(targetExec); err != nil {
			return nil, fmt.Errorf("failed to update execution: %w", err)
		}
	}

	// Record discrete measurement record
	meas := model.TestMeasurement{
		ID:               fmt.Sprintf("MEAS-%s", uuid.New().String()[:8]),
		TestExecutionID:  targetExec.ID,
		SampleIndex:      1,
		MeasuredValue:    req.MeasuredValue,
		MeasuredBoolean:  req.MeasuredBoolean,
		ObservedText:     req.ObservedText,
		Unit:             matchedCase.Unit,
		CalculatedResult: calculatedResult,
		RecordedBy:       author,
		RecordedAt:       now,
	}
	_ = s.testingRepo.CreateMeasurement(&meas)

	// Create linked engineering finding if requested or if failed with finding payload
	if req.CreateFinding && req.FindingTitle != "" {
		sev := req.FindingSeverity
		if sev == "" {
			sev = model.FindingSeverityMedium
		}
		fnd := &model.EngineeringFinding{
			ID:                     fmt.Sprintf("FND-%s", uuid.New().String()[:8]),
			PrototypeBuildID:       session.PrototypeBuildID,
			PrototypeTestSessionID: &session.ID,
			TestExecutionID:        &targetExec.ID,
			Code:                   fmt.Sprintf("FND-%s-%s", matchedCase.Code, uuid.New().String()[:4]),
			Title:                  req.FindingTitle,
			Description:            req.FindingDescription,
			Category:               matchedCase.Category,
			Severity:               sev,
			FindingDisposition:     model.FindingDispOpen,
			ChangeCandidateStatus:  model.ChangeCandidateCandidate,
			RecommendedChangeScope: model.ChangeScopeNone,
			Status:                 "OPEN",
			ReportedBy:             author,
		}
		_ = s.testingRepo.CreateFinding(fnd)
	}

	// Recalculate Session Stats
	s.recalculateSessionStats(session.ID)

	// Log activity
	badge := "emerald"
	if execStatus == model.ExecStatusFailed {
		badge = "red"
	}
	s.logActivity(user, "testing", "TestExecution", targetExec.ID, targetExec.TestCaseName,
		fmt.Sprintf("Test %s (%s)", targetExec.TestCaseCode, execStatus),
		fmt.Sprintf("Executed test '%s' (Run #%d) with result: %s", targetExec.TestCaseName, targetExec.RunNumber, execStatus), badge)

	return s.testingRepo.GetTestExecutionByID(targetExec.ID)
}

func (s *testingService) CreateNextExecutionRun(sessionID, caseSnapshotID string, user *model.User) (*model.TestExecution, error) {
	session, err := s.testingRepo.GetTestSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	author := getUserName(user)

	var snapshot model.FrozenTestPlanSnapshot
	_ = json.Unmarshal([]byte(session.TestPlanSnapshotJSON), &snapshot)

	var matchedCase *model.FrozenTestCaseItem
	for i := range snapshot.TestCases {
		if snapshot.TestCases[i].ID == caseSnapshotID {
			matchedCase = &snapshot.TestCases[i]
			break
		}
	}
	if matchedCase == nil {
		return nil, errors.New("test case not found in snapshot")
	}

	// Find max run number
	maxRun := 0
	for _, e := range session.Executions {
		if e.TestCaseSnapshotID == caseSnapshotID && e.RunNumber > maxRun {
			maxRun = e.RunNumber
		}
	}

	newExec := model.TestExecution{
		ID:                     fmt.Sprintf("EXEC-%s", uuid.New().String()[:8]),
		PrototypeTestSessionID: session.ID,
		TestCaseSnapshotID:     matchedCase.ID,
		TestCaseCode:           matchedCase.Code,
		TestCaseName:           matchedCase.Name,
		Category:               matchedCase.Category,
		TestType:               matchedCase.TestType,
		RunNumber:              maxRun + 1,
		Status:                 model.ExecStatusPending,
		Unit:                   matchedCase.Unit,
		ExecutedBy:             author,
	}

	if err := s.testingRepo.CreateTestExecution(&newExec); err != nil {
		return nil, fmt.Errorf("failed to create run iteration: %w", err)
	}

	return &newExec, nil
}

func (s *testingService) AddEvidence(executionID string, req *AddEvidenceRequest, user *model.User) (*model.TestEvidence, error) {
	exec, err := s.testingRepo.GetTestExecutionByID(executionID)
	if err != nil {
		return nil, fmt.Errorf("execution not found: %w", err)
	}

	author := getUserName(user)

	ev := &model.TestEvidence{
		ID:              fmt.Sprintf("EVID-%s", uuid.New().String()[:8]),
		TestExecutionID: exec.ID,
		FileName:        req.FileName,
		FilePath:        req.FilePath,
		FileType:        req.FileType,
		FileSize:        req.FileSize,
		Description:     req.Description,
		UploadedBy:      author,
		UploadedAt:      time.Now(),
	}

	if err := s.testingRepo.CreateEvidence(ev); err != nil {
		return nil, fmt.Errorf("failed to attach evidence: %w", err)
	}

	s.logActivity(user, "testing", "TestEvidence", ev.ID, ev.FileName, "Attached Evidence",
		fmt.Sprintf("Uploaded evidence '%s' for test %s", ev.FileName, exec.TestCaseCode), "slate")

	return ev, nil
}

func (s *testingService) evaluateTestCase(tc *model.FrozenTestCaseItem, val *float64, bVal *bool, text string) string {
	switch tc.TestType {
	case model.TestTypeNumericRange:
		if val == nil {
			return "FAILED"
		}
		if tc.MinimumValue != nil && *val < *tc.MinimumValue {
			return "FAILED"
		}
		if tc.MaximumValue != nil && *val > *tc.MaximumValue {
			return "FAILED"
		}
		return "PASSED"

	case model.TestTypeNumericMin:
		if val == nil || tc.MinimumValue == nil {
			return "FAILED"
		}
		if *val >= *tc.MinimumValue {
			return "PASSED"
		}
		return "FAILED"

	case model.TestTypeNumericMax:
		if val == nil || tc.MaximumValue == nil {
			return "FAILED"
		}
		if *val <= *tc.MaximumValue {
			return "PASSED"
		}
		return "FAILED"

	case model.TestTypeBooleanPassFail:
		if bVal == nil || tc.ExpectedBoolean == nil {
			return "FAILED"
		}
		if *bVal == *tc.ExpectedBoolean {
			return "PASSED"
		}
		return "FAILED"

	case model.TestTypeVisualInspection, model.TestTypeTextLogVerification:
		if bVal != nil && *bVal {
			return "PASSED"
		}
		return "PASSED"

	default:
		return "PASSED"
	}
}

func (s *testingService) recalculateSessionStats(sessionID string) {
	session, err := s.testingRepo.GetTestSessionByID(sessionID)
	if err != nil {
		return
	}

	// Distinct latest execution per test case
	latestMap := make(map[string]model.TestExecution)
	for _, e := range session.Executions {
		existing, ok := latestMap[e.TestCaseSnapshotID]
		if !ok || e.RunNumber > existing.RunNumber {
			latestMap[e.TestCaseSnapshotID] = e
		}
	}

	passed := 0
	failed := 0
	blocked := 0
	skipped := 0

	for _, e := range latestMap {
		switch e.Status {
		case model.ExecStatusPassed:
			passed++
		case model.ExecStatusFailed:
			failed++
		case model.ExecStatusBlocked:
			blocked++
		case model.ExecStatusSkipped:
			skipped++
		}
	}

	session.PassedCases = passed
	session.FailedCases = failed
	session.BlockedCases = blocked
	session.SkippedCases = skipped

	if session.TotalCases > 0 {
		session.PassRatePercentage = (float64(passed) / float64(session.TotalCases)) * 100.0
	}

	_ = s.testingRepo.UpdateTestSession(session)
}

// -----------------------------------------------------------------------------
// Engineering Findings
// -----------------------------------------------------------------------------
func (s *testingService) GetFindings(prototypeID, sessionID, severity, disposition, changeStatus, query string, page, limit int) ([]model.EngineeringFinding, int64, error) {
	return s.testingRepo.FindFindings(prototypeID, sessionID, severity, disposition, changeStatus, query, page, limit)
}

func (s *testingService) GetFindingByID(id string) (*model.EngineeringFinding, error) {
	return s.testingRepo.GetFindingByID(id)
}

func (s *testingService) CreateFinding(req *CreateFindingRequest, user *model.User) (*model.EngineeringFinding, error) {
	disp := req.FindingDisposition
	if disp == "" {
		disp = model.FindingDispOpen
	}
	changeStat := req.ChangeCandidateStatus
	if changeStat == "" {
		changeStat = model.ChangeCandidateNone
	}
	changeScope := req.RecommendedChangeScope
	if changeScope == "" {
		changeScope = model.ChangeScopeNone
	}

	author := getUserName(user)

	f := &model.EngineeringFinding{
		ID:                     fmt.Sprintf("FND-%s", uuid.New().String()[:8]),
		PrototypeBuildID:       req.PrototypeBuildID,
		PrototypeTestSessionID: req.PrototypeTestSessionID,
		TestExecutionID:        req.TestExecutionID,
		Code:                   fmt.Sprintf("FND-%s", uuid.New().String()[:6]),
		Title:                  req.Title,
		Description:            req.Description,
		Category:               req.Category,
		Severity:               req.Severity,
		FindingDisposition:     disp,
		ChangeCandidateStatus:  changeStat,
		RecommendedChangeScope: changeScope,
		RootCause:              req.RootCause,
		ContainmentAction:      req.ContainmentAction,
		ResolutionNotes:        req.ResolutionNotes,
		Status:                 "OPEN",
		ReportedBy:             author,
		AssignedTo:             req.AssignedTo,
	}

	if err := s.testingRepo.CreateFinding(f); err != nil {
		return nil, fmt.Errorf("failed to create finding: %w", err)
	}

	s.logActivity(user, "testing", "EngineeringFinding", f.ID, f.Title, "Created Finding",
		fmt.Sprintf("Logged %s finding '%s' (Disposition: %s, Change Candidate: %s)", f.Severity, f.Title, f.FindingDisposition, f.ChangeCandidateStatus), "amber")

	return f, nil
}

func (s *testingService) UpdateFinding(id string, req *UpdateFindingRequest, user *model.User) (*model.EngineeringFinding, error) {
	f, err := s.testingRepo.GetFindingByID(id)
	if err != nil {
		return nil, fmt.Errorf("finding not found: %w", err)
	}

	author := getUserName(user)

	if req.Title != "" {
		f.Title = req.Title
	}
	if req.Description != "" {
		f.Description = req.Description
	}
	if req.Category != "" {
		f.Category = req.Category
	}
	if req.Severity != "" {
		f.Severity = req.Severity
	}
	if req.FindingDisposition != "" {
		f.FindingDisposition = req.FindingDisposition
		if req.FindingDisposition == model.FindingDispClosed || req.FindingDisposition == model.FindingDispAcceptedAsIs {
			f.Status = "CLOSED"
			now := time.Now()
			f.ClosedAt = &now
			f.ClosedBy = author
		}
	}
	if req.ChangeCandidateStatus != "" {
		f.ChangeCandidateStatus = req.ChangeCandidateStatus
	}
	if req.RecommendedChangeScope != "" {
		f.RecommendedChangeScope = req.RecommendedChangeScope
	}
	if req.RootCause != "" {
		f.RootCause = req.RootCause
	}
	if req.ContainmentAction != "" {
		f.ContainmentAction = req.ContainmentAction
	}
	if req.ResolutionNotes != "" {
		f.ResolutionNotes = req.ResolutionNotes
	}
	if req.Status != "" {
		f.Status = req.Status
	}
	if req.AssignedTo != "" {
		f.AssignedTo = req.AssignedTo
	}

	if err := s.testingRepo.UpdateFinding(f); err != nil {
		return nil, fmt.Errorf("failed to update finding: %w", err)
	}

	s.logActivity(user, "testing", "EngineeringFinding", f.ID, f.Title, "Updated Finding",
		fmt.Sprintf("Updated finding '%s' (Severity: %s, Disposition: %s)", f.Title, f.Severity, f.FindingDisposition), "amber")

	return f, nil
}

func (s *testingService) DeleteFinding(id string, user *model.User) error {
	f, err := s.testingRepo.GetFindingByID(id)
	if err != nil {
		return fmt.Errorf("finding not found: %w", err)
	}

	if err := s.testingRepo.DeleteFinding(f.ID); err != nil {
		return fmt.Errorf("failed to delete finding: %w", err)
	}

	s.logActivity(user, "testing", "EngineeringFinding", f.ID, f.Title, "Deleted Finding",
		fmt.Sprintf("Removed engineering finding '%s'", f.Title), "red")

	return nil
}

// -----------------------------------------------------------------------------
// Validation Decision Governance Engine
// -----------------------------------------------------------------------------
func (s *testingService) SubmitValidationDecision(sessionID string, req *SubmitValidationRequest, user *model.User) (*model.PrototypeValidationDecision, error) {
	session, err := s.testingRepo.GetTestSessionByID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	author := getUserName(user)

	// Governance Rule 1: Unresolved CRITICAL finding strictly BLOCKS VALIDATED
	critCount, _ := s.testingRepo.CountUnresolvedFindings(session.PrototypeBuildID, session.ID, model.FindingSeverityCritical)
	if critCount > 0 && req.Decision == model.ValidationDecisionValidated {
		return nil, fmt.Errorf("validation blocked: %d unresolved CRITICAL engineering findings exist on this session/prototype; prototype cannot be VALIDATED", critCount)
	}

	// Governance Rule 2: Unresolved HIGH finding strictly BLOCKS VALIDATED
	highCount, _ := s.testingRepo.CountUnresolvedFindings(session.PrototypeBuildID, session.ID, model.FindingSeverityHigh)
	if highCount > 0 && req.Decision == model.ValidationDecisionValidated {
		return nil, fmt.Errorf("validation blocked: %d unresolved HIGH severity engineering findings exist; prototype requires formal disposition or CONDITIONALLY_VALIDATED with justification", highCount)
	}

	// Governance Rule 3: CONDITIONALLY_VALIDATED requires mandatory justification
	if req.Decision == model.ValidationDecisionConditionallyValidated && req.Justification == "" {
		return nil, errors.New("mandatory justification required for CONDITIONALLY_VALIDATED decision")
	}

	decision := &model.PrototypeValidationDecision{
		ID:                     fmt.Sprintf("VAL-%s", uuid.New().String()[:8]),
		PrototypeTestSessionID: session.ID,
		Decision:               req.Decision,
		Justification:          req.Justification,
		DecisionSummary:        req.DecisionSummary,
		DecidedBy:              author,
		IsCurrent:              true,
		CreatedAt:              time.Now(),
	}

	if err := s.testingRepo.CreateValidationDecision(decision); err != nil {
		return nil, fmt.Errorf("failed to record validation decision: %w", err)
	}

	// If validated, update session status to COMPLETED
	if req.Decision == model.ValidationDecisionValidated || req.Decision == model.ValidationDecisionConditionallyValidated {
		session.Status = model.SessionStatusCompleted
		now := time.Now()
		session.CompletedAt = &now
		_ = s.testingRepo.UpdateTestSession(session)
	}

	s.logActivity(user, "testing", "PrototypeValidationDecision", decision.ID, decision.Decision, "Validation Decision",
		fmt.Sprintf("Formal validation decision recorded for session %s: %s (Decided by %s)", session.SessionCode, decision.Decision, author), "emerald")

	return decision, nil
}

func (s *testingService) GetValidationDecisions(sessionID string) ([]model.PrototypeValidationDecision, error) {
	return s.testingRepo.GetValidationDecisionsBySession(sessionID)
}

func (s *testingService) logActivity(user *model.User, module, entityType, entityID, entityName, action, desc, badgeColor string) {
	userName := "System"
	userID := "USR-SYSTEM"
	avatar := ""
	if user != nil {
		userName = user.Name
		userID = user.ID
		avatar = user.Avatar
	}
	act := &model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%s", uuid.New().String()[:8]),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  avatar,
		Module:      module,
		EntityType:  entityType,
		EntityID:    entityID,
		EntityName:  entityName,
		Action:      action,
		Description: desc,
		BadgeColor:  badgeColor,
		CreatedAt:   time.Now(),
	}
	_ = s.actRepo.Create(act)
}
