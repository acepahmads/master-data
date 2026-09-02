package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) FindAll(p PaginationParams, category, status, productType string) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{})

	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("code LIKE ? OR name LIKE ? OR description LIKE ? OR target_market LIKE ?", s, s, s, s)
	}
	if category != "" && category != "All" {
		query = query.Where("category = ?", category)
	}
	if status != "" && status != "All" {
		query = query.Where("status = ?", status)
	}
	if productType != "" && productType != "All" {
		query = query.Where("product_type = ?", productType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("TradingDetail").Scopes(Paginate(p)).Order(p.Sort).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("TradingDetail").Where("id = ? OR code = ?", id, id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id string) error {
	return r.db.Where("id = ? OR code = ?", id, id).Delete(&model.Product{}).Error
}

func (r *ProductRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Product{}).Count(&count).Error
	return count, err
}
