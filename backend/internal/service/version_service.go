package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type VersionService struct {
	repo         *repository.VersionRepository
	revRepo      *repository.RevisionRepository
	productRepo  *repository.ProductRepository
	activityRepo *repository.ActivityRepository
}

func NewVersionService(
	repo *repository.VersionRepository,
	revRepo *repository.RevisionRepository,
	productRepo *repository.ProductRepository,
	activityRepo *repository.ActivityRepository,
) *VersionService {
	return &VersionService{
		repo:         repo,
		revRepo:      revRepo,
		productRepo:  productRepo,
		activityRepo: activityRepo,
	}
}

func (s *VersionService) GetAll(p repository.PaginationParams, productID, status, owner string) ([]model.ProductVersion, repository.PaginationMeta, error) {
	versions, total, err := s.repo.FindAll(p, productID, status, owner)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return versions, meta, nil
}

func (s *VersionService) GetByID(id string) (*model.ProductVersion, error) {
	return s.repo.FindByID(id)
}

func (s *VersionService) GetLifecycle(productID string) (*model.ProductLifecycleResponse, error) {
	prod, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product master not found")
	}

	versions, err := s.repo.FindByProduct(prod.ID)
	if err != nil {
		return nil, err
	}

	res := &model.ProductLifecycleResponse{
		ProductCode: prod.Code,
		ProductName: prod.Name,
		Versions:    make([]model.ProductVersionSummary, 0, len(versions)),
	}

	for _, v := range versions {
		vSummary := model.ProductVersionSummary{
			ID:              v.ID,
			VersionNumber:   v.VersionNumber,
			VersionName:     v.VersionName,
			Status:          v.Status,
			CurrentRevision: v.CurrentRevision,
			Owner:           v.Owner,
			ReleaseDate:     v.ReleaseDate,
			Revisions:       make([]model.HardwareRevisionSummary, 0, len(v.HardwareRevisions)),
		}

		for _, r := range v.HardwareRevisions {
			vSummary.Revisions = append(vSummary.Revisions, model.HardwareRevisionSummary{
				ID:          r.ID,
				Code:        r.Code,
				Name:        r.Name,
				Status:      r.Status,
				PCBRevision: r.PCBRevision,
				ReleaseDate: r.ReleaseDate,
			})
		}

		res.Versions = append(res.Versions, vSummary)
	}

	return res, nil
}

func (s *VersionService) Compare(baseID, targetID string) (map[string]interface{}, error) {
	baseVer, err := s.repo.FindByID(baseID)
	if err != nil {
		return nil, fmt.Errorf("base version not found: %w", err)
	}
	targetVer, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, fmt.Errorf("target version not found: %w", err)
	}

	return map[string]interface{}{
		"baseVersion":   baseVer,
		"targetVersion": targetVer,
		"diffParameters": map[string]bool{
			"versionNumber":   baseVer.VersionNumber != targetVer.VersionNumber,
			"versionName":     baseVer.VersionName != targetVer.VersionName,
			"status":          baseVer.Status != targetVer.Status,
			"currentRevision": baseVer.CurrentRevision != targetVer.CurrentRevision,
			"owner":           baseVer.Owner != targetVer.Owner,
			"description":     baseVer.Description != targetVer.Description,
			"releaseNotes":    baseVer.ReleaseNotes != targetVer.ReleaseNotes,
		},
	}, nil
}

func (s *VersionService) Create(v *model.ProductVersion, currentUser *model.User) (*model.ProductVersion, error) {
	if v.ProductID == "" && v.ProductCode == "" {
		return nil, errors.New("product identifier is required")
	}
	if v.VersionNumber == "" {
		return nil, errors.New("version number is required (e.g. 1.0.0)")
	}
	if v.VersionName == "" {
		return nil, errors.New("version name is required")
	}

	// Validate allowed status
	v.Status = normalizeVersionStatus(v.Status)

	// Validate product exists
	prodKey := v.ProductID
	if prodKey == "" {
		prodKey = v.ProductCode
	}
	prod, err := s.productRepo.FindByID(prodKey)
	if err != nil {
		return nil, errors.New("associated product master not found")
	}

	// Validate product type: ONLY Manufacture products are allowed to have versions
	if prod.ProductType != "" && prod.ProductType != model.ProductTypeManufacture && prod.ProductType != "RND" {
		return nil, fmt.Errorf("cannot create version baseline: product '%s' is of type %s (versioning is strictly restricted to Manufacture Products)", prod.Code, prod.ProductType)
	}

	v.ProductID = prod.ID
	v.ProductCode = prod.Code
	v.ProductName = prod.Name

	// Check uniqueness
	existing, _ := s.repo.FindByProductAndVersion(v.ProductID, v.VersionNumber)
	if existing != nil {
		return nil, fmt.Errorf("version '%s' already exists for product '%s'", v.VersionNumber, prod.Code)
	}

	// Generate ID
	if v.ID == "" {
		v.ID = fmt.Sprintf("VER-%d", time.Now().UnixNano()%1000000)
	}

	if v.Owner == "" {
		if currentUser != nil {
			v.Owner = currentUser.Name
			v.OwnerAvatar = currentUser.Avatar
		} else {
			v.Owner = "Alex Chen"
			v.OwnerAvatar = "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&auto=format&fit=crop&q=80"
		}
	}

	if v.CurrentRevision == "" {
		v.CurrentRevision = "REV-A"
	}
	v.RevisionsCount = 1
	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()

	if err := s.repo.Create(v); err != nil {
		return nil, err
	}

	// Auto-create initial REV-A revision
	initialRev := &model.HardwareRevision{
		ID:                fmt.Sprintf("REV-%d", (time.Now().UnixNano()+1)%1000000),
		ProductVersionID:  v.ID,
		VersionNumber:     v.VersionNumber,
		ProductCode:       v.ProductCode,
		ProductName:       v.ProductName,
		Code:              "REV-A",
		Name:              fmt.Sprintf("Initial Hardware Revision for v%s", v.VersionNumber),
		Status:            "Draft",
		PCBRevision:       fmt.Sprintf("PCB-v%s-A", v.VersionNumber),
		SchematicRevision: fmt.Sprintf("SCH-v%s-01", v.VersionNumber),
		Engineer:          v.Owner,
		ChangeSummary:     fmt.Sprintf("Initial engineering hardware baseline for %s v%s.", v.ProductName, v.VersionNumber),
		ChangesList:       `["Initial schematic capture and PCB footprint placement."]`,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	_ = s.revRepo.Create(initialRev)

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
		Module:      "versions",
		EntityType:  "ProductVersion",
		EntityID:    v.ID,
		EntityName:  fmt.Sprintf("%s v%s", v.ProductName, v.VersionNumber),
		Action:      "PRODUCT_VERSION_CREATED",
		Description: fmt.Sprintf("Registered new version baseline 'v%s' for '%s'", v.VersionNumber, v.ProductCode),
		BadgeColor:  "emerald",
		CreatedAt:   time.Now(),
	})

	return v, nil
}

func (s *VersionService) Update(id string, v *model.ProductVersion, currentUser *model.User) (*model.ProductVersion, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product version not found")
	}

	if v.VersionName != "" {
		existing.VersionName = v.VersionName
	}
	if v.Status != "" {
		existing.Status = normalizeVersionStatus(v.Status)
	}
	if v.Owner != "" {
		existing.Owner = v.Owner
	}
	if v.ReleaseDate != nil {
		existing.ReleaseDate = v.ReleaseDate
	}
	if v.Description != "" {
		existing.Description = v.Description
	}
	if v.ReleaseNotes != "" {
		existing.ReleaseNotes = v.ReleaseNotes
	}
	if v.Documents != "" {
		existing.Documents = v.Documents
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
		Module:      "versions",
		EntityType:  "ProductVersion",
		EntityID:    existing.ID,
		EntityName:  fmt.Sprintf("%s v%s", existing.ProductName, existing.VersionNumber),
		Action:      "PRODUCT_VERSION_UPDATED",
		Description: fmt.Sprintf("Updated version baseline 'v%s' parameters for '%s'", existing.VersionNumber, existing.ProductCode),
		BadgeColor:  "blue",
		CreatedAt:   time.Now(),
	})

	return existing, nil
}

func (s *VersionService) Delete(id string, currentUser *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product version not found")
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
		Module:      "versions",
		EntityType:  "ProductVersion",
		EntityID:    existing.ID,
		EntityName:  fmt.Sprintf("%s v%s", existing.ProductName, existing.VersionNumber),
		Action:      "PRODUCT_VERSION_DELETED",
		Description: fmt.Sprintf("Removed version baseline 'v%s' from '%s'", existing.VersionNumber, existing.ProductCode),
		BadgeColor:  "rose",
		CreatedAt:   time.Now(),
	})

	return nil
}

func normalizeVersionStatus(status string) string {
	s := strings.ToLower(strings.TrimSpace(status))
	switch s {
	case "active", "released", "production":
		return "Active"
	case "in review", "in_review", "review":
		return "In Review"
	case "draft", "new":
		return "Draft"
	case "deprecated":
		return "Deprecated"
	case "archived", "obsolete", "eol":
		return "Archived"
	default:
		return "Draft"
	}
}
