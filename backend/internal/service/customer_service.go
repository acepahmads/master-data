package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type CustomerService struct {
	repo         *repository.CustomerRepository
	activityRepo *repository.ActivityRepository
}

func NewCustomerService(repo *repository.CustomerRepository, activityRepo *repository.ActivityRepository) *CustomerService {
	return &CustomerService{repo: repo, activityRepo: activityRepo}
}

func (s *CustomerService) GetAll(p repository.PaginationParams, status, source string) ([]model.Customer, repository.PaginationMeta, error) {
	customers, total, err := s.repo.FindAll(p, status, source)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return customers, meta, nil
}

func (s *CustomerService) GetByID(id string) (*model.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *CustomerService) Create(c *model.Customer, currentUser *model.User) (*model.Customer, error) {
	if strings.TrimSpace(c.Name) == "" && strings.TrimSpace(c.Company) == "" {
		return nil, errors.New("customer name or company is required")
	}

	if c.Name == "" {
		c.Name = c.Company
	}

	count, _ := s.repo.Count()
	if c.Code == "" {
		c.Code = fmt.Sprintf("CUST-%d-%03d", time.Now().Year(), count+1)
	}

	if c.ID == "" {
		c.ID = c.Code
	}

	if c.Status == "" {
		c.Status = "Active"
	}

	if c.Source == "" {
		c.Source = "MANUAL"
	}

	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()

	if err := s.repo.Create(c); err != nil {
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

	if s.activityRepo != nil {
		_ = s.activityRepo.Create(&model.ActivityLog{
			ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
			UserID:      userID,
			UserName:    userName,
			UserAvatar:  userAvatar,
			Module:      "customers",
			EntityType:  "Customer",
			EntityID:    c.ID,
			EntityName:  c.Name,
			Action:      "Created Customer",
			Description: fmt.Sprintf("Added new customer '%s' (%s)", c.Name, c.Company),
			BadgeColor:  "indigo",
			CreatedAt:   time.Now(),
		})
	}

	return c, nil
}

func (s *CustomerService) Update(c *model.Customer, currentUser *model.User) (*model.Customer, error) {
	existing, err := s.repo.FindByID(c.ID)
	if err != nil {
		return nil, errors.New("customer not found")
	}

	if strings.TrimSpace(c.Name) != "" {
		existing.Name = strings.TrimSpace(c.Name)
	}
	existing.Company = strings.TrimSpace(c.Company)
	existing.Phone = strings.TrimSpace(c.Phone)
	existing.Email = strings.TrimSpace(c.Email)
	existing.Descr = strings.TrimSpace(c.Descr)
	if c.Status != "" {
		existing.Status = c.Status
	}
	existing.UpdatedAt = time.Now()

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

	if s.activityRepo != nil {
		_ = s.activityRepo.Create(&model.ActivityLog{
			ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
			UserID:      userID,
			UserName:    userName,
			UserAvatar:  userAvatar,
			Module:      "customers",
			EntityType:  "Customer",
			EntityID:    existing.ID,
			EntityName:  existing.Name,
			Action:      "Updated Customer",
			Description: fmt.Sprintf("Updated customer profile '%s'", existing.Name),
			BadgeColor:  "indigo",
			CreatedAt:   time.Now(),
		})
	}

	return existing, nil
}

func (s *CustomerService) Delete(id string, currentUser *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("customer not found")
	}

	if err := s.repo.Delete(id); err != nil {
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

	if s.activityRepo != nil {
		_ = s.activityRepo.Create(&model.ActivityLog{
			ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
			UserID:      userID,
			UserName:    userName,
			UserAvatar:  userAvatar,
			Module:      "customers",
			EntityType:  "Customer",
			EntityID:    existing.ID,
			EntityName:  existing.Name,
			Action:      "Deleted Customer",
			Description: fmt.Sprintf("Deleted customer profile '%s'", existing.Name),
			BadgeColor:  "rose",
			CreatedAt:   time.Now(),
		})
	}

	return nil
}

// UpsertFromImport is called during Excel import ingestion to automatically register or update customer records.
func (s *CustomerService) UpsertFromImport(name, company, phone, email, descr string) (*model.Customer, error) {
	name = strings.TrimSpace(name)
	company = strings.TrimSpace(company)
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)
	descr = strings.TrimSpace(descr)

	if name == "" && company == "" {
		return nil, nil // Nothing to upsert
	}

	if name == "" {
		name = company
	}

	existing, err := s.repo.FindByCompanyOrName(company, name)
	if err == nil && existing != nil {
		// Update missing contact details if new import has them
		changed := false
		if existing.Phone == "" && phone != "" {
			existing.Phone = phone
			changed = true
		}
		if existing.Email == "" && email != "" {
			existing.Email = email
			changed = true
		}
		if existing.Descr == "" && descr != "" {
			existing.Descr = descr
			changed = true
		}
		if changed {
			existing.UpdatedAt = time.Now()
			_ = s.repo.Update(existing)
		}
		return existing, nil
	}

	// Create new record
	count, _ := s.repo.Count()
	code := fmt.Sprintf("CUST-%d-%03d", time.Now().Year(), count+1)
	newCust := &model.Customer{
		ID:        code,
		Code:      code,
		Name:      name,
		Company:   company,
		Phone:     phone,
		Email:     email,
		Descr:     descr,
		Status:    "Active",
		Source:    "IMPORT_XLS",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(newCust); err != nil {
		return nil, err
	}
	return newCust, nil
}
