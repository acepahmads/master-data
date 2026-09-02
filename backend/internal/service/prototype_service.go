package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type PrototypeService struct {
	protoRepo   *repository.PrototypeRepository
	revRepo     *repository.RevisionRepository
	versionRepo *repository.VersionRepository
	productRepo *repository.ProductRepository
	bomRepo     *repository.BOMRepository
	activityRepo *repository.ActivityRepository
}

func NewPrototypeService(
	protoRepo *repository.PrototypeRepository,
	revRepo *repository.RevisionRepository,
	versionRepo *repository.VersionRepository,
	productRepo *repository.ProductRepository,
	bomRepo *repository.BOMRepository,
	activityRepo *repository.ActivityRepository,
) *PrototypeService {
	return &PrototypeService{
		protoRepo:    protoRepo,
		revRepo:      revRepo,
		versionRepo:  versionRepo,
		productRepo:  productRepo,
		bomRepo:      bomRepo,
		activityRepo: activityRepo,
	}
}

type CreatePrototypeRequest struct {
	Code               string  `json:"code"`
	Name               string  `json:"name" binding:"required"`
	HardwareRevisionID string  `json:"hardwareRevisionId" binding:"required"`
	BuildNumber        int     `json:"buildNumber"`
	Quantity           int     `json:"quantity" binding:"required,min=1"`
	Purpose            string  `json:"purpose" binding:"required"`
	Objective          string  `json:"objective"`
	Status             string  `json:"status"`
	EngineeringOwner   string  `json:"engineeringOwner"`
	TargetCompletionAt *string `json:"targetCompletionAt"`
	Notes              string  `json:"notes"`
	InitialNote        string  `json:"initialNote"`
}

type UpdatePrototypeRequest struct {
	Name               *string `json:"name"`
	Purpose            *string `json:"purpose"`
	Objective          *string `json:"objective"`
	Status             *string `json:"status"`
	EngineeringOwner   *string `json:"engineeringOwner"`
	TargetCompletionAt *string `json:"targetCompletionAt"`
	Notes              *string `json:"notes"`
}

type UpdateCompPrepRequest struct {
	AvailableQuantity *float64 `json:"availableQuantity"`
	Status            *string  `json:"status"`
	Notes             *string  `json:"notes"`
}

type UpdateAssemblyStageRequest struct {
	Status      string  `json:"status" binding:"required"`
	PerformedBy string  `json:"performedBy"`
	Notes       string  `json:"notes"`
}

type CreateEngineeringNoteRequest struct {
	Category string `json:"category" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

func (s *PrototypeService) GetAll(p repository.PaginationParams, filters map[string]string) ([]model.PrototypeBuild, repository.PaginationMeta, error) {
	prototypes, total, err := s.protoRepo.FindAll(p, filters)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}
	return prototypes, repository.CalculateMeta(total, p.Page, p.Limit), nil
}

func (s *PrototypeService) GetDashboardStats() (*model.PrototypeDashboardStats, error) {
	return s.protoRepo.GetDashboardStats()
}

func (s *PrototypeService) GetByID(id string) (*model.PrototypeBuild, error) {
	return s.protoRepo.FindByID(id)
}

func (s *PrototypeService) GetByRevision(revisionID string) ([]model.PrototypeBuild, error) {
	return s.protoRepo.FindByRevision(revisionID)
}

func (s *PrototypeService) Create(req CreatePrototypeRequest, userID, userName, userAvatar string) (*model.PrototypeBuild, error) {
	// 1. Validate Hardware Revision
	revision, err := s.revRepo.FindByID(req.HardwareRevisionID)
	if err != nil {
		return nil, fmt.Errorf("hardware revision not found: %w", err)
	}

	// 2. Validate Product Version
	version, err := s.versionRepo.FindByID(revision.ProductVersionID)
	if err != nil {
		return nil, fmt.Errorf("product version not found for revision: %w", err)
	}

	// 3. Validate Product & Product Type = RND
	product, err := s.productRepo.FindByID(version.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found for version: %w", err)
	}

	productType := strings.ToUpper(strings.TrimSpace(product.ProductType))
	if productType != "RND" && productType != "" {
		return nil, fmt.Errorf("only R&D Products may create Prototype Builds. Product '%s' has type '%s' and is ineligible", product.Code, product.ProductType)
	}

	// 4. Validate / Allocate Build Number for Revision
	buildNum := req.BuildNumber
	if buildNum <= 0 {
		existingBuilds, _ := s.protoRepo.FindByRevision(revision.ID)
		maxNum := 0
		for _, b := range existingBuilds {
			if b.BuildNumber > maxNum {
				maxNum = b.BuildNumber
			}
		}
		buildNum = maxNum + 1
	} else {
		existingBuilds, _ := s.protoRepo.FindByRevision(revision.ID)
		for _, b := range existingBuilds {
			if b.BuildNumber == buildNum {
				return nil, fmt.Errorf("build number #%d already exists for Hardware Revision %s (%s)", buildNum, revision.Code, revision.Name)
			}
		}
	}

	// 5. Generate Code if empty
	code := req.Code
	if code == "" {
		code = fmt.Sprintf("PROTO-%s-B%03d", revision.Code, buildNum)
	}

	// Verify Code Uniqueness
	if existing, _ := s.protoRepo.FindByCode(code); existing != nil {
		code = fmt.Sprintf("PROTO-%d-%03d", time.Now().Year(), time.Now().Unix()%1000)
	}

	// 6. Load Engineering BOM & Create Frozen Snapshot
	bom, _ := s.bomRepo.FindByRevisionID(revision.ID)
	var snapshotItems []model.FrozenBOMItemDTO
	var totalCost float64 = 0
	var totalQty float64 = 0
	currency := "USD"

	if bom != nil {
		currency = bom.Currency
		if currency == "" {
			currency = "USD"
		}
		for i, item := range bom.Items {
			lineCost := item.LineCost
			if lineCost <= 0 {
				lineCost = item.Quantity * item.UnitCost
			}
			totalCost += lineCost
			totalQty += item.Quantity

			compName := item.Name
			mpn := item.ComponentPartNumber
			category := item.ComponentCategory
			if item.Component != nil {
				if compName == "" {
					compName = item.Component.Name
				}
				if mpn == "" {
					mpn = item.Component.PartNumber
				}
				if category == "" {
					if item.Component.Category != nil {
						category = item.Component.Category.Name
					} else {
						category = item.Component.CategoryName
					}
				}
			}

			snapshotItems = append(snapshotItems, model.FrozenBOMItemDTO{
				ID:                   fmt.Sprintf("SNP-ITM-%03d", i+1),
				BOMItemID:            item.ID,
				ComponentID:          item.ComponentID,
				ComponentName:        compName,
				MPN:                  mpn,
				Category:             category,
				Quantity:             item.Quantity,
				UnitCode:             item.UnitCode,
				UnitCost:             item.UnitCost,
				LineCost:             lineCost,
				ReferenceDesignators: item.ReferenceDesignators,
				Position:             item.Position,
				Status:               "FROZEN_SNAPSHOT",
			})
		}
	}

	snapshotID := fmt.Sprintf("SNP-%s-%03d", revision.Code, buildNum)
	snapshotData := model.FrozenBOMSnapshotData{
		SnapshotID:           snapshotID,
		FrozenAt:             time.Now(),
		HardwareRevisionID:   revision.ID,
		HardwareRevisionCode: revision.Code,
		TotalItems:           len(snapshotItems),
		TotalQuantity:        totalQty,
		TotalEstimatedCost:   totalCost,
		Currency:             currency,
		Items:                snapshotItems,
	}

	snapshotJSONBytes, err := json.Marshal(snapshotData)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize BOM snapshot: %w", err)
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if status == "" {
		status = model.PrototypeStatusPlanned
	}

	owner := req.EngineeringOwner
	if owner == "" {
		owner = userName
	}

	id := fmt.Sprintf("PROTO-%d-%04d", time.Now().Year(), time.Now().UnixNano()%10000)

	proto := &model.PrototypeBuild{
		ID:                   id,
		Code:                 code,
		Name:                 req.Name,
		HardwareRevisionID:   revision.ID,
		ProductID:            product.ID,
		ProductCode:          product.Code,
		ProductName:          product.Name,
		VersionID:            version.ID,
		VersionNumber:        version.VersionNumber,
		VersionName:          version.VersionName,
		HardwareRevisionCode: revision.Code,
		HardwareRevisionName: revision.Name,
		BuildNumber:          buildNum,
		Quantity:             req.Quantity,
		Purpose:              req.Purpose,
		Objective:            req.Objective,
		Status:               status,
		EngineeringOwner:     owner,
		TargetCompletionAt:   req.TargetCompletionAt,
		Notes:                req.Notes,
		BOMSnapshotJSON:      string(snapshotJSONBytes),
		BOMSnapshotCost:      totalCost,
		BOMSnapshotID:        snapshotID,
		BOMSnapshotCurrency:  currency,
		ReadinessPercentage:  0.0,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// 7. Atomic Database Transaction
	tx := s.protoRepo.GetDB().Begin()
	if err := tx.Create(proto).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create prototype build record: %w", err)
	}

	// Create Component Preparations (Required Qty = Snapshot Qty * Build Quantity)
	for i, itm := range snapshotItems {
		required := itm.Quantity * float64(req.Quantity)
		prep := model.PrototypeComponentPreparation{
			ID:                  fmt.Sprintf("PREP-%s-%03d", id, i+1),
			PrototypeBuildID:    id,
			BOMItemID:           &itm.BOMItemID,
			ComponentID:         itm.ComponentID,
			ComponentName:       itm.ComponentName,
			ComponentPartNumber: itm.MPN,
			ComponentCategory:   itm.Category,
			RequiredQuantity:    required,
			AvailableQuantity:   0,
			MissingQuantity:     required,
			UnitCode:            itm.UnitCode,
			Status:              model.CompPrepMissing,
			Notes:               "Allocated at build creation",
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}
		if err := tx.Create(&prep).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create component preparation: %w", err)
		}
	}

	// Create 8 Standard Assembly Stages
	stageDefs := []struct {
		Stage string
		Seq   int
		Name  string
		Desc  string
	}{
		{model.AssemblyStagePCBPreparation, 1, "PCB Preparation & Solder Paste Stencil", "Laser stainless stencil alignment and Lead-Free SAC305 paste deposition."},
		{model.AssemblyStageComponentPreparation, 2, "Component Kitting & Reel Setup", "Feeder loading and component verification for high-speed SMD assembly."},
		{model.AssemblyStageSMTAssembly, 3, "SMT Pick & Place Assembly", "Automated component placement for top and bottom PCBA layers."},
		{model.AssemblyStageThroughHoleAssembly, 4, "Reflow Soldering & AOI Inspection", "9-zone nitrogen reflow oven profile and automated optical inspection."},
		{model.AssemblyStageMechanicalFitting, 5, "Through-Hole Connectors & Hand Assembly", "Terminal blocks, SMA antenna connectors, and RJ45 Ethernet jacks."},
		{model.AssemblyStageWiring, 6, "Mechanical Fitting & Enclosure Assembly", "Extruded aluminum enclosure insertion, thermal pad mounting, and gaskets."},
		{model.AssemblyStageFirmwareFlashing, 7, "Firmware Bootloader Flashing", "JTAG/UART flashing of factory bootloader and MAC allocation."},
		{model.AssemblyStageInitialPowerOn, 8, "Initial Power-On & Voltage Rail Validation", "Current limiting bench supply test: 3.3V, 5.0V, and RF LDO rails."},
	}

	for _, sDef := range stageDefs {
		stage := model.PrototypeAssemblyStage{
			ID:               fmt.Sprintf("STG-%s-%02d", id, sDef.Seq),
			PrototypeBuildID: id,
			Stage:            sDef.Stage,
			Sequence:         sDef.Seq,
			Name:             sDef.Name,
			Description:      sDef.Desc,
			Status:           model.StageStatusNotStarted,
			PerformedBy:      owner,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := tx.Create(&stage).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create assembly stage: %w", err)
		}
	}

	// Create Initial Note if provided
	if req.InitialNote != "" {
		note := model.PrototypeEngineeringNote{
			ID:               fmt.Sprintf("NOT-%s-%03d", id, 1),
			PrototypeBuildID: id,
			Category:         model.NoteCategoryGeneral,
			Title:            "Initial Build Setup",
			Content:          req.InitialNote,
			CreatedBy:        owner,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := tx.Create(&note).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create initial note: %w", err)
		}
	}

	// Log Activity
	logEntry := model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "prototypes",
		EntityType:  "PrototypeBuild",
		EntityID:    id,
		EntityName:  proto.Name,
		Action:      "Created Prototype Build",
		Description: fmt.Sprintf("Created prototype build %s (#%d, %d Units) for revision %s with frozen BOM snapshot (%d items, $%.2f)", proto.Code, proto.BuildNumber, proto.Quantity, revision.Code, len(snapshotItems), totalCost),
		BadgeColor:  "indigo",
		CreatedAt:   time.Now(),
	}
	if err := tx.Create(&logEntry).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create activity log: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("transaction commit failed: %w", err)
	}

	return s.protoRepo.FindByID(id)
}

func (s *PrototypeService) Update(id string, req UpdatePrototypeRequest, userID, userName, userAvatar string) (*model.PrototypeBuild, error) {
	proto, err := s.protoRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("prototype build not found: %w", err)
	}

	var changes []string

	if req.Name != nil && *req.Name != "" && *req.Name != proto.Name {
		changes = append(changes, fmt.Sprintf("Name changed to '%s'", *req.Name))
		proto.Name = *req.Name
	}
	if req.Purpose != nil && *req.Purpose != "" {
		proto.Purpose = *req.Purpose
	}
	if req.Objective != nil {
		proto.Objective = *req.Objective
	}
	if req.EngineeringOwner != nil && *req.EngineeringOwner != "" {
		proto.EngineeringOwner = *req.EngineeringOwner
	}
	if req.TargetCompletionAt != nil {
		proto.TargetCompletionAt = req.TargetCompletionAt
	}
	if req.Notes != nil {
		proto.Notes = *req.Notes
	}

	if req.Status != nil && *req.Status != "" {
		newStatus := strings.ToUpper(strings.TrimSpace(*req.Status))
		if newStatus != proto.Status {
			if !s.isValidStatus(newStatus) {
				return nil, fmt.Errorf("invalid prototype build status '%s'", newStatus)
			}
			changes = append(changes, fmt.Sprintf("Status transitioned from %s to %s", proto.Status, newStatus))
			proto.Status = newStatus

			now := time.Now()
			if newStatus == model.PrototypeStatusCompleted {
				proto.CompletedAt = &now
			} else if (newStatus == model.PrototypeStatusInProgress || newStatus == model.PrototypeStatusAssembly) && proto.StartedAt == nil {
				proto.StartedAt = &now
			}
		}
	}

	proto.UpdatedAt = time.Now()
	if err := s.protoRepo.Update(proto); err != nil {
		return nil, err
	}

	if len(changes) > 0 {
		_ = s.activityRepo.Create(&model.ActivityLog{
			ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
			UserID:      userID,
			UserName:    userName,
			UserAvatar:  userAvatar,
			Module:      "prototypes",
			EntityType:  "PrototypeBuild",
			EntityID:    proto.ID,
			EntityName:  proto.Name,
			Action:      "Updated Prototype Build",
			Description: strings.Join(changes, "; "),
			BadgeColor:  "blue",
			CreatedAt:   time.Now(),
		})
	}

	return s.protoRepo.FindByID(proto.ID)
}

func (s *PrototypeService) Delete(id string, userID, userName, userAvatar string) error {
	proto, err := s.protoRepo.FindByID(id)
	if err != nil {
		return fmt.Errorf("prototype build not found: %w", err)
	}

	if err := s.protoRepo.Delete(proto.ID); err != nil {
		return err
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "prototypes",
		EntityType:  "PrototypeBuild",
		EntityID:    proto.ID,
		EntityName:  proto.Name,
		Action:      "Deleted Prototype Build",
		Description: fmt.Sprintf("Deleted prototype build %s (%s)", proto.Code, proto.Name),
		BadgeColor:  "red",
		CreatedAt:   time.Now(),
	})

	return nil
}

func (s *PrototypeService) GetComponentPreparations(prototypeID string) ([]model.PrototypeComponentPreparation, float64, error) {
	proto, err := s.protoRepo.FindByID(prototypeID)
	if err != nil {
		return nil, 0, fmt.Errorf("prototype build not found: %w", err)
	}
	preps, err := s.protoRepo.FindComponentPreparations(proto.ID)
	return preps, proto.ReadinessPercentage, err
}

func (s *PrototypeService) UpdateComponentPreparation(prototypeID, compPrepID string, req UpdateCompPrepRequest, userID, userName, userAvatar string) (*model.PrototypeComponentPreparation, float64, error) {
	prep, err := s.protoRepo.FindComponentPreparationByID(compPrepID)
	if err != nil {
		return nil, 0, fmt.Errorf("component preparation record not found: %w", err)
	}

	if prep.PrototypeBuildID != prototypeID {
		proto, _ := s.protoRepo.FindByID(prototypeID)
		if proto == nil || prep.PrototypeBuildID != proto.ID {
			return nil, 0, errors.New("component preparation does not belong to specified prototype build")
		}
	}

	if req.AvailableQuantity != nil {
		prep.AvailableQuantity = *req.AvailableQuantity
		if prep.AvailableQuantity < 0 {
			prep.AvailableQuantity = 0
		}
	}
	if req.Notes != nil {
		prep.Notes = *req.Notes
	}

	// Server-Side Readiness & Missing Calculation
	if req.Status != nil && *req.Status == model.CompPrepSubstituteRequired {
		prep.Status = model.CompPrepSubstituteRequired
		prep.MissingQuantity = 0
	} else if prep.AvailableQuantity >= prep.RequiredQuantity {
		prep.Status = model.CompPrepAvailable
		prep.MissingQuantity = 0
	} else if prep.AvailableQuantity > 0 {
		prep.Status = model.CompPrepPartial
		prep.MissingQuantity = prep.RequiredQuantity - prep.AvailableQuantity
	} else {
		prep.Status = model.CompPrepMissing
		prep.MissingQuantity = prep.RequiredQuantity
	}

	prep.UpdatedAt = time.Now()
	if err := s.protoRepo.UpdateComponentPreparation(prep); err != nil {
		return nil, 0, err
	}

	// Server-Side Overall Readiness Percentage Calculation
	allPreps, _ := s.protoRepo.FindComponentPreparations(prep.PrototypeBuildID)
	var totalReq float64 = 0
	var totalAvail float64 = 0

	for _, p := range allPreps {
		totalReq += p.RequiredQuantity
		avail := p.AvailableQuantity
		if avail > p.RequiredQuantity {
			avail = p.RequiredQuantity
		}
		if p.Status == model.CompPrepSubstituteRequired {
			avail = p.RequiredQuantity
		}
		totalAvail += avail
	}

	var readinessPct float64 = 100.0
	if totalReq > 0 {
		readinessPct = (totalAvail / totalReq) * 100.0
		if readinessPct > 100.0 {
			readinessPct = 100.0
		}
	}

	// Update prototype build readiness percentage
	proto, _ := s.protoRepo.FindByID(prep.PrototypeBuildID)
	if proto != nil {
		proto.ReadinessPercentage = readinessPct
		_ = s.protoRepo.Update(proto)
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "prototypes",
		EntityType:  "ComponentPreparation",
		EntityID:    prep.ID,
		EntityName:  prep.ComponentName,
		Action:      "Updated Component Readiness",
		Description: fmt.Sprintf("Updated %s: %.0f/%.0f %s (%s). Overall readiness: %.1f%%", prep.ComponentName, prep.AvailableQuantity, prep.RequiredQuantity, prep.UnitCode, prep.Status, readinessPct),
		BadgeColor:  "amber",
		CreatedAt:   time.Now(),
	})

	return prep, readinessPct, nil
}

func (s *PrototypeService) GetAssemblyStages(prototypeID string) ([]model.PrototypeAssemblyStage, error) {
	proto, err := s.protoRepo.FindByID(prototypeID)
	if err != nil {
		return nil, fmt.Errorf("prototype build not found: %w", err)
	}
	return s.protoRepo.FindAssemblyStages(proto.ID)
}

func (s *PrototypeService) UpdateAssemblyStage(prototypeID, stageID string, req UpdateAssemblyStageRequest, userID, userName, userAvatar string) (*model.PrototypeAssemblyStage, error) {
	stage, err := s.protoRepo.FindAssemblyStageByID(stageID)
	if err != nil {
		return nil, fmt.Errorf("assembly stage not found: %w", err)
	}

	status := strings.ToUpper(strings.TrimSpace(req.Status))
	if !s.isValidStageStatus(status) {
		return nil, fmt.Errorf("invalid assembly stage status '%s'", status)
	}

	// Server-Side Sequence Validation
	if status == model.StageStatusCompleted {
		allStages, _ := s.protoRepo.FindAssemblyStages(stage.PrototypeBuildID)
		for _, prev := range allStages {
			if prev.Sequence < stage.Sequence {
				if prev.Status != model.StageStatusCompleted && prev.Status != model.StageStatusSkipped {
					return nil, fmt.Errorf("sequence violation: cannot complete Stage %d (%s) before prerequisite Stage %d (%s) is completed or skipped", stage.Sequence, stage.Name, prev.Sequence, prev.Name)
				}
			}
		}
	}

	stage.Status = status
	if req.PerformedBy != "" {
		stage.PerformedBy = req.PerformedBy
	}
	if req.Notes != "" {
		stage.Notes = req.Notes
	}

	now := time.Now()
	if status == model.StageStatusCompleted {
		stage.CompletedAt = &now
	} else if status == model.StageStatusInProgress {
		stage.StartedAt = &now
		stage.CompletedAt = nil
	} else {
		stage.CompletedAt = nil
	}

	stage.UpdatedAt = time.Now()
	if err := s.protoRepo.UpdateAssemblyStage(stage); err != nil {
		return nil, err
	}

	// Check if all 8 stages are completed to automatically transition build status
	allStages, _ := s.protoRepo.FindAssemblyStages(stage.PrototypeBuildID)
	allDone := true
	for _, st := range allStages {
		if st.Status != model.StageStatusCompleted && st.Status != model.StageStatusSkipped {
			allDone = false
			break
		}
	}

	proto, _ := s.protoRepo.FindByID(stage.PrototypeBuildID)
	if proto != nil && allDone && proto.Status != model.PrototypeStatusCompleted {
		proto.Status = "READY_FOR_TESTING"
		_ = s.protoRepo.Update(proto)
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "prototypes",
		EntityType:  "AssemblyStage",
		EntityID:    stage.ID,
		EntityName:  stage.Name,
		Action:      "Updated Assembly Stage",
		Description: fmt.Sprintf("Stage 0%d (%s) transitioned to %s by %s", stage.Sequence, stage.Name, stage.Status, stage.PerformedBy),
		BadgeColor:  "indigo",
		CreatedAt:   time.Now(),
	})

	return stage, nil
}

func (s *PrototypeService) GetEngineeringNotes(prototypeID string) ([]model.PrototypeEngineeringNote, error) {
	proto, err := s.protoRepo.FindByID(prototypeID)
	if err != nil {
		return nil, fmt.Errorf("prototype build not found: %w", err)
	}
	return s.protoRepo.FindEngineeringNotes(proto.ID)
}

func (s *PrototypeService) CreateEngineeringNote(prototypeID string, req CreateEngineeringNoteRequest, userID, userName, userAvatar string) (*model.PrototypeEngineeringNote, error) {
	proto, err := s.protoRepo.FindByID(prototypeID)
	if err != nil {
		return nil, fmt.Errorf("prototype build not found: %w", err)
	}

	category := strings.ToUpper(strings.TrimSpace(req.Category))
	if !s.isValidNoteCategory(category) {
		category = model.NoteCategoryGeneral
	}

	id := fmt.Sprintf("NOT-%s-%d", proto.ID, time.Now().UnixNano()%10000)
	note := &model.PrototypeEngineeringNote{
		ID:               id,
		PrototypeBuildID: proto.ID,
		Category:         category,
		Title:            req.Title,
		Content:          req.Content,
		CreatedBy:        userName,
		UserID:           &userID,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.protoRepo.CreateEngineeringNote(note); err != nil {
		return nil, err
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "prototypes",
		EntityType:  "EngineeringNote",
		EntityID:    note.ID,
		EntityName:  note.Title,
		Action:      "Added Engineering Note",
		Description: fmt.Sprintf("Logged %s note: '%s' for %s", note.Category, note.Title, proto.Code),
		BadgeColor:  "purple",
		CreatedAt:   time.Now(),
	})

	return note, nil
}

func (s *PrototypeService) UpdateEngineeringNote(prototypeID, noteID string, req CreateEngineeringNoteRequest, userID, userName, userAvatar string) (*model.PrototypeEngineeringNote, error) {
	note, err := s.protoRepo.FindEngineeringNoteByID(noteID)
	if err != nil {
		return nil, fmt.Errorf("engineering note not found: %w", err)
	}

	if req.Category != "" {
		category := strings.ToUpper(strings.TrimSpace(req.Category))
		if s.isValidNoteCategory(category) {
			note.Category = category
		}
	}
	if req.Title != "" {
		note.Title = req.Title
	}
	if req.Content != "" {
		note.Content = req.Content
	}

	note.UpdatedAt = time.Now()
	if err := s.protoRepo.UpdateEngineeringNote(note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *PrototypeService) DeleteEngineeringNote(prototypeID, noteID string, userID, userName, userAvatar string) error {
	note, err := s.protoRepo.FindEngineeringNoteByID(noteID)
	if err != nil {
		return fmt.Errorf("engineering note not found: %w", err)
	}

	if err := s.protoRepo.DeleteEngineeringNote(note.ID); err != nil {
		return err
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          fmt.Sprintf("ACT-%d", time.Now().UnixNano()),
		UserID:      userID,
		UserName:    userName,
		UserAvatar:  userAvatar,
		Module:      "prototypes",
		EntityType:  "EngineeringNote",
		EntityID:    note.ID,
		EntityName:  note.Title,
		Action:      "Deleted Engineering Note",
		Description: fmt.Sprintf("Removed note: '%s'", note.Title),
		BadgeColor:  "slate",
		CreatedAt:   time.Now(),
	})

	return nil
}

func (s *PrototypeService) GetBOMSnapshot(prototypeID string) (*model.FrozenBOMSnapshotData, error) {
	proto, err := s.protoRepo.FindByID(prototypeID)
	if err != nil {
		return nil, fmt.Errorf("prototype build not found: %w", err)
	}

	var snapshot model.FrozenBOMSnapshotData
	if err := json.Unmarshal([]byte(proto.BOMSnapshotJSON), &snapshot); err != nil {
		return nil, fmt.Errorf("failed to unmarshal BOM snapshot JSON: %w", err)
	}
	return &snapshot, nil
}

func (s *PrototypeService) isValidStatus(status string) bool {
	switch status {
	case model.PrototypeStatusPlanned,
		model.PrototypeStatusPreparing,
		model.PrototypeStatusInProgress,
		model.PrototypeStatusAssembly,
		model.PrototypeStatusPowerOn,
		"READY_FOR_TESTING",
		model.PrototypeStatusCompleted,
		model.PrototypeStatusOnHold,
		model.PrototypeStatusCancelled:
		return true
	default:
		return false
	}
}

func (s *PrototypeService) isValidStageStatus(status string) bool {
	switch status {
	case model.StageStatusNotStarted,
		model.StageStatusInProgress,
		model.StageStatusCompleted,
		model.StageStatusBlocked,
		model.StageStatusSkipped:
		return true
	default:
		return false
	}
}

func (s *PrototypeService) isValidNoteCategory(cat string) bool {
	switch cat {
	case model.NoteCategoryGeneral,
		model.NoteCategoryAssemblyIssue,
		model.NoteCategoryDesignObservation,
		model.NoteCategoryKnownIssue,
		model.NoteCategoryRecommendation:
		return true
	default:
		return false
	}
}
