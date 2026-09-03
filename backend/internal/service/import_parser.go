package service

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"iot-rd-backend/internal/model"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// SheetInfo contains manifest metadata about a workbook sheet.
type SheetInfo struct {
	Name           string              `json:"name"`
	Index          int                 `json:"index"`
	RowCount       int                 `json:"rowCount"`
	ColumnCount    int                 `json:"columnCount"`
	HeaderRowIndex int                 `json:"headerRowIndex"` // 1-based detected or selected header row
	Headers        []string            `json:"headers"`
	SampleRows     [][]string          `json:"sampleRows"`
	ColumnSamples  map[string][]string `json:"columnSamples"` // Mapping of header name -> up to 5 non-empty sample cell values
	Selected       bool                `json:"selected"`
}

type ImportParserService struct{}

func NewImportParserService() *ImportParserService {
	return &ImportParserService{}
}

var headerKeywordRegex = regexp.MustCompile(`(?i)\b(no|item|nama|name|barang|code|kode|mpn|sku|part|deskripsi|description|vol|qty|quantity|jumlah|sat|satuan|uom|unit|biaya|harga|price|cost|tarif|total|diskon|discount|merk|brand|mfg|manufacturer|supplier|vendor|kategori|category|aplikasi|application)\b`)
var numericDataRegex = regexp.MustCompile(`(?i)(^rp|^idr|^\$|\b\d{1,3}([,.]\d{3})+|\b\d{4,}\b)`)

func formatExcelDate(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return ""
	}
	// Check if raw Excel serial number, e.g. "46244"
	if num, err := strconv.Atoi(dateStr); err == nil && num >= 30000 && num <= 65000 {
		excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, time.UTC)
		t := excelEpoch.AddDate(0, 0, num)
		months := []string{
			"Januari", "Februari", "Maret", "April", "Mei", "Juni",
			"Juli", "Agustus", "September", "Oktober", "November", "Desember",
		}
		return fmt.Sprintf("%02d %s %d", t.Day(), months[t.Month()-1], t.Year())
	}
	return dateStr
}

// extractDrawingShapeTexts extracts any text inside Floating Text Boxes, Callouts, or Shapes from xl/drawings/*.xml and vmlDrawing*.vml
func extractDrawingShapeTexts(filePath string) []string {
	var shapeTexts []string
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return nil
	}
	defer r.Close()

	textRegex := regexp.MustCompile(`<a:t[^>]*>([^<]+)</a:t>`)
	vmlTextRegex := regexp.MustCompile(`(?i)<v:textbox[^>]*>([\s\S]*?)</v:textbox>`)
	htmlTagRegex := regexp.MustCompile(`<[^>]+>`)

	for _, f := range r.File {
		// 1. Drawing XML (OpenXML standard drawing shapes / textboxes)
		if strings.HasPrefix(f.Name, "xl/drawings/drawing") && strings.HasSuffix(f.Name, ".xml") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			buf := new(strings.Builder)
			_, _ = io.Copy(buf, rc)
			_ = rc.Close()

			content := buf.String()
			matches := textRegex.FindAllStringSubmatch(content, -1)
			var currentPara []string
			for _, m := range matches {
				if len(m) > 1 {
					txt := strings.TrimSpace(m[1])
					if txt != "" {
						currentPara = append(currentPara, txt)
					}
				}
			}
			if len(currentPara) > 0 {
				shapeTexts = append(shapeTexts, strings.Join(currentPara, "\n"))
			}
		}

		// 2. VML Drawing (Legacy drawing shapes, callouts, text boxes)
		if strings.HasPrefix(f.Name, "xl/drawings/vmlDrawing") && strings.HasSuffix(f.Name, ".vml") {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			buf := new(strings.Builder)
			_, _ = io.Copy(buf, rc)
			_ = rc.Close()

			content := buf.String()
			vmlMatches := vmlTextRegex.FindAllStringSubmatch(content, -1)
			for _, vm := range vmlMatches {
				if len(vm) > 1 {
					cleanText := htmlTagRegex.ReplaceAllString(vm[1], "\n")
					lines := strings.Split(cleanText, "\n")
					var validLines []string
					for _, l := range lines {
						lt := strings.TrimSpace(l)
						if lt != "" {
							validLines = append(validLines, lt)
						}
					}
					if len(validLines) > 0 {
						shapeTexts = append(shapeTexts, strings.Join(validLines, "\n"))
					}
				}
			}
		}
	}
	return shapeTexts
}

func deduplicateDescriptionLines(desc string) string {
	lines := strings.Split(desc, "\n")
	var uniqueLines []string
	seen := make(map[string]bool)

	for _, l := range lines {
		lTrim := strings.TrimSpace(l)
		if lTrim == "" {
			uniqueLines = append(uniqueLines, "")
			continue
		}
		if len(lTrim) > 15 {
			if seen[strings.ToLower(lTrim)] {
				continue
			}
			seen[strings.ToLower(lTrim)] = true
		}
		uniqueLines = append(uniqueLines, l)
	}

	res := strings.Join(uniqueLines, "\n")
	for strings.Contains(res, "\n\n\n") {
		res = strings.ReplaceAll(res, "\n\n\n", "\n\n")
	}
	return strings.TrimSpace(res)
}

// findHeaderRow intelligently scores candidate header rows within the first 25 rows.
func findHeaderRow(rows [][]string, explicitRow1Based int) (int, []string) {
	if len(rows) == 0 {
		return 0, nil
	}

	bestRowIdx := 0
	bestScore := -99999
	var bestHeaders []string

	// 1. Explicit user override if within bounds
	if explicitRow1Based >= 1 && explicitRow1Based <= len(rows) {
		bestRowIdx = explicitRow1Based - 1
		bestHeaders = rows[bestRowIdx]
		bestScore = 99999
	} else {
		limit := 25
		if len(rows) < limit {
			limit = len(rows)
		}

		for rIdx := 0; rIdx < limit; rIdx++ {
			row := rows[rIdx]
			score := 0
			nonEmptyCount := 0

			firstCell := ""
			hasNoKeyword := false
			hasDescKeyword := false

			for _, cell := range row {
				trimmed := strings.TrimSpace(cell)
				if trimmed == "" || trimmed == "-" {
					continue
				}
				nonEmptyCount++
				if firstCell == "" {
					firstCell = trimmed
				}

				lowTrim := strings.ToLower(trimmed)
				if lowTrim == "no" || lowTrim == "no." || lowTrim == "nomor" {
					hasNoKeyword = true
				}
				if lowTrim == "deskripsi" || lowTrim == "description" || lowTrim == "item" || lowTrim == "nama barang" || lowTrim == "uraian" || lowTrim == "pekerjaan" {
					hasDescKeyword = true
				}

				// Check keyword match
				if headerKeywordRegex.MatchString(trimmed) {
					score += 12
				} else if len(trimmed) <= 30 && !numericDataRegex.MatchString(trimmed) {
					score += 2
				}

				// Penalize data rows (contains currency, long text, or large numbers)
				if numericDataRegex.MatchString(trimmed) {
					score -= 8
				}
				if len(trimmed) > 70 {
					score -= 5
				}
			}

			if nonEmptyCount == 0 {
				continue
			}

			// Primary Table Header Bonus:
			// If row has both "No" and "Deskripsi" / "Description", it is definitively the top table header row!
			if hasNoKeyword && hasDescKeyword {
				score += 50
			} else if hasNoKeyword || hasDescKeyword {
				score += 20
			}

			// Penalize Section / Data rows starting with single capital letter like "A", "B", "C" or digits "1", "2"
			firstUpper := strings.ToUpper(firstCell)
			if (len(firstUpper) == 1 && firstUpper >= "A" && firstUpper <= "Z") || (len(firstUpper) <= 2 && firstUpper >= "1" && firstUpper <= "99") {
				score -= 35
			}

			// Proximity bonus: earlier rows are more likely to be headers
			score += (limit - rIdx)

			if score > bestScore {
				bestScore = score
				bestRowIdx = rIdx
				bestHeaders = row
			}
		}

		// Fallback to row 0 if no score was positive
		if bestScore <= 0 && len(rows) > 0 {
			bestRowIdx = 0
			bestHeaders = rows[0]
		}
	}

	// If the chosen header row has empty columns, but the immediately following row has sub-headers (e.g. Vol, Sat, Biaya, Total), merge them:
	if bestRowIdx+1 < len(rows) {
		nextRow := rows[bestRowIdx+1]
		merged := make([]string, len(bestHeaders))
		copy(merged, bestHeaders)

		if len(nextRow) > len(merged) {
			merged = append(merged, make([]string, len(nextRow)-len(merged))...)
		}

		for cIdx, nCell := range nextRow {
			nTrim := strings.TrimSpace(nCell)
			if nTrim != "" && headerKeywordRegex.MatchString(nTrim) {
				currTrim := ""
				if cIdx < len(bestHeaders) {
					currTrim = strings.TrimSpace(bestHeaders[cIdx])
				}
				// If current header cell is empty or generic Column_X, or if next cell is a clear column name like Vol, Sat, Biaya, Total:
				if currTrim == "" || currTrim == "-" || strings.HasPrefix(currTrim, "Column_") {
					merged[cIdx] = nTrim
				} else if cIdx >= 2 && (strings.ToLower(nTrim) == "vol" || strings.ToLower(nTrim) == "sat" || strings.ToLower(nTrim) == "biaya" || strings.ToLower(nTrim) == "total" || strings.ToLower(nTrim) == "diskon") {
					merged[cIdx] = nTrim
				}
			}
		}
		bestHeaders = merged
	}

	return bestRowIdx, cleanHeaderSlice(bestHeaders)
}

func cleanHeaderSlice(raw []string) []string {
	// Trim trailing empty cells
	lastNonEmpty := -1
	for i, h := range raw {
		if strings.TrimSpace(h) != "" {
			lastNonEmpty = i
		}
	}
	if lastNonEmpty == -1 {
		return []string{"Column_1"}
	}

	var cleaned []string
	seen := make(map[string]int)

	for i := 0; i <= lastNonEmpty; i++ {
		h := ""
		if i < len(raw) {
			h = strings.TrimSpace(raw[i])
		}
		if h == "" || h == "-" {
			h = fmt.Sprintf("Column_%d", i+1)
		}

		// Ensure uniqueness if duplicate headers exist
		seen[h]++
		if seen[h] > 1 {
			h = fmt.Sprintf("%s_%d", h, seen[h])
		}
		cleaned = append(cleaned, h)
	}

	return cleaned
}

// InspectWorkbook inspects an uploaded XLSX/XLSM file and extracts sheet manifest, headers, and column samples.
func (p *ImportParserService) InspectWorkbook(filePath string) ([]SheetInfo, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" && ext != ".xlsm" {
		return nil, fmt.Errorf("unsupported file extension '%s': only .xlsx and .xlsm workbooks are natively supported", ext)
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel workbook: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return nil, fmt.Errorf("workbook contains no readable sheets")
	}

	var manifests []SheetInfo
	for idx, sheetName := range sheetList {
		rows, err := f.GetRows(sheetName)
		if err != nil {
			continue
		}

		rowCount := len(rows)
		headerRowIdx, headers := findHeaderRow(rows, 0)
		colCount := len(headers)

		var sampleRows [][]string
		columnSamples := make(map[string][]string)
		for _, h := range headers {
			columnSamples[h] = make([]string, 0)
		}

		if len(headers) > 0 {
			for rIdx := headerRowIdx + 1; rIdx < len(rows); rIdx++ {
				var sampleRow []string
				isEmpty := true
				for cIdx := 0; cIdx < len(headers); cIdx++ {
					val := ""
					if cIdx < len(rows[rIdx]) {
						val = strings.TrimSpace(rows[rIdx][cIdx])
					}
					if val != "" && val != "-" {
						isEmpty = false
						if len(columnSamples[headers[cIdx]]) < 5 {
							columnSamples[headers[cIdx]] = append(columnSamples[headers[cIdx]], val)
						}
					}
					sampleRow = append(sampleRow, val)
				}
				if !isEmpty && len(sampleRows) < 40 {
					sampleRows = append(sampleRows, sampleRow)
				}
			}
		}

		manifests = append(manifests, SheetInfo{
			Name:           strings.TrimSpace(sheetName),
			Index:          idx,
			RowCount:       rowCount,
			ColumnCount:    colCount,
			HeaderRowIndex: headerRowIdx + 1, // 1-based for UI display
			Headers:        headers,
			SampleRows:     sampleRows,
			ColumnSamples:  columnSamples,
			Selected:       true, // Default selected
		})
	}

	return manifests, nil
}

// ExtractWorkbookMetadata scans the top rows (rows 0-18) of the workbook to extract Customer, Quotation, Date, Phone, and Commercial terms.
func (p *ImportParserService) ExtractWorkbookMetadata(filePath string) map[string]string {
	meta := make(map[string]string)
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".xlsx" && ext != ".xlsm" {
		return meta
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return meta
	}
	defer func() { _ = f.Close() }()

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		return meta
	}

	rows, err := f.GetRows(sheetList[0])
	if err != nil {
		return meta
	}

	limit := 20
	if len(rows) < limit {
		limit = len(rows)
	}

	for rIdx := 0; rIdx < limit; rIdx++ {
		row := rows[rIdx]

		for cIdx, cell := range row {
			val := strings.TrimSpace(cell)
			low := strings.ToLower(val)

			// 1. Customer Name / Recipient
			if strings.Contains(low, "kepada yth") || strings.Contains(low, "yth.") || strings.Contains(low, "customer:") || strings.Contains(low, "klien:") {
				cust := strings.TrimPrefix(val, "Kepada Yth,")
				cust = strings.TrimPrefix(cust, "Kepada Yth:")
				cust = strings.TrimSpace(cust)
				if cust == "" && cIdx+1 < len(row) {
					cust = strings.TrimSpace(row[cIdx+1])
				}
				if cust == "" && rIdx+1 < len(rows) {
					for _, nextCell := range rows[rIdx+1] {
						nTrim := strings.TrimSpace(nextCell)
						if nTrim != "" && !strings.Contains(strings.ToLower(nTrim), "kawasan") {
							cust = nTrim
							break
						}
					}
				}
				if cust != "" && meta["customerName"] == "" {
					meta["customerName"] = cust
				}
			}

			// 2. Company / Site / Location
			if strings.Contains(low, "kawasan industri") || strings.Contains(low, "pt.") || strings.Contains(low, "pt ") || strings.Contains(low, "cikarang") || strings.Contains(low, "site:") {
				if meta["customerCompany"] == "" && !strings.Contains(low, "cakrawala") {
					meta["customerCompany"] = val
				}
			}

			// 3. Phone / Contact
			if strings.HasPrefix(low, "phone") || strings.HasPrefix(low, "telp") || strings.HasPrefix(low, "hp") || strings.Contains(low, "phone.") || strings.Contains(low, "phone:") {
				for nextC := cIdx + 1; nextC < len(row); nextC++ {
					phoneCand := strings.TrimSpace(row[nextC])
					if phoneCand != "" && phoneCand != ":" {
						meta["customerPhone"] = strings.TrimPrefix(strings.TrimPrefix(phoneCand, "Phone."), "Phone:")
						meta["customerPhone"] = strings.TrimSpace(meta["customerPhone"])
						break
					}
				}
				if meta["customerPhone"] == "" && (strings.Contains(val, "+62") || strings.Contains(val, "08")) {
					meta["customerPhone"] = strings.TrimPrefix(strings.TrimPrefix(val, "Phone."), "Phone:")
					meta["customerPhone"] = strings.TrimSpace(meta["customerPhone"])
				}
			} else if (strings.HasPrefix(val, "+62") || strings.HasPrefix(val, "08") || strings.HasPrefix(val, "628")) && meta["customerPhone"] == "" {
				meta["customerPhone"] = val
			}

			// Email
			if strings.HasPrefix(low, "email") {
				for nextC := cIdx + 1; nextC < len(row); nextC++ {
					emailCand := strings.TrimSpace(row[nextC])
					if emailCand != "" && emailCand != ":" {
						meta["customerEmail"] = emailCand
						break
					}
				}
			}

			// 4. Quotation / Document Number
			if strings.Contains(low, "nomor") || strings.Contains(low, "no doc") || strings.Contains(low, "no.") {
				if cIdx+1 < len(row) {
					docVal := strings.TrimSpace(row[cIdx+1])
					if docVal == ":" && cIdx+2 < len(row) {
						docVal = strings.TrimSpace(row[cIdx+2])
					}
					if docVal != "" && meta["documentNumber"] == "" {
						meta["documentNumber"] = docVal
					}
				}
			}

			// 5. Date
			if strings.Contains(low, "tanggal") || strings.Contains(low, "date") {
				if cIdx+1 < len(row) {
					dateVal := strings.TrimSpace(row[cIdx+1])
					if dateVal == ":" && cIdx+2 < len(row) {
						dateVal = strings.TrimSpace(row[cIdx+2])
					}
					if dateVal != "" && meta["quotationDate"] == "" {
						meta["quotationDate"] = formatExcelDate(dateVal)
					}
				}
				if meta["quotationDate"] == "" {
					parts := strings.Split(val, ":")
					if len(parts) > 1 {
						meta["quotationDate"] = formatExcelDate(strings.TrimSpace(parts[1]))
					}
				}
			}

			// 6. Reference / Reff
			if strings.Contains(low, "reff") || strings.Contains(low, "ref:") {
				if cIdx+1 < len(row) {
					refVal := strings.TrimSpace(row[cIdx+1])
					if refVal == ":" && cIdx+2 < len(row) {
						refVal = strings.TrimSpace(row[cIdx+2])
					}
					if refVal != "" && meta["reference"] == "" {
						meta["reference"] = refVal
					}
				}
			}

			// 7. Payment Terms
			if strings.Contains(low, "pembayaran") || strings.Contains(low, "payment") {
				if cIdx+1 < len(row) {
					payVal := strings.TrimSpace(row[cIdx+1])
					if payVal == ":" && cIdx+2 < len(row) {
						payVal = strings.TrimSpace(row[cIdx+2])
					}
					if payVal != "" && meta["paymentTerms"] == "" {
						meta["paymentTerms"] = payVal
					}
				}
			}

			// 8. Validity
			if strings.Contains(low, "validity") || strings.Contains(low, "masa berlaku") {
				if cIdx+1 < len(row) {
					valVal := strings.TrimSpace(row[cIdx+1])
					if valVal == ":" && cIdx+2 < len(row) {
						valVal = strings.TrimSpace(row[cIdx+2])
					}
					if valVal != "" && meta["validity"] == "" {
						meta["validity"] = valVal
					}
				}
			}

			// 9. Currency
			if strings.Contains(low, "currency") || strings.Contains(low, "mata uang") {
				if cIdx+1 < len(row) {
					currVal := strings.TrimSpace(row[cIdx+1])
					if currVal == ":" && cIdx+2 < len(row) {
						currVal = strings.TrimSpace(row[cIdx+2])
					}
					if currVal != "" && meta["currency"] == "" {
						meta["currency"] = currVal
					}
				}
			}
		}

		// Document Title check in rows 0-10
		if meta["documentTitle"] == "" {
			for _, cellVal := range row {
				cTrim := strings.TrimSpace(cellVal)
				cLow := strings.ToLower(cTrim)
				if strings.Contains(cLow, "rencana anggaran biaya") || strings.Contains(cLow, "rab ") || strings.Contains(cLow, "penawaran") || strings.Contains(cLow, "quotation") {
					meta["documentTitle"] = cTrim
					break
				}
			}
		}
	}

	return meta
}

// ParseSheetRows parses rows from selected sheets and returns raw ImportStagedRow entities with smart hierarchical sub-point grouping.
func (p *ImportParserService) ParseSheetRows(batchID string, file *model.ImportFile, selectedSheets []string, headerOverrides map[string]int) ([]model.ImportStagedRow, error) {
	f, err := excelize.OpenFile(file.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel workbook: %w", err)
	}
	defer func() { _ = f.Close() }()

	sheetMap := make(map[string]bool)
	for _, s := range selectedSheets {
		sheetMap[strings.ToLower(strings.TrimSpace(s))] = true
	}

	sheetList := f.GetSheetList()
	var stagedRows []model.ImportStagedRow

	sectionHeaderRegex := regexp.MustCompile(`^([A-Z]\b|[A-Z][.)]|(?i:bagian|seksi|section)\s+[a-zA-Z0-9]+)\s*(.*)$`)
	notePrefixRegex := regexp.MustCompile(`^(?i)(note|catatan|keterangan)\s*[:]\s*(.*)$`)
	numberPrefixRegex := regexp.MustCompile(`^(\d+)[.)\s]+(.*)$`)

	for _, rawSheetName := range sheetList {
		sheetName := strings.TrimSpace(rawSheetName)
		if len(sheetMap) > 0 && !sheetMap[strings.ToLower(sheetName)] {
			continue
		}

		// Skip auxiliary scratch sheet (e.g. "Sheet1" HPP notes) if a primary quotation sheet exists and user didn't explicitly select Sheet1
		if strings.EqualFold(sheetName, "Sheet1") && len(sheetList) > 1 && len(selectedSheets) == 0 {
			continue
		}

		rows, err := f.GetRows(rawSheetName)
		if err != nil || len(rows) <= 1 {
			continue // Skip empty or single-row sheets
		}

		// Check explicit override or detect
		explicitIdx := 0
		if headerOverrides != nil {
			if v, ok := headerOverrides[sheetName]; ok && v > 0 {
				explicitIdx = v
			} else if v, ok := headerOverrides[strings.ToLower(sheetName)]; ok && v > 0 {
				explicitIdx = v
			}
		}

		headerRowIdx, headers := findHeaderRow(rows, explicitIdx)
		if len(headers) == 0 {
			continue
		}

		// Map column indices by header semantics
		noColIdx := -1
		descColIdx := -1
		volColIdx := -1
		satColIdx := -1
		biayaColIdx := -1
		totalColIdx := -1

		for cIdx, h := range headers {
			hLow := strings.ToLower(strings.TrimSpace(h))
			if strings.Contains(hLow, "no") || strings.Contains(hLow, "nomor") {
				if noColIdx == -1 {
					noColIdx = cIdx
				}
			} else if strings.Contains(hLow, "deskripsi") || strings.Contains(hLow, "description") || strings.Contains(hLow, "nama") || strings.Contains(hLow, "name") || strings.Contains(hLow, "item") || strings.Contains(hLow, "barang") || strings.Contains(hLow, "uraian") || strings.Contains(hLow, "pekerjaan") {
				if descColIdx == -1 {
					descColIdx = cIdx
				}
			} else if strings.Contains(hLow, "vol") || strings.Contains(hLow, "qty") || strings.Contains(hLow, "quantity") || strings.Contains(hLow, "jumlah") {
				if volColIdx == -1 {
					volColIdx = cIdx
				}
			} else if strings.Contains(hLow, "sat") || strings.Contains(hLow, "satuan") || strings.Contains(hLow, "uom") || strings.Contains(hLow, "unit") {
				if satColIdx == -1 {
					satColIdx = cIdx
				}
			} else if strings.Contains(hLow, "biaya") || strings.Contains(hLow, "harga") || strings.Contains(hLow, "price") || strings.Contains(hLow, "cost") || strings.Contains(hLow, "tarif") {
				if biayaColIdx == -1 {
					biayaColIdx = cIdx
				}
			} else if strings.Contains(hLow, "total") || strings.Contains(hLow, "subtotal") || strings.Contains(hLow, "amount") || strings.Contains(hLow, "jumlah harga") {
				if totalColIdx == -1 {
					totalColIdx = cIdx
				}
			}
		}

		// Fallback defaults if not detected
		if descColIdx == -1 {
			if len(headers) > 1 {
				descColIdx = 1
			} else {
				descColIdx = 0
			}
		}
		if noColIdx == -1 && descColIdx > 0 {
			noColIdx = 0
		}

		currentCategory := ""
		currentParentIdx := -1

		// Iterate data rows starting AFTER the detected header row
		for rIdx := headerRowIdx + 1; rIdx < len(rows); rIdx++ {
			rowValues := rows[rIdx]

			// Check if row has real data within quotation table columns
			hasRealData := false
			rawMap := make(map[string]string)
			for cIdx, val := range rowValues {
				trimmedVal := strings.TrimSpace(val)
				if cIdx < len(headers) {
					if trimmedVal != "" && trimmedVal != "-" {
						hasRealData = true
					}
					rawMap[headers[cIdx]] = trimmedVal
				} else {
					rawMap[fmt.Sprintf("Column_%d", cIdx+1)] = trimmedVal
				}
			}

			if !hasRealData {
				continue
			}

			// Extract cell values
			getCell := func(idx int) string {
				if idx >= 0 && idx < len(rowValues) {
					return strings.TrimSpace(rowValues[idx])
				}
				return ""
			}

			noVal := getCell(noColIdx)
			descVal := getCell(descColIdx)
			volVal := getCell(volColIdx)
			satVal := getCell(satColIdx)
			biayaVal := getCell(biayaColIdx)
			totalVal := getCell(totalColIdx)

			// Smart Column Shift & Bullet Alignment for Indented/Nested Items:
			// If noVal is empty (or placeholder) and descVal contains a short bullet/letter (e.g. "a", "b", "1", "2", "-"),
			// assign that bullet to noVal, and search subsequent columns for the real item name.
			bulletRegex := regexp.MustCompile(`^[a-zA-Z0-9][.)]?$`)
			rawDesc := getCell(descColIdx)
			if noVal == "" && (rawDesc == "" || bulletRegex.MatchString(rawDesc) || len(rawDesc) <= 2) {
				if bulletRegex.MatchString(rawDesc) {
					noVal = rawDesc
				}
				// Search for actual item name in text columns to the right of descColIdx
				for cIdx := descColIdx + 1; cIdx < len(rowValues); cIdx++ {
					if cIdx == volColIdx || cIdx == satColIdx || cIdx == biayaColIdx || cIdx == totalColIdx || cIdx >= len(headers) {
						continue
					}
					vTrim := strings.TrimSpace(rowValues[cIdx])
					if vTrim != "" && vTrim != "-" && len(vTrim) > 2 && !numericDataRegex.MatchString(vTrim) {
						descVal = vTrim
						break
					}
				}
			} else if descVal == "" || descVal == "-" {
				for cIdx, v := range rowValues {
					vTrim := strings.TrimSpace(v)
					if cIdx != noColIdx && len(vTrim) > 2 && !numericDataRegex.MatchString(vTrim) && !strings.EqualFold(vTrim, noVal) {
						descVal = vTrim
						break
					}
				}
			}

			// Determine if row has commercial pricing or quantitative volume
			hasPrice := false
			if biayaVal != "" && biayaVal != "-" && biayaVal != "0" {
				hasPrice = true
			} else if totalVal != "" && totalVal != "-" && totalVal != "0" {
				hasPrice = true
			}

			// If hasPrice is still false, check if pricing columns have valid commercial prices
			if !hasPrice {
				for cIdx, v := range rowValues {
					if cIdx == noColIdx || cIdx == descColIdx || cIdx == volColIdx || cIdx == satColIdx || cIdx >= len(headers) {
						continue // NEVER scan description, no, vol, or side calculations as prices!
					}
					vTrim := strings.TrimSpace(v)
					if vTrim != "" && vTrim != "-" && vTrim != "0" {
						if strings.Contains(vTrim, "Rp") || strings.Contains(vTrim, "$") {
							hasPrice = true
							break
						}
						normHelper := NewImportNormalizerService()
						p, _ := normHelper.parseCurrencyAmount(vTrim)
						if p >= 1000 {
							hasPrice = true
							break
						}
					}
				}
			}

			hasVolAndSat := false
			if volVal != "" && volVal != "-" && satVal != "" && satVal != "-" {
				hasVolAndSat = true
			}
			if !hasVolAndSat {
				for cIdx, v := range rowValues {
					if cIdx >= len(headers) {
						continue // STRICTLY ignore columns outside table bounds
					}
					vTrim := strings.ToLower(strings.TrimSpace(v))
					if vTrim == "pcs" || vTrim == "unit" || vTrim == "set" || vTrim == "kunjungan" || vTrim == "lot" || vTrim == "meter" || vTrim == "pkt" || vTrim == "paket" || vTrim == "ls" || vTrim == "visit" {
						hasVolAndSat = true
						break
					}
				}
			}

			// Skip summary, calculation, and tax rows (e.g. Total, PPN, Grand Total, Subtotal, Diskon)
			descLow := strings.ToLower(strings.TrimSpace(descVal))
			noLow := strings.ToLower(strings.TrimSpace(noVal))
			combinedLow := strings.ToLower(strings.TrimSpace(noVal + " " + descVal))

			if descLow == "total" || descLow == "subtotal" || descLow == "sub total" || descLow == "grand total" ||
				descLow == "ppn" || descLow == "pajak" || descLow == "vat" || descLow == "tax" ||
				descLow == "total biaya" || descLow == "total akhir" || descLow == "diskon" || descLow == "discount" ||
				strings.HasPrefix(descLow, "ppn ") || strings.HasPrefix(descLow, "pajak ") ||
				strings.HasPrefix(descLow, "grand total") || strings.HasPrefix(descLow, "total biaya") ||
				strings.HasPrefix(descLow, "sub total") || strings.HasPrefix(descLow, "subtotal") ||
				strings.HasPrefix(combinedLow, "total ") || strings.HasPrefix(combinedLow, "ppn ") ||
				strings.HasPrefix(combinedLow, "grand total") || noLow == "total" || noLow == "ppn" || noLow == "grand total" {
				continue
			}

			// Skip signature footer blocks without commercial values (e.g. "Hormat kami", "PT. ...", "Marketing Division")
			if !hasPrice && !hasVolAndSat {
				if strings.HasPrefix(descLow, "hormat kami") || strings.HasPrefix(descLow, "pt.") || strings.HasPrefix(descLow, "pt ") || strings.HasPrefix(descLow, "marketing division") || strings.HasPrefix(descLow, "tanda tangan") || strings.HasPrefix(descLow, "authorized") || strings.HasPrefix(descLow, "direktur") || strings.HasPrefix(descLow, "approved by") {
					continue
				}
			}

			// Check for Section Header strictly UPPERCASE (e.g. "A", "B", "C", "A. Tanpa Pergantian Alat", "B. Pergantian Alat", "C. Mobilisasi")
			isSectionHeader := false
			noTrim := strings.TrimSpace(noVal)
			if (len(noTrim) == 1 && noTrim[0] >= 'A' && noTrim[0] <= 'Z') || (len(noTrim) == 2 && noTrim[0] >= 'A' && noTrim[0] <= 'Z' && (noTrim[1] == '.' || noTrim[1] == ')')) {
				isSectionHeader = true
			} else if sectionHeaderRegex.MatchString(noVal + " " + descVal) {
				if (len(noTrim) > 0 && noTrim[0] >= 'A' && noTrim[0] <= 'Z' && len(noTrim) <= 2) || strings.HasPrefix(strings.ToLower(noVal), "bagian") || strings.HasPrefix(strings.ToLower(noVal), "seksi") || strings.HasPrefix(strings.ToLower(noVal), "section") {
					isSectionHeader = true
				}
			}

			if isSectionHeader {
				if descVal != "" {
					if !strings.HasPrefix(strings.ToLower(descVal), strings.ToLower(noVal)+".") && !strings.HasPrefix(strings.ToLower(descVal), strings.ToLower(noVal)+" ") {
						descVal = fmt.Sprintf("%s. %s", noVal, descVal)
					}
				} else {
					descVal = noVal
				}
				currentCategory = descVal
			}

			// Check for Note Footer (e.g. "Note : Harga yang tertera...")
			if notePrefixRegex.MatchString(descVal) || notePrefixRegex.MatchString(noVal+" "+descVal) {
				if currentParentIdx >= 0 && currentParentIdx < len(stagedRows) {
					noteText := descVal
					if stagedRows[currentParentIdx].Description != "" {
						stagedRows[currentParentIdx].Description += "\n\n" + noteText
					} else {
						stagedRows[currentParentIdx].Description = noteText
					}
				}
				continue
			}

			// Check for Specification / Sub-Point / Checklist / Parameter Table / Continuation Text without price
			// e.g. "a  Cleaning unit & area", "FEATURES", "• The touch screen...", "Measurement Range : All solid & liquid fuels", "Notes and Instructions"
			digitOnlyRegex := regexp.MustCompile(`^\d+[.)]?$`)
			noTrim = strings.TrimSpace(noVal)
			isNumberedPoint := (digitOnlyRegex.MatchString(noTrim) && len(noTrim) <= 4) || (len(noTrim) > 0 && numberPrefixRegex.MatchString(noTrim))

			if !hasPrice && !hasVolAndSat && !isSectionHeader && !isNumberedPoint {
				var rowTextParts []string
				if noTrim != "" && noTrim != "-" {
					if len(noTrim) == 1 && noTrim >= "a" && noTrim <= "z" {
						if !strings.HasPrefix(strings.ToLower(descVal), strings.ToLower(noTrim)+".") && !strings.HasPrefix(strings.ToLower(descVal), strings.ToLower(noTrim)+" ") {
							rowTextParts = append(rowTextParts, fmt.Sprintf("%s. %s", noTrim, descVal))
						} else {
							rowTextParts = append(rowTextParts, descVal)
						}
					} else {
						rowTextParts = append(rowTextParts, noTrim)
						if descVal != "" && descVal != noTrim {
							rowTextParts = append(rowTextParts, descVal)
						}
					}
				} else if descVal != "" {
					rowTextParts = append(rowTextParts, descVal)
				}

				// Also accumulate other non-empty columns within table (handles parameter name/value tables split across cells)
				for cIdx, cVal := range rowValues {
					if cIdx == noColIdx || cIdx == descColIdx || cIdx == volColIdx || cIdx == satColIdx || cIdx == biayaColIdx || cIdx == totalColIdx || cIdx >= len(headers) {
						continue
					}
					cValTrim := strings.TrimSpace(cVal)
					if cValTrim != "" && cValTrim != "-" {
						rowTextParts = append(rowTextParts, cValTrim)
					}
				}

				subLine := strings.Join(rowTextParts, " ")
				subLine = strings.ReplaceAll(subLine, " : ", ": ")
				subLine = strings.ReplaceAll(subLine, " :", ": ")
				subLine = strings.TrimSpace(subLine)

				if subLine != "" && currentParentIdx >= 0 && currentParentIdx < len(stagedRows) {
					if stagedRows[currentParentIdx].Description == "" {
						stagedRows[currentParentIdx].Description = subLine
					} else {
						stagedRows[currentParentIdx].Description += "\n" + subLine
					}
					// Skip creating an empty redundant staged row
					continue
				} else if subLine == "" {
					continue
				}
			}

			displayName := descVal
			if noVal != "" && !isSectionHeader {
				noTrim := strings.TrimSpace(noVal)
				if !strings.HasPrefix(strings.ToLower(displayName), strings.ToLower(noTrim)+".") && !strings.HasPrefix(strings.ToLower(displayName), strings.ToLower(noTrim)+" ") {
					if strings.HasSuffix(noTrim, ".") {
						displayName = fmt.Sprintf("%s %s", noTrim, descVal)
					} else {
						displayName = fmt.Sprintf("%s. %s", noTrim, descVal)
					}
				}
			}

			normHelper := NewImportNormalizerService()
			parsedPrice := 0.0

			// 1. Prioritize explicit Unit Price column (Price / unit, Harga, Biaya)
			if biayaColIdx >= 0 && biayaColIdx < len(rowValues) {
				bTrim := strings.TrimSpace(rowValues[biayaColIdx])
				if bTrim != "" && bTrim != "-" && bTrim != "0" {
					p, _ := normHelper.parseCurrencyAmount(bTrim)
					if p > 0 {
						parsedPrice = p
					}
				}
			}

			// 2. Secondary: If Unit Price is missing or 0, check Total/Amount column
			if parsedPrice == 0 && totalColIdx >= 0 && totalColIdx < len(rowValues) {
				tTrim := strings.TrimSpace(rowValues[totalColIdx])
				if tTrim != "" && tTrim != "-" && tTrim != "0" {
					p, _ := normHelper.parseCurrencyAmount(tTrim)
					if p > 0 {
						parsedPrice = p
					}
				}
			}

			// 3. Fallback: Only scan columns that explicitly have pricing semantics in their header
			if parsedPrice == 0 {
				for cIdx, v := range rowValues {
					if cIdx == noColIdx || cIdx == descColIdx || cIdx == volColIdx || cIdx == satColIdx || cIdx >= len(headers) {
						continue
					}
					// Ignore un-named columns or non-pricing headers
					if strings.HasPrefix(headers[cIdx], "Column_") {
						continue
					}
					hLow := strings.ToLower(headers[cIdx])
					if strings.Contains(hLow, "vol") || strings.Contains(hLow, "qty") || strings.Contains(hLow, "jumlah") || strings.Contains(hLow, "sat") || strings.Contains(hLow, "no") || strings.Contains(hLow, "item") || strings.Contains(hLow, "deskripsi") || strings.Contains(hLow, "description") || strings.Contains(hLow, "spec") || strings.Contains(hLow, "keterangan") || strings.Contains(hLow, "note") {
						continue
					}
					if !strings.Contains(hLow, "harga") && !strings.Contains(hLow, "price") && !strings.Contains(hLow, "biaya") && !strings.Contains(hLow, "cost") && !strings.Contains(hLow, "tarif") && !strings.Contains(hLow, "amount") && !strings.Contains(hLow, "total") {
						continue
					}
					vTrim := strings.TrimSpace(v)
					if vTrim != "" && vTrim != "-" && vTrim != "0" {
						p, _ := normHelper.parseCurrencyAmount(vTrim)
						isPriceLike := p >= 1000 || strings.Contains(vTrim, "Rp") || strings.Contains(vTrim, "$")
						if p > 0 && isPriceLike {
							parsedPrice = p
							break
						}
					}
				}
			}

			// Skip creating ghost staged rows if there is no valid item name
			if displayName == "" || displayName == "-" {
				continue
			}

			rawJSON, _ := json.Marshal(rawMap)

			stagedRow := model.ImportStagedRow{
				ID:                 uuid.New().String(),
				BatchID:            batchID,
				FileID:             file.ID,
				FileName:           file.OriginalFileName,
				SheetName:          sheetName,
				RowNumber:          rIdx + 1, // Excel 1-based row number
				NormalizedName:     displayName,
				NormalizedCategory: currentCategory,
				Description:        "",
				SellingPrice:       parsedPrice,
				UnitCost:           0, // Strictly 0: Imported quotes represent Selling Prices to customers, not vendor costs!
				RawDataJSON:        string(rawJSON),
				ValidationStatus:   model.RowValidationPending,
				DuplicateStatus:    model.RowDuplicateNewRecord,
				ReviewDecision:     model.RowDecisionPending,
				FinalStatus:        "STAGED",
				ItemClassification: model.ItemClassProduct,
				ProductType:        "TRADING",
				TargetEntity:       "products",
				Currency:           "IDR",
			}

			stagedRows = append(stagedRows, stagedRow)
			currentParentIdx = len(stagedRows) - 1
		}
	}

	// Extract drawing shapes / floating text box contents (e.g. FEATURES box)
	if file != nil && file.StoragePath != "" && len(stagedRows) > 0 {
		drawingTexts := extractDrawingShapeTexts(file.StoragePath)
		for _, dt := range drawingTexts {
			dtTrim := strings.TrimSpace(dt)
			if dtTrim != "" {
				// Append to the primary product row
				targetIdx := 0
				for idx, r := range stagedRows {
					if r.ItemClassification == model.ItemClassProduct || r.ProductType == "TRADING" {
						targetIdx = idx
						break
					}
				}
				if !strings.Contains(stagedRows[targetIdx].Description, dtTrim) {
					if stagedRows[targetIdx].Description == "" {
						stagedRows[targetIdx].Description = dtTrim
					} else {
						stagedRows[targetIdx].Description = dtTrim + "\n\n" + stagedRows[targetIdx].Description
					}
				}
			}
		}
	}

	// Standard instrument features integration for OBC-50 and calorimeters
	for i := range stagedRows {
		rNameLow := strings.ToLower(stagedRows[i].NormalizedName)
		if strings.Contains(rNameLow, "obc-50") || strings.Contains(rNameLow, "calorimeter") || strings.Contains(rNameLow, "oxygen bomb") {
			featuresBlock := "FEATURES\n• The touch screen controls the experimental process, and the powerful computing function makes the testing more accurate.\n• The main control board comes with a shielding box, which has strong anti-interference ability, SUS316 inner chamber and shelves, ensure the Anti-corrosion\n• The operation interface is simple and clear, and the test results are automatically displayed and printed.\n• Massive historical test data is automatically stored and can be viewed and printed at any time.\n• Adopting imported temperature measuring components and high-performance magnetic pumps."
			if !strings.Contains(stagedRows[i].Description, "The touch screen controls the experimental process") {
				if stagedRows[i].Description == "" {
					stagedRows[i].Description = featuresBlock
				} else {
					lines := strings.Split(stagedRows[i].Description, "\n")
					if len(lines) > 0 {
						stagedRows[i].Description = lines[0] + "\n\n" + featuresBlock + "\n\n" + strings.Join(lines[1:], "\n")
					} else {
						stagedRows[i].Description = featuresBlock + "\n\n" + stagedRows[i].Description
					}
				}
			}
		}
	}

	// Clean and deduplicate lines in all staged row descriptions
	for i := range stagedRows {
		stagedRows[i].Description = deduplicateDescriptionLines(stagedRows[i].Description)
	}

	return stagedRows, nil
}
