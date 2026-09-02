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

type ProductService struct {
	repo            *repository.ProductRepository
	tradingRepo     *repository.TradingRepository
	projectItemRepo *repository.ProjectItemRepository
	versionRepo     *repository.VersionRepository
	activityRepo    *repository.ActivityRepository
	bomService      *BOMService
	rateService     *ExchangeRateService
}

func NewProductService(
	repo *repository.ProductRepository,
	tradingRepo *repository.TradingRepository,
	projectItemRepo *repository.ProjectItemRepository,
	versionRepo *repository.VersionRepository,
	activityRepo *repository.ActivityRepository,
) *ProductService {
	return &ProductService{
		repo:            repo,
		tradingRepo:     tradingRepo,
		projectItemRepo: projectItemRepo,
		versionRepo:     versionRepo,
		activityRepo:    activityRepo,
	}
}

func (s *ProductService) SetBOMService(bomService *BOMService) {
	s.bomService = bomService
}

func (s *ProductService) SetExchangeRateService(rateService *ExchangeRateService) {
	s.rateService = rateService
}

func (s *ProductService) GetAll(p repository.PaginationParams, category, status, productType string) ([]model.Product, repository.PaginationMeta, error) {
	products, total, err := s.repo.FindAll(p, category, status, productType)
	if err != nil {
		return nil, repository.PaginationMeta{}, err
	}

	for i := range products {
		s.PopulatePriceSummaries(&products[i])
	}

	meta := repository.CalculateMeta(total, p.Page, p.Limit)
	return products, meta, nil
}

func (s *ProductService) GetByID(id string) (*model.Product, error) {
	prod, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if prod.ProductType == model.ProductTypeTrading && prod.TradingDetail == nil && s.tradingRepo != nil {
		detail, _ := s.tradingRepo.FindByProductID(prod.ID)
		prod.TradingDetail = detail
	}
	s.PopulatePriceSummaries(prod)
	return prod, nil
}

func (s *ProductService) PopulatePriceSummaries(prod *model.Product) {
	if prod == nil {
		return
	}

	switch prod.ProductType {
	case model.ProductTypeTrading:
		detail := prod.TradingDetail
		if detail == nil && s.tradingRepo != nil {
			detail, _ = s.tradingRepo.FindByProductID(prod.ID)
			prod.TradingDetail = detail
		}

		if detail != nil && detail.PurchasePrice > 0 {
			pCur := detail.PurchaseCurrency
			if pCur == "" {
				pCur = detail.Currency
			}
			if pCur == "" {
				pCur = "USD"
			}

			var costIDR *float64
			if detail.PurchasePriceInIDR > 0 {
				val := detail.PurchasePriceInIDR
				costIDR = &val
			} else if s.rateService != nil {
				val := s.rateService.GetIDREquivalent(detail.PurchasePrice, pCur)
				costIDR = &val
			}

			marginVal := detail.ProfitMarginPercentage
			effRate := detail.EffectiveExchangeRate
			rateMode := detail.ExchangeRateMode

			prod.CostSummary = &model.PriceSummary{
				Amount:       detail.PurchasePrice,
				Currency:     pCur,
				AmountInIDR:  costIDR,
				Label:        "Purchase Cost",
				Available:    true,
				Margin:       &marginVal,
				RateMode:     rateMode,
				ExchangeRate: &effRate,
			}
		} else {
			prod.CostSummary = &model.PriceSummary{
				Amount:    0,
				Currency:  "USD",
				Label:     "No Purchase Cost",
				Available: false,
			}
		}

		if detail != nil && detail.SellingPrice > 0 {
			sCur := detail.SellingCurrency
			if sCur == "" {
				sCur = detail.PurchaseCurrency
			}
			if sCur == "" {
				sCur = detail.Currency
			}
			if sCur == "" {
				sCur = "USD"
			}

			var sellIDR *float64
			if detail.SellingPriceInIDR > 0 {
				val := detail.SellingPriceInIDR
				sellIDR = &val
			} else if s.rateService != nil {
				val := s.rateService.GetIDREquivalent(detail.SellingPrice, sCur)
				sellIDR = &val
			}

			marginVal := detail.ProfitMarginPercentage
			effRate := detail.EffectiveExchangeRate
			rateMode := detail.ExchangeRateMode

			prod.SellingPriceSummary = &model.PriceSummary{
				Amount:       detail.SellingPrice,
				Currency:     sCur,
				AmountInIDR:  sellIDR,
				Label:        "Commercial Price",
				Available:    true,
				Margin:       &marginVal,
				RateMode:     rateMode,
				ExchangeRate: &effRate,
			}
		} else {
			prod.SellingPriceSummary = &model.PriceSummary{
				Amount:    0,
				Currency:  "USD",
				Label:     "Not Calculated",
				Available: false,
			}
		}

	case model.ProductTypeRND:
		if s.bomService != nil {
			rndCost, err := s.bomService.GetProductCostSummary(prod.ID)
			if err == nil && rndCost != nil && rndCost.EstimatedBOMCost > 0 {
				currency := rndCost.Currency
				if currency == "" {
					currency = "USD"
				}

				var costIDR, sellIDR *float64
				if s.rateService != nil {
					cVal := s.rateService.GetIDREquivalent(rndCost.EstimatedBOMCost, currency)
					costIDR = &cVal
					sVal := s.rateService.GetIDREquivalent(rndCost.EstimatedSellingPrice, currency)
					sellIDR = &sVal
				}

				marginVal := rndCost.TargetMargin

				prod.CostSummary = &model.PriceSummary{
					Amount:      rndCost.EstimatedBOMCost,
					Currency:    currency,
					AmountInIDR: costIDR,
					Label:       "BOM Cost",
					Available:   true,
					Margin:      &marginVal,
				}
				prod.SellingPriceSummary = &model.PriceSummary{
					Amount:      rndCost.EstimatedSellingPrice,
					Currency:    currency,
					AmountInIDR: sellIDR,
					Label:       "Estimated Price",
					Available:   true,
					Margin:      &marginVal,
				}
			} else {
				prod.CostSummary = &model.PriceSummary{
					Amount:    0,
					Currency:  "USD",
					Label:     "No BOM Cost",
					Available: false,
				}
				prod.SellingPriceSummary = &model.PriceSummary{
					Amount:    0,
					Currency:  "USD",
					Label:     "Not Calculated",
					Available: false,
				}
			}
		} else {
			prod.CostSummary = &model.PriceSummary{
				Amount:    0,
				Currency:  "USD",
				Label:     "No BOM Cost",
				Available: false,
			}
			prod.SellingPriceSummary = &model.PriceSummary{
				Amount:    0,
				Currency:  "USD",
				Label:     "Not Calculated",
				Available: false,
			}
		}

	case model.ProductTypeProject:
		if s.bomService != nil {
			prjCost, err := s.bomService.GetProjectCostSummary(prod.ID)
			if err == nil && prjCost != nil && prjCost.TotalEstimatedCost > 0 {
				currency := prjCost.Currency
				if currency == "" {
					currency = "USD"
				}

				var costIDR, sellIDR *float64
				if s.rateService != nil {
					cVal := s.rateService.GetIDREquivalent(prjCost.TotalEstimatedCost, currency)
					costIDR = &cVal
					sVal := s.rateService.GetIDREquivalent(prjCost.EstimatedSellingPrice, currency)
					sellIDR = &sVal
				}

				marginVal := prjCost.TargetMargin

				prod.CostSummary = &model.PriceSummary{
					Amount:      prjCost.TotalEstimatedCost,
					Currency:    currency,
					AmountInIDR: costIDR,
					Label:       "Project Cost",
					Available:   true,
					Margin:      &marginVal,
				}
				prod.SellingPriceSummary = &model.PriceSummary{
					Amount:      prjCost.EstimatedSellingPrice,
					Currency:    currency,
					AmountInIDR: sellIDR,
					Label:       "Solution Price",
					Available:   true,
					Margin:      &marginVal,
				}
			} else {
				prod.CostSummary = &model.PriceSummary{
					Amount:    0,
					Currency:  "USD",
					Label:     "Not Calculated",
					Available: false,
				}
				prod.SellingPriceSummary = &model.PriceSummary{
					Amount:    0,
					Currency:  "USD",
					Label:     "Not Calculated",
					Available: false,
				}
			}
		} else {
			prod.CostSummary = &model.PriceSummary{
				Amount:    0,
				Currency:  "USD",
				Label:     "Not Calculated",
				Available: false,
			}
			prod.SellingPriceSummary = &model.PriceSummary{
				Amount:    0,
				Currency:  "USD",
				Label:     "Not Calculated",
				Available: false,
			}
		}
	}
}

func (s *ProductService) Create(p *model.Product, currentUser *model.User) (*model.Product, error) {
	if p.Name == "" {
		return nil, errors.New("product name is required")
	}

	// Normalize ProductType
	p.ProductType = normalizeProductType(p.ProductType)

	if p.Code == "" {
		count, _ := s.repo.Count()
		prefix := "PRD"
		if p.ProductType == model.ProductTypeTrading {
			prefix = "PRD-TRD"
		} else if p.ProductType == model.ProductTypeProject {
			prefix = "PRD-PRJ"
		}
		p.Code = fmt.Sprintf("%s-2024-%03d", prefix, count+1)
	}

	if p.ID == "" {
		p.ID = p.Code
	}

	if p.Status == "" {
		p.Status = "Active"
	}

	if p.Images != "" {
		var imgList []string
		if err := json.Unmarshal([]byte(p.Images), &imgList); err == nil && len(imgList) > 0 {
			p.Image = imgList[0]
		}
	} else if p.Image != "" {
		p.Images = fmt.Sprintf("[\"%s\"]", p.Image)
	}

	if p.Image == "" {
		p.Image = "https://images.unsplash.com/photo-1518770660439-4636190af475?w=400&auto=format&fit=crop&q=80"
		p.Images = fmt.Sprintf("[\"%s\"]", p.Image)
	}

	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	if err := s.repo.Create(p); err != nil {
		return nil, err
	}

	// If Trading Detail was provided during creation
	if p.ProductType == model.ProductTypeTrading && p.TradingDetail != nil && s.tradingRepo != nil {
		p.TradingDetail.ID = fmt.Sprintf("TRD-DET-%d", time.Now().UnixNano()%100000)
		p.TradingDetail.ProductID = p.ID
		if s.rateService != nil {
			_ = s.rateService.CalculateTradingDetailPricing(p.TradingDetail)
		}
		p.TradingDetail.CreatedAt = time.Now()
		p.TradingDetail.UpdatedAt = time.Now()
		_ = s.tradingRepo.Upsert(p.TradingDetail)
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
		Module:      "products",
		EntityType:  "Product",
		EntityID:    p.ID,
		EntityName:  p.Name,
		Action:      fmt.Sprintf("Created %s Product", p.ProductType),
		Description: fmt.Sprintf("Registered new %s product '%s' (%s)", p.ProductType, p.Name, p.Code),
		BadgeColor:  "emerald",
		CreatedAt:   time.Now(),
	})

	return p, nil
}

func (s *ProductService) Update(id string, p *model.Product, currentUser *model.User) (*model.Product, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	newType := normalizeProductType(p.ProductType)
	if newType != "" && newType != existing.ProductType {
		// Type Conversion Safety Rules
		if existing.ProductType == model.ProductTypeRND && s.versionRepo != nil {
			versions, _ := s.versionRepo.FindByProduct(existing.ID)
			if len(versions) > 0 {
				return nil, fmt.Errorf("cannot convert product type from R&D to %s: product has %d existing engineering versions and hardware revisions", newType, len(versions))
			}
		}
		if existing.ProductType == model.ProductTypeProject && s.projectItemRepo != nil {
			items, _ := s.projectItemRepo.FindByProjectID(existing.ID)
			if len(items) > 0 {
				return nil, fmt.Errorf("cannot convert product type from PROJECT to %s: product has %d existing project composition items", newType, len(items))
			}
		}
		existing.ProductType = newType
	}

	existing.Name = p.Name
	existing.Category = p.Category
	existing.CategoryID = p.CategoryID
	existing.Status = p.Status
	existing.Description = p.Description
	existing.LongDescription = p.LongDescription
	existing.ProjectLead = p.ProjectLead
	existing.TargetMarket = p.TargetMarket
	existing.PowerRating = p.PowerRating
	existing.IngressProtection = p.IngressProtection
	existing.OperatingTemp = p.OperatingTemp
	if p.Certifications != "" {
		existing.Certifications = p.Certifications
	}
	if p.Images != "" {
		existing.Images = p.Images
		var imgList []string
		if err := json.Unmarshal([]byte(p.Images), &imgList); err == nil && len(imgList) > 0 {
			existing.Image = imgList[0]
		}
	} else if p.Image != "" {
		existing.Image = p.Image
		existing.Images = fmt.Sprintf("[\"%s\"]", p.Image)
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
		Module:      "products",
		EntityType:  "Product",
		EntityID:    existing.ID,
		EntityName:  existing.Name,
		Action:      "Updated Product Data",
		Description: fmt.Sprintf("Updated specifications and master attributes for '%s'", existing.Name),
		BadgeColor:  "blue",
		CreatedAt:   time.Now(),
	})

	return existing, nil
}

func (s *ProductService) Delete(id string, currentUser *model.User) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	if err := s.repo.Delete(existing.ID); err != nil {
		return err
	}

	// Also clean up trading detail or project items if any
	if s.tradingRepo != nil {
		_ = s.tradingRepo.Delete(existing.ID)
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
		Module:      "products",
		EntityType:  "Product",
		EntityID:    existing.ID,
		EntityName:  existing.Name,
		Action:      "Deleted Product",
		Description: fmt.Sprintf("Removed product master '%s'", existing.Name),
		BadgeColor:  "rose",
		CreatedAt:   time.Now(),
	})

	return nil
}

// Trading Detail Methods
func (s *ProductService) GetTradingDetail(productID string) (*model.TradingProductDetail, error) {
	prod, err := s.repo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if prod.ProductType != model.ProductTypeTrading {
		return nil, errors.New("product is not a TRADING product")
	}
	return s.tradingRepo.FindByProductID(prod.ID)
}

func (s *ProductService) UpdateTradingDetail(productID string, detail *model.TradingProductDetail, currentUser *model.User) (*model.TradingProductDetail, error) {
	prod, err := s.repo.FindByID(productID)
	if err != nil {
		return nil, errors.New("product not found")
	}
	if prod.ProductType != model.ProductTypeTrading {
		return nil, errors.New("product is not a TRADING product")
	}

	detail.ProductID = prod.ID
	if detail.ID == "" {
		existing, _ := s.tradingRepo.FindByProductID(prod.ID)
		if existing != nil {
			detail.ID = existing.ID
		} else {
			detail.ID = fmt.Sprintf("TRD-DET-%d", time.Now().UnixNano()%100000)
		}
	}

	if s.rateService != nil {
		_ = s.rateService.CalculateTradingDetailPricing(detail)
	}

	detail.UpdatedAt = time.Now()

	if err := s.tradingRepo.Upsert(detail); err != nil {
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
		EntityType:  "Product",
		EntityID:    prod.ID,
		EntityName:  prod.Name,
		Action:      "Updated Commercial Sourcing Info",
		Description: fmt.Sprintf("Updated trading commercial details and pricing for '%s'", prod.Name),
		BadgeColor:  "amber",
		CreatedAt:   time.Now(),
	})

	return detail, nil
}

func normalizeProductType(t string) string {
	s := strings.ToUpper(strings.TrimSpace(t))
	switch s {
	case "TRADING", "COMMERCIAL":
		return model.ProductTypeTrading
	case "PROJECT", "SOLUTION", "SYSTEM":
		return model.ProductTypeProject
	case "RND", "R&D", "ENGINEERING":
		return model.ProductTypeRND
	default:
		return model.ProductTypeRND
	}
}
