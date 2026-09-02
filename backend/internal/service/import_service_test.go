package service

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"iot-rd-backend/internal/model"

	"github.com/xuri/excelize/v2"
)

func TestSmartMPNGeneration(t *testing.T) {
	svc := &ImportService{}

	t.Run("Project Batch - Strictly PRJ-[DOC_ID]-[PROJECT_TAG]-[SEQUENCE]", func(t *testing.T) {
		batch := &model.ImportBatch{
			ID:              "batch-proj-1",
			SourceYear:      "2026",
			SourceType:      "EXCEL_WORKBOOK",
			SourceReference: "CQC08260429 - RAB Maintenance SPARING - Lipo R1",
		}

		rows := []struct {
			name        string
			itemClass   string
			prodType    string
			category    string
			rowNumber   int
			expected    string
		}{
			{
				name:        "Jasa Preventive Maintenance & Kalibrasi Online",
				itemClass:   model.ItemClassService,
				prodType:    "",
				category:    "Engineering Services",
				rowNumber:   20,
				expected:    "PRJ-CQC08260429-LIPO-0020",
			},
			{
				name:        "Cleaning unit & area",
				itemClass:   model.ItemClassService,
				prodType:    "",
				category:    "Engineering Services",
				rowNumber:   21,
				expected:    "PRJ-CQC08260429-LIPO-0021",
			},
			{
				name:        "Cleaning sensor",
				itemClass:   model.ItemClassService,
				prodType:    "",
				category:    "Engineering Services",
				rowNumber:   22,
				expected:    "PRJ-CQC08260429-LIPO-0022",
			},
			{
				name:        "- pH",
				itemClass:   model.ItemClassProduct,
				prodType:    "PROJECT",
				category:    "General Hardware",
				rowNumber:   31,
				expected:    "PRJ-CQC08260429-LIPO-0031",
			},
			{
				name:        "- Ammonium",
				itemClass:   model.ItemClassProduct,
				prodType:    "PROJECT",
				category:    "General Hardware",
				rowNumber:   32,
				expected:    "PRJ-CQC08260429-LIPO-0032",
			},
			{
				name:        "- COD & TSS",
				itemClass:   model.ItemClassProduct,
				prodType:    "PROJECT",
				category:    "General Hardware",
				rowNumber:   34,
				expected:    "PRJ-CQC08260429-LIPO-0034",
			},
		}

		for _, tc := range rows {
			row := &model.ImportStagedRow{
				NormalizedName:     tc.name,
				ItemClassification: tc.itemClass,
				ProductType:        tc.prodType,
				NormalizedCategory: tc.category,
				RowNumber:          tc.rowNumber,
			}
			code := svc.GenerateSmartMPN(row, batch)
			if code != tc.expected {
				t.Errorf("For %s (Row #%d): expected %s, got %s", tc.name, tc.rowNumber, tc.expected, code)
			}
			// Strict assertion: No SEN or SRV in the new format
			if strings.Contains(code, "-SEN-") || strings.Contains(code, "-SRV-") {
				t.Errorf("Code must not contain -SEN- or -SRV-, got %s", code)
			}
		}
	})

	t.Run("Trading Batch - Strictly TRD prefix with customer extraction", func(t *testing.T) {
		batch := &model.ImportBatch{
			ID:              "batch-trd-1",
			SourceYear:      "2026",
			BatchCode:       "TRD-2026-0001",
			SourceType:      "TRADING",
			SourceReference: "PO-TRADING PT Indofood CBP",
		}

		row := &model.ImportStagedRow{
			NormalizedName:     "Sensor Turbidity Optical",
			ItemClassification: model.ItemClassProduct,
			ProductType:        "TRADING",
			NormalizedCategory: "Sensors",
			RowNumber:          5,
		}

		code := svc.GenerateSmartMPN(row, batch)
		expected := "TRD-TRD20260001-INDOFOOD-0005"
		if code != expected {
			t.Errorf("expected %s, got %s", expected, code)
		}
	})

	t.Run("RnD Batch - Strictly DEV prefix with CBI fallback", func(t *testing.T) {
		batch := &model.ImportBatch{
			ID:              "batch-dev-1",
			SourceYear:      "2026",
			BatchCode:       "DEV-2026-0001",
			SourceType:      "RND_INTERNAL",
			SourceReference: "DEV - IoT Edge Gateway Board Rev 3",
		}

		row := &model.ImportStagedRow{
			NormalizedName:     "Gateway Edge Controller PCB",
			ItemClassification: model.ItemClassProduct,
			ProductType:        "RND",
			NormalizedCategory: "Controller Hardware",
			RowNumber:          1,
		}

		code := svc.GenerateSmartMPN(row, batch)
		expected := "DEV-DEV20260001-CBI-0001"
		if code != expected {
			t.Errorf("expected %s, got %s", expected, code)
		}
	})
}

func TestCurrencyParsing(t *testing.T) {
	norm := NewImportNormalizerService()

	testCases := []struct {
		input       string
		expectedVal float64
		expectedCur string
	}{
		{"Rp 500.000", 500000, "IDR"},
		{"Rp 150.000", 150000, "IDR"},
		{"Rp 350.000", 350000, "IDR"},
		{"Rp500,000", 500000, "IDR"},
		{"Rp150,000", 150000, "IDR"},
		{"Rp3,500,000", 3500000, "IDR"},
		{"Rp 15.625.311", 15625311, "IDR"},
		{"Rp 3.500.000", 3500000, "IDR"},
		{"500.000", 500000, "IDR"},
		{"150.000", 150000, "IDR"},
		{"500,000", 500000, "IDR"},
		{"150,000", 150000, "IDR"},
		{"4.500.000,50", 4500000.50, "IDR"},
		{"$12.50", 12.50, "USD"},
		{"$1,500.00", 1500.00, "USD"},
		{"$500,000", 500000.00, "USD"},
	}

	for _, tc := range testCases {
		val, cur := norm.parseCurrencyAmount(tc.input)
		if val != tc.expectedVal {
			t.Errorf("For %s: expected value %f, got %f", tc.input, tc.expectedVal, val)
		}
		if cur != tc.expectedCur {
			t.Errorf("For %s: expected currency %s, got %s", tc.input, tc.expectedCur, cur)
		}
	}
}

func TestSmartHierarchicalSubPointParsing(t *testing.T) {
	parser := NewImportParserService()
	svc := &ImportService{}

	// Create in-memory workbook to test Excel sub-point extraction
	f := excelize.NewFile()
	sheet := "RAB Maintenance"
	f.NewSheet(sheet)
	f.DeleteSheet("Sheet1")

	// Set headers
	headers := []string{"No", "Deskripsi", "Vol", "Sat", "Biaya", "Total"}
	for i, h := range headers {
		col, _ := excelize.CoordinatesToCellName(i+1, 1)
		_ = f.SetCellValue(sheet, col, h)
	}

	// Data rows
	data := [][]string{
		{"A", "Tanpa Pergantian Alat (sensor)", "", "", "", ""},
		{"1.", "Jasa Preventive Maintenance (Lingkup Pekerjaan)", "1", "Kunjungan", "5.000.000", "5.000.000"},
		{"a", "Cleaning unit & area", "", "", "", ""},
		{"b", "Cleaning sensor", "", "", "", ""},
		{"c", "Pengecekkan & pengukuran power sistem", "", "", "", ""},
		{"2.", "Validasi Sensor+CRM ( Teknisi CBI )", "1", "Kunjungan", "10.000.000", "10.000.000"},
		{"", "Pekerjaan kalibrasi dilakukan terhadap Sensor Aquas:", "", "", "", ""},
		{"", "- pH", "", "", "", ""},
		{"", "- Ammonium", "", "", "", ""},
		{"B", "Pergantian Alat (Sensor)", "", "", "", ""},
		{"1.", "Penggantian & repair sparepart yang rusak", "", "", "", ""},
		{"a", "Sensor Aquas (pH)", "1", "Pcs", "15.625.311", "opsional"},
		{"b", "Sensor Photonic (COD, TSS)", "1", "Pcs", "157.616.701", "opsional"},
		{"2", "Perbaikan data pembacaan dilakukan secara remote", "", "", "", ""},
		{"", "Layanan support perbaikan pembacaan data sparing ( secara remote ) - ( request di luar jam kerja & hari kerja )", "1", "Kejadian", "350.000", "Opsional"},
		{"", "Total", "", "", "20.650.000", "20.650.000"},
		{"", "PPN 11%", "", "", "2.271.500", "2.271.500"},
		{"", "Grand Total", "", "", "22.921.500", "22.921.500"},
		{"", "Hormat kami,", "", "", "", ""},
		{"", "PT. Cakrawala Bima Instrument", "", "", "", ""},
	}

	for rIdx, row := range data {
		for cIdx, val := range row {
			cell, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+2)
			_ = f.SetCellValue(sheet, cell, val)
		}
	}

	tmpFile := filepath.Join(os.TempDir(), "test_smart_parse.xlsx")
	if err := f.SaveAs(tmpFile); err != nil {
		t.Fatalf("Failed to save test xlsx: %v", err)
	}
	defer os.Remove(tmpFile)

	impFile := &model.ImportFile{
		ID:               "f-1",
		StoragePath:      tmpFile,
		OriginalFileName: "test_smart_parse.xlsx",
	}

	stagedRows, err := parser.ParseSheetRows("batch-1", impFile, []string{sheet}, nil)
	if err != nil {
		t.Fatalf("ParseSheetRows failed: %v", err)
	}

	// Expect 9 staged rows including Item B. 2 and Remote Support
	if len(stagedRows) != 9 {
		t.Fatalf("Expected 9 staged rows with full section headers & merged subpoints, got %d", len(stagedRows))
	}

	// Assert Row 0: Section A
	if stagedRows[0].NormalizedName != "A. Tanpa Pergantian Alat (sensor)" {
		t.Errorf("Row 0 name mismatch: %s", stagedRows[0].NormalizedName)
	}

	// Assert Row 1: Item 1 with sub-bullets in description
	if !strings.Contains(stagedRows[1].NormalizedName, "Jasa Preventive Maintenance") {
		t.Errorf("Row 1 name mismatch: %s", stagedRows[1].NormalizedName)
	}
	if !strings.Contains(stagedRows[1].Description, "Cleaning unit & area") || !strings.Contains(stagedRows[1].Description, "Cleaning sensor") {
		t.Errorf("Row 1 description missing subpoints: %s", stagedRows[1].Description)
	}

	// Assert Row 2: Item 2 with calibration notes in description
	if !strings.Contains(stagedRows[2].NormalizedName, "Validasi Sensor+CRM") {
		t.Errorf("Row 2 name mismatch: %s", stagedRows[2].NormalizedName)
	}
	if !strings.Contains(stagedRows[2].Description, "- pH") || !strings.Contains(stagedRows[2].Description, "- Ammonium") {
		t.Errorf("Row 2 description missing subpoints: %s", stagedRows[2].Description)
	}

	// Assert Row 3: Section B
	if stagedRows[3].NormalizedName != "B. Pergantian Alat (Sensor)" {
		t.Errorf("Row 3 name mismatch: %s", stagedRows[3].NormalizedName)
	}

	// Assert Row 4: Subgroup 1
	if !strings.Contains(stagedRows[4].NormalizedName, "Penggantian & repair sparepart") {
		t.Errorf("Row 4 name mismatch: %s", stagedRows[4].NormalizedName)
	}

	// Assert Row 5 & 6: Priced spareparts in Staging (contains numbering) and Master (clean without numbering)
	if !strings.Contains(stagedRows[5].NormalizedName, "Sensor Aquas (pH)") {
		t.Errorf("Row 5 name mismatch: %s", stagedRows[5].NormalizedName)
	}
	if !strings.Contains(stagedRows[6].NormalizedName, "Sensor Photonic (COD, TSS)") {
		t.Errorf("Row 6 name mismatch: %s", stagedRows[6].NormalizedName)
	}

	// Assert stripLeadingNumberPrefix cleans for Master Product database:
	cleanProd5 := svc.stripLeadingNumberPrefix(stagedRows[5].NormalizedName)
	if cleanProd5 != "Sensor Aquas (pH)" {
		t.Errorf("Expected clean master name 'Sensor Aquas (pH)', got '%s'", cleanProd5)
	}

	// Assert Row 7: Item B. 2 (Plain number 2)
	if !strings.Contains(stagedRows[7].NormalizedName, "Perbaikan data pembacaan") {
		t.Errorf("Row 7 name mismatch: %s", stagedRows[7].NormalizedName)
	}

	// Assert Row 8: Remote Support Service with price
	if !strings.Contains(stagedRows[8].NormalizedName, "Layanan support perbaikan") {
		t.Errorf("Row 8 name mismatch: %s", stagedRows[8].NormalizedName)
	}
}
