package service

import (
	"errors"
	"fmt"
	"time"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type ProjectItemService struct {
	repo         *repository.ProjectItemRepository
	productRepo  *repository.ProductRepository
	compRepo     *repository.ComponentRepository
	activityRepo *repository.ActivityRepository
}

func NewProjectItemService(
	repo *repository.ProjectItemRepository,
	productRepo *repository.ProductRepository,
	compRepo *repository.ComponentRepository,
	activityRepo *repository.ActivityRepository,
) *ProjectItemService {
	return &ProjectItemService{
		repo:         repo,
		productRepo:  productRepo,
		compRepo:     compRepo,
		activityRepo: activityRepo,
	}
}

// GetProjectItems returns flat and hierarchical items for a project product
func (s *ProjectItemService) GetProjectItems(projectID string) ([]model.ProjectItem, error) {
	prod, err := s.productRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("project product not found")
	}

	items, err := s.repo.FindByProjectID(prod.ID)
	if err != nil {
		return nil, err
	}

	// Enrich pricing if empty from linked product/trading details or component
	for i := range items {
		if items[i].SellingPrice == 0 && items[i].ProductID != nil && *items[i].ProductID != "" {
			if linkedProd, err := s.productRepo.FindByID(*items[i].ProductID); err == nil && linkedProd != nil {
				if linkedProd.TradingDetail != nil {
					items[i].CostPrice = linkedProd.TradingDetail.PurchasePrice
					items[i].SellingPrice = linkedProd.TradingDetail.SellingPrice
					if linkedProd.TradingDetail.Currency != "" {
						items[i].Currency = linkedProd.TradingDetail.Currency
					}
				}
			}
		} else if items[i].CostPrice == 0 && items[i].ComponentID != nil && *items[i].ComponentID != "" {
			if linkedComp, err := s.compRepo.FindByID(*items[i].ComponentID); err == nil && linkedComp != nil {
				items[i].CostPrice = linkedComp.EstimatedUnitCost
				if linkedComp.Currency != "" {
					items[i].Currency = linkedComp.Currency
				}
			}
		}
		if items[i].Currency == "" {
			items[i].Currency = "IDR"
		}
	}

	return items, nil
}

// GetProjectTree returns nested tree structure of project items
func (s *ProjectItemService) GetProjectTree(projectID string) ([]model.ProjectItem, error) {
	items, err := s.GetProjectItems(projectID)
	if err != nil {
		return nil, err
	}

	// Build map and tree
	itemMap := make(map[string]*model.ProjectItem)
	var rootItems []model.ProjectItem

	for i := range items {
		itemCopy := items[i]
		itemCopy.SubItems = []model.ProjectItem{}
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

	return rootItems, nil
}

func (s *ProjectItemService) CreateItem(projectID string, item *model.ProjectItem, currentUser *model.User) (*model.ProjectItem, error) {
	prod, err := s.productRepo.FindByID(projectID)
	if err != nil {
		return nil, errors.New("parent project product not found")
	}
	if prod.ProductType != model.ProductTypeProject {
		return nil, fmt.Errorf("cannot add project items: product '%s' is not a PROJECT product", prod.Name)
	}

	item.ProjectProductID = prod.ID

	if item.Name == "" {
		return nil, errors.New("item name is required")
	}

	// Validate item type
	switch item.ItemType {
	case model.ProjectItemTypeProduct:
		if item.ProductID == nil || *item.ProductID == "" {
			return nil, errors.New("referenced product ID is required for PRODUCT item type")
		}
		refProd, err := s.productRepo.FindByID(*item.ProductID)
		if err != nil {
			return nil, errors.New("referenced product master not found")
		}
		item.ProductID = &refProd.ID
		item.ProductCode = refProd.Code
		item.ProductType = refProd.ProductType
		if item.Name == "" {
			item.Name = refProd.Name
		}
	case model.ProjectItemTypeComponent:
		if item.ComponentID != nil && *item.ComponentID != "" {
			comp, err := s.compRepo.FindByID(*item.ComponentID)
			if err == nil && comp != nil {
				item.ComponentMPN = comp.PartNumber
				if item.Name == "" {
					item.Name = comp.Name
				}
			}
		}
	case model.ProjectItemTypeSubAssembly:
		// Grouping node
	default:
		item.ItemType = model.ProjectItemTypeSubAssembly
	}

	if item.Quantity <= 0 {
		item.Quantity = 1.0
	}
	if item.UnitName == "" {
		item.UnitName = "pcs"
	}
	if item.Status == "" {
		item.Status = "Active"
	}

	if item.ID == "" {
		count, _ := s.repo.CountByProjectID(prod.ID)
		item.ID = fmt.Sprintf("PRJ-ITM-%d", time.Now().UnixNano()%1000000)
		item.SortOrder = int(count) + 1
	}

	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

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
		Module:      "products",
		EntityType:  "ProjectItem",
		EntityID:    item.ID,
		EntityName:  item.Name,
		Action:      "Added Project Structure Item",
		Description: fmt.Sprintf("Added %s '%s' to project '%s'", item.ItemType, item.Name, prod.Name),
		BadgeColor:  "indigo",
		CreatedAt:   time.Now(),
	})

	return item, nil
}

func (s *ProjectItemService) UpdateItem(id string, item *model.ProjectItem, currentUser *model.User) (*model.ProjectItem, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("project item not found")
	}

	existing.Name = item.Name
	existing.ItemType = item.ItemType
	existing.ParentItemID = item.ParentItemID
	existing.ProductID = item.ProductID
	existing.ProductCode = item.ProductCode
	existing.ProductType = item.ProductType
	existing.ComponentID = item.ComponentID
	existing.ComponentMPN = item.ComponentMPN
	existing.Quantity = item.Quantity
	existing.UnitID = item.UnitID
	existing.UnitName = item.UnitName
	existing.SortOrder = item.SortOrder
	existing.Notes = item.Notes
	existing.Status = item.Status
	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

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
		Module:      "products",
		EntityType:  "ProjectItem",
		EntityID:    existing.ID,
		EntityName:  existing.Name,
		Action:      "Updated Project Item",
		Description: fmt.Sprintf("Updated project item specifications for '%s'", existing.Name),
		BadgeColor:  "blue",
		CreatedAt:   time.Now(),
	})

	return existing, nil
}

func (s *ProjectItemService) DeleteItem(id string, currentUser *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("project item not found")
	}

	if err := s.repo.Delete(existing.ID); err != nil {
		return err
	}

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
		Module:      "products",
		EntityType:  "ProjectItem",
		EntityID:    existing.ID,
		EntityName:  existing.Name,
		Action:      "Removed Project Item",
		Description: fmt.Sprintf("Deleted item '%s' from project structure", existing.Name),
		BadgeColor:  "rose",
		CreatedAt:   time.Now(),
	})

	return nil
}

func (s *ProjectItemService) ReorderItems(projectID string, itemIDs []string, currentUser *model.User) error {
	return s.repo.Reorder(projectID, itemIDs)
}
