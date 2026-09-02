package service

import (
	"errors"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type CategoryService struct {
	repo         *repository.CategoryRepository
	activityRepo *repository.ActivityRepository
}

func NewCategoryService(repo *repository.CategoryRepository, activityRepo *repository.ActivityRepository) *CategoryService {
	return &CategoryService{repo: repo, activityRepo: activityRepo}
}

func (s *CategoryService) GetTree() ([]model.ComponentCategory, error) {
	return s.repo.FindAllTree()
}

func (s *CategoryService) GetFlat() ([]model.ComponentCategory, error) {
	return s.repo.FindAllFlat()
}

func (s *CategoryService) GetByID(id string) (*model.ComponentCategory, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) Create(cat *model.ComponentCategory, currentUser *model.User) (*model.ComponentCategory, error) {
	if cat.Name == "" {
		return nil, errors.New("category name is required")
	}

	if cat.Code == "" {
		count, _ := s.repo.Count()
		cat.Code = fmt.Sprintf("CAT-%03d", count+1)
	}

	if cat.ID == "" {
		cat.ID = fmt.Sprintf("CAT-%s", cat.Code)
	}

	if cat.Status == "" {
		cat.Status = "Active"
	}

	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()

	if err := s.repo.Create(cat); err != nil {
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
		Module:      "categories",
		EntityType:  "Category",
		EntityID:    cat.ID,
		EntityName:  cat.Name,
		Action:      "Created Category Node",
		Description: fmt.Sprintf("Created category branch '%s' (%s)", cat.Name, cat.Code),
		BadgeColor:  "purple",
		CreatedAt:   time.Now(),
	})

	return cat, nil
}

func (s *CategoryService) Update(id string, cat *model.ComponentCategory, currentUser *model.User) (*model.ComponentCategory, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	existing.Name = cat.Name
	existing.Code = cat.Code
	existing.Description = cat.Description
	existing.Status = cat.Status
	existing.ParentID = cat.ParentID
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *CategoryService) Delete(id string) error {
	return s.repo.Delete(id)
}
