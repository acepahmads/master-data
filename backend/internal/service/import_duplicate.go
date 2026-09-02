package service

import (
	"encoding/json"
	"strings"

	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/repository"
)

type DuplicateCandidate struct {
	MasterID   string `json:"masterId"`
	MasterType string `json:"masterType"` // PRODUCT, COMPONENT
	Code       string `json:"code"`
	Name       string `json:"name"`
	MatchType  string `json:"matchType"`  // EXACT_CODE, EXACT_NAME, EXACT_MPN, FUZZY_NAME, BATCH_DUPLICATE
}

type ImportDuplicateService struct {
	productRepo   *repository.ProductRepository
	componentRepo *repository.ComponentRepository
}

func NewImportDuplicateService(prodRepo *repository.ProductRepository, compRepo *repository.ComponentRepository) *ImportDuplicateService {
	return &ImportDuplicateService{
		productRepo:   prodRepo,
		componentRepo: compRepo,
	}
}

// AnalyzeDuplicates evaluates a set of staged rows for collisions against Master Data and intra-batch duplicates.
func (d *ImportDuplicateService) AnalyzeDuplicates(rows []model.ImportStagedRow) error {
	// Preload existing products and components for fast in-memory matching
	existingProds, _, err := d.productRepo.FindAll(repository.PaginationParams{Page: 1, Limit: 10000}, "", "", "")
	if err != nil {
		existingProds = []model.Product{}
	}
	existingComps, _, err := d.componentRepo.FindAll(repository.PaginationParams{Page: 1, Limit: 10000}, repository.ComponentFilters{})
	if err != nil {
		existingComps = []model.Component{}
	}

	// Build lookup maps
	prodCodeMap := make(map[string]model.Product)
	prodNameMap := make(map[string]model.Product)
	for _, p := range existingProds {
		if p.Code != "" {
			prodCodeMap[strings.ToUpper(strings.TrimSpace(p.Code))] = p
		}
		if p.Name != "" {
			prodNameMap[strings.ToUpper(strings.TrimSpace(p.Name))] = p
		}
	}

	compMPNMap := make(map[string]model.Component)
	compNameMap := make(map[string]model.Component)
	for _, c := range existingComps {
		if c.PartNumber != "" {
			compMPNMap[strings.ToUpper(strings.TrimSpace(c.PartNumber))] = c
		}
		if c.Name != "" {
			compNameMap[strings.ToUpper(strings.TrimSpace(c.Name))] = c
		}
	}

	// Intra-batch tracker
	batchCodeSeen := make(map[string]string) // code -> rowId
	batchNameSeen := make(map[string]string) // name -> rowId

	for i := range rows {
		row := &rows[i]
		var candidates []DuplicateCandidate
		rowCode := strings.ToUpper(strings.TrimSpace(row.NormalizedCode))
		rowName := strings.ToUpper(strings.TrimSpace(row.NormalizedName))

		// 1. Check Exact Product Code Match
		if rowCode != "" {
			if matchedP, exists := prodCodeMap[rowCode]; exists {
				candidates = append(candidates, DuplicateCandidate{
					MasterID:   matchedP.ID,
					MasterType: "PRODUCT",
					Code:       matchedP.Code,
					Name:       matchedP.Name,
					MatchType:  "EXACT_CODE",
				})
			}
		}

		// 2. Check Exact Product Name Match
		if rowName != "" {
			if matchedP, exists := prodNameMap[rowName]; exists {
				candidates = append(candidates, DuplicateCandidate{
					MasterID:   matchedP.ID,
					MasterType: "PRODUCT",
					Code:       matchedP.Code,
					Name:       matchedP.Name,
					MatchType:  "EXACT_NAME",
				})
			}
		}

		// 3. Check Exact Component MPN Match
		if rowCode != "" {
			if matchedC, exists := compMPNMap[rowCode]; exists {
				candidates = append(candidates, DuplicateCandidate{
					MasterID:   matchedC.ID,
					MasterType: "COMPONENT",
					Code:       matchedC.PartNumber,
					Name:       matchedC.Name,
					MatchType:  "EXACT_MPN",
				})
			}
		}

		// 4. Check Exact Component Name Match
		if rowName != "" {
			if matchedC, exists := compNameMap[rowName]; exists {
				candidates = append(candidates, DuplicateCandidate{
					MasterID:   matchedC.ID,
					MasterType: "COMPONENT",
					Code:       matchedC.PartNumber,
					Name:       matchedC.Name,
					MatchType:  "EXACT_NAME",
				})
			}
		}

		// 5. Check Intra-Batch Duplicates
		if rowCode != "" {
			if prevRowID, exists := batchCodeSeen[rowCode]; exists && prevRowID != row.ID {
				candidates = append(candidates, DuplicateCandidate{
					MasterID:   prevRowID,
					MasterType: "STAGED_ROW",
					Code:       rowCode,
					Name:       row.NormalizedName,
					MatchType:  "BATCH_DUPLICATE",
				})
			} else {
				batchCodeSeen[rowCode] = row.ID
			}
		}
		if rowName != "" {
			if prevRowID, exists := batchNameSeen[rowName]; exists && prevRowID != row.ID {
				candidates = append(candidates, DuplicateCandidate{
					MasterID:   prevRowID,
					MasterType: "STAGED_ROW",
					Code:       rowCode,
					Name:       row.NormalizedName,
					MatchType:  "BATCH_DUPLICATE",
				})
			} else {
				batchNameSeen[rowName] = row.ID
			}
		}

		// Assign duplicate diagnostics
		if len(candidates) > 0 {
			first := candidates[0]
			if first.MatchType == "BATCH_DUPLICATE" {
				row.DuplicateStatus = model.RowDuplicateBatchCollision
			} else if first.MatchType == "EXACT_CODE" || first.MatchType == "EXACT_MPN" {
				row.DuplicateStatus = model.RowDuplicateMatchedExisting
				row.MatchedMasterID = first.MasterID
				row.MatchedMasterType = first.MasterType
			} else {
				row.DuplicateStatus = model.RowDuplicatePossibleMatch
				row.MatchedMasterID = first.MasterID
				row.MatchedMasterType = first.MasterType
			}
			candJSON, _ := json.Marshal(candidates)
			row.DuplicateCandidatesJSON = string(candJSON)
		} else {
			row.DuplicateStatus = model.RowDuplicateNewRecord
			row.DuplicateCandidatesJSON = "[]"
			row.MatchedMasterID = ""
			row.MatchedMasterType = ""
		}
	}

	return nil
}
