package service

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"iot-rd-backend/internal/model"
)

type ImportNormalizerService struct{}

func NewImportNormalizerService() *ImportNormalizerService {
	return &ImportNormalizerService{}
}

// AutoDetectMapping creates heuristic and sample-aware column mapping based on headers and sample cell values.
func (n *ImportNormalizerService) AutoDetectMapping(headers []string, samplesMap map[string][]string) map[string]string {
	mapping := make(map[string]string)
	mappedFields := make(map[string]bool)

	for _, h := range headers {
		lower := strings.ToLower(strings.TrimSpace(h))
		var samples []string
		if samplesMap != nil {
			samples = samplesMap[h]
		}

		// Skip summary columns explicitly
		if strings.Contains(lower, "total") || strings.Contains(lower, "diskon") || strings.Contains(lower, "discount") || strings.Contains(lower, "ppn") || strings.Contains(lower, "pajak") || strings.Contains(lower, "grand") {
			mapping[h] = "IGNORE"
			continue
		}

		// 1. Selling Price / Cost check by Header or Samples
		if strings.Contains(lower, "biaya") || strings.Contains(lower, "harga jual") || strings.Contains(lower, "selling") || strings.Contains(lower, "msrp") || strings.Contains(lower, "retail") || strings.Contains(lower, "pricelist") {
			if !mappedFields["selling_price"] {
				mapping[h] = "selling_price"
				mappedFields["selling_price"] = true
				continue
			}
		}

		if strings.Contains(lower, "harga beli") || strings.Contains(lower, "purchase") || strings.Contains(lower, "cost") || strings.Contains(lower, "modal") || lower == "harga" || lower == "price" {
			if !mappedFields["unit_cost"] {
				mapping[h] = "unit_cost"
				mappedFields["unit_cost"] = true
				continue
			}
		}

		// Check sample values for Currency / Prices if header is generic
		if !mappedFields["selling_price"] && len(samples) > 0 {
			currencyMatchCount := 0
			for _, s := range samples {
				if strings.Contains(s, "Rp") || strings.Contains(s, "$") || strings.Contains(s, "IDR") || strings.Contains(s, "USD") {
					currencyMatchCount++
				}
			}
			if currencyMatchCount >= 1 {
				mapping[h] = "selling_price"
				mappedFields["selling_price"] = true
				continue
			}
		}

		// 2. Unit of Measure (UoM)
		if strings.Contains(lower, "satuan") || lower == "sat" || strings.Contains(lower, "unit") || strings.Contains(lower, "uom") {
			if !mappedFields["unit"] {
				mapping[h] = "unit"
				mappedFields["unit"] = true
				continue
			}
		}
		if !mappedFields["unit"] && len(samples) > 0 {
			uomMatch := 0
			for _, s := range samples {
				sLow := strings.ToLower(strings.TrimSpace(s))
				if sLow == "pcs" || sLow == "unit" || sLow == "kunjungan" || sLow == "set" || sLow == "box" || sLow == "roll" || sLow == "meter" || sLow == "bh" || sLow == "buah" {
					uomMatch++
				}
			}
			if uomMatch >= 1 {
				mapping[h] = "unit"
				mappedFields["unit"] = true
				continue
			}
		}

		// 3. Product / Item Name
		if strings.Contains(lower, "nama") || strings.Contains(lower, "item name") || strings.Contains(lower, "product name") || strings.Contains(lower, "deskripsi") || strings.Contains(lower, "description") || lower == "item" || lower == "barang" || strings.Contains(lower, "pekerjaan") {
			if !strings.Contains(lower, "supplier") && !strings.Contains(lower, "merk") && !strings.Contains(lower, "pabrik") {
				if !mappedFields["name"] {
					mapping[h] = "name"
					mappedFields["name"] = true
					continue
				}
			}
		}

		// Check sample values for descriptive text if name is not yet mapped
		if !mappedFields["name"] && len(samples) > 0 {
			textMatchCount := 0
			for _, s := range samples {
				sTrim := strings.TrimSpace(s)
				if len(sTrim) > 3 && !numericDataRegex.MatchString(sTrim) && sTrim != "-" && sTrim != "." {
					textMatchCount++
				}
			}
			if textMatchCount >= 1 && (strings.HasPrefix(lower, "column_") || strings.Contains(lower, "alat") || strings.Contains(lower, "uraian") || strings.Contains(lower, "keterangan")) {
				mapping[h] = "name"
				mappedFields["name"] = true
				continue
			}
		}

		// 4. Code / MPN / Part Number
		if strings.Contains(lower, "kode") || strings.Contains(lower, "code") || strings.Contains(lower, "mpn") || strings.Contains(lower, "part number") || strings.Contains(lower, "sku") || lower == "pn" {
			if !mappedFields["code"] {
				mapping[h] = "code"
				mappedFields["code"] = true
				continue
			}
		}

		// 5. Category
		if strings.Contains(lower, "kategori") || strings.Contains(lower, "category") || strings.Contains(lower, "kelompok") || strings.Contains(lower, "group") || strings.Contains(lower, "tipe barang") {
			if !mappedFields["category"] {
				mapping[h] = "category"
				mappedFields["category"] = true
				continue
			}
		}

		// 6. Manufacturer / Brand
		if strings.Contains(lower, "merk") || strings.Contains(lower, "brand") || strings.Contains(lower, "manufacturer") || strings.Contains(lower, "pabrikan") || strings.Contains(lower, "mfg") {
			if !mappedFields["manufacturer"] {
				mapping[h] = "manufacturer"
				mappedFields["manufacturer"] = true
				continue
			}
		}

		// 7. Supplier
		if strings.Contains(lower, "supplier") || strings.Contains(lower, "vendor") || strings.Contains(lower, "distributor") || strings.Contains(lower, "pemasok") || strings.Contains(lower, "toko") {
			if !mappedFields["supplier"] {
				mapping[h] = "supplier"
				mappedFields["supplier"] = true
				continue
			}
		}

		// 8. Target Market / Application
		if strings.Contains(lower, "aplikasi") || strings.Contains(lower, "application") || strings.Contains(lower, "target") || strings.Contains(lower, "kegunaan") {
			if !mappedFields["target_market"] {
				mapping[h] = "target_market"
				mappedFields["target_market"] = true
				continue
			}
		}

		mapping[h] = "IGNORE"
	}

	return mapping
}

// NormalizeStagedRow applies mapping and cleans data into canonical fields.
func (n *ImportNormalizerService) NormalizeStagedRow(row *model.ImportStagedRow, mapping map[string]string) error {
	var rawMap map[string]string
	if err := json.Unmarshal([]byte(row.RawDataJSON), &rawMap); err != nil {
		return err
	}

	normMap := make(map[string]interface{})

	for sourceHeader, targetField := range mapping {
		rawVal, exists := rawMap[sourceHeader]
		if !exists || strings.TrimSpace(rawVal) == "" {
			continue
		}
		rawValTrimmed := strings.TrimSpace(rawVal)

		switch targetField {
		case "name":
			cleaned := n.cleanString(rawValTrimmed)
			if cleaned != "-" && cleaned != "." && cleaned != "—" && cleaned != "_" {
				row.NormalizedName = cleaned
				normMap["name"] = row.NormalizedName
			}

		case "code":
			row.NormalizedCode = strings.ToUpper(n.cleanString(rawValTrimmed))
			normMap["code"] = row.NormalizedCode

		case "category":
			row.NormalizedCategory = n.cleanString(rawValTrimmed)
			normMap["category"] = row.NormalizedCategory

		case "manufacturer":
			row.NormalizedMfg = n.cleanString(rawValTrimmed)
			normMap["manufacturer"] = row.NormalizedMfg

		case "supplier":
			row.NormalizedSupplier = n.cleanString(rawValTrimmed)
			normMap["supplier"] = row.NormalizedSupplier

		case "unit":
			row.NormalizedUnit = n.normalizeUoM(rawValTrimmed)
			normMap["unit"] = row.NormalizedUnit

		case "unit_cost":
			cost, curr := n.parseCurrencyAmount(rawValTrimmed)
			row.UnitCost = cost
			if curr != "" {
				row.Currency = curr
			}
			normMap["unitCost"] = row.UnitCost
			normMap["currency"] = row.Currency

		case "selling_price":
			price, curr := n.parseCurrencyAmount(rawValTrimmed)
			row.SellingPrice = price
			if curr != "" {
				row.Currency = curr
			}
			normMap["sellingPrice"] = row.SellingPrice

		case "target_market":
			row.TargetMarket = n.cleanString(rawValTrimmed)
			normMap["targetMarket"] = row.TargetMarket

		case "description":
			cleaned := n.cleanMultilineString(rawValTrimmed)
			if row.Description != "" && len(row.Description) > len(cleaned) {
				// Retain rich multi-line specification accumulated by parser
				normMap["description"] = row.Description
			} else {
				row.Description = cleaned
				normMap["description"] = row.Description
			}
		}
	}

	// Smart Fallback for Price / Cost (e.g. when "Rp" is separated from numeric column "5.000.000" or generic column):
	if row.SellingPrice == 0 {
		for hKey, cellVal := range rawMap {
			// Skip un-named columns (Column_X) and non-price headers
			if strings.HasPrefix(hKey, "Column_") {
				continue
			}
			hLow := strings.ToLower(hKey)
			if strings.Contains(hLow, "diskon") || strings.Contains(hLow, "ppn") || strings.Contains(hLow, "vol") || strings.Contains(hLow, "qty") || strings.Contains(hLow, "jumlah") || strings.Contains(hLow, "no") || strings.Contains(hLow, "sat") || strings.Contains(hLow, "deskripsi") || strings.Contains(hLow, "description") || strings.Contains(hLow, "item") || strings.Contains(hLow, "nama") {
				continue
			}
			if !strings.Contains(hLow, "harga") && !strings.Contains(hLow, "price") && !strings.Contains(hLow, "biaya") && !strings.Contains(hLow, "cost") && !strings.Contains(hLow, "tarif") && !strings.Contains(hLow, "amount") && !strings.Contains(hLow, "total") {
				continue
			}
			cTrim := strings.TrimSpace(cellVal)
			if cTrim != "" && cTrim != "-" && cTrim != "0" {
				price, curr := n.parseCurrencyAmount(cTrim)
				isPriceLike := price >= 100 || strings.Contains(cTrim, "Rp") || strings.Contains(cTrim, "$") || strings.Contains(cTrim, ".") || strings.Contains(cTrim, ",")
				if price > 0 && isPriceLike {
					row.SellingPrice = price
					row.UnitCost = 0 // Strictly 0 for sales quotations!
					if curr != "" {
						row.Currency = curr
					}
					normMap["sellingPrice"] = row.SellingPrice
					normMap["unitCost"] = 0
					normMap["currency"] = row.Currency
					break
				}
			}
		}
	}

	// Smart Fallback for Missing or Single-Bullet Name (e.g. "a", "b", "c", "1."):
	singleCharRegex := regexp.MustCompile(`^(?i)^[a-z0-9][.)]?$`)
	if row.NormalizedName == "" || row.NormalizedName == "Untitled" || row.NormalizedName == "-" || len(row.NormalizedName) <= 2 || singleCharRegex.MatchString(row.NormalizedName) {
		var bestFallback string
		var fallbackScore int

		for headerKey, cellVal := range rawMap {
			cTrim := strings.TrimSpace(cellVal)
			cLow := strings.ToLower(cTrim)
			hLow := strings.ToLower(headerKey)

			// Skip numerical, formula, or summary columns
			if len(cTrim) < 2 || cTrim == "-" || cTrim == "." || cTrim == "—" || cTrim == "_" || singleCharRegex.MatchString(cTrim) {
				continue
			}
			if strings.Contains(hLow, "total") || strings.Contains(hLow, "diskon") || strings.Contains(hLow, "sat") || strings.Contains(hLow, "vol") || strings.Contains(hLow, "ppn") || strings.Contains(hLow, "biaya") || strings.Contains(hLow, "harga") {
				continue
			}
			if numericDataRegex.MatchString(cTrim) || strings.Contains(cTrim, "Rp") || strings.Contains(cTrim, "$") {
				continue
			}

			// Score candidate
			score := len(cTrim)
			if strings.Contains(cLow, "sensor") || strings.Contains(cLow, "aquas") || strings.Contains(cLow, "photonic") || strings.Contains(cLow, "ammonium") || strings.Contains(cLow, "logger") || strings.Contains(cLow, "mekanik") || strings.Contains(cLow, "cable") || strings.Contains(cLow, "kabel") || strings.Contains(cLow, "debit") {
				score += 80
			}
			if strings.Contains(cLow, "jasa") || strings.Contains(cLow, "maintenance") || strings.Contains(cLow, "validasi") || strings.Contains(cLow, "alat") || strings.Contains(cLow, "pekerjaan") || strings.Contains(cLow, "cleaning") || strings.Contains(cLow, "teknisi") || strings.Contains(cLow, "layanan") || strings.Contains(cLow, "support") || strings.Contains(cLow, "cost") {
				score += 50
			}
			if strings.Contains(hLow, "deskripsi") || strings.Contains(hLow, "nama") || strings.Contains(hLow, "uraian") || strings.Contains(hLow, "item") {
				score += 30
			}

			if score > fallbackScore {
				fallbackScore = score
				bestFallback = cTrim
			}
		}

		if bestFallback != "" {
			row.NormalizedName = n.cleanString(bestFallback)
			normMap["name"] = row.NormalizedName
		}
	}

	// Ensure row.NormalizedName maintains the original document numbering for Staging Review:
	noVal := ""
	for hKey, cVal := range rawMap {
		hLow := strings.ToLower(hKey)
		if hLow == "no" || hLow == "no." || hLow == "nomor" || hLow == "column_1" || hLow == "column_2" {
			cTrim := strings.TrimSpace(cVal)
			if (len(cTrim) <= 3 && singleCharRegex.MatchString(cTrim)) || regexp.MustCompile(`^\d+[.)]?$`).MatchString(cTrim) {
				noVal = cTrim
				break
			}
		}
	}
	if noVal != "" && row.NormalizedName != "" {
		noTrim := strings.TrimSpace(noVal)
		if !strings.HasPrefix(strings.ToLower(row.NormalizedName), strings.ToLower(noTrim)+".") && !strings.HasPrefix(strings.ToLower(row.NormalizedName), strings.ToLower(noTrim)+" ") {
			if strings.HasSuffix(noTrim, ".") {
				row.NormalizedName = fmt.Sprintf("%s %s", noTrim, row.NormalizedName)
			} else {
				row.NormalizedName = fmt.Sprintf("%s. %s", noTrim, row.NormalizedName)
			}
			normMap["name"] = row.NormalizedName
		}
	}

	if row.Description != "" {
		normMap["description"] = row.Description
	}

	if row.NormalizedUnit == "" {
		row.NormalizedUnit = "PCS"
	}
	if row.Currency == "" {
		row.Currency = "IDR"
	}

	normJSON, _ := json.Marshal(normMap)
	row.NormalizedDataJSON = string(normJSON)

	return nil
}

func (n *ImportNormalizerService) cleanString(val string) string {
	// Remove excessive whitespaces
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(strings.TrimSpace(val), " ")
}

func (n *ImportNormalizerService) cleanMultilineString(val string) string {
	lines := strings.Split(val, "\n")
	var cleanedLines []string
	hSpace := regexp.MustCompile(`[ \t\r]+`)
	for _, l := range lines {
		trimmed := hSpace.ReplaceAllString(strings.TrimSpace(l), " ")
		if trimmed != "" {
			cleanedLines = append(cleanedLines, trimmed)
		}
	}
	return strings.Join(cleanedLines, "\n")
}

func (n *ImportNormalizerService) normalizeUoM(raw string) string {
	lower := strings.ToLower(strings.TrimSpace(raw))
	switch lower {
	case "pcs", "pc", "piece", "pieces", "buah", "bh", "unit", "un":
		return "PCS"
	case "set", "st", "paket", "kit":
		return "SET"
	case "box", "bx", "dus":
		return "BOX"
	case "roll", "rl", "gulung":
		return "ROLL"
	case "meter", "m", "mtr":
		return "METER"
	case "cm", "centimeter":
		return "CM"
	case "kg", "kilogram":
		return "KG"
	default:
		return strings.ToUpper(raw)
	}
}

func (n *ImportNormalizerService) parseCurrencyAmount(raw string) (float64, string) {
	currency := "IDR"
	cleaned := strings.TrimSpace(raw)
	if cleaned == "" || cleaned == "-" || cleaned == "0" {
		return 0, currency
	}

	// Multi-line cells represent descriptions/specifications, never a single currency amount
	if strings.Contains(cleaned, "\n") {
		return 0, currency
	}

	rawLow := strings.ToLower(cleaned)

	// Filter out technical specifications and engineering measurement units
	if strings.Contains(rawLow, "j/g") || strings.Contains(rawLow, "kcal") || strings.Contains(rawLow, "°c") ||
		strings.Contains(rawLow, "ͦc") || strings.Contains(rawLow, "kg") || strings.Contains(rawLow, "cm") ||
		strings.Contains(rawLow, "mm") || strings.Contains(rawLow, "±") || strings.Contains(rawLow, "rh") ||
		strings.Contains(rawLow, "%") || strings.Contains(rawLow, "ppm") || strings.Contains(rawLow, "m/s") ||
		strings.Contains(rawLow, "db") || strings.Contains(rawLow, "μg") || strings.Contains(rawLow, "ug") ||
		strings.Contains(rawLow, "mg/") || strings.Contains(rawLow, "220v") || strings.Contains(rawLow, "360w") ||
		strings.Contains(rawLow, "watt") || strings.Contains(rawLow, "volt") || strings.Contains(rawLow, "minutes") ||
		strings.Contains(rawLow, "menit") || strings.Contains(rawLow, "jam") || strings.Contains(rawLow, "quality") ||
		strings.Contains(rawLow, "weight") || strings.Contains(rawLow, "hz") || strings.Contains(rawLow, "vac") ||
		strings.Contains(rawLow, "vdc") || strings.Contains(rawLow, "ampere") || strings.Contains(rawLow, "bar") ||
		strings.Contains(rawLow, "psi") || strings.Contains(rawLow, "kpa") || strings.Contains(rawLow, "mpa") ||
		strings.Contains(rawLow, "ip66") || strings.Contains(rawLow, "ip67") || strings.Contains(rawLow, "ip68") ||
		strings.Contains(rawLow, "carrier board") || strings.Contains(rawLow, "storage") || strings.Contains(rawLow, "psu") {
		return 0, currency
	}

	// Filter engineering rating values like "24V to 12V / 5V" or "50Hz" or "500mA"
	specRatingRegex := regexp.MustCompile(`(?i)\b\d+\s*(v|w|a|va|kv|kw|ma|rpm|fps|kbps|mbps)\b`)
	if specRatingRegex.MatchString(cleaned) && !strings.Contains(strings.ToUpper(cleaned), "RP") && !strings.Contains(cleaned, "$") {
		return 0, currency
	}

	if strings.HasPrefix(strings.ToUpper(cleaned), "USD") || strings.HasPrefix(cleaned, "$") {
		currency = "USD"
	} else if strings.HasPrefix(strings.ToUpper(cleaned), "EUR") || strings.HasPrefix(cleaned, "€") {
		currency = "EUR"
	} else if strings.HasPrefix(strings.ToUpper(cleaned), "SGD") || strings.HasPrefix(cleaned, "S$") {
		currency = "SGD"
	} else if strings.HasPrefix(strings.ToUpper(cleaned), "CNY") || strings.HasPrefix(strings.ToUpper(cleaned), "RMB") || strings.HasPrefix(cleaned, "¥") {
		currency = "CNY"
	}

	// Remove non-numeric symbols except decimal dots/commas
	re := regexp.MustCompile(`[^0-9.,]`)
	numericStr := re.ReplaceAllString(cleaned, "")

	if numericStr == "" {
		return 0, currency
	}

	// Handle Indonesian thousand separator "4.500.000,00" or English "4,500,000.00"
	if strings.Contains(numericStr, ".") && strings.Contains(numericStr, ",") {
		lastDot := strings.LastIndex(numericStr, ".")
		lastComma := strings.LastIndex(numericStr, ",")
		if lastComma > lastDot {
			// Indonesian style: 4.500.000,50 -> 4500000.50
			numericStr = strings.ReplaceAll(numericStr, ".", "")
			numericStr = strings.ReplaceAll(numericStr, ",", ".")
		} else {
			// English style: 4,500,000.50 -> 4500000.50
			numericStr = strings.ReplaceAll(numericStr, ",", "")
		}
	} else if strings.Contains(numericStr, ".") && strings.Count(numericStr, ".") > 1 {
		// Indonesian multiple dots: 15.625.311 -> 15625311
		numericStr = strings.ReplaceAll(numericStr, ".", "")
	} else if strings.Contains(numericStr, ",") && strings.Count(numericStr, ",") > 1 {
		// English multiple commas: 4,500,000 -> 4500000
		numericStr = strings.ReplaceAll(numericStr, ",", "")
	} else if strings.Contains(numericStr, ".") && !strings.Contains(numericStr, ",") {
		// Exactly one dot: check if it's a thousand separator (e.g. 500.000, 150.000) or decimal (e.g. 12.50)
		dotIdx := strings.Index(numericStr, ".")
		decimals := numericStr[dotIdx+1:]
		if len(decimals) == 3 && (currency == "IDR" || strings.Contains(strings.ToUpper(raw), "RP") || (!strings.Contains(strings.ToUpper(raw), "USD") && !strings.Contains(strings.ToUpper(raw), "EUR") && !strings.Contains(raw, "$"))) {
			// Indonesian thousand separator with single dot: 500.000 -> 500000, 150.000 -> 150000
			numericStr = strings.ReplaceAll(numericStr, ".", "")
		}
	} else if strings.Contains(numericStr, ",") && !strings.Contains(numericStr, ".") {
		commaIdx := strings.Index(numericStr, ",")
		decimals := numericStr[commaIdx+1:]
		if len(decimals) == 3 {
			// Single comma thousand separator (e.g. Rp500,000 -> 500000, 150,000 -> 150000, $500,000 -> 500000)
			numericStr = strings.ReplaceAll(numericStr, ",", "")
		} else {
			// Decimal comma: 450,50 -> 450.50, 12,5 -> 12.5
			numericStr = strings.ReplaceAll(numericStr, ",", ".")
		}
	}

	val, err := strconv.ParseFloat(numericStr, 64)
	if err != nil {
		return 0, currency
	}

	// Reject unreasonable runaway numbers from concatenated text numbers (e.g. 400504080045) unless explicitly marked with currency symbols
	if val > 99999999999 && !strings.Contains(strings.ToUpper(raw), "RP") && !strings.Contains(raw, "$") && !strings.Contains(strings.ToUpper(raw), "USD") && !strings.Contains(strings.ToUpper(raw), "EUR") {
		return 0, currency
	}

	return val, currency
}
