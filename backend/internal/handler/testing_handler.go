package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/service"
)

func getUser(c *gin.Context) *model.User {
	if u, exists := c.Get("user"); exists {
		if user, ok := u.(*model.User); ok {
			return user
		}
	}
	if u, exists := c.Get("currentUser"); exists {
		if user, ok := u.(*model.User); ok {
			return user
		}
	}
	return nil
}

type TestingHandler struct {
	testingService service.TestingService
}

func NewTestingHandler(testingService service.TestingService) *TestingHandler {
	return &TestingHandler{testingService: testingService}
}

// -----------------------------------------------------------------------------
// Testing Dashboard
// -----------------------------------------------------------------------------
func (h *TestingHandler) GetDashboard(c *gin.Context) {
	stats, err := h.testingService.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": stats})
}

// -----------------------------------------------------------------------------
// Test Plans
// -----------------------------------------------------------------------------
func (h *TestingHandler) GetTestPlans(c *gin.Context) {
	category := c.Query("category")
	status := c.Query("status")
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	plans, total, err := h.testingService.GetTestPlans(category, status, query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    plans,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *TestingHandler) GetTestPlan(c *gin.Context) {
	id := c.Param("id")
	tp, err := h.testingService.GetTestPlanByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "test plan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": tp})
}

func (h *TestingHandler) CreateTestPlan(c *gin.Context) {
	var req service.CreateTestPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	tp, err := h.testingService.CreateTestPlan(&req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "test plan created successfully", "data": tp})
}

func (h *TestingHandler) UpdateTestPlan(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTestPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	tp, err := h.testingService.UpdateTestPlan(id, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "test plan updated successfully", "data": tp})
}

func (h *TestingHandler) DeleteTestPlan(c *gin.Context) {
	id := c.Param("id")
	user := getUser(c)

	if err := h.testingService.DeleteTestPlan(id, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "test plan deleted successfully"})
}

// -----------------------------------------------------------------------------
// Test Plan Versions & Immutability
// -----------------------------------------------------------------------------
func (h *TestingHandler) GetTestPlanVersion(c *gin.Context) {
	id := c.Param("id")
	ver, err := h.testingService.GetTestPlanVersionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "test plan version not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": ver})
}

func (h *TestingHandler) CreateTestPlanVersion(c *gin.Context) {
	testPlanID := c.Param("id")
	var req service.CreateTestPlanVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	ver, err := h.testingService.CreateTestPlanVersion(testPlanID, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "version created successfully", "data": ver})
}

func (h *TestingHandler) UpdateTestPlanVersion(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTestPlanVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	ver, err := h.testingService.UpdateTestPlanVersion(id, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "version updated successfully", "data": ver})
}

func (h *TestingHandler) ReleaseTestPlanVersion(c *gin.Context) {
	id := c.Param("id")
	var req service.ReleaseTestPlanVersionRequest
	_ = c.ShouldBindJSON(&req)

	user := getUser(c)

	ver, err := h.testingService.ReleaseTestPlanVersion(id, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "version released successfully", "data": ver})
}

func (h *TestingHandler) CloneTestPlanVersion(c *gin.Context) {
	id := c.Param("id")
	var req service.CloneTestPlanVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	ver, err := h.testingService.CloneTestPlanVersion(id, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "version cloned successfully", "data": ver})
}

// -----------------------------------------------------------------------------
// Test Cases
// -----------------------------------------------------------------------------
func (h *TestingHandler) CreateTestCase(c *gin.Context) {
	versionID := c.Param("id")
	var req service.CreateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	tc, err := h.testingService.CreateTestCase(versionID, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "test case created successfully", "data": tc})
}

func (h *TestingHandler) UpdateTestCase(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	tc, err := h.testingService.UpdateTestCase(id, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "test case updated successfully", "data": tc})
}

func (h *TestingHandler) DeleteTestCase(c *gin.Context) {
	id := c.Param("id")
	user := getUser(c)

	if err := h.testingService.DeleteTestCase(id, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "test case deleted successfully"})
}

// -----------------------------------------------------------------------------
// Prototype Test Sessions & Snapshots
// -----------------------------------------------------------------------------
func (h *TestingHandler) GetPrototypeTestSessions(c *gin.Context) {
	prototypeID := c.Param("id")
	sessions, err := h.testingService.GetTestSessionsByPrototype(prototypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sessions})
}

func (h *TestingHandler) CreatePrototypeTestSession(c *gin.Context) {
	prototypeID := c.Param("id")
	var req service.CreateTestSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	session, err := h.testingService.CreatePrototypeTestSession(prototypeID, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "test session initiated successfully with frozen snapshot", "data": session})
}

func (h *TestingHandler) GetTestSession(c *gin.Context) {
	id := c.Param("id")
	session, err := h.testingService.GetTestSessionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "test session not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": session})
}

func (h *TestingHandler) UpdateTestSession(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateTestSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	session, err := h.testingService.UpdateTestSession(id, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "test session updated successfully", "data": session})
}

// -----------------------------------------------------------------------------
// Test Executions & Measurements
// -----------------------------------------------------------------------------
func (h *TestingHandler) GetSessionExecutions(c *gin.Context) {
	sessionID := c.Param("id")
	execs, err := h.testingService.GetExecutionsBySession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": execs})
}

func (h *TestingHandler) RecordExecution(c *gin.Context) {
	sessionID := c.Param("id")
	var req service.RecordExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	exec, err := h.testingService.RecordExecutionRun(sessionID, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "execution recorded successfully", "data": exec})
}

func (h *TestingHandler) CreateNextRun(c *gin.Context) {
	sessionID := c.Param("id")
	caseSnapshotID := c.Param("caseId")

	user := getUser(c)

	exec, err := h.testingService.CreateNextExecutionRun(sessionID, caseSnapshotID, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "new execution run iteration created", "data": exec})
}

func (h *TestingHandler) AddEvidence(c *gin.Context) {
	executionID := c.Param("id")
	var req service.AddEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	ev, err := h.testingService.AddEvidence(executionID, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "evidence attached successfully", "data": ev})
}

// -----------------------------------------------------------------------------
// Engineering Findings
// -----------------------------------------------------------------------------
func (h *TestingHandler) GetFindings(c *gin.Context) {
	prototypeID := c.Query("prototypeId")
	sessionID := c.Query("sessionId")
	severity := c.Query("severity")
	disposition := c.Query("disposition")
	changeStatus := c.Query("changeStatus")
	query := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	findings, total, err := h.testingService.GetFindings(prototypeID, sessionID, severity, disposition, changeStatus, query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    findings,
		"meta": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

func (h *TestingHandler) GetFinding(c *gin.Context) {
	id := c.Param("id")
	f, err := h.testingService.GetFindingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "finding not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": f})
}

func (h *TestingHandler) CreateFinding(c *gin.Context) {
	var req service.CreateFindingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	f, err := h.testingService.CreateFinding(&req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "engineering finding created successfully", "data": f})
}

func (h *TestingHandler) UpdateFinding(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdateFindingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	f, err := h.testingService.UpdateFinding(id, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "finding updated successfully", "data": f})
}

func (h *TestingHandler) DeleteFinding(c *gin.Context) {
	id := c.Param("id")
	user := getUser(c)

	if err := h.testingService.DeleteFinding(id, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "finding deleted successfully"})
}

// -----------------------------------------------------------------------------
// Validation Decisions
// -----------------------------------------------------------------------------
func (h *TestingHandler) GetValidationDecisions(c *gin.Context) {
	sessionID := c.Param("id")
	decisions, err := h.testingService.GetValidationDecisions(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": decisions})
}

func (h *TestingHandler) SubmitValidationDecision(c *gin.Context) {
	sessionID := c.Param("id")
	var req service.SubmitValidationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	user := getUser(c)

	decision, err := h.testingService.SubmitValidationDecision(sessionID, &req, user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "validation decision recorded successfully", "data": decision})
}
