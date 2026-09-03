package service

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type BOMService struct {
	repo            *repository.BOMRepository
	revRepo         *repository.RevisionRepository
	versionRepo     *repository.VersionRepository
	productRepo     *repository.ProductRepository
	compRepo        *repository.ComponentRepository
	tradingRepo     *repository.TradingRepository
	projectItemRepo *repository.ProjectItemRepository
	activityRepo    *repository.ActivityRepository
}

func NewBOMService(
	repo *repository.BOMRepository,
	revRepo *repository.RevisionRepository,
	versionRepo *repository.VersionRepository,
	productRepo *repository.ProductRepository,
	compRepo *repository.ComponentRepository,
	tradingRepo *repository.TradingRepository,
	projectItemRepo *repository.ProjectItemRepository,
	activityRepo *repository.ActivityRepository,
) *BOMService {
	return &BOMService{
		repo:            repo,
		revRepo:         revRepo,
		versionRepo:     versionRepo,
		productRepo:     productRepo,
		compRepo:        compRepo,
		tradingRepo:     tradingRepo,
		projectItemRepo: projectItemRepo,
		activityRepo:    activityRepo,
	}
}

func (s *BOMService) GetBOMByRevisionID(revisionID string) (*model.EngineeringBOM, error) {
	bom, err := s.repo.FindByRevisionID(revisionID)
	if err != nil {
		return nil, err
	}
	// Build nested hierarchy
	bom.Items = s.buildItemHierarchy(bom.Items)
	return bom, nil
}

func (s *BOMService) GetBOMByID(id string) (*model.EngineeringBOM, error) {
	bom, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	bom.Items = s.buildItemHierarchy(bom.Items)
	return bom, nil
}

func (s *BOMService) CreateBOM(revisionID string, bom *model.EngineeringBOM, currentUser *model.User) (*model.EngineeringBOM, error) {
	// 1. Verify Revision exists
	rev, err := s.revRepo.FindByID(revisionID)
	if err != nil {
		return nil, errors.New("associated hardware revision not found")
	}

	// 2. Verify Product Version exists
	ver, err := s.versionRepo.FindByID(rev.ProductVersionID)
	if err != nil {
		return nil, errors.New("associated product version not found")
	}

	// 3. Verify Product Master exists and is strictly R&D
	prod, err := s.productRepo.FindByID(ver.ProductID)
	if err != nil {
		return nil, errors.New("associated product master not found")
	}
	if prod.ProductType != "" && prod.ProductType != model.ProductTypeManufacture && prod.ProductType != "RND" {
		return nil, fmt.Errorf("cannot create engineering BOM: product '%s' is of type %s (BOM management is strictly restricted to Manufacture Products)", prod.Code, prod.ProductType)
	}

	// 4. Check if BOM already exists for this revision
	existing, _ := s.repo.FindByRevisionID(rev.ID)
	if existing != nil {
		return existing, nil
	}

	bom.HardwareRevisionID = rev.ID
	if bom.BOMCode == "" {
		bom.BOMCode = fmt.Sprintf("BOM-%s-%s-%s", prod.Code, ver.VersionNumber, rev.Code)
	}
	if bom.Name == "" {
		bom.Name = fmt.Sprintf("%s Engineering BOM (%s %s)", prod.Name, ver.VersionNumber, rev.Code)
	}
	if bom.ID == "" {
		bom.ID = fmt.Sprintf("BOM-%d", time.Now().UnixNano()%1000000)
	}
	if bom.Status == "" {
		bom.Status = model.BOMStatusActive
	}
	if bom.Currency == "" {
		bom.Currency = "USD"
	}
	if bom.TargetMargin <= 0 {
		bom.TargetMargin = 35.00
	}
	if currentUser != nil {
		bom.CreatedBy = currentUser.Name
	} else {
		bom.CreatedBy = "System Engineer"
	}

	bom.CreatedAt = time.Now()
	bom.UpdatedAt = time.Now()

	if err := s.repo.Create(bom); err != nil {
		return nil, err
	}

	// Log Activity
	s.logActivity(currentUser, "Created Engineering BOM", fmt.Sprintf("Created engineering BOM '%s' (%s) for revision %s", bom.Name, bom.BOMCode, rev.Code), bom.ID, "emerald")

	return bom, nil
}

func (s *BOMService) UpdateBOM(id string, b *model.EngineeringBOM, currentUser *model.User) (*model.EngineeringBOM, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("engineering BOM not found")
	}

	if b.Name != "" {
		existing.Name = b.Name
	}
	if b.Description != "" {
		existing.Description = b.Description
	}
	if b.Status != "" {
		existing.Status = b.Status
	}
	if b.TargetMargin > 0 && b.TargetMargin < 100 {
		existing.TargetMargin = b.TargetMargin
	}
	if b.Notes != "" {
		existing.Notes = b.Notes
	}

	// Recalculate selling price based on margin
	if existing.TargetMargin > 0 && existing.TargetMargin < 100 {
		existing.EstimatedSellingPrice = math.Round((existing.TotalCost/(1.0-(existing.TargetMargin/100.0)))*100) / 100
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	s.logActivity(currentUser, "Updated Engineering BOM", fmt.Sprintf("Updated metadata and target margin for BOM '%s'", existing.BOMCode), existing.ID, "blue")

	return existing, nil
}

func (s *BOMService) AddItem(bomID string, item *model.BOMItem, currentUser *model.User) (*model.BOMItem, error) {
	bom, err := s.repo.FindByID(bomID)
	if err != nil {
		return nil, errors.New("engineering BOM not found")
	}

	item.EngineeringBOMID = bom.ID

	if item.ItemType == "" {
		item.ItemType = model.BOMItemTypeComponent
	}

	if item.Quantity <= 0 {
		item.Quantity = 1.0
	}
	if item.UnitCode == "" {
		item.UnitCode = "PCS"
	}

	if item.ItemType == model.BOMItemTypeComponent {
		if item.ComponentID == nil || *item.ComponentID == "" {
			return nil, errors.New("component selection is required for component BOM item")
		}
		comp, err := s.compRepo.FindByID(*item.ComponentID)
		if err != nil || comp == nil {
			return nil, errors.New("referenced component master not found")
		}

		item.ComponentPartNumber = comp.PartNumber
		item.ComponentCategory = comp.CategoryName
		if item.Name == "" {
			item.Name = comp.Name
		}
		if item.UnitCost <= 0 {
			item.UnitCost = comp.EstimatedUnitCost
		}
		item.LineCost = math.Round((item.Quantity*item.UnitCost)*10000) / 10000
	} else if item.ItemType == model.BOMItemTypeSubAssembly {
		if item.Name == "" {
			return nil, errors.New("sub-assembly name is required")
		}
		item.ComponentID = nil
		item.UnitCost = 0
		item.LineCost = 0
	}

	if item.ID == "" {
		item.ID = fmt.Sprintf("BOM-ITM-%d", time.Now().UnixNano()%1000000)
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	if err := s.repo.CreateItem(item); err != nil {
		return nil, err
	}

	// Trigger Cost Recalculation
	_ = s.RecalculateBOMCost(bom.ID)

	s.logActivity(currentUser, "Added BOM Item", fmt.Sprintf("Added %s '%s' (Qty: %.2f %s) to BOM %s", item.ItemType, item.Name, item.Quantity, item.UnitCode, bom.BOMCode), item.ID, "indigo")

	return item, nil
}

func (s *BOMService) UpdateItem(id string, item *model.BOMItem, currentUser *model.User) (*model.BOMItem, error) {
	existing, err := s.repo.FindItemByID(id)
	if err != nil {
		return nil, errors.New("BOM item not found")
	}

	// Cycle prevention: cannot set parent to itself
	if item.ParentItemID != nil && *item.ParentItemID == existing.ID {
		return nil, errors.New("circular hierarchy detected: item cannot be its own parent")
	}

	existing.Name = item.Name
	existing.ParentItemID = item.ParentItemID
	existing.ItemType = item.ItemType
	if item.Quantity > 0 {
		existing.Quantity = item.Quantity
	}
	if item.UnitCode != "" {
		existing.UnitCode = item.UnitCode
	}
	if item.UnitCost >= 0 {
		existing.UnitCost = item.UnitCost
	}
	if existing.ItemType == model.BOMItemTypeComponent {
		existing.LineCost = math.Round((existing.Quantity*existing.UnitCost)*10000) / 10000
	}
	existing.ReferenceDesignators = item.ReferenceDesignators
	existing.Notes = item.Notes
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateItem(existing); err != nil {
		return nil, err
	}

	// Recalculate Total BOM Cost
	_ = s.RecalculateBOMCost(existing.EngineeringBOMID)

	s.logActivity(currentUser, "Updated BOM Item", fmt.Sprintf("Updated specifications for item '%s'", existing.Name), existing.ID, "blue")

	return existing, nil
}

func (s *BOMService) DeleteItem(id string, currentUser *model.User) error {
	existing, err := s.repo.FindItemByID(id)
	if err != nil {
		return errors.New("BOM item not found")
	}

	bomID := existing.EngineeringBOMID
	name := existing.Name

	if err := s.repo.DeleteItem(existing.ID); err != nil {
		return err
	}

	_ = s.RecalculateBOMCost(bomID)

	s.logActivity(currentUser, "Removed BOM Item", fmt.Sprintf("Deleted item '%s' from engineering BOM", name), existing.ID, "rose")

	return nil
}

func (s *BOMService) ReorderItems(bomID string, itemIDs []string, currentUser *model.User) error {
	return s.repo.ReorderItems(bomID, itemIDs)
}

// RecalculateBOMCost computes recursive sub-assembly line costs, total BOM cost, and estimated selling price.
func (s *BOMService) RecalculateBOMCost(bomID string) error {
	bom, err := s.repo.FindByID(bomID)
	if err != nil {
		return err
	}

	items, err := s.repo.FindItemsByBOMID(bomID)
	if err != nil {
		return err
	}

	// 1. Group items by ID and map children
	itemMap := make(map[string]*model.BOMItem)
	childrenMap := make(map[string][]*model.BOMItem)

	for i := range items {
		itemMap[items[i].ID] = &items[i]
		if items[i].ParentItemID != nil && *items[i].ParentItemID != "" {
			pID := *items[i].ParentItemID
			childrenMap[pID] = append(childrenMap[pID], &items[i])
		}
	}

	// 2. Recursive function to compute Sub-Assembly Line Costs
	var computeSubAssemblyCost func(item *model.BOMItem) float64
	computeSubAssemblyCost = func(item *model.BOMItem) float64 {
		if item.ItemType == model.BOMItemTypeComponent {
			item.LineCost = math.Round((item.Quantity*item.UnitCost)*10000) / 10000
			return item.LineCost
		}

		// Sub-Assembly: sum children
		var subTotal float64
		for _, child := range childrenMap[item.ID] {
			subTotal += computeSubAssemblyCost(child)
		}
		item.UnitCost = math.Round(subTotal*10000) / 10000
		item.LineCost = math.Round((item.Quantity*item.UnitCost)*10000) / 10000
		_ = s.repo.UpdateItem(item)
		return item.LineCost
	}

	// 3. Compute Top Level Total
	var totalBOMCost float64
	for i := range items {
		if items[i].ParentItemID == nil || *items[i].ParentItemID == "" {
			totalBOMCost += computeSubAssemblyCost(&items[i])
		}
	}

	bom.TotalCost = math.Round(totalBOMCost*10000) / 10000
	if bom.TargetMargin > 0 && bom.TargetMargin < 100 {
		bom.EstimatedSellingPrice = math.Round((bom.TotalCost/(1.0-(bom.TargetMargin/100.0)))*100) / 100
	} else {
		bom.EstimatedSellingPrice = bom.TotalCost
	}
	bom.UpdatedAt = time.Now()

	return s.repo.Update(bom)
}

// GetBOMCostSummary provides cost category breakdown and top 5 expensive components.
func (s *BOMService) GetBOMCostSummary(bomID string) (*model.BOMCostSummary, error) {
	bom, err := s.repo.FindByID(bomID)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.FindItemsByBOMID(bom.ID)
	if err != nil {
		return nil, err
	}

	summary := &model.BOMCostSummary{
		TotalBOMCost:          bom.TotalCost,
		Currency:              bom.Currency,
		TargetMargin:          bom.TargetMargin,
		EstimatedSellingPrice: bom.EstimatedSellingPrice,
		CostByCategory:        make(map[string]float64),
		TopCostDrivers:        []model.TopCostItem{},
	}

	var components []model.BOMItem
	for _, item := range items {
		if item.ItemType == model.BOMItemTypeComponent {
			summary.TotalComponentsCount++
			cat := item.ComponentCategory
			if cat == "" {
				cat = "General Electronic"
			}
			summary.CostByCategory[cat] += item.LineCost
			components = append(components, item)
		} else {
			summary.TotalAssembliesCount++
		}
	}

	// Sort components by line cost descending
	sort.Slice(components, func(i, j int) bool {
		return components[i].LineCost > components[j].LineCost
	})

	limit := 5
	if len(components) < limit {
		limit = len(components)
	}

	for i := 0; i < limit; i++ {
		c := components[i]
		pct := 0.0
		if bom.TotalCost > 0 {
			pct = math.Round((c.LineCost/bom.TotalCost)*1000) / 10
		}
		summary.TopCostDrivers = append(summary.TopCostDrivers, model.TopCostItem{
			ItemName:          c.Name,
			PartNumber:        c.ComponentPartNumber,
			Quantity:          c.Quantity,
			UnitCost:          c.UnitCost,
			LineCost:          c.LineCost,
			PercentageOfTotal: pct,
		})
	}

	return summary, nil
}

// GetProductCostSummary rolls up estimated BOM cost for an R&D product based on its latest active hardware revision.
func (s *BOMService) GetProductCostSummary(productID string) (*model.ProductCostSummary, error) {
	prod, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	summary := &model.ProductCostSummary{
		ProductID:   prod.ID,
		ProductCode: prod.Code,
		ProductName: prod.Name,
		ProductType: prod.ProductType,
		Currency:    "USD",
	}

	if prod.ProductType != model.ProductTypeManufacture && prod.ProductType != "RND" {
		return summary, nil
	}

	// Fetch versions and locate active BOM baseline
	versions, err := s.versionRepo.FindByProduct(prod.ID)
	if err == nil && len(versions) > 0 {
		for _, ver := range versions {
			revs, _ := s.revRepo.FindByVersion(ver.ID)
			for _, rev := range revs {
				bom, _ := s.repo.FindByRevisionID(rev.ID)
				if bom != nil && bom.TotalCost > 0 {
					summary.CurrentVersion = ver.VersionNumber
					summary.CurrentRevision = rev.Code
					summary.RevisionID = rev.ID
					summary.BOMID = bom.ID
					summary.BOMCode = bom.BOMCode
					summary.EstimatedBOMCost = bom.TotalCost
					summary.Currency = bom.Currency
					summary.TargetMargin = bom.TargetMargin
					summary.EstimatedSellingPrice = bom.EstimatedSellingPrice
					return summary, nil
				}
			}
		}

		// Fallback to latest version if no BOM populated yet
		latestVer := versions[0]
		summary.CurrentVersion = latestVer.VersionNumber
		revs, _ := s.revRepo.FindByVersion(latestVer.ID)
		if len(revs) > 0 {
			summary.CurrentRevision = revs[0].Code
			summary.RevisionID = revs[0].ID
		}
	}

	return summary, nil
}

// GetProjectCostSummary aggregates costs for a Project Solution without copying internal R&D BOMs.
func (s *BOMService) GetProjectCostSummary(projectID string) (*model.ProjectCostSummary, error) {
	prod, err := s.productRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project product not found")
	}

	summary := &model.ProjectCostSummary{
		ProjectID:    prod.ID,
		ProjectCode:  prod.Code,
		ProjectName:  prod.Name,
		ProductType:  prod.ProductType,
		Currency:     "USD",
		TargetMargin: 30.00, // standard project solution margin
	}

	items, err := s.projectItemRepo.FindByProjectID(prod.ID)
	if err != nil {
		return summary, nil
	}
	summary.ItemsCount = len(items)

	for _, item := range items {
		if item.ItemType == model.ProjectItemTypeProduct {
			cost := s.computeProjectItemCost(&item)
			if item.ProductType == model.ProductTypeTrading {
				summary.TradingProductsCost += cost
			} else {
				summary.RNDProductsCost += cost
			}
		} else if item.ItemType == model.ProjectItemTypeComponent {
			summary.DirectComponentsCost += s.computeProjectItemCost(&item)
		}
	}

	summary.TradingProductsCost = math.Round(summary.TradingProductsCost*100) / 100
	summary.RNDProductsCost = math.Round(summary.RNDProductsCost*100) / 100
	summary.DirectComponentsCost = math.Round(summary.DirectComponentsCost*100) / 100
	summary.TotalEstimatedCost = math.Round((summary.TradingProductsCost+summary.RNDProductsCost+summary.DirectComponentsCost)*100) / 100

	if summary.TargetMargin > 0 && summary.TargetMargin < 100 {
		summary.EstimatedSellingPrice = math.Round((summary.TotalEstimatedCost/(1.0-(summary.TargetMargin/100.0)))*100) / 100
	} else {
		summary.EstimatedSellingPrice = summary.TotalEstimatedCost
	}

	return summary, nil
}

func (s *BOMService) computeProjectItemCost(item *model.ProjectItem) float64 {
	switch item.ItemType {
	case model.ProjectItemTypeProduct:
		if item.ProductID == nil {
			return 0
		}
		refProd, err := s.productRepo.FindByID(*item.ProductID)
		if err != nil || refProd == nil {
			return 0
		}
		if refProd.ProductType == model.ProductTypeTrading {
			detail, _ := s.tradingRepo.FindByProductID(refProd.ID)
			if detail != nil {
				return math.Round((item.Quantity*detail.PurchasePrice)*100) / 100
			}
		} else if refProd.ProductType == model.ProductTypeManufacture || refProd.ProductType == "RND" {
			costSum, _ := s.GetProductCostSummary(refProd.ID)
			if costSum != nil && costSum.EstimatedBOMCost > 0 {
				return math.Round((item.Quantity*costSum.EstimatedBOMCost)*100) / 100
			}
		}
	case model.ProjectItemTypeComponent:
		if item.ComponentID != nil {
			comp, _ := s.compRepo.FindByID(*item.ComponentID)
			if comp != nil {
				return math.Round((item.Quantity*comp.EstimatedUnitCost)*100) / 100
			}
		}
	case model.ProjectItemTypeSubAssembly:
		// Recursive sum of sub-items
		items, _ := s.projectItemRepo.FindByProjectID(item.ProjectProductID)
		var subSum float64
		for _, child := range items {
			if child.ParentItemID != nil && *child.ParentItemID == item.ID {
				subSum += s.computeProjectItemCost(&child)
			}
		}
		return math.Round((item.Quantity*subSum)*100) / 100
	}
	return 0
}

func (s *BOMService) buildItemHierarchy(items []model.BOMItem) []model.BOMItem {
	itemMap := make(map[string]*model.BOMItem)
	var rootItems []model.BOMItem

	for i := range items {
		itemCopy := items[i]
		itemCopy.SubItems = []model.BOMItem{}
		itemMap[itemCopy.ID] = &itemCopy
	}

	for i := range items {
		curr := itemMap[items[i].ID]
		if curr.ParentItemID != nil && *curr.ParentItemID != "" {
			if parent, exists := itemMap[*curr.ParentItemID]; exists {
				parent.SubItems = append(parent.SubItems, *curr)
			} else {
				rootItems = append(rootItems, *curr)
			}
		} else {
			rootItems = append(rootItems, *curr)
		}
	}

	return rootItems
}

func (s *BOMService) logActivity(currentUser *model.User, action, desc, entityID, color string) {
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
		Module:      "bom",
		EntityType:  "EngineeringBOM",
		EntityID:    entityID,
		EntityName:  "BOM Specification",
		Action:      action,
		Description: desc,
		BadgeColor:  color,
		CreatedAt:   time.Now(),
	})
}
