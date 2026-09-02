package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"

	"gorm.io/gorm"
)

type ChangeService struct {
	repo         *repository.ChangeRepository
	activityRepo *repository.ActivityRepository
}

func NewChangeService(repo *repository.ChangeRepository, activityRepo *repository.ActivityRepository) *ChangeService {
	return &ChangeService{
		repo:         repo,
		activityRepo: activityRepo,
	}
}

func genShortID(prefix string) string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(b))
}

func (s *ChangeService) logActivity(actorID, actorName, entityType, entityID, entityName, action, desc, badgeColor string) {
	if s.activityRepo != nil {
		_ = s.activityRepo.Create(&model.ActivityLog{
			ID:          genShortID("ACT"),
			UserID:      actorID,
			UserName:    actorName,
			Module:      "changes",
			EntityType:  entityType,
			EntityID:    entityID,
			EntityName:  entityName,
			Action:      action,
			Description: desc,
			BadgeColor:  badgeColor,
		})
	}
}

func genCode(prefix string) string {
	b := make([]byte, 2)
	rand.Read(b)
	return fmt.Sprintf("%s-%d-%s", prefix, time.Now().Year(), strings.ToUpper(hex.EncodeToString(b)))
}

// -----------------------------------------------------------------------------
// DASHBOARD STATS
// -----------------------------------------------------------------------------

func (s *ChangeService) GetDashboard() (*repository.ChangeDashboardStats, error) {
	return s.repo.GetDashboardStats()
}

// -----------------------------------------------------------------------------
// ECR CRUD & LIFECYCLE
// -----------------------------------------------------------------------------

type CreateECRRequest struct {
	Title                    string     `json:"title" binding:"required"`
	Origin                   string     `json:"origin"`
	SourceFindingID          *string    `json:"sourceFindingId"`
	ProductID                string     `json:"productId"`
	ProductVersionID         string     `json:"productVersionId"`
	SourceHardwareRevisionID string     `json:"sourceHardwareRevisionId" binding:"required"`
	ChangeType               string     `json:"changeType"`
	Priority                 string     `json:"priority"`
	Description              string     `json:"description" binding:"required"`
	BusinessReason           string     `json:"businessReason"`
	TechnicalJustification   string     `json:"technicalJustification"`
	AssignedLead             string     `json:"assignedLead"`
	EstimatedCost            float64    `json:"estimatedCost"`
	TargetReleaseDate        *time.Time `json:"targetReleaseDate"`
}

func (s *ChangeService) CreateECR(req CreateECRRequest, actorID, actorName string) (*model.EngineeringChangeRequest, error) {
	// If source finding is provided, verify finding exists
	if req.SourceFindingID != nil && *req.SourceFindingID != "" {
		var finding model.EngineeringFinding
		err := s.repo.DB().Where("id = ? OR code = ?", *req.SourceFindingID, *req.SourceFindingID).First(&finding).Error
		if err != nil {
			return nil, fmt.Errorf("source finding '%s' not found: %w", *req.SourceFindingID, err)
		}
		req.SourceFindingID = &finding.ID
		// Ingest context if empty
		if req.ChangeType == "" {
			req.ChangeType = string(finding.RecommendedChangeScope)
		}
		if req.Priority == "" {
			req.Priority = string(finding.Severity)
		}
	}

	// 1. Verify revision exists
	var rev model.HardwareRevision
	if err := s.repo.DB().Where("id = ? OR code = ?", req.SourceHardwareRevisionID, req.SourceHardwareRevisionID).First(&rev).Error; err != nil {
		return nil, fmt.Errorf("hardware revision '%s' not found: %w", req.SourceHardwareRevisionID, err)
	}
	req.SourceHardwareRevisionID = rev.ID

	// 2. Resolve version
	var ver model.ProductVersion
	verID := req.ProductVersionID
	if verID == "" {
		verID = rev.ProductVersionID
	}
	if err := s.repo.DB().Where("id = ? OR version_number = ?", verID, verID).First(&ver).Error; err != nil {
		return nil, fmt.Errorf("product version '%s' not found: %w", verID, err)
	}
	req.ProductVersionID = ver.ID

	// 3. Resolve product
	var prod model.Product
	prodID := req.ProductID
	if prodID == "" {
		prodID = ver.ProductID
	}
	if err := s.repo.DB().Where("id = ? OR code = ?", prodID, prodID).First(&prod).Error; err != nil {
		return nil, fmt.Errorf("product '%s' not found: %w", prodID, err)
	}
	req.ProductID = prod.ID

	origin := req.Origin
	if origin == "" {
		if req.SourceFindingID != nil && *req.SourceFindingID != "" {
			origin = "TEST_FINDING"
		} else {
			origin = "ENGINEERING_REVIEW"
		}
	}

	changeType := req.ChangeType
	if changeType == "" {
		changeType = "SCHEMATIC_CIRCUIT"
	}

	priority := req.Priority
	if priority == "" {
		priority = "MEDIUM"
	}

	ecr := &model.EngineeringChangeRequest{
		ID:                       genShortID("ECR"),
		Code:                     genCode("ECR"),
		Title:                    req.Title,
		Origin:                   origin,
		SourceFindingID:          req.SourceFindingID,
		ProductID:                req.ProductID,
		ProductVersionID:         req.ProductVersionID,
		SourceHardwareRevisionID: req.SourceHardwareRevisionID,
		ChangeType:               changeType,
		Priority:                 priority,
		Description:              req.Description,
		BusinessReason:           req.BusinessReason,
		TechnicalJustification:   req.TechnicalJustification,
		Status:                   "DRAFT",
		RequestedBy:              actorName,
		AssignedLead:             req.AssignedLead,
		EstimatedCost:            req.EstimatedCost,
		TargetReleaseDate:        req.TargetReleaseDate,
	}

	if err := s.repo.CreateECR(ecr); err != nil {
		return nil, fmt.Errorf("failed to create ECR: %w", err)
	}

	s.logActivity(actorID, actorName, "EngineeringChangeRequest", ecr.ID, ecr.Code, "Created ECR", fmt.Sprintf("Created change request %s: %s (Priority: %s)", ecr.Code, ecr.Title, ecr.Priority), "blue")

	return s.repo.GetECRByID(ecr.ID)
}

func (s *ChangeService) GetECRs(filter repository.ECRFilter) ([]model.EngineeringChangeRequest, int64, error) {
	return s.repo.GetECRs(filter)
}

func (s *ChangeService) GetECR(id string) (*model.EngineeringChangeRequest, error) {
	return s.repo.GetECRByID(id)
}

func (s *ChangeService) UpdateECR(id string, req CreateECRRequest, actorID, actorName string) (*model.EngineeringChangeRequest, error) {
	ecr, err := s.repo.GetECRByID(id)
	if err != nil {
		return nil, err
	}

	if ecr.Status != "DRAFT" && ecr.Status != "UNDER_REVIEW" {
		return nil, errors.New("cannot edit ECR after it has been submitted for approval or converted to ECO")
	}

	ecr.Title = req.Title
	if req.ChangeType != "" {
		ecr.ChangeType = req.ChangeType
	}
	if req.Priority != "" {
		ecr.Priority = req.Priority
	}
	ecr.Description = req.Description
	ecr.BusinessReason = req.BusinessReason
	ecr.TechnicalJustification = req.TechnicalJustification
	ecr.AssignedLead = req.AssignedLead
	ecr.EstimatedCost = req.EstimatedCost
	ecr.TargetReleaseDate = req.TargetReleaseDate

	if err := s.repo.UpdateECR(ecr); err != nil {
		return nil, err
	}

	s.logActivity(actorID, actorName, "EngineeringChangeRequest", ecr.ID, ecr.Code, "Updated ECR", fmt.Sprintf("Updated metadata on %s: %s", ecr.Code, ecr.Title), "slate")

	return s.repo.GetECRByID(ecr.ID)
}

func (s *ChangeService) DeleteECR(id string, actorID, actorName string) error {
	ecr, err := s.repo.GetECRByID(id)
	if err != nil {
		return err
	}

	if ecr.Status != "DRAFT" {
		return errors.New("only DRAFT ECRs may be deleted")
	}

	return s.repo.DeleteECR(ecr.ID)
}

func (s *ChangeService) SubmitECR(id string, actorID, actorName string) (*model.EngineeringChangeRequest, error) {
	ecr, err := s.repo.GetECRByID(id)
	if err != nil {
		return nil, err
	}

	if ecr.Status != "DRAFT" {
		return nil, fmt.Errorf("cannot submit ECR in status '%s'; must be DRAFT", ecr.Status)
	}

	ecr.Status = "UNDER_REVIEW"
	if err := s.repo.UpdateECR(ecr); err != nil {
		return nil, err
	}

	s.logActivity(actorID, actorName, "EngineeringChangeRequest", ecr.ID, ecr.Code, "Submitted ECR", fmt.Sprintf("Submitted %s for technical review & impact analysis", ecr.Code), "amber")

	return s.repo.GetECRByID(ecr.ID)
}

// -----------------------------------------------------------------------------
// CHANGE ITEMS
// -----------------------------------------------------------------------------

type AddChangeItemRequest struct {
	Sequence               int     `json:"sequence"`
	ChangeType             string  `json:"changeType" binding:"required"`
	Title                  string  `json:"title" binding:"required"`
	AffectedReference      string  `json:"affectedReference"`
	AffectedComponentID    *string `json:"affectedComponentId"`
	CurrentState           string  `json:"currentState"`
	ProposedState          string  `json:"proposedState"`
	CurrentValue           string  `json:"currentValue"`
	ProposedValue          string  `json:"proposedValue"`
	TechnicalJustification string  `json:"technicalJustification"`
	RiskLevel              string  `json:"riskLevel"`
	Notes                  string  `json:"notes"`
}

func (s *ChangeService) AddChangeItem(ecrID string, req AddChangeItemRequest, actorID, actorName string) (*model.ECRChangeItem, error) {
	ecr, err := s.repo.GetECRByID(ecrID)
	if err != nil {
		return nil, err
	}

	if ecr.Status != "DRAFT" && ecr.Status != "UNDER_REVIEW" {
		return nil, errors.New("cannot add change items to locked or approved ECR")
	}

	risk := req.RiskLevel
	if risk == "" {
		risk = "MEDIUM"
	}

	seq := req.Sequence
	if seq <= 0 {
		seq = len(ecr.ChangeItems) + 1
	}

	curr := req.CurrentState
	if curr == "" {
		curr = req.CurrentValue
	}
	prop := req.ProposedState
	if prop == "" {
		prop = req.ProposedValue
	}

	item := &model.ECRChangeItem{
		ID:                     genShortID("CHG"),
		ECRID:                  ecr.ID,
		Sequence:               seq,
		ChangeType:             req.ChangeType,
		Title:                  req.Title,
		AffectedReference:      req.AffectedReference,
		AffectedComponentID:    req.AffectedComponentID,
		CurrentState:           curr,
		ProposedState:          prop,
		TechnicalJustification: req.TechnicalJustification,
		RiskLevel:              risk,
		Notes:                  req.Notes,
	}

	if err := s.repo.AddChangeItem(item); err != nil {
		return nil, err
	}

	return item, nil
}

func (s *ChangeService) DeleteChangeItem(itemID string, actorID, actorName string) error {
	return s.repo.DeleteChangeItem(itemID)
}

// -----------------------------------------------------------------------------
// IMPACT ANALYSIS
// -----------------------------------------------------------------------------

type SaveImpactRequest struct {
	Domain              string  `json:"domain" binding:"required"`
	ImpactLevel         string  `json:"impactLevel" binding:"required"`
	Description         string  `json:"description" binding:"required"`
	Mitigation          string  `json:"mitigation"`
	CostDelta           float64 `json:"costDelta"`
	ScheduleImpactWeeks float64 `json:"scheduleImpactWeeks"`
	ScrapDisposition    string  `json:"scrapDisposition"`
}

func (s *ChangeService) SaveImpactAnalysis(ecrID string, req SaveImpactRequest, actorID, actorName string) (*model.ECRImpactAnalysis, error) {
	ecr, err := s.repo.GetECRByID(ecrID)
	if err != nil {
		return nil, err
	}

	scrap := req.ScrapDisposition
	if scrap == "" {
		scrap = "RUN_OUT"
	}

	impact := &model.ECRImpactAnalysis{
		ID:                  genShortID("IMP"),
		ECRID:               ecr.ID,
		Domain:              req.Domain,
		ImpactLevel:         req.ImpactLevel,
		Description:         req.Description,
		Mitigation:          req.Mitigation,
		CostDelta:           req.CostDelta,
		ScheduleImpactWeeks: req.ScheduleImpactWeeks,
		ScrapDisposition:    scrap,
		ReviewerID:          actorID,
		ReviewerName:        actorName,
		ReviewStatus:        "REVIEWED",
	}

	if err := s.repo.SaveImpactAnalysis(impact); err != nil {
		return nil, err
	}

	s.logActivity(actorID, actorName, "ECRImpactAnalysis", ecr.ID, ecr.Code, "Updated Impact Analysis", fmt.Sprintf("Logged %s impact (%s) on %s", impact.Domain, impact.ImpactLevel, ecr.Code), "blue")

	return impact, nil
}

// -----------------------------------------------------------------------------
// REVIEWS
// -----------------------------------------------------------------------------

type SubmitReviewRequest struct {
	Discipline string `json:"discipline" binding:"required"` // HARDWARE, FIRMWARE, MECHANICAL, SOFTWARE, QA, PROJECT_MANAGEMENT
	Decision   string `json:"decision" binding:"required"`   // APPROVE, REJECT, REQUEST_CHANGES
	Comments   string `json:"comments"`
}

func (s *ChangeService) SubmitReview(ecrID string, req SubmitReviewRequest, actorID, actorName string) (*model.ECRReview, error) {
	ecr, err := s.repo.GetECRByID(ecrID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	review := &model.ECRReview{
		ID:           genShortID("REV"),
		ECRID:        ecr.ID,
		Discipline:   req.Discipline,
		ReviewerID:   actorID,
		ReviewerName: actorName,
		Decision:     req.Decision,
		Comments:     req.Comments,
		ReviewedAt:   &now,
	}

	if err := s.repo.SaveReview(review); err != nil {
		return nil, err
	}

	// Update ECR status if appropriate
	if req.Decision == "REQUEST_CHANGES" {
		ecr.Status = "UNDER_REVIEW"
		_ = s.repo.UpdateECR(ecr)
	} else if req.Decision == "APPROVE" && ecr.Status == "UNDER_REVIEW" {
		ecr.Status = "PENDING_APPROVAL"
		_ = s.repo.UpdateECR(ecr)
	}

	s.logActivity(actorID, actorName, "ECRReview", ecr.ID, ecr.Code, "Submitted Technical Review", fmt.Sprintf("%s Review on %s: %s", review.Discipline, ecr.Code, review.Decision), "indigo")

	return review, nil
}

// -----------------------------------------------------------------------------
// APPROVALS
// -----------------------------------------------------------------------------

type SubmitApprovalRequest struct {
	Stage         string `json:"stage" binding:"required"` // TIER_1_TECH_LEAD, TIER_2_ENG_MANAGER, TIER_3_QUALITY_DIRECTOR
	Decision      string `json:"decision" binding:"required"` // APPROVED, REJECTED, REQUEST_CHANGES
	Justification string `json:"justification"`
}

func (s *ChangeService) SubmitApproval(ecrID string, req SubmitApprovalRequest, actorID, actorName, actorRole string) (*model.ECRApproval, error) {
	ecr, err := s.repo.GetECRByID(ecrID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	approval := &model.ECRApproval{
		ID:            genShortID("APPR"),
		ECRID:         ecr.ID,
		Stage:         req.Stage,
		ApproverID:    actorID,
		ApproverName:  actorName,
		ApproverRole:  actorRole,
		Decision:      req.Decision,
		Justification: req.Justification,
		DecidedAt:     &now,
	}

	if err := s.repo.SaveApproval(approval); err != nil {
		return nil, err
	}

	// Check if all CCB stages approved
	approvals, _ := s.repo.GetApprovalsByECR(ecr.ID)
	allApproved := true
	anyRejected := false
	anyChanges := false

	for _, a := range approvals {
		if a.Decision == "REJECTED" {
			anyRejected = true
		}
		if a.Decision == "REQUEST_CHANGES" {
			anyChanges = true
		}
		if a.Decision != "APPROVED" {
			allApproved = false
		}
	}

	if anyRejected {
		ecr.Status = "REJECTED"
		_ = s.repo.UpdateECR(ecr)
	} else if anyChanges {
		ecr.Status = "UNDER_REVIEW"
		_ = s.repo.UpdateECR(ecr)
	} else if allApproved && len(approvals) >= 1 {
		ecr.Status = "APPROVED"
		_ = s.repo.UpdateECR(ecr)
	}

	s.logActivity(actorID, actorName, "ECRApproval", ecr.ID, ecr.Code, "CCB Approval Decision", fmt.Sprintf("%s on %s: %s (Status: %s)", approval.Stage, ecr.Code, approval.Decision, ecr.Status), "emerald")

	return approval, nil
}

// -----------------------------------------------------------------------------
// ECO CREATION & WORKFLOW
// -----------------------------------------------------------------------------

type CreateECORequest struct {
	ECRID               string `json:"ecrId" binding:"required"`
	Title               string `json:"title"`
	ImplementationOwner string `json:"implementationOwner"`
	ImplementationNotes string `json:"implementationNotes"`
	ScrapDisposition    string `json:"scrapDisposition"`
}

func (s *ChangeService) CreateECO(req CreateECORequest, actorID, actorName string) (*model.EngineeringChangeOrder, error) {
	ecr, err := s.repo.GetECRByID(req.ECRID)
	if err != nil {
		return nil, fmt.Errorf("source ECR '%s' not found: %w", req.ECRID, err)
	}

	// CRITICAL GUARDRAIL #1: Must be explicitly APPROVED ECR
	if ecr.Status != "APPROVED" && ecr.Status != "CONVERTED_TO_ECO" {
		return nil, fmt.Errorf("cannot create ECO from ECR in status '%s'; ECR must be APPROVED first", ecr.Status)
	}

	// Check if ECO already exists
	if ecr.ECO != nil {
		return nil, fmt.Errorf("ECO '%s' already exists for ECR '%s'", ecr.ECO.Code, ecr.Code)
	}

	title := req.Title
	if title == "" {
		title = fmt.Sprintf("Implement %s: %s", ecr.Code, ecr.Title)
	}

	owner := req.ImplementationOwner
	if owner == "" {
		owner = ecr.AssignedLead
		if owner == "" {
			owner = actorName
		}
	}

	scrap := req.ScrapDisposition
	if scrap == "" {
		scrap = "RUN_OUT"
	}

	// Lookup source BOM from source revision
	var sourceBOM model.EngineeringBOM
	var sourceBOMID *string
	if err := s.repo.DB().Where("hardware_revision_id = ?", ecr.SourceHardwareRevisionID).First(&sourceBOM).Error; err == nil {
		sourceBOMID = &sourceBOM.ID
	}

	eco := &model.EngineeringChangeOrder{
		ID:                        genShortID("ECO"),
		Code:                      genCode("ECO"),
		ECRID:                     ecr.ID,
		Title:                     title,
		SourceHardwareRevisionID:  ecr.SourceHardwareRevisionID,
		SourceBOMID:               sourceBOMID,
		ImplementationOwner:       owner,
		Status:                    "DRAFT",
		ImplementationNotes:       req.ImplementationNotes,
		ScrapDisposition:          scrap,
		RequiresRegressionTesting: true,
		CreatedBy:                 actorName,
	}

	if err := s.repo.CreateECO(eco); err != nil {
		return nil, fmt.Errorf("failed to create ECO: %w", err)
	}

	s.logActivity(actorID, actorName, "EngineeringChangeOrder", eco.ID, eco.Code, "Created ECO", fmt.Sprintf("Created Change Order %s from approved %s", eco.Code, ecr.Code), "blue")

	return s.repo.GetECOByID(eco.ID)
}

func (s *ChangeService) GetECOs(filter repository.ECOFilter) ([]model.EngineeringChangeOrder, int64, error) {
	return s.repo.GetECOs(filter)
}

func (s *ChangeService) GetECO(id string) (*model.EngineeringChangeOrder, error) {
	return s.repo.GetECOByID(id)
}

func (s *ChangeService) ApproveECO(id string, actorID, actorName string) (*model.EngineeringChangeOrder, error) {
	eco, err := s.repo.GetECOByID(id)
	if err != nil {
		return nil, err
	}

	if eco.Status != "DRAFT" {
		return nil, fmt.Errorf("cannot approve ECO in status '%s'; must be DRAFT", eco.Status)
	}

	now := time.Now()
	eco.Status = "APPROVED"
	eco.ApprovedAt = &now

	if err := s.repo.UpdateECO(eco); err != nil {
		return nil, err
	}

	s.logActivity(actorID, actorName, "EngineeringChangeOrder", eco.ID, eco.Code, "Approved ECO", fmt.Sprintf("Approved %s for implementation baseline generation", eco.Code), "emerald")

	return s.repo.GetECOByID(eco.ID)
}

// -----------------------------------------------------------------------------
// BOM CHANGE PREVIEW
// -----------------------------------------------------------------------------

func (s *ChangeService) GenerateBOMChangePreview(ecoID string) (*model.BOMChangePreviewSummary, error) {
	eco, err := s.repo.GetECOByID(ecoID)
	if err != nil {
		return nil, err
	}

	ecr, err := s.repo.GetECRByID(eco.ECRID)
	if err != nil {
		return nil, err
	}

	// Fetch source BOM items
	var sourceBOMItems []model.BOMItem
	if eco.SourceBOMID != nil && *eco.SourceBOMID != "" {
		_ = s.repo.DB().Where("engineering_bom_id = ?", *eco.SourceBOMID).Preload("Component").Find(&sourceBOMItems)
	}

	summary := &model.BOMChangePreviewSummary{
		SourceRevisionCode: "REV-A0",
		TargetRevisionCode: "REV-B0 (Target)",
	}

	if eco.SourceHardwareRevision != nil {
		summary.SourceRevisionCode = eco.SourceHardwareRevision.Code
	}

	// Map of proposed changes by reference designator
	changeByRef := make(map[string]model.ECRChangeItem)
	for _, ci := range ecr.ChangeItems {
		if ci.AffectedReference != "" {
			changeByRef[ci.AffectedReference] = ci
		}
	}

	var sourceTotal, targetTotal float64

	for _, sItem := range sourceBOMItems {
		compName := "Component"
		mpn := ""
		pkg := ""
		unitCost := sItem.UnitCost
		compId := ""
		if sItem.ComponentID != nil {
			compId = *sItem.ComponentID
		}
		if sItem.Component != nil {
			compName = sItem.Component.Name
			mpn = sItem.Component.PartNumber
			pkg = sItem.Component.PackageType
			if sItem.Component.EstimatedUnitCost > 0 {
				unitCost = sItem.Component.EstimatedUnitCost
			}
		}

		lineTotal := sItem.Quantity * unitCost
		sourceTotal += lineTotal

		pItem := model.BOMChangePreviewItem{
			ReferenceDesignator: sItem.ReferenceDesignators,
			ComponentID:         compId,
			ComponentName:       compName,
			MPN:                 mpn,
			Package:             pkg,
			SourceQuantity:      int(sItem.Quantity),
			TargetQuantity:      int(sItem.Quantity),
			SourceUnitCost:      unitCost,
			TargetUnitCost:      unitCost,
			SourceTotalCost:     lineTotal,
			TargetTotalCost:     lineTotal,
			CostDelta:           0.0,
			ChangeStatus:        "UNCHANGED",
		}

		// Check if reference has change item
		if ci, hasChange := changeByRef[sItem.ReferenceDesignators]; hasChange {
			if ci.ChangeType == "BOM_COMPONENT" || ci.ChangeType == "COMPONENT_VALUE" {
				pItem.ChangeStatus = "MODIFIED"
				pItem.ChangeReason = ci.Title
				if ci.AffectedComponent != nil {
					pItem.ComponentName = ci.AffectedComponent.Name
					pItem.MPN = ci.AffectedComponent.PartNumber
					pItem.Package = ci.AffectedComponent.PackageType
					pItem.TargetUnitCost = ci.AffectedComponent.EstimatedUnitCost
					pItem.TargetTotalCost = float64(pItem.TargetQuantity) * ci.AffectedComponent.EstimatedUnitCost
					pItem.CostDelta = pItem.TargetTotalCost - pItem.SourceTotalCost
				} else {
					// Assume minor delta from description
					pItem.TargetUnitCost = unitCost + 0.06
					pItem.TargetTotalCost = float64(pItem.TargetQuantity) * pItem.TargetUnitCost
					pItem.CostDelta = pItem.TargetTotalCost - pItem.SourceTotalCost
				}
				summary.TotalModified++
			}
		} else {
			summary.TotalUnchanged++
		}

		targetTotal += pItem.TargetTotalCost
		summary.Items = append(summary.Items, pItem)
	}

	summary.SourceBOMCost = sourceTotal
	summary.TargetBOMCost = targetTotal
	summary.NetCostDelta = targetTotal - sourceTotal
	if sourceTotal > 0 {
		summary.PercentCostDelta = (summary.NetCostDelta / sourceTotal) * 100.0
	}

	return summary, nil
}

// -----------------------------------------------------------------------------
// ATOMIC ECO IMPLEMENTATION ENGINE (PART 10)
// -----------------------------------------------------------------------------

func (s *ChangeService) ImplementECO(ecoID string, actorID, actorName string) (*model.EngineeringChangeOrder, error) {
	var implementedECO *model.EngineeringChangeOrder

	// Execute inside atomic transaction
	err := s.repo.DB().Transaction(func(tx *gorm.DB) error {
		// 1. Validate ECO
		var eco model.EngineeringChangeOrder
		if err := tx.Where("id = ? OR code = ?", ecoID, ecoID).
			Preload("ECR.ChangeItems.AffectedComponent").
			Preload("SourceHardwareRevision").
			Preload("SourceBOM.Items").
			First(&eco).Error; err != nil {
			return fmt.Errorf("ECO '%s' not found: %w", ecoID, err)
		}

		if eco.Status != "APPROVED" && eco.Status != "IN_IMPLEMENTATION" {
			return fmt.Errorf("cannot implement ECO in status '%s'; must be APPROVED", eco.Status)
		}

		// 2. Validate source revision
		if eco.SourceHardwareRevision == nil {
			return errors.New("source hardware revision is missing")
		}

		// 3. Increment revision code (e.g. REV-001 -> REV-002 or REV-B0)
		sourceRev := eco.SourceHardwareRevision
		var revCount int64
		tx.Model(&model.HardwareRevision{}).Where("product_version_id = ?", sourceRev.ProductVersionID).Count(&revCount)

		targetRevCode := fmt.Sprintf("REV-%03d", revCount+1)
		targetRevName := fmt.Sprintf("%s (ECO %s Baseline)", sourceRev.Name, eco.Code)

		// 4. Create Target Hardware Revision (In DRAFT status)
		targetRev := model.HardwareRevision{
			ID:                genShortID("REV"),
			ProductVersionID:  sourceRev.ProductVersionID,
			Code:              targetRevCode,
			Name:              targetRevName,
			Status:            "Draft",
			ChangeSummary:     fmt.Sprintf("Baseline generated from %s (Source: %s). Includes ECO change items.", eco.Code, sourceRev.Code),
			PCBRevision:       "B0",
			SchematicRevision: "v1.1",
			Engineer:          actorName,
		}

		if err := tx.Create(&targetRev).Error; err != nil {
			return fmt.Errorf("failed to create target hardware revision: %w", err)
		}

		// 5. Clone Source BOM to Target BOM
		var targetBOMID *string
		if eco.SourceBOMID != nil && *eco.SourceBOMID != "" {
			var srcBOM model.EngineeringBOM
			if err := tx.Where("id = ?", *eco.SourceBOMID).Preload("Items").First(&srcBOM).Error; err == nil {
				targetBOMCode := fmt.Sprintf("BOM-%s", targetRev.Code)
				targetBOM := model.EngineeringBOM{
					ID:                 genShortID("BOM"),
					BOMCode:            targetBOMCode,
					HardwareRevisionID: targetRev.ID,
					Name:               fmt.Sprintf("BOM for %s", targetRev.Code),
					Status:             "Draft",
					Currency:           "USD",
					CreatedBy:          actorName,
				}

				if err := tx.Create(&targetBOM).Error; err != nil {
					return fmt.Errorf("failed to clone target BOM: %w", err)
				}
				targetBOMID = &targetBOM.ID

				// Clone all BOM items and apply ECR change items
				changeByRef := make(map[string]model.ECRChangeItem)
				if eco.ECR != nil {
					for _, ci := range eco.ECR.ChangeItems {
						if ci.AffectedReference != "" {
							changeByRef[ci.AffectedReference] = ci
						}
					}
				}

				var totalCost float64
				for _, sItem := range srcBOM.Items {
					targetItem := model.BOMItem{
						ID:                   genShortID("BMI"),
						EngineeringBOMID:     targetBOM.ID,
						ItemType:             sItem.ItemType,
						Name:                 sItem.Name,
						ComponentID:          sItem.ComponentID,
						ComponentPartNumber:  sItem.ComponentPartNumber,
						ComponentCategory:    sItem.ComponentCategory,
						Quantity:             sItem.Quantity,
						UnitID:               sItem.UnitID,
						UnitCode:             sItem.UnitCode,
						UnitCost:             sItem.UnitCost,
						LineCost:             sItem.LineCost,
						ReferenceDesignators: sItem.ReferenceDesignators,
						Position:             sItem.Position,
						Notes:                sItem.Notes,
					}

					// Apply Change Item delta if matching ref designator
					if ci, hasChange := changeByRef[sItem.ReferenceDesignators]; hasChange {
						if ci.AffectedComponentID != nil && *ci.AffectedComponentID != "" {
							targetItem.ComponentID = ci.AffectedComponentID
							if ci.AffectedComponent != nil && ci.AffectedComponent.EstimatedUnitCost > 0 {
								targetItem.UnitCost = ci.AffectedComponent.EstimatedUnitCost
								targetItem.LineCost = targetItem.Quantity * targetItem.UnitCost
							}
						}
						targetItem.Notes = fmt.Sprintf("Modified by %s: %s", eco.Code, ci.Title)
					}

					if err := tx.Create(&targetItem).Error; err != nil {
						return fmt.Errorf("failed to clone BOM item '%s': %w", targetItem.ReferenceDesignators, err)
					}

					totalCost += targetItem.LineCost
				}

				targetBOM.TotalCost = totalCost
				_ = tx.Save(&targetBOM)
			}
		}

		// 6. Update ECO status to IMPLEMENTED & Set Regression Mandate
		now := time.Now()
		eco.TargetHardwareRevisionID = &targetRev.ID
		eco.TargetBOMID = targetBOMID
		eco.Status = "IMPLEMENTED"
		eco.RequiresRegressionTesting = true
		eco.ImplementedAt = &now

		if err := tx.Save(&eco).Error; err != nil {
			return fmt.Errorf("failed to update ECO status: %w", err)
		}

		// 7. Update ECR status to CONVERTED_TO_ECO
		if eco.ECR != nil {
			eco.ECR.Status = "CONVERTED_TO_ECO"
			_ = tx.Save(eco.ECR)
		}

		// 8. Log Activity
		s.logActivity(actorID, actorName, "EngineeringChangeOrder", eco.ID, eco.Code, "Implemented ECO", fmt.Sprintf("Implemented %s: Spawned target revision %s and target BOM. Regression testing required.", eco.Code, targetRev.Code), "emerald")

		implementedECO = &eco
		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.repo.GetECOByID(implementedECO.ID)
}

func (s *ChangeService) VerifyECO(id string, notes string, actorID, actorName string) (*model.EngineeringChangeOrder, error) {
	eco, err := s.repo.GetECOByID(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	eco.Status = "VERIFIED"
	eco.VerifiedAt = &now
	if notes != "" {
		eco.ImplementationNotes = fmt.Sprintf("%s\n[Verification Note]: %s", eco.ImplementationNotes, notes)
	}

	if err := s.repo.UpdateECO(eco); err != nil {
		return nil, err
	}

	s.logActivity(actorID, actorName, "EngineeringChangeOrder", eco.ID, eco.Code, "Verified ECO", fmt.Sprintf("Verified %s post-implementation regression results", eco.Code), "purple")

	return s.repo.GetECOByID(eco.ID)
}

func (s *ChangeService) CloseECO(id string, actorID, actorName string) (*model.EngineeringChangeOrder, error) {
	eco, err := s.repo.GetECOByID(id)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	eco.Status = "CLOSED"
	eco.ClosedAt = &now

	if err := s.repo.UpdateECO(eco); err != nil {
		return nil, err
	}

	s.logActivity(actorID, actorName, "EngineeringChangeOrder", eco.ID, eco.Code, "Closed ECO", fmt.Sprintf("Formally closed change order %s", eco.Code), "slate")

	return s.repo.GetECOByID(eco.ID)
}

// -----------------------------------------------------------------------------
// TRACEABILITY
// -----------------------------------------------------------------------------

func (s *ChangeService) GetTraceability(productID string) ([]model.ChangeTraceabilityChain, error) {
	return s.repo.GetTraceabilityChains(productID)
}
