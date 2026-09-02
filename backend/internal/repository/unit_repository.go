package repository

import (
	"iot-rd-backend/internal/model"

	"gorm.io/gorm"
)

type UnitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) *UnitRepository {
	return &UnitRepository{db: db}
}

func (r *UnitRepository) FindAll(p PaginationParams) ([]model.Unit, int64, error) {
	var units []model.Unit
	var total int64

	query := r.db.Model(&model.Unit{})
	if p.Search != "" {
		s := "%" + p.Search + "%"
		query = query.Where("name LIKE ? OR code LIKE ? OR symbol LIKE ? OR category LIKE ?", s, s, s, s)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Scopes(Paginate(p)).Order("code asc").Find(&units).Error
	return units, total, err
}

func (r *UnitRepository) FindByID(id string) (*model.Unit, error) {
	var unit model.Unit
	err := r.db.Where("id = ? OR code = ?", id, id).First(&unit).Error
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

func (r *UnitRepository) Create(unit *model.Unit) error {
	return r.db.Create(unit).Error
}

func (r *UnitRepository) Update(unit *model.Unit) error {
	return r.db.Save(unit).Error
}

func (r *UnitRepository) Delete(id string) error {
	return r.db.Where("id = ? OR code = ?", id, id).Delete(&model.Unit{}).Error
}

func (r *UnitRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Unit{}).Count(&count).Error
	return count, err
}
