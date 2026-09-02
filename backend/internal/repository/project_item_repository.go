package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type ProjectItemRepository struct {
	db *gorm.DB
}

func NewProjectItemRepository(db *gorm.DB) *ProjectItemRepository {
	return &ProjectItemRepository{db: db}
}

func (r *ProjectItemRepository) FindByProjectID(projectID string) ([]model.ProjectItem, error) {
	var items []model.ProjectItem
	err := r.db.Where("project_product_id = ?", projectID).Order("sort_order ASC, created_at ASC").Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ProjectItemRepository) FindByID(id string) (*model.ProjectItem, error) {
	var item model.ProjectItem
	err := r.db.Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ProjectItemRepository) Create(item *model.ProjectItem) error {
	return r.db.Create(item).Error
}

func (r *ProjectItemRepository) Update(item *model.ProjectItem) error {
	return r.db.Save(item).Error
}

func (r *ProjectItemRepository) Delete(id string) error {
	// First delete any sub-items
	_ = r.db.Where("parent_item_id = ?", id).Delete(&model.ProjectItem{}).Error
	return r.db.Where("id = ?", id).Delete(&model.ProjectItem{}).Error
}

func (r *ProjectItemRepository) CountByProjectID(projectID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.ProjectItem{}).Where("project_product_id = ?", projectID).Count(&count).Error
	return count, err
}

func (r *ProjectItemRepository) Reorder(projectID string, itemIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range itemIDs {
			if err := tx.Model(&model.ProjectItem{}).Where("id = ? AND project_product_id = ?", id, projectID).Update("sort_order", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
