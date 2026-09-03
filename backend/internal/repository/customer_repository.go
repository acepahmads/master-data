package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindAll(p PaginationParams, status, source string) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64

	query := r.db.Model(&model.Customer{})
	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR company LIKE ? OR phone LIKE ? OR email LIKE ? OR descr LIKE ?", s, s, s, s, s, s)
	}
	if status != "" && status != "All" {
		query = query.Where("status = ?", status)
	}
	if source != "" && source != "All" {
		query = query.Where("source = ?", source)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order(p.Sort).Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepository) FindByID(id string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("id = ? OR code = ?", id, id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) FindByCompanyOrName(company, name string) (*model.Customer, error) {
	var customer model.Customer
	q := r.db.Model(&model.Customer{})
	if company != "" && name != "" {
		q = q.Where("company = ? AND name = ?", company, name)
	} else if company != "" {
		q = q.Where("company = ?", company)
	} else if name != "" {
		q = q.Where("name = ?", name)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	if err := q.First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) Create(c *model.Customer) error {
	return r.db.Create(c).Error
}

func (r *CustomerRepository) Update(c *model.Customer) error {
	return r.db.Save(c).Error
}

func (r *CustomerRepository) Delete(id string) error {
	return r.db.Where("id = ? OR code = ?", id, id).Delete(&model.Customer{}).Error
}

func (r *CustomerRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Customer{}).Count(&count).Error
	return count, err
}
