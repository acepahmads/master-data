package repository

import (
	"strings"

	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type PrototypeRepository struct {
	db *gorm.DB
}

func NewPrototypeRepository(db *gorm.DB) *PrototypeRepository {
	return &PrototypeRepository{db: db}
}

func (r *PrototypeRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *PrototypeRepository) FindAll(p PaginationParams, filters map[string]string) ([]model.PrototypeBuild, int64, error) {
	var prototypes []model.PrototypeBuild
	var total int64

	query := r.db.Model(&model.PrototypeBuild{}).
		Preload("ComponentPreparations").
		Preload("AssemblyStages", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence asc")
		}).
		Preload("EngineeringNotes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		})

	if productID, ok := filters["productId"]; ok && productID != "" {
		query = query.Where("product_id = ? OR product_code = ?", productID, productID)
	}
	if versionID, ok := filters["versionId"]; ok && versionID != "" {
		query = query.Where("version_id = ? OR version_number = ?", versionID, versionID)
	}
	if revisionID, ok := filters["revisionId"]; ok && revisionID != "" {
		query = query.Where("hardware_revision_id = ? OR hardware_revision_code = ?", revisionID, revisionID)
	}
	if status, ok := filters["status"]; ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if owner, ok := filters["owner"]; ok && owner != "" {
		query = query.Where("engineering_owner LIKE ?", "%"+owner+"%")
	}
	if search, ok := filters["search"]; ok && search != "" {
		s := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(code) LIKE ? OR LOWER(name) LIKE ? OR LOWER(purpose) LIKE ? OR LOWER(engineering_owner) LIKE ?", s, s, s, s)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order("created_at desc").Find(&prototypes).Error
	return prototypes, total, err
}

func (r *PrototypeRepository) FindByID(id string) (*model.PrototypeBuild, error) {
	var proto model.PrototypeBuild
	err := r.db.
		Preload("ComponentPreparations").
		Preload("AssemblyStages", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence asc")
		}).
		Preload("EngineeringNotes", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).
		Where("id = ? OR code = ?", id, id).
		First(&proto).Error
	if err != nil {
		return nil, err
	}
	return &proto, nil
}

func (r *PrototypeRepository) FindByCode(code string) (*model.PrototypeBuild, error) {
	var proto model.PrototypeBuild
	err := r.db.Where("code = ?", code).First(&proto).Error
	if err != nil {
		return nil, err
	}
	return &proto, nil
}

func (r *PrototypeRepository) FindByRevision(revisionID string) ([]model.PrototypeBuild, error) {
	var prototypes []model.PrototypeBuild
	err := r.db.
		Preload("ComponentPreparations").
		Preload("AssemblyStages", func(db *gorm.DB) *gorm.DB {
			return db.Order("sequence asc")
		}).
		Preload("EngineeringNotes").
		Where("hardware_revision_id = ? OR hardware_revision_code = ?", revisionID, revisionID).
		Order("build_number asc, created_at desc").
		Find(&prototypes).Error
	return prototypes, err
}

func (r *PrototypeRepository) Create(proto *model.PrototypeBuild) error {
	return r.db.Create(proto).Error
}

func (r *PrototypeRepository) Update(proto *model.PrototypeBuild) error {
	return r.db.Save(proto).Error
}

func (r *PrototypeRepository) Delete(id string) error {
	_ = r.db.Unscoped().Where("prototype_build_id = ? OR prototype_build_id IN (SELECT id FROM prototype_builds WHERE code = ?)", id, id).Delete(&model.PrototypeComponentPreparation{}).Error
	_ = r.db.Unscoped().Where("prototype_build_id = ? OR prototype_build_id IN (SELECT id FROM prototype_builds WHERE code = ?)", id, id).Delete(&model.PrototypeAssemblyStage{}).Error
	_ = r.db.Unscoped().Where("prototype_build_id = ? OR prototype_build_id IN (SELECT id FROM prototype_builds WHERE code = ?)", id, id).Delete(&model.PrototypeEngineeringNote{}).Error
	return r.db.Unscoped().Where("id = ? OR code = ?", id, id).Delete(&model.PrototypeBuild{}).Error
}

func (r *PrototypeRepository) GetDashboardStats() (*model.PrototypeDashboardStats, error) {
	var stats model.PrototypeDashboardStats

	r.db.Model(&model.PrototypeBuild{}).Count(&stats.TotalBuilds)
	r.db.Model(&model.PrototypeBuild{}).Where("status = ?", model.PrototypeStatusPlanned).Count(&stats.PlannedCount)
	r.db.Model(&model.PrototypeBuild{}).Where("status = ?", model.PrototypeStatusPreparing).Count(&stats.ComponentPrepCount)
	r.db.Model(&model.PrototypeBuild{}).Where("status IN ?", []string{model.PrototypeStatusInProgress, model.PrototypeStatusAssembly}).Count(&stats.InAssemblyCount)
	r.db.Model(&model.PrototypeBuild{}).Where("status = ?", model.PrototypeStatusPowerOn).Count(&stats.PowerOnCount)
	r.db.Model(&model.PrototypeBuild{}).Where("status = ?", "READY_FOR_TESTING").Count(&stats.ReadyForTestingCount)
	r.db.Model(&model.PrototypeBuild{}).Where("status = ?", model.PrototypeStatusCompleted).Count(&stats.CompletedCount)
	r.db.Model(&model.PrototypeBuild{}).Where("status = ?", model.PrototypeStatusOnHold).Count(&stats.OnHoldCount)

	return &stats, nil
}

// Component Preparation repository methods
func (r *PrototypeRepository) FindComponentPreparations(prototypeBuildID string) ([]model.PrototypeComponentPreparation, error) {
	var items []model.PrototypeComponentPreparation
	err := r.db.Where("prototype_build_id = ?", prototypeBuildID).Order("created_at asc").Find(&items).Error
	return items, err
}

func (r *PrototypeRepository) FindComponentPreparationByID(id string) (*model.PrototypeComponentPreparation, error) {
	var item model.PrototypeComponentPreparation
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PrototypeRepository) UpdateComponentPreparation(comp *model.PrototypeComponentPreparation) error {
	return r.db.Save(comp).Error
}

// Assembly Stages repository methods
func (r *PrototypeRepository) FindAssemblyStages(prototypeBuildID string) ([]model.PrototypeAssemblyStage, error) {
	var stages []model.PrototypeAssemblyStage
	err := r.db.Where("prototype_build_id = ?", prototypeBuildID).Order("sequence asc").Find(&stages).Error
	return stages, err
}

func (r *PrototypeRepository) FindAssemblyStageByID(id string) (*model.PrototypeAssemblyStage, error) {
	var stage model.PrototypeAssemblyStage
	err := r.db.Where("id = ?", id).First(&stage).Error
	if err != nil {
		return nil, err
	}
	return &stage, nil
}

func (r *PrototypeRepository) UpdateAssemblyStage(stage *model.PrototypeAssemblyStage) error {
	return r.db.Save(stage).Error
}

// Engineering Notes repository methods
func (r *PrototypeRepository) FindEngineeringNotes(prototypeBuildID string) ([]model.PrototypeEngineeringNote, error) {
	var notes []model.PrototypeEngineeringNote
	err := r.db.Where("prototype_build_id = ?", prototypeBuildID).Order("created_at desc").Find(&notes).Error
	return notes, err
}

func (r *PrototypeRepository) FindEngineeringNoteByID(id string) (*model.PrototypeEngineeringNote, error) {
	var note model.PrototypeEngineeringNote
	err := r.db.Where("id = ?", id).First(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *PrototypeRepository) CreateEngineeringNote(note *model.PrototypeEngineeringNote) error {
	return r.db.Create(note).Error
}

func (r *PrototypeRepository) UpdateEngineeringNote(note *model.PrototypeEngineeringNote) error {
	return r.db.Save(note).Error
}

func (r *PrototypeRepository) DeleteEngineeringNote(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.PrototypeEngineeringNote{}).Error
}
