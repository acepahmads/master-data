package handler

import (
	"strconv"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ChangeHandler struct {
	service *service.ChangeService
}

func NewChangeHandler(service *service.ChangeService) *ChangeHandler {
	return &ChangeHandler{service: service}
}

func getActor(c *gin.Context) (string, string, string) {
	actorID := "USR-SYSTEM"
	actorName := "System"
	actorRole := "Admin"

	if uVal, exists := c.Get("user"); exists {
		if user, ok := uVal.(*model.User); ok && user != nil {
			actorID = user.ID
			actorName = user.Name
		}
	}
	if rVal, exists := c.Get("role"); exists {
		if roleStr, ok := rVal.(string); ok {
			actorRole = roleStr
		}
	}

	return actorID, actorName, actorRole
}

// -----------------------------------------------------------------------------
// DASHBOARD
// -----------------------------------------------------------------------------

func (h *ChangeHandler) GetDashboard(c *gin.Context) {
	stats, err := h.service.GetDashboard()
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, stats, "Dashboard retrieved successfully")
}

// -----------------------------------------------------------------------------
// ECR HANDLERS
// -----------------------------------------------------------------------------

func (h *ChangeHandler) GetECRs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := repository.ECRFilter{
		Status:             c.Query("status"),
		Priority:           c.Query("priority"),
		ProductID:          c.Query("productId"),
		HardwareRevisionID: c.Query("hardwareRevisionId"),
		Origin:             c.Query("origin"),
		Search:             c.Query("search"),
		Page:               page,
		Limit:              limit,
	}

	ecrs, total, err := h.service.GetECRs(filter)
	if err != nil {
		ServerError(c, err)
		return
	}

	Success(c, gin.H{
		"items": ecrs,
		"total": total,
		"page":  page,
		"limit": limit,
	}, "ECRs retrieved successfully")
}

func (h *ChangeHandler) GetECR(c *gin.Context) {
	id := c.Param("id")
	ecr, err := h.service.GetECR(id)
	if err != nil {
		NotFound(c, "Engineering change request not found")
		return
	}
	Success(c, ecr, "ECR retrieved successfully")
}

func (h *ChangeHandler) CreateECR(c *gin.Context) {
	var req service.CreateECRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid ECR payload", err.Error())
		return
	}

	actorID, actorName, _ := getActor(c)
	ecr, err := h.service.CreateECR(req, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to create ECR", err.Error())
		return
	}

	Created(c, ecr, "ECR created successfully")
}

func (h *ChangeHandler) UpdateECR(c *gin.Context) {
	id := c.Param("id")
	var req service.CreateECRRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid ECR payload", err.Error())
		return
	}

	actorID, actorName, _ := getActor(c)
	ecr, err := h.service.UpdateECR(id, req, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to update ECR", err.Error())
		return
	}

	Success(c, ecr, "ECR updated successfully")
}

func (h *ChangeHandler) DeleteECR(c *gin.Context) {
	id := c.Param("id")
	actorID, actorName, _ := getActor(c)
	if err := h.service.DeleteECR(id, actorID, actorName); err != nil {
		BadRequest(c, "Failed to delete ECR", err.Error())
		return
	}
	Success(c, nil, "ECR deleted successfully")
}

func (h *ChangeHandler) SubmitECR(c *gin.Context) {
	id := c.Param("id")
	actorID, actorName, _ := getActor(c)
	ecr, err := h.service.SubmitECR(id, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to submit ECR", err.Error())
		return
	}
	Success(c, ecr, "ECR submitted for review successfully")
}

// -----------------------------------------------------------------------------
// CHANGE ITEMS
// -----------------------------------------------------------------------------

func (h *ChangeHandler) AddChangeItem(c *gin.Context) {
	ecrID := c.Param("id")
	var req service.AddChangeItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid change item payload", err.Error())
		return
	}

	actorID, actorName, _ := getActor(c)
	item, err := h.service.AddChangeItem(ecrID, req, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to add change item", err.Error())
		return
	}

	Created(c, item, "Change item added successfully")
}

func (h *ChangeHandler) DeleteChangeItem(c *gin.Context) {
	itemID := c.Param("id")
	actorID, actorName, _ := getActor(c)
	if err := h.service.DeleteChangeItem(itemID, actorID, actorName); err != nil {
		BadRequest(c, "Failed to delete change item", err.Error())
		return
	}
	Success(c, nil, "Change item deleted successfully")
}

// -----------------------------------------------------------------------------
// IMPACT & REVIEWS & APPROVALS
// -----------------------------------------------------------------------------

func (h *ChangeHandler) SaveImpact(c *gin.Context) {
	ecrID := c.Param("id")
	var req service.SaveImpactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid impact analysis payload", err.Error())
		return
	}

	actorID, actorName, _ := getActor(c)
	impact, err := h.service.SaveImpactAnalysis(ecrID, req, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to save impact analysis", err.Error())
		return
	}

	Created(c, impact, "Impact analysis saved successfully")
}

func (h *ChangeHandler) SubmitReview(c *gin.Context) {
	ecrID := c.Param("id")
	var req service.SubmitReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid review payload", err.Error())
		return
	}

	actorID, actorName, _ := getActor(c)
	review, err := h.service.SubmitReview(ecrID, req, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to submit review", err.Error())
		return
	}

	Created(c, review, "Technical review recorded successfully")
}

func (h *ChangeHandler) SubmitApproval(c *gin.Context) {
	ecrID := c.Param("id")
	var req service.SubmitApprovalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid approval payload", err.Error())
		return
	}

	actorID, actorName, actorRole := getActor(c)
	approval, err := h.service.SubmitApproval(ecrID, req, actorID, actorName, actorRole)
	if err != nil {
		BadRequest(c, "Failed to submit approval", err.Error())
		return
	}

	Created(c, approval, "Approval decision recorded successfully")
}

// -----------------------------------------------------------------------------
// ECO HANDLERS
// -----------------------------------------------------------------------------

func (h *ChangeHandler) GetECOs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	filter := repository.ECOFilter{
		Status:             c.Query("status"),
		ProductID:          c.Query("productId"),
		HardwareRevisionID: c.Query("hardwareRevisionId"),
		Search:             c.Query("search"),
		Page:               page,
		Limit:              limit,
	}

	ecos, total, err := h.service.GetECOs(filter)
	if err != nil {
		ServerError(c, err)
		return
	}

	Success(c, gin.H{
		"items": ecos,
		"total": total,
		"page":  page,
		"limit": limit,
	}, "ECOs retrieved successfully")
}

func (h *ChangeHandler) GetECO(c *gin.Context) {
	id := c.Param("id")
	eco, err := h.service.GetECO(id)
	if err != nil {
		NotFound(c, "Engineering change order not found")
		return
	}
	Success(c, eco, "ECO retrieved successfully")
}

func (h *ChangeHandler) CreateECO(c *gin.Context) {
	var req service.CreateECORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid ECO payload", err.Error())
		return
	}

	actorID, actorName, _ := getActor(c)
	eco, err := h.service.CreateECO(req, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to create ECO", err.Error())
		return
	}

	Created(c, eco, "ECO created successfully")
}

func (h *ChangeHandler) ApproveECO(c *gin.Context) {
	id := c.Param("id")
	actorID, actorName, _ := getActor(c)
	eco, err := h.service.ApproveECO(id, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to approve ECO", err.Error())
		return
	}
	Success(c, eco, "ECO approved for implementation")
}

func (h *ChangeHandler) GetChangePreview(c *gin.Context) {
	id := c.Param("id")
	preview, err := h.service.GenerateBOMChangePreview(id)
	if err != nil {
		BadRequest(c, "Failed to generate BOM change preview", err.Error())
		return
	}
	Success(c, preview, "BOM change preview generated successfully")
}

func (h *ChangeHandler) ImplementECO(c *gin.Context) {
	id := c.Param("id")
	actorID, actorName, _ := getActor(c)
	eco, err := h.service.ImplementECO(id, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to implement ECO", err.Error())
		return
	}
	Success(c, eco, "ECO implemented successfully; target revision and BOM baseline created")
}

func (h *ChangeHandler) VerifyECO(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Notes string `json:"notes"`
	}
	_ = c.ShouldBindJSON(&req)

	actorID, actorName, _ := getActor(c)
	eco, err := h.service.VerifyECO(id, req.Notes, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to verify ECO", err.Error())
		return
	}
	Success(c, eco, "ECO verified successfully")
}

func (h *ChangeHandler) CloseECO(c *gin.Context) {
	id := c.Param("id")
	actorID, actorName, _ := getActor(c)
	eco, err := h.service.CloseECO(id, actorID, actorName)
	if err != nil {
		BadRequest(c, "Failed to close ECO", err.Error())
		return
	}
	Success(c, eco, "ECO closed successfully")
}

// -----------------------------------------------------------------------------
// TRACEABILITY
// -----------------------------------------------------------------------------

func (h *ChangeHandler) GetTraceability(c *gin.Context) {
	productID := c.Query("productId")
	chains, err := h.service.GetTraceability(productID)
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, chains, "Traceability chains retrieved successfully")
}
