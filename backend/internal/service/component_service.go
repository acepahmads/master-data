package service

import (
	"errors"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"

	"github.com/google/uuid"
)

type ComponentService struct {
	repo         *repository.ComponentRepository
	activityRepo *repository.ActivityRepository
}

func NewComponentService(repo *repository.ComponentRepository, activityRepo *repository.ActivityRepository) *ComponentService {
	return &ComponentService{repo: repo, activityRepo: activityRepo}
}

func (s *ComponentService) GetAll(p repository.PaginationParams, filters repository.ComponentFilters) ([]model.Component, repository.PaginationMeta, error) {
	components, total, err := s.repo.FindAll(p, filters)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return components, meta, nil
}

func (s *ComponentService) GetByID(id string) (*model.Component, error) {
	return s.repo.FindByID(id)
}

func (s *ComponentService) Create(c *model.Component, currentUser *model.User) (*model.Component, error) {
	if c.PartNumber == "" || c.Name == "" {
		return nil, errors.New("part number and name are required")
	}

	if c.ID == "" {
		count, _ := s.repo.Count()
		c.ID = fmt.Sprintf("CMP-%03d", count+1)
	}

	if c.InternalCode == "" {
		c.InternalCode = fmt.Sprintf("IPN-ENG-%03d", time.Now().Unix()%1000)
	}

	if c.Status == "" {
		c.Status = "Active"
	}

	if c.Image == "" {
		c.Image = "https://images.unsplash.com/photo-1518770660439-4636190af475?w=300&auto=format&fit=crop&q=80"
	}

	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	// Assign IDs to specs
	for i := range c.Specs {
		if c.Specs[i].ID == "" {
			c.Specs[i].ID = fmt.Sprintf("SPEC-%s", uuid.New().String()[:8])
		}
		c.Specs[i].ComponentID = c.ID
		c.Specs[i].DisplayOrder = i + 1
		c.Specs[i].CreatedAt = time.Now()
		c.Specs[i].UpdatedAt = time.Now()
	}

	if err := s.repo.Create(c); err != nil {
		return nil, err
	}

	// Auto Log Activity
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
		Module:      "components",
		EntityType:  "Component",
		EntityID:    c.ID,
		EntityName:  c.PartNumber,
		Action:      "Registered New Component",
		Description: fmt.Sprintf("Registered component '%s' (%s)", c.PartNumber, c.Name),
		BadgeColor:  "emerald",
		CreatedAt:   time.Now(),
	})

	return s.repo.FindByID(c.ID)
}

func (s *ComponentService) Update(id string, c *model.Component, currentUser *model.User) (*model.Component, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("component not found")
	}

	existing.PartNumber = c.PartNumber
	existing.InternalCode = c.InternalCode
	existing.Name = c.Name
	existing.CategoryID = c.CategoryID
	existing.ManufacturerID = c.ManufacturerID
	existing.SupplierID = c.SupplierID
	existing.UnitID = c.UnitID
	existing.PackageType = c.PackageType
	existing.MountingType = c.MountingType
	existing.LeadTime = c.LeadTime
	existing.RoHSCompliant = c.RoHSCompliant
	existing.REACHCompliant = c.REACHCompliant
	existing.EstimatedUnitCost = c.EstimatedUnitCost
	existing.Currency = c.Currency
	existing.PurchaseLink = c.PurchaseLink
	existing.ShortDescription = c.ShortDescription
	existing.DetailedDescription = c.DetailedDescription
	existing.Status = c.Status
	if c.Image != "" {
		existing.Image = c.Image
	}
	existing.UpdatedAt = time.Now()

	// Update specs list
	for i := range c.Specs {
		if c.Specs[i].ID == "" {
			c.Specs[i].ID = fmt.Sprintf("SPEC-%s", uuid.New().String()[:8])
		}
		c.Specs[i].ComponentID = existing.ID
		c.Specs[i].DisplayOrder = i + 1
		c.Specs[i].CreatedAt = time.Now()
		c.Specs[i].UpdatedAt = time.Now()
	}
	existing.Specs = c.Specs

	if err := s.repo.Update(existing); err != nil {
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
		Module:      "components",
		EntityType:  "Component",
		EntityID:    existing.ID,
		EntityName:  existing.PartNumber,
		Action:      "Updated Component Specs",
		Description: fmt.Sprintf("Updated specifications and attributes for '%s'", existing.PartNumber),
		BadgeColor:  "blue",
		CreatedAt:   time.Now(),
	})

	return s.repo.FindByID(existing.ID)
}

func (s *ComponentService) Delete(id string, currentUser *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("component not found")
	}

	if err := s.repo.Delete(existing.ID); err != nil {
		return err
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
		Module:      "components",
		EntityType:  "Component",
		EntityID:    existing.ID,
		EntityName:  existing.PartNumber,
		Action:      "Deleted Component",
		Description: fmt.Sprintf("Removed component '%s'", existing.PartNumber),
		BadgeColor:  "rose",
		CreatedAt:   time.Now(),
	})

	return nil
}

func (s *ComponentService) AddDocument(componentID string, doc *model.ComponentDocument, currentUser *model.User) (*model.ComponentDocument, error) {
	comp, err := s.repo.FindByID(componentID)
	if err != nil {
		return nil, errors.New("component not found")
	}

	doc.ID = fmt.Sprintf("DOC-%d", time.Now().UnixNano())
	doc.ComponentID = comp.ID
	doc.CreatedAt = time.Now()

	if err := s.repo.AddDocument(doc); err != nil {
		return nil, err
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      currentUser.ID,
		UserName:    currentUser.Name,
		UserAvatar:  currentUser.Avatar,
		Module:      "components",
		EntityType:  "Document",
		EntityID:    doc.ID,
		EntityName:  doc.FileName,
		Action:      "Uploaded Component Document",
		Description: fmt.Sprintf("Attached %s '%s' to component '%s'", doc.DocumentType, doc.FileName, comp.PartNumber),
		BadgeColor:  "amber",
		CreatedAt:   time.Now(),
	})

	return doc, nil
}

func (s *ComponentService) DeleteDocument(docID string) error {
	return s.repo.DeleteDocument(docID)
}
