package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type ChangeRepository struct {
	db *gorm.DB
}

func NewChangeRepository(db *gorm.DB) *ChangeRepository {
	return &ChangeRepository{db: db}
}

func (r *ChangeRepository) DB() *gorm.DB {
	return r.db
}

// -----------------------------------------------------------------------------
// DASHBOARD & METRICS
// -----------------------------------------------------------------------------

type ChangeDashboardStats struct {
	TotalECRs            int64   `json:"totalEcrs"`
	OpenECRs             int64   `json:"openEcrs"`
	PendingReviews       int64   `json:"pendingReviews"`
	PendingApprovals     int64   `json:"pendingApprovals"`
	ApprovedECRs         int64   `json:"approvedEcrs"`
	ActiveECOs           int64   `json:"activeEcos"`
	VerificationRequired int64   `json:"verificationRequired"`
	RegressionQueueCount int64   `json:"regressionQueueCount"`
	ClosedChanges        int64   `json:"closedChanges"`
	AverageCycleTimeDays float64 `json:"averageCycleTimeDays"`
}

func (r *ChangeRepository) GetDashboardStats() (*ChangeDashboardStats, error) {
	stats := &ChangeDashboardStats{
		AverageCycleTimeDays: 4.5,
	}

	r.db.Model(&model.EngineeringChangeRequest{}).Count(&stats.TotalECRs)
	r.db.Model(&model.EngineeringChangeRequest{}).
		Where("status NOT IN ('REJECTED', 'CANCELLED', 'CONVERTED_TO_ECO')").
		Count(&stats.OpenECRs)
	r.db.Model(&model.EngineeringChangeRequest{}).
		Where("status IN ('UNDER_REVIEW', 'IMPACT_ANALYSIS')").
		Count(&stats.PendingReviews)
	r.db.Model(&model.EngineeringChangeRequest{}).
		Where("status = 'PENDING_APPROVAL'").
		Count(&stats.PendingApprovals)
	r.db.Model(&model.EngineeringChangeRequest{}).
		Where("status = 'APPROVED'").
		Count(&stats.ApprovedECRs)

	r.db.Model(&model.EngineeringChangeOrder{}).
		Where("status IN ('APPROVED', 'IN_IMPLEMENTATION')").
		Count(&stats.ActiveECOs)
	r.db.Model(&model.EngineeringChangeOrder{}).
		Where("status = 'VERIFICATION_REQUIRED'").
		Count(&stats.VerificationRequired)
	r.db.Model(&model.EngineeringChangeOrder{}).
		Where("requires_regression_testing = ? AND status != 'CLOSED'", true).
		Count(&stats.RegressionQueueCount)
	r.db.Model(&model.EngineeringChangeOrder{}).
		Where("status = 'CLOSED'").
		Count(&stats.ClosedChanges)

	return stats, nil
}

// -----------------------------------------------------------------------------
// ECR OPERATIONS
// -----------------------------------------------------------------------------

type ECRFilter struct {
	Status             string
	Priority           string
	ProductID          string
	HardwareRevisionID string
	Origin             string
	Search             string
	Page               int
	Limit              int
}

func (r *ChangeRepository) GetECRs(filter ECRFilter) ([]model.EngineeringChangeRequest, int64, error) {
	var ecrs []model.EngineeringChangeRequest
	var total int64

	query := r.db.Model(&model.EngineeringChangeRequest{}).
		Preload("Product").
		Preload("ProductVersion").
		Preload("SourceHardwareRevision").
		Preload("ChangeItems").
		Preload("ECO")

	if filter.Status != "" && filter.Status != "ALL" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Priority != "" && filter.Priority != "ALL" {
		query = query.Where("priority = ?", filter.Priority)
	}
	if filter.ProductID != "" && filter.ProductID != "ALL" {
		query = query.Where("product_id = ?", filter.ProductID)
	}
	if filter.HardwareRevisionID != "" && filter.HardwareRevisionID != "ALL" {
		query = query.Where("source_hardware_revision_id = ?", filter.HardwareRevisionID)
	}
	if filter.Origin != "" && filter.Origin != "ALL" {
		query = query.Where("origin = ?", filter.Origin)
	}
	if filter.Search != "" {
		s := "%" + filter.Search + "%"
		query = query.Where("code LIKE ? OR title LIKE ? OR description LIKE ?", s, s, s)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&ecrs).Error
	return ecrs, total, err
}

func (r *ChangeRepository) GetECRByID(id string) (*model.EngineeringChangeRequest, error) {
	var ecr model.EngineeringChangeRequest
	err := r.db.Where("id = ? OR code = ?", id, id).
		Preload("Product").
		Preload("ProductVersion").
		Preload("SourceHardwareRevision").
		Preload("SourceFinding").
		Preload("ChangeItems.AffectedComponent").
		Preload("ImpactAnalyses").
		Preload("Reviews").
		Preload("Approvals").
		Preload("ECO").
		First(&ecr).Error
	if err != nil {
		return nil, err
	}
	return &ecr, nil
}

func (r *ChangeRepository) CreateECR(ecr *model.EngineeringChangeRequest) error {
	return r.db.Create(ecr).Error
}

func (r *ChangeRepository) UpdateECR(ecr *model.EngineeringChangeRequest) error {
	return r.db.Save(ecr).Error
}

func (r *ChangeRepository) DeleteECR(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.EngineeringChangeRequest{}).Error
}

// -----------------------------------------------------------------------------
// CHANGE ITEMS
// -----------------------------------------------------------------------------

func (r *ChangeRepository) AddChangeItem(item *model.ECRChangeItem) error {
	return r.db.Create(item).Error
}

func (r *ChangeRepository) UpdateChangeItem(item *model.ECRChangeItem) error {
	return r.db.Save(item).Error
}

func (r *ChangeRepository) DeleteChangeItem(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.ECRChangeItem{}).Error
}

func (r *ChangeRepository) GetChangeItemsByECR(ecrID string) ([]model.ECRChangeItem, error) {
	var items []model.ECRChangeItem
	err := r.db.Where("ecr_id = ?", ecrID).Preload("AffectedComponent").Order("sequence ASC").Find(&items).Error
	return items, err
}

// -----------------------------------------------------------------------------
// IMPACT ANALYSES
// -----------------------------------------------------------------------------

func (r *ChangeRepository) SaveImpactAnalysis(impact *model.ECRImpactAnalysis) error {
	var existing model.ECRImpactAnalysis
	err := r.db.Where("ecr_id = ? AND domain = ?", impact.ECRID, impact.Domain).First(&existing).Error
	if err == nil {
		impact.ID = existing.ID
		return r.db.Save(impact).Error
	}
	return r.db.Create(impact).Error
}

func (r *ChangeRepository) GetImpactsByECR(ecrID string) ([]model.ECRImpactAnalysis, error) {
	var impacts []model.ECRImpactAnalysis
	err := r.db.Where("ecr_id = ?", ecrID).Find(&impacts).Error
	return impacts, err
}

// -----------------------------------------------------------------------------
// REVIEWS
// -----------------------------------------------------------------------------

func (r *ChangeRepository) SaveReview(review *model.ECRReview) error {
	var existing model.ECRReview
	err := r.db.Where("ecr_id = ? AND discipline = ?", review.ECRID, review.Discipline).First(&existing).Error
	if err == nil {
		review.ID = existing.ID
		return r.db.Save(review).Error
	}
	return r.db.Create(review).Error
}

func (r *ChangeRepository) GetReviewsByECR(ecrID string) ([]model.ECRReview, error) {
	var reviews []model.ECRReview
	err := r.db.Where("ecr_id = ?", ecrID).Find(&reviews).Error
	return reviews, err
}

// -----------------------------------------------------------------------------
// APPROVALS
// -----------------------------------------------------------------------------

func (r *ChangeRepository) SaveApproval(approval *model.ECRApproval) error {
	var existing model.ECRApproval
	err := r.db.Where("ecr_id = ? AND stage = ?", approval.ECRID, approval.Stage).First(&existing).Error
	if err == nil {
		approval.ID = existing.ID
		return r.db.Save(approval).Error
	}
	return r.db.Create(approval).Error
}

func (r *ChangeRepository) GetApprovalsByECR(ecrID string) ([]model.ECRApproval, error) {
	var approvals []model.ECRApproval
	err := r.db.Where("ecr_id = ?", ecrID).Find(&approvals).Error
	return approvals, err
}

// -----------------------------------------------------------------------------
// ECO OPERATIONS
// -----------------------------------------------------------------------------

type ECOFilter struct {
	Status             string
	ProductID          string
	HardwareRevisionID string
	Search             string
	Page               int
	Limit              int
}

func (r *ChangeRepository) GetECOs(filter ECOFilter) ([]model.EngineeringChangeOrder, int64, error) {
	var ecos []model.EngineeringChangeOrder
	var total int64

	query := r.db.Model(&model.EngineeringChangeOrder{}).
		Preload("ECR.Product").
		Preload("SourceHardwareRevision").
		Preload("TargetHardwareRevision").
		Preload("SourceBOM").
		Preload("TargetBOM")

	if filter.Status != "" && filter.Status != "ALL" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.HardwareRevisionID != "" && filter.HardwareRevisionID != "ALL" {
		query = query.Where("source_hardware_revision_id = ? OR target_hardware_revision_id = ?", filter.HardwareRevisionID, filter.HardwareRevisionID)
	}
	if filter.Search != "" {
		s := "%" + filter.Search + "%"
		query = query.Where("code LIKE ? OR title LIKE ?", s, s)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&ecos).Error
	return ecos, total, err
}

func (r *ChangeRepository) GetECOByID(id string) (*model.EngineeringChangeOrder, error) {
	var eco model.EngineeringChangeOrder
	err := r.db.Where("id = ? OR code = ?", id, id).
		Preload("ECR.Product").
		Preload("ECR.ProductVersion").
		Preload("ECR.ChangeItems.AffectedComponent").
		Preload("ECR.ImpactAnalyses").
		Preload("SourceHardwareRevision").
		Preload("TargetHardwareRevision").
		Preload("SourceBOM.Items.Component").
		Preload("TargetBOM.Items.Component").
		First(&eco).Error
	if err != nil {
		return nil, err
	}
	return &eco, nil
}

func (r *ChangeRepository) CreateECO(eco *model.EngineeringChangeOrder) error {
	return r.db.Create(eco).Error
}

func (r *ChangeRepository) UpdateECO(eco *model.EngineeringChangeOrder) error {
	return r.db.Save(eco).Error
}

// -----------------------------------------------------------------------------
// TRACEABILITY
// -----------------------------------------------------------------------------

func (r *ChangeRepository) GetTraceabilityChains(productID string) ([]model.ChangeTraceabilityChain, error) {
	var chains []model.ChangeTraceabilityChain

	query := r.db.Model(&model.EngineeringChangeRequest{}).
		Preload("SourceFinding").
		Preload("SourceHardwareRevision").
		Preload("ECO.TargetHardwareRevision").
		Preload("ECO.TargetBOM")

	if productID != "" && productID != "ALL" {
		query = query.Where("product_id = ?", productID)
	}

	var ecrs []model.EngineeringChangeRequest
	if err := query.Order("created_at DESC").Find(&ecrs).Error; err != nil {
		return nil, err
	}

	for _, ecr := range ecrs {
		chain := model.ChangeTraceabilityChain{
			ECRID:              ecr.ID,
			ECRCode:            ecr.Code,
			ECRTitle:           ecr.Title,
			ECRStatus:          ecr.Status,
			SourceRevisionID:   ecr.SourceHardwareRevisionID,
			SourceRevisionCode: ecr.SourceHardwareRevisionID,
		}
		if ecr.SourceHardwareRevision != nil {
			chain.SourceRevisionCode = ecr.SourceHardwareRevision.Code
		}

		if ecr.SourceFinding != nil {
			chain.FindingID = &ecr.SourceFinding.ID
			chain.FindingCode = &ecr.SourceFinding.Code
			chain.FindingTitle = &ecr.SourceFinding.Title
		}

		if ecr.ECO != nil {
			chain.ECOID = &ecr.ECO.ID
			chain.ECOCode = &ecr.ECO.Code
			chain.ECOStatus = &ecr.ECO.Status

			if ecr.ECO.TargetHardwareRevision != nil {
				chain.TargetRevisionID = &ecr.ECO.TargetHardwareRevision.ID
				chain.TargetRevisionCode = &ecr.ECO.TargetHardwareRevision.Code

				// Look up downstream Prototype on target revision
				var proto model.PrototypeBuild
				if err := r.db.Where("hardware_revision_id = ?", ecr.ECO.TargetHardwareRevision.ID).First(&proto).Error; err == nil {
					chain.PrototypeID = &proto.ID
					chain.PrototypeCode = &proto.Code

					// Look up test session on prototype
					var sess model.PrototypeTestSession
					if err := r.db.Where("prototype_build_id = ?", proto.ID).First(&sess).Error; err == nil {
						chain.TestSessionID = &sess.ID
						chain.TestSessionCode = &sess.SessionCode

						// Look up validation decision on session
						var valDec model.PrototypeValidationDecision
						if err := r.db.Where("prototype_test_session_id = ? AND is_current = ?", sess.ID, true).First(&valDec).Error; err == nil {
							chain.ValidationDecisionID = &valDec.ID
							decStr := string(valDec.Decision)
							chain.ValidationDecisionState = &decStr
						}
					}
				}
			}

			if ecr.ECO.TargetBOM != nil {
				chain.TargetBOMID = &ecr.ECO.TargetBOM.ID
				chain.TargetBOMCode = &ecr.ECO.TargetBOM.BOMCode
			}

			chain.RegressionBuildRequired = ecr.ECO.RequiresRegressionTesting
			if ecr.ECO.Status == "VERIFIED" || ecr.ECO.Status == "CLOSED" || (chain.ValidationDecisionState != nil && *chain.ValidationDecisionState == "VALIDATED") {
				chain.ClosedLoopVerified = true
			}
		}

		chain.CCBApprovalCount = len(ecr.Approvals)

		chains = append(chains, chain)
	}

	return chains, nil
}
