package service

import (
	"encoding/json"
	"strings"

	"iot-rd-backend/internal/model"
)

type ImportValidatorService struct{}

func NewImportValidatorService() *ImportValidatorService {
	return &ImportValidatorService{}
}

// ValidateRow performs deterministic rule validation on a normalized staged row.
func (v *ImportValidatorService) ValidateRow(row *model.ImportStagedRow) {
	var errs []string
	var warns []string

	// 1. Mandatory Item Name
	if strings.TrimSpace(row.NormalizedName) == "" || strings.TrimSpace(row.NormalizedName) == "Untitled" || strings.TrimSpace(row.NormalizedName) == "-" {
		errs = append(errs, "Missing required item or product name")
	}

	// 2. Numeric Range Validations
	if row.UnitCost < 0 {
		errs = append(errs, "Purchase cost cannot be negative")
	}
	if row.SellingPrice < 0 {
		errs = append(errs, "Selling price cannot be negative")
	}

	// 3. Currency Validation
	validCurrencies := map[string]bool{"IDR": true, "USD": true, "EUR": true, "SGD": true, "CNY": true, "JPY": true, "GBP": true}
	if row.Currency != "" && !validCurrencies[strings.ToUpper(row.Currency)] {
		errs = append(errs, "Unsupported currency code: "+row.Currency)
	}

	// 4. Missing Warnings
	if strings.TrimSpace(row.NormalizedCategory) == "" {
		warns = append(warns, "Category is unassigned; will default to General Hardware")
	}
	if strings.TrimSpace(row.NormalizedMfg) == "" && row.ItemClassification != model.ItemClassService {
		warns = append(warns, "Manufacturer / Brand not specified")
	}
	if strings.TrimSpace(row.NormalizedCode) == "" {
		warns = append(warns, "Product Code / MPN missing; system will auto-generate code upon import")
	}

	// Determine final validation state
	if len(errs) > 0 {
		row.ValidationStatus = model.RowValidationError
	} else if len(warns) > 0 {
		row.ValidationStatus = model.RowValidationWarning
	} else {
		row.ValidationStatus = model.RowValidationValid
	}

	errJSON, _ := json.Marshal(errs)
	warnJSON, _ := json.Marshal(warns)
	row.ValidationErrorsJSON = string(errJSON)
	row.ValidationWarningsJSON = string(warnJSON)
}
