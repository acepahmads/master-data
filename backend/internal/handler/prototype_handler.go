package handler

import (
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
	"iot-rd-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PrototypeHandler struct {
	service *service.PrototypeService
}

func NewPrototypeHandler(service *service.PrototypeService) *PrototypeHandler {
	return &PrototypeHandler{service: service}
}

func (h *PrototypeHandler) GetAll(c *gin.Context) {
	p := repository.GetPaginationParams(c)
	filters := map[string]string{
		"productId":  c.Query("productId"),
		"versionId":  c.Query("versionId"),
		"revisionId": c.Query("revisionId"),
		"status":     c.Query("status"),
		"owner":      c.Query("owner"),
		"search":     c.Query("search"),
	}

	prototypes, meta, err := h.service.GetAll(p, filters)
	if err != nil {
		ServerError(c, err)
		return
	}

	SuccessWithMeta(c, prototypes, meta)
}

func (h *PrototypeHandler) GetDashboard(c *gin.Context) {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, stats)
}

func (h *PrototypeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	proto, err := h.service.GetByID(id)
	if err != nil {
		NotFound(c, "Prototype build not found")
		return
	}
	Success(c, proto)
}

func (h *PrototypeHandler) GetByRevision(c *gin.Context) {
	revisionID := c.Param("revisionId")
	if revisionID == "" {
		revisionID = c.Param("id")
	}
	prototypes, err := h.service.GetByRevision(revisionID)
	if err != nil {
		ServerError(c, err)
		return
	}
	Success(c, prototypes)
}

func (h *PrototypeHandler) Create(c *gin.Context) {
	var req service.CreatePrototypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid prototype creation data: "+err.Error())
		return
	}

	userID, userName, userAvatar := getUserContext(c)

	proto, err := h.service.Create(req, userID, userName, userAvatar)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, proto, "Prototype build created successfully with frozen BOM snapshot")
}

func (h *PrototypeHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req service.UpdatePrototypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid update data: "+err.Error())
		return
	}

	userID, userName, userAvatar := getUserContext(c)

	proto, err := h.service.Update(id, req, userID, userName, userAvatar)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, proto, "Prototype build updated successfully")
}

func (h *PrototypeHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, userName, userAvatar := getUserContext(c)

	if err := h.service.Delete(id, userID, userName, userAvatar); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"deleted": true}, "Prototype build deleted successfully")
}

func (h *PrototypeHandler) GetBOMSnapshot(c *gin.Context) {
	id := c.Param("id")
	snapshot, err := h.service.GetBOMSnapshot(id)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, snapshot)
}

func (h *PrototypeHandler) GetComponents(c *gin.Context) {
	id := c.Param("id")
	preps, readinessPct, err := h.service.GetComponentPreparations(id)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, gin.H{
		"components":          preps,
		"readinessPercentage": readinessPct,
	})
}

func (h *PrototypeHandler) UpdateComponent(c *gin.Context) {
	id := c.Param("id")
	componentID := c.Param("componentId")

	var req service.UpdateCompPrepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid component preparation data: "+err.Error())
		return
	}

	userID, userName, userAvatar := getUserContext(c)

	prep, readinessPct, err := h.service.UpdateComponentPreparation(id, componentID, req, userID, userName, userAvatar)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{
		"component":           prep,
		"readinessPercentage": readinessPct,
	}, "Component preparation status updated")
}

func (h *PrototypeHandler) GetAssembly(c *gin.Context) {
	id := c.Param("id")
	stages, err := h.service.GetAssemblyStages(id)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, stages)
}

func (h *PrototypeHandler) UpdateAssemblyStage(c *gin.Context) {
	id := c.Param("id")
	stageID := c.Param("stageId")

	var req service.UpdateAssemblyStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid assembly stage data: "+err.Error())
		return
	}

	userID, userName, userAvatar := getUserContext(c)

	stage, err := h.service.UpdateAssemblyStage(id, stageID, req, userID, userName, userAvatar)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, stage, "Assembly stage updated")
}

func (h *PrototypeHandler) GetNotes(c *gin.Context) {
	id := c.Param("id")
	notes, err := h.service.GetEngineeringNotes(id)
	if err != nil {
		NotFound(c, err.Error())
		return
	}
	Success(c, notes)
}

func (h *PrototypeHandler) CreateNote(c *gin.Context) {
	id := c.Param("id")
	var req service.CreateEngineeringNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid note data: "+err.Error())
		return
	}

	userID, userName, userAvatar := getUserContext(c)

	note, err := h.service.CreateEngineeringNote(id, req, userID, userName, userAvatar)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Created(c, note, "Engineering note added")
}

func (h *PrototypeHandler) UpdateNote(c *gin.Context) {
	id := c.Param("id")
	noteID := c.Param("noteId")

	var req service.CreateEngineeringNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid note data: "+err.Error())
		return
	}

	userID, userName, userAvatar := getUserContext(c)

	note, err := h.service.UpdateEngineeringNote(id, noteID, req, userID, userName, userAvatar)
	if err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, note, "Engineering note updated")
}

func (h *PrototypeHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	noteID := c.Param("noteId")
	userID, userName, userAvatar := getUserContext(c)

	if err := h.service.DeleteEngineeringNote(id, noteID, userID, userName, userAvatar); err != nil {
		BadRequest(c, err.Error())
		return
	}

	Success(c, gin.H{"deleted": true}, "Engineering note deleted")
}

func getUserContext(c *gin.Context) (string, string, string) {
	userID := "USR-SYSTEM"
	userName := "System User"
	userAvatar := ""

	if u, exists := c.Get("user"); exists {
		if user, ok := u.(*model.User); ok && user != nil {
			userID = user.ID
			userName = user.Name
			userAvatar = user.Avatar
		}
	}
	return userID, userName, userAvatar
}
