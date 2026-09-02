package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type RevisionRepository struct {
	db *gorm.DB
}

func NewRevisionRepository(db *gorm.DB) *RevisionRepository {
	return &RevisionRepository{db: db}
}

func (r *RevisionRepository) FindAll(p PaginationParams, versionID, productID, status string) ([]model.HardwareRevision, int64, error) {
	var revisions []model.HardwareRevision
	var total int64

	query := r.db.Model(&model.HardwareRevision{})

	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("code LIKE ? OR name LIKE ? OR product_name LIKE ? OR product_code LIKE ? OR pcb_revision LIKE ? OR change_summary LIKE ?", s, s, s, s, s, s)
	}
	if versionID != "" && versionID != "All" {
		query = query.Where("product_version_id = ?", versionID)
	}
	if productID != "" && productID != "All" {
		query = query.Where("product_code = ?", productID)
	}
	if status != "" && status != "All" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order(p.Sort).Find(&revisions).Error
	return revisions, total, err
}

func (r *RevisionRepository) FindByID(id string) (*model.HardwareRevision, error) {
	var revision model.HardwareRevision
	err := r.db.Where("id = ? OR code = ?", id, id).First(&revision).Error
	if err != nil {
		return nil, err
	}
	return &revision, nil
}

func (r *RevisionRepository) FindByVersionAndCode(versionID, code string) (*model.HardwareRevision, error) {
	var revision model.HardwareRevision
	err := r.db.Where("product_version_id = ? AND code = ?", versionID, code).First(&revision).Error
	if err != nil {
		return nil, err
	}
	return &revision, nil
}

func (r *RevisionRepository) FindByVersion(versionID string) ([]model.HardwareRevision, error) {
	var revisions []model.HardwareRevision
	err := r.db.Where("product_version_id = ?", versionID).Order("created_at ASC").Find(&revisions).Error
	return revisions, err
}

func (r *RevisionRepository) Create(revision *model.HardwareRevision) error {
	return r.db.Create(revision).Error
}

func (r *RevisionRepository) Update(revision *model.HardwareRevision) error {
	return r.db.Save(revision).Error
}

func (r *RevisionRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.HardwareRevision{}).Error
}

func (r *RevisionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.HardwareRevision{}).Count(&count).Error
	return count, err
}
