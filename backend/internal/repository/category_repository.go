package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) FindAllTree() ([]model.ComponentCategory, error) {
	var roots []model.ComponentCategory
	err := r.db.Where("parent_id IS NULL").
		Preload("Children.Children").
		Order("code asc").
		Find(&roots).Error
	return roots, err
}

func (r *CategoryRepository) FindAllFlat() ([]model.ComponentCategory, error) {
	var categories []model.ComponentCategory
	err := r.db.Order("name asc").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByID(id string) (*model.ComponentCategory, error) {
	var cat model.ComponentCategory
	err := r.db.Where("id = ? OR code = ?", id, id).
		Preload("Parent").
		Preload("Children").
		First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *CategoryRepository) Create(cat *model.ComponentCategory) error {
	return r.db.Create(cat).Error
}

func (r *CategoryRepository) Update(cat *model.ComponentCategory) error {
	return r.db.Save(cat).Error
}

func (r *CategoryRepository) Delete(id string) error {
	return r.db.Where("id = ? OR code = ?", id, id).Delete(&model.ComponentCategory{}).Error
}

func (r *CategoryRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.ComponentCategory{}).Count(&count).Error
	return count, err
}
