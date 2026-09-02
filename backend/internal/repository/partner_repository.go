package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type ManufacturerRepository struct {
	db *gorm.DB
}

func NewManufacturerRepository(db *gorm.DB) *ManufacturerRepository {
	return &ManufacturerRepository{db: db}
}

func (r *ManufacturerRepository) FindAll(p PaginationParams, status string) ([]model.Manufacturer, int64, error) {
	var mfgs []model.Manufacturer
	var total int64

	query := r.db.Model(&model.Manufacturer{})
	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR country LIKE ? OR contact_person LIKE ?", s, s, s, s)
	}
	if status != "" && status != "All" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order(p.Sort).Find(&mfgs).Error
	return mfgs, total, err
}

func (r *ManufacturerRepository) FindByID(id string) (*model.Manufacturer, error) {
	var mfg model.Manufacturer
	err := r.db.Where("id = ? OR code = ?", id, id).First(&mfg).Error
	if err != nil {
		return nil, err
	}
	return &mfg, nil
}

func (r *ManufacturerRepository) Create(mfg *model.Manufacturer) error {
	return r.db.Create(mfg).Error
}

func (r *ManufacturerRepository) Update(mfg *model.Manufacturer) error {
	return r.db.Save(mfg).Error
}

func (r *ManufacturerRepository) Delete(id string) error {
	return r.db.Where("id = ? OR code = ?", id, id).Delete(&model.Manufacturer{}).Error
}

func (r *ManufacturerRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Manufacturer{}).Count(&count).Error
	return count, err
}

// SUPPLIER REPOSITORY

type SupplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) *SupplierRepository {
	return &SupplierRepository{db: db}
}

func (r *SupplierRepository) FindAll(p PaginationParams, status string) ([]model.Supplier, int64, error) {
	var suppliers []model.Supplier
	var total int64

	query := r.db.Model(&model.Supplier{})
	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR contact_person LIKE ? OR email LIKE ?", s, s, s, s)
	}
	if status != "" && status != "All" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order(p.Sort).Find(&suppliers).Error
	return suppliers, total, err
}

func (r *SupplierRepository) FindByID(id string) (*model.Supplier, error) {
	var supplier model.Supplier
	err := r.db.Where("id = ? OR code = ?", id, id).First(&supplier).Error
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

func (r *SupplierRepository) Create(supplier *model.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *SupplierRepository) Update(supplier *model.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *SupplierRepository) Delete(id string) error {
	return r.db.Where("id = ? OR code = ?", id, id).Delete(&model.Supplier{}).Error
}

func (r *SupplierRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Supplier{}).Count(&count).Error
	return count, err
}
