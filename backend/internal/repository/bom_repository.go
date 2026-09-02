package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type BOMRepository struct {
	db *gorm.DB
}

func NewBOMRepository(db *gorm.DB) *BOMRepository {
	return &BOMRepository{db: db}
}

func (r *BOMRepository) FindByRevisionID(revisionID string) (*model.EngineeringBOM, error) {
	var bom model.EngineeringBOM
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("bom_items.position ASC, bom_items.created_at ASC")
	}).Preload("Items.Component").Where("hardware_revision_id = ?", revisionID).First(&bom).Error
	if err != nil {
		return nil, err
	}
	return &bom, nil
}

func (r *BOMRepository) FindByID(id string) (*model.EngineeringBOM, error) {
	var bom model.EngineeringBOM
	err := r.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("bom_items.position ASC, bom_items.created_at ASC")
	}).Preload("Items.Component").Where("id = ? OR bom_code = ?", id, id).First(&bom).Error
	if err != nil {
		return nil, err
	}
	return &bom, nil
}

func (r *BOMRepository) Create(bom *model.EngineeringBOM) error {
	return r.db.Create(bom).Error
}

func (r *BOMRepository) Update(bom *model.EngineeringBOM) error {
	return r.db.Save(bom).Error
}

func (r *BOMRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.EngineeringBOM{}).Error
}

func (r *BOMRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.EngineeringBOM{}).Count(&count).Error
	return count, err
}

// BOM Item Methods
func (r *BOMRepository) FindItemsByBOMID(bomID string) ([]model.BOMItem, error) {
	var items []model.BOMItem
	err := r.db.Preload("Component").Where("engineering_bom_id = ?", bomID).Order("position ASC, created_at ASC").Find(&items).Error
	return items, err
}

func (r *BOMRepository) FindItemByID(id string) (*model.BOMItem, error) {
	var item model.BOMItem
	err := r.db.Preload("Component").Where("id = ?", id).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *BOMRepository) CreateItem(item *model.BOMItem) error {
	return r.db.Create(item).Error
}

func (r *BOMRepository) UpdateItem(item *model.BOMItem) error {
	return r.db.Save(item).Error
}

func (r *BOMRepository) DeleteItem(id string) error {
	// Also delete child items
	_ = r.db.Where("parent_item_id = ?", id).Delete(&model.BOMItem{}).Error
	return r.db.Where("id = ?", id).Delete(&model.BOMItem{}).Error
}

func (r *BOMRepository) ReorderItems(bomID string, itemIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range itemIDs {
			if err := tx.Model(&model.BOMItem{}).Where("id = ? AND engineering_bom_id = ?", id, bomID).Update("position", i+1).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
