package service

import (
	"errors"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type RoleService struct {
	repo         *repository.RoleRepository
	activityRepo *repository.ActivityRepository
}

func NewRoleService(repo *repository.RoleRepository, activityRepo *repository.ActivityRepository) *RoleService {
	return &RoleService{repo: repo, activityRepo: activityRepo}
}

func (s *RoleService) GetAll() ([]model.Role, error) {
	return s.repo.FindAll()
}

func (s *RoleService) GetByID(id string) (*model.Role, error) {
	return s.repo.FindByID(id)
}

func (s *RoleService) GetAllPermissions() ([]model.Permission, error) {
	return s.repo.FindAllPermissions()
}

func (s *RoleService) Create(role *model.Role) (*model.Role, error) {
	if role.Name == "" || role.Code == "" {
		return nil, errors.New("role name and code are required")
	}

	if role.ID == "" {
		role.ID = fmt.Sprintf("ROLE-%d", time.Now().Unix()%10000)
	}

	if role.Status == "" {
		role.Status = "Active"
	}

	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

	if err := s.repo.Create(role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) Update(id string, role *model.Role) (*model.Role, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("role not found")
	}

	existing.Name = role.Name
	existing.Code = role.Code
	existing.Description = role.Description
	existing.Status = role.Status
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *RoleService) UpdatePermissions(roleID string, permissionCodes []string, currentUser *model.User) error {
	role, err := s.repo.FindByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	if err := s.repo.UpdatePermissions(role.ID, permissionCodes); err != nil {
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
		Module:      "administration",
		EntityType:  "Role",
		EntityID:    role.ID,
		EntityName:  role.Name,
		Action:      "Modified RBAC Permissions",
		Description: fmt.Sprintf("Updated granular security permissions matrix for '%s'", role.Name),
		BadgeColor:  "amber",
		CreatedAt:   time.Now(),
	})

	return nil
}
