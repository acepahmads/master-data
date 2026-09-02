package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type VersionRepository struct {
	db *gorm.DB
}

func NewVersionRepository(db *gorm.DB) *VersionRepository {
	return &VersionRepository{db: db}
}

func (r *VersionRepository) FindAll(p PaginationParams, productID, status, owner string) ([]model.ProductVersion, int64, error) {
	var versions []model.ProductVersion
	var total int64

	query := r.db.Model(&model.ProductVersion{})

	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("version_number LIKE ? OR version_name LIKE ? OR product_name LIKE ? OR product_code LIKE ? OR owner LIKE ?", s, s, s, s, s)
	}
	if productID != "" && productID != "All" {
		query = query.Where("product_id = ? OR product_code = ?", productID, productID)
	}
	if status != "" && status != "All" {
		query = query.Where("status = ?", status)
	}
	if owner != "" {
		query = query.Where("owner = ?", owner)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order(p.Sort).Find(&versions).Error
	return versions, total, err
}

func (r *VersionRepository) FindByID(id string) (*model.ProductVersion, error) {
	var version model.ProductVersion
	err := r.db.Where("id = ?", id).Preload("HardwareRevisions").First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *VersionRepository) FindByProductAndVersion(productID, versionNumber string) (*model.ProductVersion, error) {
	var version model.ProductVersion
	err := r.db.Where("(product_id = ? OR product_code = ?) AND version_number = ?", productID, productID, versionNumber).First(&version).Error
	if err != nil {
		return nil, err
	}
	return &version, nil
}

func (r *VersionRepository) FindByProduct(productIDOrCode string) ([]model.ProductVersion, error) {
	var versions []model.ProductVersion
	err := r.db.Where("product_id = ? OR product_code = ?", productIDOrCode, productIDOrCode).
		Preload("HardwareRevisions").
		Order("version_number ASC").
		Find(&versions).Error
	return versions, err
}

func (r *VersionRepository) Create(version *model.ProductVersion) error {
	return r.db.Create(version).Error
}

func (r *VersionRepository) Update(version *model.ProductVersion) error {
	return r.db.Save(version).Error
}

func (r *VersionRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.ProductVersion{}).Error
}

func (r *VersionRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.ProductVersion{}).Count(&count).Error
	return count, err
}
