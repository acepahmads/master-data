package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type RevisionService struct {
	repo         *repository.RevisionRepository
	verRepo      *repository.VersionRepository
	activityRepo *repository.ActivityRepository
}

func NewRevisionService(
	repo *repository.RevisionRepository,
	verRepo *repository.VersionRepository,
	activityRepo *repository.ActivityRepository,
) *RevisionService {
	return &RevisionService{
		repo:         repo,
		verRepo:      verRepo,
		activityRepo: activityRepo,
	}
}

func (s *RevisionService) GetAll(p repository.PaginationParams, versionID, productID, status string) ([]model.HardwareRevision, repository.PaginationMeta, error) {
	revisions, total, err := s.repo.FindAll(p, versionID, productID, status)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return revisions, meta, nil
}

func (s *RevisionService) GetByID(id string) (*model.HardwareRevision, error) {
	return s.repo.FindByID(id)
}

func (s *RevisionService) Create(r *model.HardwareRevision, currentUser *model.User) (*model.HardwareRevision, error) {
	if r.ProductVersionID == "" {
		return nil, errors.New("parent product version ID is required")
	}
	if r.Code == "" {
		return nil, errors.New("revision code is required (e.g. REV-B)")
	}
	if r.Name == "" {
		return nil, errors.New("revision name is required")
	}

	r.Code = strings.ToUpper(strings.TrimSpace(r.Code))
	r.Status = normalizeRevisionStatus(r.Status)

	// Validate parent version exists
	ver, err := s.verRepo.FindByID(r.ProductVersionID)
	if err != nil {
		return nil, errors.New("associated product version not found")
	}

	r.VersionNumber = ver.VersionNumber
	r.ProductCode = ver.ProductCode
	r.ProductName = ver.ProductName

	// Check uniqueness
	existing, _ := s.repo.FindByVersionAndCode(r.ProductVersionID, r.Code)
	if existing != nil {
		return nil, fmt.Errorf("revision code '%s' already exists for version v%s", r.Code, ver.VersionNumber)
	}

	// Generate ID
	if r.ID == "" {
		r.ID = fmt.Sprintf("REV-%d", time.Now().UnixNano()%1000000)
	}

	if r.Engineer == "" {
		if currentUser != nil {
			r.Engineer = currentUser.Name
		} else {
			r.Engineer = "Alex Chen"
		}
	}

	if r.PCBRevision == "" {
		r.PCBRevision = fmt.Sprintf("PCB-v%s-%s", ver.VersionNumber, strings.ReplaceAll(r.Code, "REV-", ""))
	}
	if r.SchematicRevision == "" {
		r.SchematicRevision = fmt.Sprintf("SCH-v%s-01", ver.VersionNumber)
	}

	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()

	if err := s.repo.Create(r); err != nil {
		return nil, err
	}

	// Update parent version current revision & count
	ver.CurrentRevision = r.Code
	ver.RevisionsCount++
	_ = s.verRepo.Update(ver)

	// Activity log
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
		Module:      "revisions",
		EntityType:  "HardwareRevision",
		EntityID:    r.ID,
		EntityName:  fmt.Sprintf("%s %s (v%s)", r.ProductCode, r.Code, r.VersionNumber),
		Action:      "HARDWARE_REVISION_CREATED",
		Description: fmt.Sprintf("Created hardware revision '%s' for '%s v%s'", r.Code, r.ProductCode, r.VersionNumber),
		BadgeColor:  "emerald",
		CreatedAt:   time.Now(),
	})

	return r, nil
}

func (s *RevisionService) Update(id string, r *model.HardwareRevision, currentUser *model.User) (*model.HardwareRevision, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("hardware revision not found")
	}

	if r.Name != "" {
		existing.Name = r.Name
	}
	if r.Status != "" {
		existing.Status = normalizeRevisionStatus(r.Status)
	}
	if r.PCBRevision != "" {
		existing.PCBRevision = r.PCBRevision
	}
	if r.SchematicRevision != "" {
		existing.SchematicRevision = r.SchematicRevision
	}
	if r.Engineer != "" {
		existing.Engineer = r.Engineer
	}
	if r.ApprovedBy != nil {
		existing.ApprovedBy = r.ApprovedBy
	}
	if r.ReleaseDate != nil {
		existing.ReleaseDate = r.ReleaseDate
	}
	if r.ChangeSummary != "" {
		existing.ChangeSummary = r.ChangeSummary
	}
	if r.ChangesList != "" {
		existing.ChangesList = r.ChangesList
	}
	if r.Documents != "" {
		existing.Documents = r.Documents
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

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "revisions",
		EntityType:  "HardwareRevision",
		EntityID:    existing.ID,
		EntityName:  fmt.Sprintf("%s %s", existing.ProductCode, existing.Code),
		Action:      "HARDWARE_REVISION_UPDATED",
		Description: fmt.Sprintf("Updated hardware revision '%s' for '%s'", existing.Code, existing.ProductCode),
		BadgeColor:  "blue",
		CreatedAt:   time.Now(),
	})

	return existing, nil
}

func (s *RevisionService) Delete(id string, currentUser *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("hardware revision not found")
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
		Module:      "revisions",
		EntityType:  "HardwareRevision",
		EntityID:    existing.ID,
		EntityName:  fmt.Sprintf("%s %s", existing.ProductCode, existing.Code),
		Action:      "HARDWARE_REVISION_DELETED",
		Description: fmt.Sprintf("Removed hardware revision '%s' from '%s'", existing.Code, existing.ProductCode),
		BadgeColor:  "rose",
		CreatedAt:   time.Now(),
	})

	return nil
}

func normalizeRevisionStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "released", "production":
		return "Released"
	case "approved":
		return "Approved"
	case "in review", "in_review", "review", "pending":
		return "In Review"
	case "draft", "new":
		return "Draft"
	case "superseded":
		return "Superseded"
	case "obsolete", "deprecated":
		return "Obsolete"
	default:
		return "Draft"
	}
}
