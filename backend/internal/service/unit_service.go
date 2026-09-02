package service

import (
	"errors"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type UnitService struct {
	repo *repository.UnitRepository
}

func NewUnitService(repo *repository.UnitRepository) *UnitService {
	return &UnitService{repo: repo}
}

func (s *UnitService) GetAll(p repository.PaginationParams) ([]model.Unit, repository.PaginationMeta, error) {
	units, total, err := s.repo.FindAll(p)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return units, meta, nil
}

func (s *UnitService) GetByID(id string) (*model.Unit, error) {
	return s.repo.FindByID(id)
}

func (s *UnitService) Create(unit *model.Unit) (*model.Unit, error) {
	if unit.Name == "" || unit.Code == "" {
		return nil, errors.New("unit code and name are required")
	}

	if unit.ID == "" {
		unit.ID = fmt.Sprintf("UOM-%d", time.Now().Unix()%10000)
	}

	if unit.Status == "" {
		unit.Status = "Active"
	}

	unit.CreatedAt = time.Now()
	unit.UpdatedAt = time.Now()

	if err := s.repo.Create(unit); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *UnitService) Update(id string, unit *model.Unit) (*model.Unit, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("unit not found")
	}

	existing.Name = unit.Name
	existing.Code = unit.Code
	existing.Symbol = unit.Symbol
	existing.Category = unit.Category
	existing.Description = unit.Description
	existing.Status = unit.Status
	existing.IsDefault = unit.IsDefault
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *UnitService) Delete(id string) error {
	return s.repo.Delete(id)
}
