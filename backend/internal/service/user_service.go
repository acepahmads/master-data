package service

import (
	"errors"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo         *repository.UserRepository
	activityRepo *repository.ActivityRepository
}

func NewUserService(repo *repository.UserRepository, activityRepo *repository.ActivityRepository) *UserService {
	return &UserService{repo: repo, activityRepo: activityRepo}
}

func (s *UserService) GetAll(p repository.PaginationParams) ([]model.User, repository.PaginationMeta, error) {
	users, total, err := s.repo.FindAll(p)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return users, meta, nil
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) Create(u *model.User, password string) (*model.User, error) {
	if u.Email == "" || u.Name == "" {
		return nil, errors.New("name and email are required")
	}

	if password == "" {
		password = "ChangeMe123!"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.PasswordHash = string(hash)

	if u.ID == "" {
		count, _ := s.repo.Count()
		u.ID = fmt.Sprintf("USR-%03d", count+1)
	}

	if u.Username == "" {
		u.Username = u.Email
	}

	if u.Status == "" {
		u.Status = "Active"
	}

	if u.Avatar == "" {
		u.Avatar = "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80"
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	if err := s.repo.Create(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) Update(id string, u *model.User) (*model.User, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	existing.Name = u.Name
	existing.Email = u.Email
	existing.RoleID = u.RoleID
	existing.Department = u.Department
	existing.Status = u.Status
	existing.TwoFactor = u.TwoFactor
	if u.Avatar != "" {
		existing.Avatar = u.Avatar
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *UserService) ToggleStatus(id string) (*model.User, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if existing.Status == "Active" {
		existing.Status = "Suspended"
	} else {
		existing.Status = "Active"
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

func (s *UserService) ResetPassword(id string, newPassword string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if newPassword == "" {
		newPassword = "ChangeMe123!"
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	existing.PasswordHash = string(hash)
	existing.UpdatedAt = time.Now()

	return s.repo.Update(existing)
}

func (s *UserService) Delete(id string) error {
	return s.repo.Delete(id)
}
