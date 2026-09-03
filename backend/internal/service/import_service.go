package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"iot-rd-backend/internal/config"
	"iot-rd-backend/internal/model"
	"iot-rd-backend/internal/provider"
	"iot-rd-backend/internal/repository"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ImportService struct {
	importRepo    *repository.ImportRepository
	productRepo   *repository.ProductRepository
	tradingRepo   *repository.TradingRepository
	componentRepo *repository.ComponentRepository
	catRepo       *repository.CategoryRepository
	mfgRepo       *repository.ManufacturerRepository
	supplierRepo  *repository.SupplierRepository
	activityRepo  *repository.ActivityRepository
	parser        *ImportParserService
	normalizer    *ImportNormalizerService
	validator     *ImportValidatorService
	duplicate     *ImportDuplicateService
	aiProvider    provider.AIClassificationProvider
	config        *config.Config
}

func NewImportService(
	impRepo *repository.ImportRepository,
	prodRepo *repository.ProductRepository,
	tradingRepo *repository.TradingRepository,
	compRepo *repository.ComponentRepository,
	catRepo *repository.CategoryRepository,
	mfgRepo *repository.ManufacturerRepository,
	supplierRepo *repository.SupplierRepository,
	activityRepo *repository.ActivityRepository,
	aiProvider provider.AIClassificationProvider,
	cfg *config.Config,
) *ImportService {
	return &ImportService{
		importRepo:    impRepo,
		productRepo:   prodRepo,
		tradingRepo:   tradingRepo,
		componentRepo: compRepo,
		catRepo:       catRepo,
		mfgRepo:       mfgRepo,
		supplierRepo:  supplierRepo,
		activityRepo:  activityRepo,
		parser:        NewImportParserService(),
		normalizer:    NewImportNormalizerService(),
		validator:     NewImportValidatorService(),
		duplicate:     NewImportDuplicateService(prodRepo, compRepo),
		aiProvider:    aiProvider,
		config:        cfg,
	}
}

// CreateBatch creates a new import batch container.
func (s *ImportService) CreateBatch(sourceYear, sourceReference, sourceNotes, user string) (*model.ImportBatch, error) {
	if sourceYear == "" {
		sourceYear = "2026"
	}
	
	// Generate unique sequential batch code
	batchCode := fmt.Sprintf("IMP-%s-%04d", sourceYear, time.Now().Unix()%10000)

	batch := &model.ImportBatch{
		ID:              uuid.New().String(),
		BatchCode:       batchCode,
		SourceYear:      sourceYear,
		SourceType:      "EXCEL_WORKBOOK",
		SourceReference: sourceReference,
		SourceNotes:     sourceNotes,
		Status:          model.BatchStatusUploaded,
		UploadedBy:      user,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.importRepo.CreateBatch(batch); err != nil {
		return nil, err
	}

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          uuid.New().String(),
		UserID:      user,
		Module:      "dataimport",
		Action:      "CREATE",
		EntityType:  "IMPORT_BATCH",
		EntityID:    batch.ID,
		Description: fmt.Sprintf("Created import batch %s (Source Year: %s)", batch.BatchCode, batch.SourceYear),
		CreatedAt:   time.Now(),
	})

	return batch, nil
}

// GetBatch retrieves a batch with its files.
func (s *ImportService) GetBatch(id string) (*model.ImportBatch, error) {
	return s.importRepo.GetBatchByID(id)
}

// GetBatches returns batches filtered by year and status.
func (s *ImportService) GetBatches(year string, status string) ([]model.ImportBatch, error) {
	return s.importRepo.GetBatches(year, status)
}

// UploadFiles stores uploaded Excel files into controlled local storage and attaches them to the batch.
func (s *ImportService) UploadFiles(batchID string, files []*multipart.FileHeader, user string) ([]model.ImportFile, error) {
	batch, err := s.importRepo.GetBatchByID(batchID)
	if err != nil {
		return nil, fmt.Errorf("import batch not found: %w", err)
	}

	storageDir := filepath.Join(s.config.UploadPath, "imports", batch.SourceYear, batch.BatchCode)
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	var createdFiles []model.ImportFile

	for _, fileHeader := range files {
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		if ext != ".xlsx" && ext != ".xlsm" {
			continue // Skip non-excel
		}

		srcFile, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer func() { _ = srcFile.Close() }()

		// Compute SHA256 and save file
		hasher := sha256.New()
		storedFilename := fmt.Sprintf("%s_%s%s", uuid.New().String()[:8], strings.ReplaceAll(fileHeader.Filename, " ", "_"), ext)
		dstPath := filepath.Join(storageDir, storedFilename)

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return nil, err
		}
		defer func() { _ = dstFile.Close() }()

		multiWriter := io.MultiWriter(dstFile, hasher)
		if _, err := io.Copy(multiWriter, srcFile); err != nil {
			return nil, err
		}

		sha256Hex := hex.EncodeToString(hasher.Sum(nil))

		// Inspect workbook sheets
		sheetManifest, err := s.parser.InspectWorkbook(dstPath)
		manifestJSON := "[]"
		sheetCount := 1
		totalRows := 0
		if err == nil && len(sheetManifest) > 0 {
			mBytes, _ := json.Marshal(sheetManifest)
			manifestJSON = string(mBytes)
			sheetCount = len(sheetManifest)
			for _, sm := range sheetManifest {
				totalRows += sm.RowCount
			}
		}

		importFile := model.ImportFile{
			ID:                uuid.New().String(),
			BatchID:           batch.ID,
			OriginalFileName:  fileHeader.Filename,
			StoredFileName:    storedFilename,
			StoragePath:       dstPath,
			FileSHA256:        sha256Hex,
			FileSizeBytes:     fileHeader.Size,
			MimeType:          "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			SheetCount:        sheetCount,
			RowCount:          totalRows,
			SheetManifestJSON: manifestJSON,
			Status:            "UPLOADED",
			SourceYear:        batch.SourceYear,
			UploadedBy:        user,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := s.importRepo.CreateFile(&importFile); err != nil {
			return nil, err
		}
		createdFiles = append(createdFiles, importFile)
	}

	_ = s.importRepo.RecalculateBatchStats(batch.ID)

	return createdFiles, nil
}

// ParseBatch extracts rows from uploaded files into imp_staged_rows.
func (s *ImportService) ParseBatch(batchID string, selectedSheetsByFile map[string][]string, headerOverrides map[string]int) (*model.ImportBatch, error) {
	batch, err := s.importRepo.GetBatchByID(batchID)
	if err != nil {
		return nil, err
	}

	// Clean existing staged rows before re-parsing
	_ = s.importRepo.DeleteStagedRowsByBatch(batch.ID)

	files, err := s.importRepo.GetFilesByBatchID(batch.ID)
	if err != nil {
		return nil, err
	}

	var allStagedRows []model.ImportStagedRow
	autoMapping := make(map[string]string)

	for _, file := range files {
		selectedSheets := selectedSheetsByFile[file.ID]
		rows, err := s.parser.ParseSheetRows(batch.ID, &file, selectedSheets, headerOverrides)
		if err != nil {
			continue
		}
		allStagedRows = append(allStagedRows, rows...)

		// Collect headers and column samples for auto-mapping per file from selected sheets
		selectedSheetMap := make(map[string]bool)
		for _, s := range selectedSheets {
			selectedSheetMap[strings.ToLower(strings.TrimSpace(s))] = true
		}

		var fileHeaders []string
		fileSamples := make(map[string][]string)

		var manifest []SheetInfo
		if err := json.Unmarshal([]byte(file.SheetManifestJSON), &manifest); err == nil {
			for _, m := range manifest {
				// Only inspect headers from selected sheets
				if len(selectedSheetMap) > 0 && !selectedSheetMap[strings.ToLower(strings.TrimSpace(m.Name))] {
					continue
				}
				fileHeaders = append(fileHeaders, m.Headers...)
				if m.ColumnSamples != nil {
					for k, v := range m.ColumnSamples {
						if fileSamples[k] == nil {
							fileSamples[k] = v
						}
					}
				}
			}
		}

		// Run AutoDetectMapping for this file's headers and samples
		fileMapping := s.normalizer.AutoDetectMapping(fileHeaders, fileSamples)
		for k, v := range fileMapping {
			if v != "IGNORE" || autoMapping[k] == "" {
				autoMapping[k] = v
			}
		}
	}

	mapJSON, _ := json.Marshal(autoMapping)
	batch.ColumnMappingJSON = string(mapJSON)
	batch.Status = model.BatchStatusParsed

	// Insert staged rows
	if len(allStagedRows) > 0 {
		if err := s.importRepo.CreateStagedRows(allStagedRows); err != nil {
			return nil, err
		}
	}

	_ = s.importRepo.UpdateBatch(batch)
	_ = s.importRepo.RecalculateBatchStats(batch.ID)

	return s.importRepo.GetBatchByID(batch.ID)
}

// ApplyColumnMapping updates column mapping and executes normalization on all staged rows.
func (s *ImportService) ApplyColumnMapping(batchID string, mapping map[string]string) (*model.ImportBatch, error) {
	batch, err := s.importRepo.GetBatchByID(batchID)
	if err != nil {
		return nil, err
	}

	mapJSON, _ := json.Marshal(mapping)
	batch.ColumnMappingJSON = string(mapJSON)

	rows, err := s.importRepo.GetAllStagedRowsInBatch(batch.ID)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	for i := range rows {
		_ = s.normalizer.NormalizeStagedRow(&rows[i], mapping)

		// Instant AI classification for seamless zero-config onboarding
		if aiOut, err := s.aiProvider.SuggestClassification(ctx, provider.RowClassificationInput{
			RawName:         rows[i].NormalizedName,
			RawCode:         rows[i].NormalizedCode,
			RawCategory:     rows[i].NormalizedCategory,
			RawManufacturer: rows[i].NormalizedMfg,
			RawSupplier:     rows[i].NormalizedSupplier,
			UnitCost:        rows[i].UnitCost,
			Currency:        rows[i].Currency,
			FileName:        rows[i].FileName,
			SheetName:       rows[i].SheetName,
		}); err == nil && aiOut != nil {
			rows[i].AISuggestedType = aiOut.ItemClassification
			rows[i].AISuggestedProdType = aiOut.ProductType
			rows[i].AISuggestedCategory = aiOut.SuggestedCategory
			rows[i].AIConfidence = aiOut.ConfidenceScore
			payloadJSON, _ := json.Marshal(aiOut)
			rows[i].AISuggestionsJSON = string(payloadJSON)

			if rows[i].ReviewDecision == model.RowDecisionPending {
				rows[i].ItemClassification = aiOut.ItemClassification
				if aiOut.ItemClassification == model.ItemClassProduct {
					rows[i].ProductType = aiOut.ProductType
					rows[i].TargetEntity = "products"
				} else if aiOut.ItemClassification == model.ItemClassComponent {
					rows[i].TargetEntity = "components"
				} else if aiOut.ItemClassification == model.ItemClassService {
					rows[i].TargetEntity = "service_items"
				}
				if rows[i].NormalizedCategory == "" && aiOut.SuggestedCategory != "" {
					rows[i].NormalizedCategory = aiOut.SuggestedCategory
				}
			}
		}

		// Smart Auto-Generated MPN / Code:
		if strings.TrimSpace(rows[i].NormalizedCode) == "" || strings.TrimSpace(rows[i].NormalizedCode) == "-" {
			rows[i].NormalizedCode = s.GenerateSmartMPN(&rows[i], batch)
		}

		s.validator.ValidateRow(&rows[i])

		// Smart Auto-Approve Rule:
		// Automatically approve items that have a concrete Name and positive Price/Cost (> 0)
		hasValidName := strings.TrimSpace(rows[i].NormalizedName) != "" &&
			strings.TrimSpace(rows[i].NormalizedName) != "Untitled" &&
			strings.TrimSpace(rows[i].NormalizedName) != "-"
		hasCommercialValue := rows[i].SellingPrice > 0 || rows[i].UnitCost > 0

		if hasValidName && hasCommercialValue && (rows[i].ReviewDecision == model.RowDecisionPending || rows[i].ReviewDecision == "") {
			rows[i].ReviewDecision = model.RowDecisionApproved
		}
	}

	// Analyze duplicates
	_ = s.duplicate.AnalyzeDuplicates(rows)

	if err := s.importRepo.BulkUpdateRows(rows); err != nil {
		return nil, err
	}

	batch.Status = model.BatchStatusValidating
	_ = s.importRepo.UpdateBatch(batch)
	_ = s.importRepo.RecalculateBatchStats(batch.ID)

	return s.importRepo.GetBatchByID(batch.ID)
}

// RunAIClassification runs heuristic and AI suggestions on staged rows.
func (s *ImportService) RunAIClassification(batchID string) (*model.ImportBatch, error) {
	batch, err := s.importRepo.GetBatchByID(batchID)
	if err != nil {
		return nil, err
	}

	rows, err := s.importRepo.GetAllStagedRowsInBatch(batch.ID)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	for i := range rows {
		row := &rows[i]
		aiOut, err := s.aiProvider.SuggestClassification(ctx, provider.RowClassificationInput{
			RawName:         row.NormalizedName,
			RawCode:         row.NormalizedCode,
			RawCategory:     row.NormalizedCategory,
			RawManufacturer: row.NormalizedMfg,
			RawSupplier:     row.NormalizedSupplier,
			UnitCost:        row.UnitCost,
			Currency:        row.Currency,
			FileName:        row.FileName,
			SheetName:       row.SheetName,
		})

		if err == nil && aiOut != nil {
			row.AISuggestedType = aiOut.ItemClassification
			row.AISuggestedProdType = aiOut.ProductType
			row.AISuggestedCategory = aiOut.SuggestedCategory
			row.AIConfidence = aiOut.ConfidenceScore
			payloadJSON, _ := json.Marshal(aiOut)
			row.AISuggestionsJSON = string(payloadJSON)

			// If row is not yet reviewed, prefill AI recommendation safely
			if row.ReviewDecision == model.RowDecisionPending {
				row.ItemClassification = aiOut.ItemClassification
				if aiOut.ItemClassification == model.ItemClassProduct {
					row.ProductType = aiOut.ProductType
					row.TargetEntity = "products"
				} else if aiOut.ItemClassification == model.ItemClassComponent {
					row.TargetEntity = "components"
				} else if aiOut.ItemClassification == model.ItemClassService {
					row.TargetEntity = "service_items"
				}
				if row.NormalizedCategory == "" && aiOut.SuggestedCategory != "" {
					row.NormalizedCategory = aiOut.SuggestedCategory
				}
			}
		}

		// Smart Auto-Generated MPN / Code:
		if strings.TrimSpace(row.NormalizedCode) == "" || strings.TrimSpace(row.NormalizedCode) == "-" {
			row.NormalizedCode = s.GenerateSmartMPN(row, batch)
		}

		// Smart Auto-Approve Rule on AI run as well
		hasValidName := strings.TrimSpace(row.NormalizedName) != "" &&
			strings.TrimSpace(row.NormalizedName) != "Untitled" &&
			strings.TrimSpace(row.NormalizedName) != "-"
		hasCommercialValue := row.SellingPrice > 0 || row.UnitCost > 0

		if hasValidName && hasCommercialValue && (row.ReviewDecision == model.RowDecisionPending || row.ReviewDecision == "") {
			row.ReviewDecision = model.RowDecisionApproved
		}
	}

	if err := s.importRepo.BulkUpdateRows(rows); err != nil {
		return nil, err
	}

	batch.Status = model.BatchStatusReadyForReview
	_ = s.importRepo.UpdateBatch(batch)
	_ = s.importRepo.RecalculateBatchStats(batch.ID)

	return s.importRepo.GetBatchByID(batch.ID)
}

// GenerateMissingCodes auto-generates missing or placeholder codes for all staged rows in the batch.
func (s *ImportService) GenerateMissingCodes(batchID string, parentProjectName string) (*model.ImportBatch, error) {
	batch, err := s.importRepo.GetBatchByID(batchID)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(parentProjectName) != "" {
		batch.SourceReference = parentProjectName
	}

	rows, err := s.importRepo.GetAllStagedRowsInBatch(batch.ID)
	if err != nil {
		return nil, err
	}

	var mapping map[string]string
	if batch.ColumnMappingJSON != "" {
		_ = json.Unmarshal([]byte(batch.ColumnMappingJSON), &mapping)
	}

	for i := range rows {
		// Re-normalize row if mapping is present to pick up corrected prices (e.g. 500.000, 150.000)
		if mapping != nil {
			_ = s.normalizer.NormalizeStagedRow(&rows[i], mapping)
		}

		if strings.TrimSpace(rows[i].NormalizedCode) == "" ||
			strings.TrimSpace(rows[i].NormalizedCode) == "-" ||
			strings.Contains(rows[i].NormalizedCode, "-SRV-") ||
			strings.Contains(rows[i].NormalizedCode, "-SEN-") ||
			strings.Contains(rows[i].NormalizedCode, "-SOL-") ||
			strings.HasPrefix(rows[i].NormalizedCode, "PROJ-") {
			rows[i].NormalizedCode = s.GenerateSmartMPN(&rows[i], batch)
		}
		s.validator.ValidateRow(&rows[i])

		// Live sync to master data if row was already committed/imported
		if rows[i].CreatedMasterID != "" {
			tx := s.importRepo.DB()
			if rows[i].TargetMasterTable == "products" || rows[i].ItemClassification == model.ItemClassProduct || rows[i].ItemClassification == model.ItemClassAccessory {
				var prod model.Product
				if err := tx.Where("id = ?", rows[i].CreatedMasterID).First(&prod).Error; err == nil {
					prod.Name = rows[i].NormalizedName
					prod.Code = rows[i].NormalizedCode
					prod.Category = rows[i].NormalizedCategory
					prod.UpdatedAt = time.Now()
					_ = tx.Save(&prod).Error

					if prod.ProductType == "TRADING" {
						var td model.TradingProductDetail
						if err := tx.Where("product_id = ?", prod.ID).First(&td).Error; err == nil {
							td.PurchasePrice = rows[i].UnitCost
							td.SellingPrice = rows[i].SellingPrice
							td.Currency = rows[i].Currency
							td.PurchaseCurrency = rows[i].Currency
							td.SellingCurrency = rows[i].Currency
							td.UpdatedAt = time.Now()
							_ = tx.Save(&td).Error
						}
					}
				}
			} else if rows[i].TargetMasterTable == "components" || rows[i].ItemClassification == model.ItemClassComponent {
				var comp model.Component
				if err := tx.Where("id = ?", rows[i].CreatedMasterID).First(&comp).Error; err == nil {
					comp.Name = rows[i].NormalizedName
					comp.PartNumber = rows[i].NormalizedCode
					comp.UpdatedAt = time.Now()
					_ = tx.Save(&comp).Error
				}
			}
		}
	}

	_ = s.duplicate.AnalyzeDuplicates(rows)

	if err := s.importRepo.BulkUpdateRows(rows); err != nil {
		return nil, err
	}

	_ = s.importRepo.RecalculateBatchStats(batch.ID)
	return s.importRepo.GetBatchByID(batch.ID)
}

// GenerateSmartMPN creates a standardized Code/MPN following: PRJ-[DOC_ID]-[PROJECT_TAG]-[SEQUENCE]
// Rules:
// 1. If Project -> "PRJ" prefix.
// 2. If Trading -> "TRD" prefix.
// 3. If RnD / Development -> "DEV" prefix.
// 4. Document code (e.g. CQC08260429) extracted from source file/reference.
// 5. Project/Customer tag (e.g. LIPO) extracted from Parent Project Name or reference.
// 6. No category tag (SEN/SRV removed).
// 7. Sequence: 4-digit row number (%04d).
func (s *ImportService) GenerateSmartMPN(row *model.ImportStagedRow, batch *model.ImportBatch) string {
	streamPrefix := s.determineStreamPrefix(batch, row)
	docCode := s.extractDocCode(batch, row)
	projTag := s.extractProjectTag(batch, row)

	rowNum := row.RowNumber
	if rowNum <= 0 {
		rowNum = 1
	}

	return fmt.Sprintf("%s-%s-%s-%04d", streamPrefix, docCode, projTag, rowNum)
}

func (s *ImportService) determineStreamPrefix(batch *model.ImportBatch, row *model.ImportStagedRow) string {
	combined := strings.ToUpper(batch.SourceReference + " " + batch.SourceNotes + " " + batch.SourceType + " " + row.FileName + " " + row.SheetName)

	// 1. If source is detected as Project or has service/project indicators:
	isProject := strings.Contains(combined, "PROJECT") ||
		strings.Contains(combined, "PROJ") ||
		strings.Contains(combined, "PRJ") ||
		strings.Contains(combined, "RAB") ||
		strings.Contains(combined, "SPK") ||
		strings.Contains(combined, "KONTRAK") ||
		strings.Contains(combined, "CONTRACT") ||
		strings.Contains(combined, "MAINTENANCE") ||
		strings.Contains(combined, "CQC") ||
		strings.Contains(combined, "INSTALLATION") ||
		strings.Contains(combined, "PEMASANGAN") ||
		strings.Contains(combined, "INTEGRATED SOLUTION") ||
		row.ProductType == "PROJECT" ||
		row.ItemClassification == model.ItemClassService

	if isProject {
		return "PRJ"
	}

	// 2. If source is detected as Trading:
	isTrading := strings.Contains(combined, "TRADING") ||
		strings.Contains(combined, "TRD") ||
		strings.Contains(combined, "RESELLER") ||
		strings.Contains(combined, "BUY-SELL") ||
		strings.Contains(combined, "PO-TRADING") ||
		row.ProductType == "TRADING"

	if isTrading {
		return "TRD"
	}

	// 3. If source is detected as RnD / Development:
	isRnD := strings.Contains(combined, "DEV") ||
		strings.Contains(combined, "RND") ||
		strings.Contains(combined, "R&D") ||
		strings.Contains(combined, "PROTOTYPE") ||
		strings.Contains(combined, "INTERNAL") ||
		row.ProductType == "RND"

	if isRnD {
		return "DEV"
	}

	return "PRJ"
}

func (s *ImportService) extractDocCode(batch *model.ImportBatch, row *model.ImportStagedRow) string {
	text := batch.SourceReference + " " + row.FileName + " " + batch.SourceNotes

	// Look for patterns like CQC08260429, SPK-001, RAB-2026, etc.
	docRegex := regexp.MustCompile(`(?i)\b(CQC[0-9A-Z_-]{4,14}|SPK[0-9A-Z_-]{3,12}|RAB[0-9A-Z_-]{3,12}|[A-Z]{2,4}[0-9]{4,12})\b`)
	if match := docRegex.FindStringSubmatch(text); len(match) > 1 {
		cleaned := sanitizeTag(match[1])
		if cleaned != "" {
			return cleaned
		}
	}

	if batch.BatchCode != "" {
		cleaned := sanitizeTag(batch.BatchCode)
		if cleaned != "" {
			return cleaned
		}
	}

	return fmt.Sprintf("CQC%s", batch.SourceYear)
}

func (s *ImportService) extractProjectTag(batch *model.ImportBatch, row *model.ImportStagedRow) string {
	text := batch.SourceReference + " " + batch.SourceNotes + " " + row.FileName + " " + row.TargetMarket

	// Pattern 1: Explicit PT / CV / PDAM / TBK entities
	ptRegex := regexp.MustCompile(`(?i)\b(?:PT\.?|CV\.?|PDAM)\s+([A-Za-z0-9]{2,16})`)
	if match := ptRegex.FindStringSubmatch(text); len(match) > 1 {
		cand := strings.ToUpper(strings.TrimSpace(match[1]))
		if !isStopWord(cand) && len(cand) >= 2 {
			return sanitizeTag(cand)
		}
	}

	tbkRegex := regexp.MustCompile(`(?i)\b([A-Za-z0-9]{2,16})\s+(?:TBK|PERSERO)\b`)
	if match := tbkRegex.FindStringSubmatch(text); len(match) > 1 {
		cand := strings.ToUpper(strings.TrimSpace(match[1]))
		if !isStopWord(cand) && len(cand) >= 2 {
			return sanitizeTag(cand)
		}
	}

	// Pattern 2: Customer / Client keyword
	custRegex := regexp.MustCompile(`(?i)\b(?:CUSTOMER|CLIENT|KLIEN|PROYEK|PROJECT)\s*[:=-]\s*([A-Za-z0-9]{2,16})`)
	if match := custRegex.FindStringSubmatch(text); len(match) > 1 {
		cand := strings.ToUpper(strings.TrimSpace(match[1]))
		if !isStopWord(cand) && len(cand) >= 2 {
			return sanitizeTag(cand)
		}
	}

	// Pattern 3: Look for tokens separated by dash/underscore/space (e.g. "RAB Maintenance SPARING - Lipo R1" -> "LIPO")
	parts := strings.FieldsFunc(text, func(r rune) bool {
		return r == '-' || r == '_' || r == ':' || r == '(' || r == ')' || r == '[' || r == ']' || r == '.'
	})
	for _, part := range parts {
		words := strings.Fields(part)
		for _, w := range words {
			upperW := strings.ToUpper(strings.TrimSpace(w))
			if !isStopWord(upperW) && len(upperW) >= 3 && len(upperW) <= 14 {
				hasLetter := false
				for _, ch := range upperW {
					if ch >= 'A' && ch <= 'Z' {
						hasLetter = true
						break
					}
				}
				// Make sure it is not the docCode itself (e.g. CQC08260429)
				if hasLetter && !strings.HasPrefix(upperW, "CQC") && !strings.HasPrefix(upperW, "IMP") {
					return sanitizeTag(upperW)
				}
			}
		}
	}

	return "CBI"
}

func isStopWord(w string) bool {
	stopWords := map[string]bool{
		"PT": true, "CV": true, "TBK": true, "PERSERO": true, "INDONESIA": true,
		"RAB": true, "MAINTENANCE": true, "SPARING": true, "PROJECT": true, "PROYEK": true,
		"PENGADAAN": true, "CONTRACT": true, "KONTRAK": true, "REVISI": true, "REV": true,
		"FINAL": true, "DRAFT": true, "COPY": true, "FILE": true, "EXCEL": true,
		"XLSX": true, "XLS": true, "CSV": true, "DATA": true, "IMPORT": true,
		"HARGA": true, "BIAYA": true, "COST": true, "PRICE": true, "OFFER": true,
		"QUOTATION": true, "PENGAJUAN": true, "JASA": true, "SERVICE": true,
		"CBI": true, "IOT": true, "UNIT": true, "AREA": true, "ITEM": true,
		"R1": true, "R2": true, "V1": true, "V2": true, "V3": true, "THE": true, "AND": true, "FOR": true,
		"DEV": true, "TRD": true, "PROJ": true, "PRJ": true, "RND": true, "INTERNAL": true,
		"EDGE": true, "GATEWAY": true, "BOARD": true, "HARDWARE": true, "SOFTWARE": true,
		"SYSTEM": true, "MODULE": true, "CONTROLLER": true, "SENSOR": true,
	}
	return stopWords[w]
}

func sanitizeTag(t string) string {
	reg := regexp.MustCompile(`[^A-Z0-9]`)
	cleaned := reg.ReplaceAllString(strings.ToUpper(t), "")
	if len(cleaned) > 16 {
		cleaned = cleaned[:16]
	}
	if cleaned == "" {
		return "CBI"
	}
	return cleaned
}

func (s *ImportService) determineCategoryTag(row *model.ImportStagedRow) string {
	nameUpper := strings.ToUpper(row.NormalizedName)
	catUpper := strings.ToUpper(row.NormalizedCategory)

	// 1. Service / Jasa
	if row.ItemClassification == model.ItemClassService ||
		strings.Contains(nameUpper, "JASA") ||
		strings.Contains(nameUpper, "SERVICE") ||
		strings.Contains(nameUpper, "MAINTENANCE") ||
		strings.Contains(nameUpper, "CLEANING") ||
		strings.Contains(nameUpper, "KALIBRASI") ||
		strings.Contains(nameUpper, "PENGECEKKAN") ||
		strings.Contains(nameUpper, "PEKERJAAN") ||
		strings.Contains(nameUpper, "INSTALLATION") ||
		strings.Contains(nameUpper, "PEMASANGAN") ||
		strings.Contains(nameUpper, "TRAINING") {
		return "SRV"
	}

	// 2. Sensor
	if strings.Contains(nameUpper, "PH") ||
		strings.Contains(nameUpper, "AMMONIUM") ||
		strings.Contains(nameUpper, "COD") ||
		strings.Contains(nameUpper, "TSS") ||
		strings.Contains(nameUpper, "SENSOR") ||
		strings.Contains(nameUpper, "TURBIDITY") ||
		strings.Contains(nameUpper, "PROBE") ||
		strings.Contains(nameUpper, "TRANSMITTER") ||
		strings.Contains(nameUpper, "FLOW") ||
		strings.Contains(nameUpper, "LEVEL") ||
		strings.Contains(nameUpper, "DO ") ||
		strings.Contains(nameUpper, "ORP") ||
		strings.Contains(catUpper, "SENSOR") {
		return "SEN"
	}

	// 3. Solution / Integrated hardware
	if strings.Contains(nameUpper, "SOLUTION") ||
		strings.Contains(catUpper, "SOLUTION") ||
		strings.Contains(nameUpper, "GATEWAY") ||
		strings.Contains(nameUpper, "CONTROLLER") ||
		strings.Contains(nameUpper, "PANEL") ||
		strings.Contains(nameUpper, "RACK") ||
		strings.Contains(nameUpper, "DATALOGGER") ||
		strings.Contains(nameUpper, "HARDWARE") {
		return "SOL"
	}

	// 4. Accessory / Spare part / Consumable
	if row.ItemClassification == model.ItemClassAccessory || strings.Contains(catUpper, "ACCESSORY") {
		return "ACC"
	}
	if row.ItemClassification == model.ItemClassSparePart || strings.Contains(catUpper, "SPARE") {
		return "SPR"
	}
	if row.ItemClassification == model.ItemClassConsumable || strings.Contains(catUpper, "MATERIAL") || strings.Contains(catUpper, "REAGENT") {
		return "MAT"
	}
	if row.ItemClassification == model.ItemClassComponent {
		return "CMP"
	}

	return "ITM"
}

// GetStagedRows returns paginated rows with filter and search.
func (s *ImportService) GetStagedRows(batchID string, status string, decision string, classification string, search string, page int, limit int) ([]model.ImportStagedRow, int64, error) {
	return s.importRepo.GetStagedRows(batchID, status, decision, classification, search, page, limit)
}

// TriageRow updates review decisions and edits for a single staged row.
func (s *ImportService) TriageRow(
	rowID string,
	itemClass string,
	prodType string,
	name string,
	code string,
	category string,
	cost float64,
	price float64,
	currency string,
	decision string,
	reviewer string,
	notes string,
	description string,
) (*model.ImportStagedRow, error) {
	row, err := s.importRepo.GetStagedRowByID(rowID)
	if err != nil {
		return nil, err
	}

	if name != "" {
		row.NormalizedName = name
	}
	if code != "" {
		row.NormalizedCode = code
	}
	if category != "" {
		row.NormalizedCategory = category
	}
	if description != "" {
		row.Description = description
	}
	if cost >= 0 {
		row.UnitCost = cost
	}
	if price >= 0 {
		row.SellingPrice = price
	}
	if currency != "" {
		row.Currency = currency
	}
	if itemClass != "" {
		row.ItemClassification = itemClass
		if itemClass == model.ItemClassProduct {
			row.TargetEntity = "products"
			row.ProductType = prodType
			if row.ProductType == "" {
				row.ProductType = "TRADING"
			}
		} else if itemClass == model.ItemClassComponent {
			row.TargetEntity = "components"
			row.ProductType = ""
		} else if itemClass == model.ItemClassService {
			row.TargetEntity = "service_items"
			row.ProductType = ""
		}
	}

	now := time.Now()
	row.ReviewDecision = decision
	row.ReviewedBy = reviewer
	row.ReviewedAt = &now
	row.ReviewNotes = notes

	// Re-validate after edits
	s.validator.ValidateRow(row)

	if err := s.importRepo.UpdateStagedRow(row); err != nil {
		return nil, err
	}

	// Live Sync to Master Data if row was already committed/imported
	if row.CreatedMasterID != "" {
		tx := s.importRepo.DB()
		if row.TargetMasterTable == "products" || row.ItemClassification == model.ItemClassProduct || row.ItemClassification == model.ItemClassAccessory {
			var prod model.Product
			if err := tx.Where("id = ?", row.CreatedMasterID).First(&prod).Error; err == nil {
				prod.Name = row.NormalizedName
				prod.Code = row.NormalizedCode
				prod.Category = row.NormalizedCategory
				if row.ProductType != "" {
					prod.ProductType = row.ProductType
				}
				prod.UpdatedAt = time.Now()
				_ = tx.Save(&prod).Error

				// If trading product, sync price
				if prod.ProductType == "TRADING" {
					var td model.TradingProductDetail
					if err := tx.Where("product_id = ?", prod.ID).First(&td).Error; err == nil {
						td.PurchasePrice = row.UnitCost
						td.SellingPrice = row.SellingPrice
						td.Currency = row.Currency
						td.PurchaseCurrency = row.Currency
						td.SellingCurrency = row.Currency
						td.UpdatedAt = time.Now()
						_ = tx.Save(&td).Error
					}
				}
			}
		} else if row.TargetMasterTable == "components" || row.ItemClassification == model.ItemClassComponent {
			var comp model.Component
			if err := tx.Where("id = ?", row.CreatedMasterID).First(&comp).Error; err == nil {
				comp.Name = row.NormalizedName
				comp.PartNumber = row.NormalizedCode
				comp.UpdatedAt = time.Now()
				_ = tx.Save(&comp).Error
			}
		}
	}

	_ = s.importRepo.RecalculateBatchStats(row.BatchID)

	return row, nil
}

// ApproveBatch marks all valid / warned rows in the batch as approved for commit.
func (s *ImportService) ApproveBatch(batchID string, user string) (*model.ImportBatch, error) {
	batch, err := s.importRepo.GetBatchByID(batchID)
	if err != nil {
		return nil, err
	}

	rows, err := s.importRepo.GetAllStagedRowsInBatch(batch.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	for i := range rows {
		row := &rows[i]
		// Rows with ERROR cannot be auto-approved
		if row.ValidationStatus != model.RowValidationError && row.ReviewDecision != model.RowDecisionRejected {
			row.ReviewDecision = model.RowDecisionApproved
			row.ReviewedBy = user
			row.ReviewedAt = &now
		}
	}

	if err := s.importRepo.BulkUpdateRows(rows); err != nil {
		return nil, err
	}

	batch.Status = model.BatchStatusApproved
	batch.ApprovedBy = user
	batch.ApprovedAt = &now
	_ = s.importRepo.UpdateBatch(batch)
	_ = s.importRepo.RecalculateBatchStats(batch.ID)

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          uuid.New().String(),
		UserID:      user,
		Module:      "dataimport",
		Action:      "APPROVE",
		EntityType:  "IMPORT_BATCH",
		EntityID:    batch.ID,
		Description: fmt.Sprintf("Formally approved import batch %s for Master Data migration", batch.BatchCode),
		CreatedAt:   time.Now(),
	})

	return s.importRepo.GetBatchByID(batch.ID)
}

// CommitApprovedBatch dispatches approved staged rows into production Master tables atomically.
func (s *ImportService) CommitApprovedBatch(batchID string, user string, parentProjectName string) (map[string]interface{}, error) {
	batch, err := s.importRepo.GetBatchByID(batchID)
	if err != nil {
		return nil, err
	}

	approvedRows, err := s.importRepo.GetApprovedRowsByBatch(batch.ID)
	if err != nil {
		return nil, err
	}

	if len(approvedRows) == 0 {
		return nil, fmt.Errorf("no approved rows available to commit in batch %s", batch.BatchCode)
	}

	var createdProds, createdComps, createdServices int
	var committedRowIDs []string
	var parentProjectID string

	// Execute transactional master data creation
	txErr := s.importRepo.DB().Transaction(func(tx *gorm.DB) error {
		// If Parent Project Name is specified, create or resolve the parent Project in master products
		if strings.TrimSpace(parentProjectName) != "" {
			parentName := strings.TrimSpace(parentProjectName)
			parentCode := fmt.Sprintf("PRJ-%s-%s", batch.SourceYear, strings.ToUpper(uuid.New().String()[:6]))

			// Extract metadata from source workbook (Customer, Company, Phone, Quotation No, Payment Terms, Notes)
			files, _ := s.importRepo.GetFilesByBatchID(batch.ID)
			var meta map[string]string
			if len(files) > 0 {
				meta = s.parser.ExtractWorkbookMetadata(files[0].StoragePath)
			}

			clientName := user
			targetMarket := "Industrial IoT"
			if meta != nil && meta["customerName"] != "" {
				clientName = meta["customerName"]
			}
			if meta != nil && meta["customerCompany"] != "" {
				targetMarket = meta["customerCompany"]
			}

			docTitle := "RENCANA ANGGARAN BIAYA MAINTENANCE & VALIDASI SPARING"
			if meta != nil && meta["documentTitle"] != "" {
				docTitle = meta["documentTitle"]
			}

			customerInfo := "Kepada Yth:\n"
			if meta != nil && meta["customerName"] != "" {
				customerInfo += fmt.Sprintf("  %s\n", meta["customerName"])
			} else {
				customerInfo += fmt.Sprintf("  %s\n", user)
			}
			if meta != nil && meta["customerCompany"] != "" {
				customerInfo += fmt.Sprintf("  %s\n", meta["customerCompany"])
			}
			if meta != nil && meta["customerPhone"] != "" {
				customerInfo += fmt.Sprintf("  Phone. : %s\n", meta["customerPhone"])
			} else {
				customerInfo += "  Phone. : -\n"
			}
			if meta != nil && meta["customerEmail"] != "" {
				customerInfo += fmt.Sprintf("  Email  : %s\n", meta["customerEmail"])
			} else {
				customerInfo += "  Email  : -\n"
			}

			docCommercial := "Informasi Dokumen & Ketentuan:\n"
			if meta != nil && meta["documentNumber"] != "" {
				docCommercial += fmt.Sprintf("  Nomor       : %s\n", meta["documentNumber"])
			}
			if meta != nil && meta["quotationDate"] != "" {
				docCommercial += fmt.Sprintf("  Tanggal     : %s\n", meta["quotationDate"])
			}
			if meta != nil && meta["reference"] != "" {
				docCommercial += fmt.Sprintf("  Reff        : %s\n", meta["reference"])
			}
			if meta != nil && meta["paymentTerms"] != "" {
				docCommercial += fmt.Sprintf("  Pembayaran  : %s\n", meta["paymentTerms"])
			}
			if meta != nil && meta["validity"] != "" {
				docCommercial += fmt.Sprintf("  Validity    : %s\n", meta["validity"])
			}
			if meta != nil && meta["currency"] != "" {
				docCommercial += fmt.Sprintf("  Currency    : %s\n", meta["currency"])
			}

			fullDesc := fmt.Sprintf("%s\n\n%s\n\n%s", docTitle, strings.TrimSpace(customerInfo), strings.TrimSpace(docCommercial))

			var existingParent model.Product
			if err := tx.Unscoped().Where("name = ? AND product_type = ?", parentName, "PROJECT").First(&existingParent).Error; err == nil {
				existingParent.DeletedAt = gorm.DeletedAt{}
				existingParent.Status = "Active"
				if clientName != "" {
					existingParent.ProjectLead = clientName
				}
				if targetMarket != "" {
					existingParent.TargetMarket = targetMarket
				}
				if fullDesc != "" {
					existingParent.Description = fullDesc
					existingParent.LongDescription = fullDesc
				}
				existingParent.UpdatedAt = time.Now()
				_ = tx.Save(&existingParent).Error
				parentProjectID = existingParent.ID
			} else {
				newParent := model.Product{
					ID:              fmt.Sprintf("PRD-%s", uuid.New().String()[:8]),
					Code:            parentCode,
					Name:            parentName,
					Category:        "Project Solutions",
					Status:          "Active",
					ProductType:     "PROJECT",
					ProjectLead:     clientName,
					TargetMarket:    targetMarket,
					Description:     fullDesc,
					LongDescription: fullDesc,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				if err := tx.Create(&newParent).Error; err == nil {
					parentProjectID = newParent.ID
					createdProds++
				}
			}
		}

		for i := range approvedRows {
			row := &approvedRows[i]
			cleanMasterName := s.stripLeadingNumberPrefix(row.NormalizedName)

			// Skip if already imported (idempotency guarantee)
			if row.FinalStatus == "IMPORTED" {
				continue
			}

			// 1. Resolve / Create Category
			catID := "CAT-001"
			if row.NormalizedCategory != "" {
				var cat model.ComponentCategory
				if err := tx.Unscoped().Where("name = ?", row.NormalizedCategory).First(&cat).Error; err != nil {
					// Create category
					cat = model.ComponentCategory{
						ID:          fmt.Sprintf("CAT-%s", uuid.New().String()[:8]),
						Name:        row.NormalizedCategory,
						Code:        strings.ToUpper(strings.ReplaceAll(row.NormalizedCategory, " ", "_")),
						Description: "Imported Category from " + batch.BatchCode,
					}
					_ = tx.Create(&cat).Error
				} else if cat.DeletedAt.Valid {
					cat.DeletedAt = gorm.DeletedAt{}
					_ = tx.Save(&cat).Error
				}
				catID = cat.ID
			} else {
				var defaultCat model.ComponentCategory
				if err := tx.First(&defaultCat).Error; err == nil {
					catID = defaultCat.ID
				}
			}

			// 2. Resolve / Create Manufacturer
			mfgID := "MFG-001"
			if row.NormalizedMfg != "" {
				var mfg model.Manufacturer
				if err := tx.Unscoped().Where("name = ?", row.NormalizedMfg).First(&mfg).Error; err != nil {
					mfg = model.Manufacturer{
						ID:   fmt.Sprintf("MFG-%s", uuid.New().String()[:8]),
						Name: row.NormalizedMfg,
						Code: strings.ToUpper(strings.ReplaceAll(row.NormalizedMfg, " ", "_")),
					}
					_ = tx.Create(&mfg).Error
				} else if mfg.DeletedAt.Valid {
					mfg.DeletedAt = gorm.DeletedAt{}
					_ = tx.Save(&mfg).Error
				}
				mfgID = mfg.ID
			} else {
				var defaultMfg model.Manufacturer
				if err := tx.First(&defaultMfg).Error; err == nil {
					mfgID = defaultMfg.ID
				}
			}

			// 3. Resolve / Create Supplier
			supplierID := "SUP-001"
			if row.NormalizedSupplier != "" {
				var sup model.Supplier
				if err := tx.Unscoped().Where("name = ?", row.NormalizedSupplier).First(&sup).Error; err != nil {
					sup = model.Supplier{
						ID:   fmt.Sprintf("SUP-%s", uuid.New().String()[:8]),
						Name: row.NormalizedSupplier,
						Code: strings.ToUpper(strings.ReplaceAll(row.NormalizedSupplier, " ", "_")),
					}
					_ = tx.Create(&sup).Error
				} else if sup.DeletedAt.Valid {
					sup.DeletedAt = gorm.DeletedAt{}
					_ = tx.Save(&sup).Error
				}
				supplierID = sup.ID
			} else {
				var defaultSup model.Supplier
				if err := tx.First(&defaultSup).Error; err == nil {
					supplierID = defaultSup.ID
				}
			}

			// 4. Resolve Unit
			unitID := "UNIT-PCS"
			var defaultUnit model.Unit
			if err := tx.Where("code = ? OR name = ?", row.NormalizedUnit, row.NormalizedUnit).First(&defaultUnit).Error; err == nil {
				unitID = defaultUnit.ID
			} else if err := tx.First(&defaultUnit).Error; err == nil {
				unitID = defaultUnit.ID
			}

			// 5. Dispatch to Target Table
			switch row.ItemClassification {
			case model.ItemClassProduct, model.ItemClassAccessory, model.ItemClassSparePart:
				prodCode := row.NormalizedCode
				if prodCode == "" {
					prodCode = s.GenerateSmartMPN(row, batch)
				}

				prodType := row.ProductType
				if prodType == "" || prodType == "PROJECT" {
					prodType = "TRADING"
				}

				// Check if product exists (idempotent update, restore, or insert)
				var existingProd model.Product
				if err := tx.Unscoped().Where("code = ?", prodCode).First(&existingProd).Error; err == nil {
					// Update and restore existing product
					existingProd.DeletedAt = gorm.DeletedAt{}
					existingProd.Name = cleanMasterName
					existingProd.TargetMarket = row.TargetMarket
					existingProd.Category = row.NormalizedCategory
					existingProd.CategoryID = &catID
					existingProd.ProductType = prodType
					existingProd.Status = "Active"
					if row.Description != "" {
						existingProd.Description = row.Description
						existingProd.LongDescription = row.Description
					}
					existingProd.UpdatedAt = time.Now()
					if err := tx.Save(&existingProd).Error; err != nil {
						return err
					}
					row.CreatedMasterID = existingProd.ID
					createdProds++

					// Upsert trading detail if TRADING product
					if prodType == "TRADING" && (row.UnitCost > 0 || row.SellingPrice > 0) {
						var tradingDetail model.TradingProductDetail
						if err := tx.Unscoped().Where("product_id = ?", existingProd.ID).First(&tradingDetail).Error; err == nil {
							tradingDetail.DeletedAt = gorm.DeletedAt{}
							tradingDetail.SupplierID = &supplierID
							tradingDetail.ManufacturerID = &mfgID
							tradingDetail.PurchasePrice = row.UnitCost
							tradingDetail.SellingPrice = row.SellingPrice
							tradingDetail.Currency = row.Currency
							tradingDetail.PurchaseCurrency = row.Currency
							tradingDetail.SellingCurrency = row.Currency
							tradingDetail.UpdatedAt = time.Now()
							_ = tx.Save(&tradingDetail).Error
						} else {
							newTrading := model.TradingProductDetail{
								ID:               fmt.Sprintf("TRD-%s", uuid.New().String()[:8]),
								ProductID:        existingProd.ID,
								SupplierID:       &supplierID,
								ManufacturerID:   &mfgID,
								PurchasePrice:    row.UnitCost,
								SellingPrice:     row.SellingPrice,
								Currency:         row.Currency,
								PurchaseCurrency: row.Currency,
								SellingCurrency:  row.Currency,
								CreatedAt:        time.Now(),
								UpdatedAt:        time.Now(),
							}
							_ = tx.Create(&newTrading).Error
						}
					}
				} else {
					// Insert new product
					newProd := model.Product{
						ID:              fmt.Sprintf("PRD-%s", uuid.New().String()[:8]),
						Code:            prodCode,
						Name:            cleanMasterName,
						TargetMarket:    row.TargetMarket,
						Category:        row.NormalizedCategory,
						CategoryID:      &catID,
						Status:          "Active",
						ProductType:     prodType,
						ProjectLead:     user,
						Description:     row.Description,
						LongDescription: row.Description,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}
					if err := tx.Create(&newProd).Error; err != nil {
						return err
					}
					row.CreatedMasterID = newProd.ID
					createdProds++

					// If TRADING product, create TradingProductDetail with purchase & selling price
					if prodType == "TRADING" && (row.UnitCost > 0 || row.SellingPrice > 0) {
						newTrading := model.TradingProductDetail{
							ID:               fmt.Sprintf("TRD-%s", uuid.New().String()[:8]),
							ProductID:        newProd.ID,
							SupplierID:       &supplierID,
							ManufacturerID:   &mfgID,
							PurchasePrice:    row.UnitCost,
							SellingPrice:     row.SellingPrice,
							Currency:         row.Currency,
							PurchaseCurrency: row.Currency,
							SellingCurrency:  row.Currency,
							CreatedAt:        time.Now(),
							UpdatedAt:        time.Now(),
						}
						_ = tx.Create(&newTrading).Error
					}
				}
				row.TargetMasterTable = "products"

			case model.ItemClassComponent, model.ItemClassConsumable:
				compPN := row.NormalizedCode
				if compPN == "" {
					compPN = s.GenerateSmartMPN(row, batch)
				}

				ipn := fmt.Sprintf("IPN-%s-%04d", batch.SourceYear, row.RowNumber)

				var existingComp model.Component
				if err := tx.Unscoped().Where("part_number = ?", compPN).First(&existingComp).Error; err == nil {
					existingComp.DeletedAt = gorm.DeletedAt{}
					existingComp.Name = cleanMasterName
					existingComp.CategoryID = catID
					existingComp.ManufacturerID = mfgID
					existingComp.SupplierID = supplierID
					existingComp.UnitID = unitID
					existingComp.Status = "Approved"
					existingComp.EstimatedUnitCost = row.UnitCost
					existingComp.Currency = row.Currency
					existingComp.ShortDescription = row.Description
					existingComp.DetailedDescription = row.Description
					existingComp.UpdatedAt = time.Now()
					if err := tx.Save(&existingComp).Error; err != nil {
						return err
					}
					row.CreatedMasterID = existingComp.ID
					createdComps++
				} else {
					newComp := model.Component{
						ID:                  fmt.Sprintf("CMP-%s", uuid.New().String()[:8]),
						PartNumber:          compPN,
						InternalCode:        ipn,
						Name:                cleanMasterName,
						CategoryID:          catID,
						ManufacturerID:      mfgID,
						SupplierID:          supplierID,
						UnitID:              unitID,
						Status:              "Approved",
						EstimatedUnitCost:   row.UnitCost,
						Currency:            row.Currency,
						ShortDescription:    row.Description,
						DetailedDescription: row.Description,
						CreatedAt:           time.Now(),
						UpdatedAt:           time.Now(),
					}
					if err := tx.Create(&newComp).Error; err != nil {
						return err
					}
					row.CreatedMasterID = newComp.ID
					createdComps++
				}
				row.TargetMasterTable = "components"

			default: // Service or other
				row.CreatedMasterID = fmt.Sprintf("SRV-%s", uuid.New().String()[:8])
				row.TargetMasterTable = "service_items"
				createdServices++
			}

			row.FinalStatus = "IMPORTED"
			committedRowIDs = append(committedRowIDs, row.ID)
			_ = tx.Save(row).Error
		}

		// Build / Sync Project Hierarchy (project_items) under parentProjectID
		if parentProjectID != "" {
			// Clear existing project items for idempotency
			tx.Where("project_product_id = ?", parentProjectID).Delete(&model.ProjectItem{})

			var currentSectionItemID *string
			var currentSubGroupID *string
			var lastSectionName string

			for i := range approvedRows {
				row := &approvedRows[i]
				cleanMasterName := s.stripLeadingNumberPrefix(row.NormalizedName)
				rawNameTrim := strings.TrimSpace(row.NormalizedName)

				notesVal := row.NormalizedCategory
				if row.Description != "" {
					notesVal = row.Description
				}

				// Check if this row is a Section Header
				isSecHeader := false
				if (len(rawNameTrim) >= 2 && rawNameTrim[0] >= 'A' && rawNameTrim[0] <= 'Z' && (rawNameTrim[1] == '.' || rawNameTrim[1] == ' ')) ||
					(len(row.NormalizedCategory) > 0 && row.NormalizedCategory != lastSectionName && (strings.HasPrefix(row.NormalizedCategory, "A.") || strings.HasPrefix(row.NormalizedCategory, "B.") || strings.HasPrefix(row.NormalizedCategory, "C."))) {
					if len(rawNameTrim) >= 2 && rawNameTrim[0] >= 'A' && rawNameTrim[0] <= 'Z' && (rawNameTrim[1] == '.' || rawNameTrim[1] == ' ') {
						isSecHeader = true
					}
				}

				if row.NormalizedCategory != "" && row.NormalizedCategory != lastSectionName && (strings.HasPrefix(row.NormalizedCategory, "A.") || strings.HasPrefix(row.NormalizedCategory, "B.") || strings.HasPrefix(row.NormalizedCategory, "C.")) {
					lastSectionName = row.NormalizedCategory
					if !isSecHeader {
						secItem := model.ProjectItem{
							ID:               fmt.Sprintf("PRJ-SEC-%s", uuid.New().String()[:8]),
							ProjectProductID: parentProjectID,
							ItemType:         model.ProjectItemTypeSubAssembly,
							Name:             row.NormalizedCategory,
							Quantity:         1.0,
							CostPrice:        0,
							SellingPrice:     0,
							Currency:         "IDR",
							SortOrder:        row.RowNumber - 1,
							Notes:            "Project Section",
							Status:           "Active",
							CreatedAt:        time.Now(),
							UpdatedAt:        time.Now(),
						}
						_ = tx.Create(&secItem).Error
						secID := secItem.ID
						currentSectionItemID = &secID
						currentSubGroupID = nil
					}
				}

				if isSecHeader {
					secItem := model.ProjectItem{
						ID:               fmt.Sprintf("PRJ-SEC-%s", uuid.New().String()[:8]),
						ProjectProductID: parentProjectID,
						ItemType:         model.ProjectItemTypeSubAssembly,
						Name:             rawNameTrim,
						Quantity:         1.0,
						CostPrice:        row.UnitCost,
						SellingPrice:     row.SellingPrice,
						Currency:         "IDR",
						SortOrder:        row.RowNumber,
						Notes:            notesVal,
						Status:           "Active",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					}
					_ = tx.Create(&secItem).Error
					secID := secItem.ID
					currentSectionItemID = &secID
					currentSubGroupID = nil
					continue
				}

				isSubGroupHeader := false
				digitPrefixRegex := regexp.MustCompile(`^\d+[.)\s]+`)
				if digitPrefixRegex.MatchString(rawNameTrim) && row.UnitCost == 0 && row.SellingPrice == 0 && row.ItemClassification == model.ItemClassService {
					isSubGroupHeader = true
				}

				if isSubGroupHeader {
					subGroupItem := model.ProjectItem{
						ID:               fmt.Sprintf("PRJ-GRP-%s", uuid.New().String()[:8]),
						ProjectProductID: parentProjectID,
						ParentItemID:     currentSectionItemID,
						ItemType:         model.ProjectItemTypeSubAssembly,
						Name:             rawNameTrim,
						Quantity:         1.0,
						CostPrice:        row.UnitCost,
						SellingPrice:     row.SellingPrice,
						Currency:         "IDR",
						SortOrder:        row.RowNumber,
						Notes:            notesVal,
						Status:           "Active",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					}
					_ = tx.Create(&subGroupItem).Error
					subID := subGroupItem.ID
					currentSubGroupID = &subID
					continue
				}

				var targetParentID *string
				if currentSubGroupID != nil {
					targetParentID = currentSubGroupID
				} else if currentSectionItemID != nil {
					targetParentID = currentSectionItemID
				}

				var prjItem model.ProjectItem
				if row.TargetMasterTable == "products" && row.CreatedMasterID != "" {
					prodCode := row.NormalizedCode
					if prodCode == "" {
						prodCode = s.GenerateSmartMPN(row, batch)
					}
					prodType := row.ProductType
					if prodType == "" || prodType == "PROJECT" {
						prodType = "TRADING"
					}
					pID := row.CreatedMasterID
					uID := row.NormalizedUnit
					curr := row.Currency
					if curr == "" {
						curr = "IDR"
					}
					prjItem = model.ProjectItem{
						ID:               fmt.Sprintf("PRJ-ITM-%s", uuid.New().String()[:8]),
						ProjectProductID: parentProjectID,
						ParentItemID:     targetParentID,
						ItemType:         model.ProjectItemTypeProduct,
						Name:             cleanMasterName,
						ProductID:        &pID,
						ProductCode:      prodCode,
						ProductType:      prodType,
						Quantity:         1.0,
						CostPrice:        row.UnitCost,
						SellingPrice:     row.SellingPrice,
						Currency:         curr,
						UnitID:           &uID,
						UnitName:         row.NormalizedUnit,
						SortOrder:        row.RowNumber,
						Notes:            notesVal,
						Status:           "Active",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					}
				} else if row.TargetMasterTable == "components" && row.CreatedMasterID != "" {
					compPN := row.NormalizedCode
					if compPN == "" {
						compPN = s.GenerateSmartMPN(row, batch)
					}
					cID := row.CreatedMasterID
					uID := row.NormalizedUnit
					curr := row.Currency
					if curr == "" {
						curr = "IDR"
					}
					prjItem = model.ProjectItem{
						ID:               fmt.Sprintf("PRJ-ITM-%s", uuid.New().String()[:8]),
						ProjectProductID: parentProjectID,
						ParentItemID:     targetParentID,
						ItemType:         model.ProjectItemTypeComponent,
						Name:             cleanMasterName,
						ComponentID:      &cID,
						ComponentMPN:     compPN,
						Quantity:         1.0,
						CostPrice:        row.UnitCost,
						SellingPrice:     row.SellingPrice,
						Currency:         curr,
						UnitID:           &uID,
						UnitName:         row.NormalizedUnit,
						SortOrder:        row.RowNumber,
						Notes:            notesVal,
						Status:           "Active",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					}
				} else if row.TargetMasterTable == "service_items" || row.ItemClassification == model.ItemClassService {
					uID := row.NormalizedUnit
					curr := row.Currency
					if curr == "" {
						curr = "IDR"
					}
					notes := "Engineering Service / Maintenance Task"
					if row.Description != "" {
						notes = row.Description
					} else if row.SellingPrice > 0 {
						notes = fmt.Sprintf("Service Fee: %s %.2f", curr, row.SellingPrice)
					}
					prjItem = model.ProjectItem{
						ID:               fmt.Sprintf("PRJ-ITM-%s", uuid.New().String()[:8]),
						ProjectProductID: parentProjectID,
						ParentItemID:     targetParentID,
						ItemType:         model.ProjectItemTypeSubAssembly,
						Name:             row.NormalizedName,
						Quantity:         1.0,
						CostPrice:        row.UnitCost,
						SellingPrice:     row.SellingPrice,
						Currency:         curr,
						UnitID:           &uID,
						UnitName:         row.NormalizedUnit,
						SortOrder:        row.RowNumber,
						Notes:            notes,
						Status:           "Active",
						CreatedAt:        time.Now(),
						UpdatedAt:        time.Now(),
					}
				}

				if prjItem.ID != "" {
					_ = tx.Create(&prjItem).Error
				}
			}
		}

		now := time.Now()
		batch.Status = model.BatchStatusImported
		batch.ImportedAt = &now
		return tx.Save(batch).Error
	})

	if txErr != nil {
		return nil, fmt.Errorf("master dispatch transaction failed: %w", txErr)
	}

	_ = s.importRepo.RecalculateBatchStats(batch.ID)

	_ = s.activityRepo.Create(&model.ActivityLog{
		ID:          uuid.New().String(),
		UserID:      user,
		Module:      "dataimport",
		Action:      "IMPORT_EXECUTE",
		EntityType:  "IMPORT_BATCH",
		EntityID:    batch.ID,
		Description: fmt.Sprintf("Committed %d items from batch %s (Products: %d, Components: %d, Services: %d)", len(committedRowIDs), batch.BatchCode, createdProds, createdComps, createdServices),
		CreatedAt:   time.Now(),
	})

	return map[string]interface{}{
		"batchCode":         batch.BatchCode,
		"status":            "IMPORTED",
		"totalCommitted":    len(committedRowIDs),
		"createdProducts":   createdProds,
		"createdComponents": createdComps,
		"createdServices":   createdServices,
		"parentProjectId":   parentProjectID,
		"parentProjectName": parentProjectName,
		"importedAt":        time.Now(),
	}, nil
}

// DeleteBatch cleans up uploaded physical files and deletes batch records.
func (s *ImportService) DeleteBatch(batchID string) error {
	files, _ := s.importRepo.GetFilesByBatchID(batchID)
	for _, f := range files {
		if f.StoragePath != "" {
			_ = os.Remove(f.StoragePath)
		}
	}
	return s.importRepo.DeleteBatch(batchID)
}

func (s *ImportService) stripLeadingNumberPrefix(name string) string {
	prefixRegex := regexp.MustCompile(`^(?i)^(\d+|[a-z])[.)\s]+(.*)$`)
	if matches := prefixRegex.FindStringSubmatch(strings.TrimSpace(name)); len(matches) >= 3 && len(strings.TrimSpace(matches[2])) > 2 {
		return strings.TrimSpace(matches[2])
	}
	return strings.TrimSpace(name)
}
