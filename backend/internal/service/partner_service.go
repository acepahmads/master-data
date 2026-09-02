package service

import (
	"errors"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type ManufacturerService struct {
	repo         *repository.ManufacturerRepository
	activityRepo *repository.ActivityRepository
}

func NewManufacturerService(repo *repository.ManufacturerRepository, activityRepo *repository.ActivityRepository) *ManufacturerService {
	return &ManufacturerService{repo: repo, activityRepo: activityRepo}
}

func (s *ManufacturerService) GetAll(p repository.PaginationParams, status string) ([]model.Manufacturer, repository.PaginationMeta, error) {
	mfgs, total, err := s.repo.FindAll(p, status)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return mfgs, meta, nil
}

func (s *ManufacturerService) GetByID(id string) (*model.Manufacturer, error) {
	return s.repo.FindByID(id)
}

func (s *ManufacturerService) Create(mfg *model.Manufacturer, currentUser *model.User) (*model.Manufacturer, error) {
	if mfg.Name == "" {
		return nil, errors.New("manufacturer name is required")
	}

	if mfg.ID == "" {
		mfg.ID = fmt.Sprintf("MFG-%d", time.Now().Unix()%10000)
	}

	if mfg.Status == "" {
		mfg.Status = "Approved"
	}

	mfg.CreatedAt = time.Now()
	mfg.UpdatedAt = time.Now()

	if err := s.repo.Create(mfg); err != nil {
		return nil, err
	}

	userName := "System"
	userID := ""
	userAvatar := ""
	if currentUser != nil {
		userName = currentUser.Name
		userID = currentUser.ID
		userAvatar = currentUser.Avatar
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "manufacturers",
		EntityType:  "Manufacturer",
		EntityID:    mfg.ID,
		EntityName:  mfg.Name,
		Action:      "Added Manufacturer",
		Description: fmt.Sprintf("Added semiconductor/hardware manufacturer '%s'", mfg.Name),
		BadgeColor:  "emerald",
		CreatedAt:   time.Now(),
	})

	return mfg, nil
}

func (s *ManufacturerService) Update(id string, mfg *model.Manufacturer, currentUser *model.User) (*model.Manufacturer, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("manufacturer not found")
	}

	existing.Name = mfg.Name
	existing.Code = mfg.Code
	existing.Country = mfg.Country
	existing.CountryCode = mfg.CountryCode
	existing.Website = mfg.Website
	existing.Status = mfg.Status
	existing.Rating = mfg.Rating
	existing.ContactPerson = mfg.ContactPerson
	existing.ContactEmail = mfg.ContactEmail
	existing.Description = mfg.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *ManufacturerService) Delete(id string) error {
	return s.repo.Delete(id)
}

// SUPPLIER SERVICE

type SupplierService struct {
	repo         *repository.SupplierRepository
	activityRepo *repository.ActivityRepository
}

func NewSupplierService(repo *repository.SupplierRepository, activityRepo *repository.ActivityRepository) *SupplierService {
	return &SupplierService{repo: repo, activityRepo: activityRepo}
}

func (s *SupplierService) GetAll(p repository.PaginationParams, status string) ([]model.Supplier, repository.PaginationMeta, error) {
	suppliers, total, err := s.repo.FindAll(p, status)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return suppliers, meta, nil
}

func (s *SupplierService) GetByID(id string) (*model.Supplier, error) {
	return s.repo.FindByID(id)
}

func (s *SupplierService) Create(supplier *model.Supplier, currentUser *model.User) (*model.Supplier, error) {
	if supplier.Name == "" {
		return nil, errors.New("supplier name is required")
	}

	if supplier.ID == "" {
		supplier.ID = fmt.Sprintf("SUP-%d", time.Now().Unix()%10000)
	}

	if supplier.Status == "" {
		supplier.Status = "Active"
	}

	supplier.CreatedAt = time.Now()
	supplier.UpdatedAt = time.Now()

	if err := s.repo.Create(supplier); err != nil {
		return nil, err
	}

	userName := "System"
	userID := ""
	userAvatar := ""
	if currentUser != nil {
		userName = currentUser.Name
		userID = currentUser.ID
		userAvatar = currentUser.Avatar
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "suppliers",
		EntityType:  "Supplier",
		EntityID:    supplier.ID,
		EntityName:  supplier.Name,
		Action:      "Added Supplier",
		Description: fmt.Sprintf("Added authorized supplier '%s'", supplier.Name),
		BadgeColor:  "indigo",
		CreatedAt:   time.Now(),
	})

	return supplier, nil
}

func (s *SupplierService) Update(id string, supplier *model.Supplier, currentUser *model.User) (*model.Supplier, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("supplier not found")
	}

	existing.Name = supplier.Name
	existing.Code = supplier.Code
	existing.Country = supplier.Country
	existing.CountryCode = supplier.CountryCode
	existing.Address = supplier.Address
	existing.ContactPerson = supplier.ContactPerson
	existing.Email = supplier.Email
	existing.Phone = supplier.Phone
	existing.Website = supplier.Website
	existing.Tier = supplier.Tier
	existing.Status = supplier.Status
	existing.Terms = supplier.Terms
	existing.LeadTimeAvg = supplier.LeadTimeAvg
	existing.Description = supplier.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *SupplierService) Delete(id string) error {
	return s.repo.Delete(id)
}
