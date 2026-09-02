package repository

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"iot-rd-backend/internal/model"
)

type TestingRepository interface {
	GetDashboardStats() (*model.TestingDashboardStats, error)

	// Master Test Plans
	FindTestPlans(category, status, query string, page, limit int) ([]model.TestPlan, int64, error)
	GetTestPlanByID(id string) (*model.TestPlan, error)
	CreateTestPlan(tp *model.TestPlan) error
	UpdateTestPlan(tp *model.TestPlan) error
	DeleteTestPlan(id string) error

	// Test Plan Versions
	FindTestPlanVersions(testPlanID string) ([]model.TestPlanVersion, error)
	GetTestPlanVersionByID(id string) (*model.TestPlanVersion, error)
	CreateTestPlanVersion(v *model.TestPlanVersion) error
	UpdateTestPlanVersion(v *model.TestPlanVersion) error

	// Test Cases
	FindTestCasesByVersion(versionID string) ([]model.TestCase, error)
	GetTestCaseByID(id string) (*model.TestCase, error)
	CreateTestCase(tc *model.TestCase) error
	UpdateTestCase(tc *model.TestCase) error
	DeleteTestCase(id string) error

	// Prototype Test Sessions
	FindTestSessionsByPrototype(prototypeID string) ([]model.PrototypeTestSession, error)
	GetTestSessionByID(id string) (*model.PrototypeTestSession, error)
	CreateTestSession(session *model.PrototypeTestSession) error
	UpdateTestSession(session *model.PrototypeTestSession) error

	// Test Executions
	FindExecutionsBySession(sessionID string) ([]model.TestExecution, error)
	GetTestExecutionByID(id string) (*model.TestExecution, error)
	CreateTestExecution(exec *model.TestExecution) error
	UpdateTestExecution(exec *model.TestExecution) error

	// Measurements & Evidence
	CreateMeasurement(meas *model.TestMeasurement) error
	CreateEvidence(ev *model.TestEvidence) error

	// Engineering Findings
	FindFindings(prototypeID, sessionID, severity, disposition, changeStatus, query string, page, limit int) ([]model.EngineeringFinding, int64, error)
	GetFindingByID(id string) (*model.EngineeringFinding, error)
	CreateFinding(f *model.EngineeringFinding) error
	UpdateFinding(f *model.EngineeringFinding) error
	DeleteFinding(id string) error

	// Validation Decisions
	CreateValidationDecision(dec *model.PrototypeValidationDecision) error
	GetValidationDecisionsBySession(sessionID string) ([]model.PrototypeValidationDecision, error)
	CountUnresolvedFindings(prototypeID, sessionID string, minSeverity string) (int64, error)
}

type testingRepository struct {
	db *gorm.DB
}

func NewTestingRepository(db *gorm.DB) TestingRepository {
	return &testingRepository{db: db}
}

func (r *testingRepository) GetDashboardStats() (*model.TestingDashboardStats, error) {
	var stats model.TestingDashboardStats

	// Total test plans
	r.db.Model(&model.TestPlan{}).Count(&stats.TotalTestPlans)

	// Active test sessions
	r.db.Model(&model.PrototypeTestSession{}).Where("status != ?", model.SessionStatusCompleted).Count(&stats.ActiveTestSessions)

	// Total executions & breakdown
	r.db.Model(&model.TestExecution{}).Count(&stats.TotalExecutions)
	r.db.Model(&model.TestExecution{}).Where("status = ?", model.ExecStatusPassed).Count(&stats.TotalPassed)
	r.db.Model(&model.TestExecution{}).Where("status = ?", model.ExecStatusFailed).Count(&stats.TotalFailed)
	r.db.Model(&model.TestExecution{}).Where("status = ?", model.ExecStatusBlocked).Count(&stats.TotalBlocked)

	if stats.TotalExecutions > 0 {
		stats.OverallPassRate = float64(stats.TotalPassed) / float64(stats.TotalExecutions) * 100.0
	}

	// Findings
	r.db.Model(&model.EngineeringFinding{}).Where("status != 'CLOSED'").Count(&stats.OpenFindingsCount)
	r.db.Model(&model.EngineeringFinding{}).Where("status != 'CLOSED' AND severity = ?", model.FindingSeverityCritical).Count(&stats.CriticalFindingsCount)
	r.db.Model(&model.EngineeringFinding{}).Where("status != 'CLOSED' AND severity = ?", model.FindingSeverityHigh).Count(&stats.HighFindingsCount)
	r.db.Model(&model.EngineeringFinding{}).Where("change_candidate_status = ?", model.ChangeCandidateCandidate).Count(&stats.ChangeCandidatesCount)

	// Validation Decisions
	r.db.Model(&model.PrototypeTestSession{}).Where("status = ? AND id NOT IN (SELECT prototype_test_session_id FROM prototype_validation_decisions WHERE is_current = 1)", model.SessionStatusCompleted).Count(&stats.PendingValidations)
	r.db.Model(&model.PrototypeValidationDecision{}).Where("is_current = 1 AND decision = ?", model.ValidationDecisionValidated).Count(&stats.TotalValidatedSessions)
	r.db.Model(&model.PrototypeValidationDecision{}).Where("is_current = 1 AND decision = ?", model.ValidationDecisionRejected).Count(&stats.TotalRejectedSessions)

	return &stats, nil
}

// Master Test Plans
func (r *testingRepository) FindTestPlans(category, status, query string, page, limit int) ([]model.TestPlan, int64, error) {
	var plans []model.TestPlan
	var total int64

	db := r.db.Model(&model.TestPlan{}).Preload("Versions", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("created_at DESC")
	})

	if category != "" {
		db = db.Where("category = ?", category)
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}
	if query != "" {
		q := "%" + query + "%"
		db = db.Where("code LIKE ? OR name LIKE ? OR description LIKE ?", q, q, q)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Order("created_at DESC").Find(&plans).Error; err != nil {
		return nil, 0, err
	}
	return plans, total, nil
}

func (r *testingRepository) GetTestPlanByID(id string) (*model.TestPlan, error) {
	var tp model.TestPlan
	err := r.db.Preload("Versions.TestCases", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sequence ASC")
	}).Where("id = ? OR code = ?", id, id).First(&tp).Error
	if err != nil {
		return nil, err
	}
	return &tp, nil
}

func (r *testingRepository) CreateTestPlan(tp *model.TestPlan) error {
	if tp.ID == "" {
		tp.ID = fmt.Sprintf("TP-%s", uuid.New().String()[:8])
	}
	return r.db.Create(tp).Error
}

func (r *testingRepository) UpdateTestPlan(tp *model.TestPlan) error {
	return r.db.Save(tp).Error
}

func (r *testingRepository) DeleteTestPlan(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.TestPlan{}).Error
}

// Test Plan Versions
func (r *testingRepository) FindTestPlanVersions(testPlanID string) ([]model.TestPlanVersion, error) {
	var versions []model.TestPlanVersion
	err := r.db.Preload("TestCases", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sequence ASC")
	}).Where("test_plan_id = ?", testPlanID).Order("created_at DESC").Find(&versions).Error
	return versions, err
}

func (r *testingRepository) GetTestPlanVersionByID(id string) (*model.TestPlanVersion, error) {
	var ver model.TestPlanVersion
	err := r.db.Preload("TestCases", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("sequence ASC")
	}).Where("id = ?", id).First(&ver).Error
	if err != nil {
		return nil, err
	}
	return &ver, nil
}

func (r *testingRepository) CreateTestPlanVersion(v *model.TestPlanVersion) error {
	if v.ID == "" {
		v.ID = fmt.Sprintf("TPV-%s", uuid.New().String()[:8])
	}
	return r.db.Create(v).Error
}

func (r *testingRepository) UpdateTestPlanVersion(v *model.TestPlanVersion) error {
	return r.db.Save(v).Error
}

// Test Cases
func (r *testingRepository) FindTestCasesByVersion(versionID string) ([]model.TestCase, error) {
	var cases []model.TestCase
	err := r.db.Where("test_plan_version_id = ?", versionID).Order("sequence ASC").Find(&cases).Error
	return cases, err
}

func (r *testingRepository) GetTestCaseByID(id string) (*model.TestCase, error) {
	var tc model.TestCase
	err := r.db.Where("id = ? OR code = ?", id, id).First(&tc).Error
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *testingRepository) CreateTestCase(tc *model.TestCase) error {
	if tc.ID == "" {
		tc.ID = fmt.Sprintf("TC-%s", uuid.New().String()[:8])
	}
	return r.db.Create(tc).Error
}

func (r *testingRepository) UpdateTestCase(tc *model.TestCase) error {
	return r.db.Save(tc).Error
}

func (r *testingRepository) DeleteTestCase(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.TestCase{}).Error
}

// Prototype Test Sessions
func (r *testingRepository) FindTestSessionsByPrototype(prototypeID string) ([]model.PrototypeTestSession, error) {
	var sessions []model.PrototypeTestSession
	err := r.db.Preload("Executions.Measurements").Preload("Executions.Evidences").Preload("Decisions", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("created_at DESC")
	}).Where("prototype_build_id = ?", prototypeID).Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func (r *testingRepository) GetTestSessionByID(id string) (*model.PrototypeTestSession, error) {
	var session model.PrototypeTestSession
	err := r.db.Preload("Executions.Measurements").Preload("Executions.Evidences").Preload("Decisions", func(tx *gorm.DB) *gorm.DB {
		return tx.Order("created_at DESC")
	}).Where("id = ? OR session_code = ?", id, id).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *testingRepository) CreateTestSession(session *model.PrototypeTestSession) error {
	if session.ID == "" {
		session.ID = fmt.Sprintf("SESS-%s", uuid.New().String()[:8])
	}
	return r.db.Create(session).Error
}

func (r *testingRepository) UpdateTestSession(session *model.PrototypeTestSession) error {
	return r.db.Save(session).Error
}

// Test Executions
func (r *testingRepository) FindExecutionsBySession(sessionID string) ([]model.TestExecution, error) {
	var execs []model.TestExecution
	err := r.db.Preload("Measurements").Preload("Evidences").Where("prototype_test_session_id = ?", sessionID).Order("test_case_snapshot_id ASC, run_number ASC").Find(&execs).Error
	return execs, err
}

func (r *testingRepository) GetTestExecutionByID(id string) (*model.TestExecution, error) {
	var exec model.TestExecution
	err := r.db.Preload("Measurements").Preload("Evidences").Where("id = ?", id).First(&exec).Error
	if err != nil {
		return nil, err
	}
	return &exec, nil
}

func (r *testingRepository) CreateTestExecution(exec *model.TestExecution) error {
	if exec.ID == "" {
		exec.ID = fmt.Sprintf("EXEC-%s", uuid.New().String()[:8])
	}
	return r.db.Create(exec).Error
}

func (r *testingRepository) UpdateTestExecution(exec *model.TestExecution) error {
	return r.db.Save(exec).Error
}

// Measurements & Evidence
func (r *testingRepository) CreateMeasurement(meas *model.TestMeasurement) error {
	if meas.ID == "" {
		meas.ID = fmt.Sprintf("MEAS-%s", uuid.New().String()[:8])
	}
	return r.db.Create(meas).Error
}

func (r *testingRepository) CreateEvidence(ev *model.TestEvidence) error {
	if ev.ID == "" {
		ev.ID = fmt.Sprintf("EVID-%s", uuid.New().String()[:8])
	}
	return r.db.Create(ev).Error
}

// Engineering Findings
func (r *testingRepository) FindFindings(prototypeID, sessionID, severity, disposition, changeStatus, query string, page, limit int) ([]model.EngineeringFinding, int64, error) {
	var findings []model.EngineeringFinding
	var total int64

	db := r.db.Model(&model.EngineeringFinding{})

	if prototypeID != "" {
		db = db.Where("prototype_build_id = ?", prototypeID)
	}
	if sessionID != "" {
		db = db.Where("prototype_test_session_id = ?", sessionID)
	}
	if severity != "" {
		db = db.Where("severity = ?", severity)
	}
	if disposition != "" {
		db = db.Where("finding_disposition = ?", disposition)
	}
	if changeStatus != "" {
		db = db.Where("change_candidate_status = ?", changeStatus)
	}
	if query != "" {
		q := "%" + query + "%"
		db = db.Where("code LIKE ? OR title LIKE ? OR description LIKE ?", q, q, q)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Order("created_at DESC").Find(&findings).Error; err != nil {
		return nil, 0, err
	}
	return findings, total, nil
}

func (r *testingRepository) GetFindingByID(id string) (*model.EngineeringFinding, error) {
	var f model.EngineeringFinding
	err := r.db.Where("id = ? OR code = ?", id, id).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *testingRepository) CreateFinding(f *model.EngineeringFinding) error {
	if f.ID == "" {
		f.ID = fmt.Sprintf("FND-%s", uuid.New().String()[:8])
	}
	return r.db.Create(f).Error
}

func (r *testingRepository) UpdateFinding(f *model.EngineeringFinding) error {
	return r.db.Save(f).Error
}

func (r *testingRepository) DeleteFinding(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.EngineeringFinding{}).Error
}

// Validation Decisions
func (r *testingRepository) CreateValidationDecision(dec *model.PrototypeValidationDecision) error {
	if dec.ID == "" {
		dec.ID = fmt.Sprintf("VAL-%s", uuid.New().String()[:8])
	}
	// Mark previous decisions for this session as not current
	if err := r.db.Model(&model.PrototypeValidationDecision{}).
		Where("prototype_test_session_id = ?", dec.PrototypeTestSessionID).
		Update("is_current", false).Error; err != nil {
		return err
	}
	dec.IsCurrent = true
	return r.db.Create(dec).Error
}

func (r *testingRepository) GetValidationDecisionsBySession(sessionID string) ([]model.PrototypeValidationDecision, error) {
	var decs []model.PrototypeValidationDecision
	err := r.db.Where("prototype_test_session_id = ?", sessionID).Order("created_at DESC").Find(&decs).Error
	return decs, err
}

func (r *testingRepository) CountUnresolvedFindings(prototypeID, sessionID string, minSeverity string) (int64, error) {
	var count int64
	db := r.db.Model(&model.EngineeringFinding{}).Where("status != 'CLOSED' AND finding_disposition NOT IN ('ACCEPTED_AS_IS', 'CLOSED')")
	if sessionID != "" {
		db = db.Where("prototype_test_session_id = ?", sessionID)
	} else if prototypeID != "" {
		db = db.Where("prototype_build_id = ?", prototypeID)
	}

	if minSeverity == model.FindingSeverityCritical {
		db = db.Where("severity = ?", model.FindingSeverityCritical)
	} else if minSeverity == model.FindingSeverityHigh {
		db = db.Where("severity IN (?, ?)", model.FindingSeverityCritical, model.FindingSeverityHigh)
	}

	err := db.Count(&count).Error
	return count, err
}
