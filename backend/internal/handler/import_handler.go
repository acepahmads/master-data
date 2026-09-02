package handler

import (
	"net/http"
	"strconv"

	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ImportHandler struct {
	importService *service.ImportService
}

func NewImportHandler(impService *service.ImportService) *ImportHandler {
	return &ImportHandler{importService: impService}
}

// CreateBatchRequest represents payload for creating an import batch container.
type CreateBatchRequest struct {
	SourceYear      string `json:"sourceYear"`
	SourceReference string `json:"sourceReference"`
	SourceNotes     string `json:"sourceNotes"`
}

// CreateBatch handles POST /api/v1/data-import/batches
func (h *ImportHandler) CreateBatch(c *gin.Context) {
	var req CreateBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	user, _ := c.Get("user_email")
	userStr, _ := user.(string)
	if userStr == "" {
		userStr = "admin@iotcontrol.io"
	}

	batch, err := h.importService.CreateBatch(req.SourceYear, req.SourceReference, req.SourceNotes, userStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create import batch: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": batch})
}

// GetBatches handles GET /api/v1/data-import/batches
func (h *ImportHandler) GetBatches(c *gin.Context) {
	year := c.Query("year")
	status := c.Query("status")

	batches, err := h.importService.GetBatches(year, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch import batches: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": batches})
}

// GetBatchByID handles GET /api/v1/data-import/batches/:id
func (h *ImportHandler) GetBatchByID(c *gin.Context) {
	id := c.Param("id")
	batch, err := h.importService.GetBatch(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Import batch not found: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": batch})
}

// UploadFiles handles POST /api/v1/data-import/batches/:id/files
func (h *ImportHandler) UploadFiles(c *gin.Context) {
	id := c.Param("id")
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read multipart form: " + err.Error()})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files attached in 'files' field"})
		return
	}

	user, _ := c.Get("user_email")
	userStr, _ := user.(string)
	if userStr == "" {
		userStr = "admin@iotcontrol.io"
	}

	createdFiles, err := h.importService.UploadFiles(id, files, userStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload files: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Files uploaded and inspected successfully",
		"data":    createdFiles,
	})
}

// ParseBatchRequest contains sheet selections and optional header row index overrides per sheet.
type ParseBatchRequest struct {
	SelectedSheetsByFile map[string][]string `json:"selectedSheetsByFile"`
	HeaderRowOverrides   map[string]int      `json:"headerRowOverrides"`
}

// ParseBatch handles POST /api/v1/data-import/batches/:id/parse
func (h *ImportHandler) ParseBatch(c *gin.Context) {
	id := c.Param("id")
	var req ParseBatchRequest
	_ = c.ShouldBindJSON(&req)

	batch, err := h.importService.ParseBatch(id, req.SelectedSheetsByFile, req.HeaderRowOverrides)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse batch workbooks: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": batch})
}

// ApplyMappingRequest contains mapping dictionary.
type ApplyMappingRequest struct {
	Mapping map[string]string `json:"mapping"`
}

// ApplyMapping handles POST /api/v1/data-import/batches/:id/mapping
func (h *ImportHandler) ApplyMapping(c *gin.Context) {
	id := c.Param("id")
	var req ApplyMappingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid mapping payload: " + err.Error()})
		return
	}

	batch, err := h.importService.ApplyColumnMapping(id, req.Mapping)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply column mapping: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": batch})
}

// RunAIClassification handles POST /api/v1/data-import/batches/:id/ai-classify
func (h *ImportHandler) RunAIClassification(c *gin.Context) {
	id := c.Param("id")
	batch, err := h.importService.RunAIClassification(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to run AI classification: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": batch})
}

// GenerateCodesRequest represents optional payload for generating codes.
type GenerateCodesRequest struct {
	ParentProjectName string `json:"parentProjectName"`
}

// GenerateMissingCodes handles POST /api/v1/data-import/batches/:id/generate-codes
func (h *ImportHandler) GenerateMissingCodes(c *gin.Context) {
	id := c.Param("id")
	var req GenerateCodesRequest
	_ = c.ShouldBindJSON(&req)

	batch, err := h.importService.GenerateMissingCodes(id, req.ParentProjectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate missing codes: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Missing Codes/MPN generated successfully",
		"data":    batch,
	})
}

// GetStagedRows handles GET /api/v1/data-import/batches/:id/rows
func (h *ImportHandler) GetStagedRows(c *gin.Context) {
	id := c.Param("id")
	status := c.Query("status")
	decision := c.Query("decision")
	classification := c.Query("classification")
	search := c.Query("search")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	rows, total, err := h.importService.GetStagedRows(id, status, decision, classification, search, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch staged rows: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  rows,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// TriageRowRequest contains row triage edits and decisions.
type TriageRowRequest struct {
	ItemClassification string  `json:"itemClassification"`
	ProductType        string  `json:"productType"`
	NormalizedName     string  `json:"normalizedName"`
	NormalizedCode     string  `json:"normalizedCode"`
	NormalizedCategory string  `json:"normalizedCategory"`
	Description        string  `json:"description"`
	UnitCost           float64 `json:"unitCost"`
	SellingPrice       float64 `json:"sellingPrice"`
	Currency           string  `json:"currency"`
	Decision           string  `json:"decision"`
	ReviewNotes        string  `json:"reviewNotes"`
}

// TriageRow handles PUT /api/v1/data-import/rows/:rowId
func (h *ImportHandler) TriageRow(c *gin.Context) {
	rowID := c.Param("rowId")
	var req TriageRowRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	user, _ := c.Get("user_email")
	userStr, _ := user.(string)
	if userStr == "" {
		userStr = "reviewer@iotcontrol.io"
	}

	row, err := h.importService.TriageRow(
		rowID,
		req.ItemClassification,
		req.ProductType,
		req.NormalizedName,
		req.NormalizedCode,
		req.NormalizedCategory,
		req.UnitCost,
		req.SellingPrice,
		req.Currency,
		req.Decision,
		userStr,
		req.ReviewNotes,
		req.Description,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to triage row: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": row})
}

// ApproveBatch handles POST /api/v1/data-import/batches/:id/approve
func (h *ImportHandler) ApproveBatch(c *gin.Context) {
	id := c.Param("id")
	user, _ := c.Get("user_email")
	userStr, _ := user.(string)
	if userStr == "" {
		userStr = "manager@iotcontrol.io"
	}

	batch, err := h.importService.ApproveBatch(id, userStr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve import batch: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import batch approved successfully for Master Data migration",
		"data":    batch,
	})
}

// CommitBatchRequest contains optional parent project name for project imports.
type CommitBatchRequest struct {
	ParentProjectName string `json:"parentProjectName"`
}

// CommitBatch handles POST /api/v1/data-import/batches/:id/commit
func (h *ImportHandler) CommitBatch(c *gin.Context) {
	id := c.Param("id")
	user, _ := c.Get("user_email")
	userStr, _ := user.(string)
	if userStr == "" {
		userStr = "manager@iotcontrol.io"
	}

	var req CommitBatchRequest
	_ = c.ShouldBindJSON(&req)

	result, err := h.importService.CommitApprovedBatch(id, userStr, req.ParentProjectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit batch to Master Data: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Approved records committed to Master Data successfully",
		"data":    result,
	})
}

// DeleteBatch handles DELETE /api/v1/data-import/batches/:id
func (h *ImportHandler) DeleteBatch(c *gin.Context) {
	id := c.Param("id")
	if err := h.importService.DeleteBatch(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete import batch: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Import batch and associated staging records deleted successfully",
	})
}
