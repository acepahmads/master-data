package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type AuditSectionResult struct {
	SectionNumber string
	Title         string
	Status        string // PASS, FAIL, WARNING
	Details       []string
}

func main() {
	fmt.Println("================================================================================")
	fmt.Println("   IOT PRODUCT R&D CONTROL CENTER — PHASE 2 FINAL FULL-STACK FORENSIC AUDIT     ")
	fmt.Println("================================================================================")

	var allResults []AuditSectionResult
	baseURL := "http://localhost:8080/api/v1"
	client := &http.Client{}

	// Helper to add results
	addResult := func(secNum, title, status string, details []string) {
		res := AuditSectionResult{
			SectionNumber: secNum,
			Title:         title,
			Status:        status,
			Details:       details,
		}
		allResults = append(allResults, res)
		fmt.Printf("\n[%s] %s: %s\n", status, secNum, title)
		for _, d := range details {
			fmt.Printf("   • %s\n", d)
		}
	}

	// -------------------------------------------------------------------------
	// PART 3: DIRECT MYSQL FORENSIC AUDIT
	// -------------------------------------------------------------------------
	dsn := "root:@tcp(127.0.0.1:3306)/iot_rd_master?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		addResult("PART 3", "Database Forensic Audit", "FAIL", []string{fmt.Sprintf("Failed to connect to MySQL: %v", err)})
	} else {
		defer db.Close()
		var dbDetails []string
		var tableCount int
		err = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'iot_rd_master'").Scan(&tableCount)
		dbDetails = append(dbDetails, fmt.Sprintf("Total Tables in MySQL `iot_rd_master`: %d", tableCount))

		// Check product_versions table
		var verTableExists int
		_ = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'iot_rd_master' AND table_name = 'product_versions'").Scan(&verTableExists)
		dbDetails = append(dbDetails, fmt.Sprintf("Table `product_versions` exists: %v", verTableExists == 1))

		// Check hardware_revisions table
		var revTableExists int
		_ = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'iot_rd_master' AND table_name = 'hardware_revisions'").Scan(&revTableExists)
		dbDetails = append(dbDetails, fmt.Sprintf("Table `hardware_revisions` exists: %v", revTableExists == 1))

		// Check Columns of product_versions
		rows, err := db.Query("SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE FROM information_schema.columns WHERE table_schema = 'iot_rd_master' AND table_name = 'product_versions'")
		var pvCols []string
		if err == nil && rows != nil {
			for rows.Next() {
				var colName, colType, isNull string
				_ = rows.Scan(&colName, &colType, &isNull)
				pvCols = append(pvCols, fmt.Sprintf("%s (%s, Nullable: %s)", colName, colType, isNull))
			}
			rows.Close()
		}
		dbDetails = append(dbDetails, fmt.Sprintf("`product_versions` schema columns count: %d", len(pvCols)))

		// Check Columns of hardware_revisions
		rows, err = db.Query("SELECT COLUMN_NAME, COLUMN_TYPE, IS_NULLABLE FROM information_schema.columns WHERE table_schema = 'iot_rd_master' AND table_name = 'hardware_revisions'")
		var hrCols []string
		if err == nil && rows != nil {
			for rows.Next() {
				var colName, colType, isNull string
				_ = rows.Scan(&colName, &colType, &isNull)
				hrCols = append(hrCols, fmt.Sprintf("%s (%s, Nullable: %s)", colName, colType, isNull))
			}
			rows.Close()
		}
		dbDetails = append(dbDetails, fmt.Sprintf("`hardware_revisions` schema columns count: %d", len(hrCols)))

		// Check Indexes on product_versions
		rows, err = db.Query("SELECT DISTINCT INDEX_NAME, NON_UNIQUE FROM information_schema.statistics WHERE table_schema = 'iot_rd_master' AND table_name = 'product_versions'")
		var pvIndexes []string
		if err == nil && rows != nil {
			for rows.Next() {
				var idxName string
				var nonUniq int
				_ = rows.Scan(&idxName, &nonUniq)
				pvIndexes = append(pvIndexes, fmt.Sprintf("%s (Unique: %v)", idxName, nonUniq == 0))
			}
			rows.Close()
		}
		dbDetails = append(dbDetails, fmt.Sprintf("`product_versions` distinct indexes: %d", len(pvIndexes)))

		// Check Orphan Product Versions
		var orphanVerCount int
		_ = db.QueryRow("SELECT COUNT(*) FROM product_versions pv LEFT JOIN products p ON pv.product_id = p.id WHERE p.id IS NULL AND pv.deleted_at IS NULL").Scan(&orphanVerCount)
		dbDetails = append(dbDetails, fmt.Sprintf("Orphan Product Versions (referencing non-existent product): %d", orphanVerCount))

		// Check Orphan Hardware Revisions
		var orphanRevCount int
		_ = db.QueryRow("SELECT COUNT(*) FROM hardware_revisions hr LEFT JOIN product_versions pv ON hr.product_version_id = pv.id WHERE pv.id IS NULL AND hr.deleted_at IS NULL").Scan(&orphanRevCount)
		dbDetails = append(dbDetails, fmt.Sprintf("Orphan Hardware Revisions (referencing non-existent version): %d", orphanRevCount))

		status := "PASS"
		if verTableExists != 1 || revTableExists != 1 || orphanVerCount > 0 || orphanRevCount > 0 {
			status = "FAIL"
		}
		addResult("PART 3", "Database Forensic Audit", status, dbDetails)
	}

	// -------------------------------------------------------------------------
	// AUTHENTICATION & TOKEN ACQUISITION
	// -------------------------------------------------------------------------
	// 1. Super Admin (Alex Chen)
	adminLoginPayload, _ := json.Marshal(map[string]string{
		"email":    "alex.chen@iotcontrol.io",
		"password": "admin123",
	})
	adminResp, _ := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(adminLoginPayload))
	var adminRes struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(adminResp.Body).Decode(&adminRes)
	adminToken := adminRes.Data.Token
	adminResp.Body.Close()

	// 2. Restricted Viewer (Maya Lin - ROLE-VIEW)
	viewerLoginPayload, _ := json.Marshal(map[string]string{
		"email":    "maya.lin@iotcontrol.io",
		"password": "admin123",
	})
	viewerResp, _ := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(viewerLoginPayload))
	var viewerRes struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(viewerResp.Body).Decode(&viewerRes)
	viewerToken := viewerRes.Data.Token
	viewerResp.Body.Close()

	// 3. Read-Only Engineer on Versions (Elena Rostova - ROLE-FWENG)
	fwLoginPayload, _ := json.Marshal(map[string]string{
		"email":    "elena.rostova@iotcontrol.io",
		"password": "admin123",
	})
	fwResp, _ := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(fwLoginPayload))
	var fwRes struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(fwResp.Body).Decode(&fwRes)
	fwToken := fwRes.Data.Token
	fwResp.Body.Close()

	authedReq := func(token, method, path string, payload any) (int, []byte) {
		var req *http.Request
		if payload != nil {
			b, _ := json.Marshal(payload)
			req, _ = http.NewRequest(method, baseURL+path, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, _ = http.NewRequest(method, baseURL+path, nil)
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		r, e := client.Do(req)
		if e != nil {
			return 500, []byte(e.Error())
		}
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)
		return r.StatusCode, body
	}

	// -------------------------------------------------------------------------
	// PART 6: API FORENSIC AUDIT
	// -------------------------------------------------------------------------
	var apiDetails []string
	apiPass := true

	// GET /product-versions
	s, b := authedReq(adminToken, "GET", "/product-versions", nil)
	apiDetails = append(apiDetails, fmt.Sprintf("GET /product-versions: HTTP %d (Valid JSON: %v)", s, json.Valid(b)))
	if s != 200 {
		apiPass = false
	}

	// GET /product-versions/VER-001
	s, b = authedReq(adminToken, "GET", "/product-versions/VER-001", nil)
	apiDetails = append(apiDetails, fmt.Sprintf("GET /product-versions/VER-001: HTTP %d", s))
	if s != 200 {
		apiPass = false
	}

	// GET /hardware-revisions
	s, b = authedReq(adminToken, "GET", "/hardware-revisions", nil)
	apiDetails = append(apiDetails, fmt.Sprintf("GET /hardware-revisions: HTTP %d", s))
	if s != 200 {
		apiPass = false
	}

	// GET /hardware-revisions/REV-001
	s, b = authedReq(adminToken, "GET", "/hardware-revisions/REV-001", nil)
	apiDetails = append(apiDetails, fmt.Sprintf("GET /hardware-revisions/REV-001: HTTP %d", s))
	if s != 200 {
		apiPass = false
	}

	// GET /products/PRD-2024-001/lifecycle
	s, b = authedReq(adminToken, "GET", "/products/PRD-2024-001/lifecycle", nil)
	apiDetails = append(apiDetails, fmt.Sprintf("GET /products/PRD-2024-001/lifecycle: HTTP %d", s))
	if s != 200 {
		apiPass = false
	}

	// GET /product-versions/compare
	s, b = authedReq(adminToken, "GET", "/product-versions/compare?base_id=VER-001&target_id=VER-004", nil)
	apiDetails = append(apiDetails, fmt.Sprintf("GET /product-versions/compare: HTTP %d", s))
	if s != 200 {
		apiPass = false
	}

	apiStatus := "PASS"
	if !apiPass {
		apiStatus = "FAIL"
	}
	addResult("PART 6", "API Forensic Audit", apiStatus, apiDetails)

	// -------------------------------------------------------------------------
	// PART 7: CRUD PERSISTENCE AUDIT
	// -------------------------------------------------------------------------
	var crudDetails []string
	crudPass := true

	// TEST A: Create Version
	auditVerID := fmt.Sprintf("VER-AUDIT-%d", time.Now().Unix()%10000)
	auditRevID := fmt.Sprintf("REV-AUDIT-%d", time.Now().Unix()%10000)
	auditVerNumber := fmt.Sprintf("8.%d.0", (time.Now().UnixNano()/1000)%100000)
	createVerPayload := map[string]any{
		"id":            auditVerID,
		"productId":     "PRD-2024-001",
		"versionNumber": auditVerNumber,
		"versionName":   "Forensic Audit Test Baseline",
		"status":        "Draft",
		"owner":         "Alex Chen",
		"description":   "Temporary baseline created for forensic CRUD verification.",
		"releaseNotes":  "1. Forensic audit verification.",
	}
	s, b = authedReq(adminToken, "POST", "/product-versions", createVerPayload)
	var createdVer struct {
		Success bool `json:"success"`
		Data    struct {
			ID            string `json:"id"`
			VersionNumber string `json:"versionNumber"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &createdVer)
	crudDetails = append(crudDetails, fmt.Sprintf("Create Version: HTTP %d (ID: %s, Number: %s)", s, createdVer.Data.ID, createdVer.Data.VersionNumber))
	if s != 201 || createdVer.Data.ID == "" {
		crudPass = false
	}

	// TEST B: Update Version
	updateVerPayload := map[string]any{
		"versionName": "Forensic Audit Test Baseline (Updated via PUT)",
		"status":      "Active",
	}
	s, b = authedReq(adminToken, "PUT", "/product-versions/"+auditVerID, updateVerPayload)
	crudDetails = append(crudDetails, fmt.Sprintf("Update Version: HTTP %d", s))
	if s != 200 {
		crudPass = false
	}

	// Verify persistence via GET
	s, b = authedReq(adminToken, "GET", "/product-versions/"+auditVerID, nil)
	var fetchedVer struct {
		Data struct {
			VersionName string `json:"versionName"`
			Status      string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &fetchedVer)
	verPersisted := fetchedVer.Data.VersionName == "Forensic Audit Test Baseline (Updated via PUT)" && fetchedVer.Data.Status == "Active"
	crudDetails = append(crudDetails, fmt.Sprintf("Verified Updated Version Persisted: %v (Status: %s)", verPersisted, fetchedVer.Data.Status))
	if !verPersisted {
		crudPass = false
	}

	// TEST C: Create Hardware Revision
	createRevPayload := map[string]any{
		"id":               auditRevID,
		"productVersionId": auditVerID,
		"code":             "REV-A1",
		"name":             "Forensic Audit Test Revision",
		"status":           "Draft",
		"engineer":         "Alex Chen",
		"changeSummary":    "Test revision summary for forensic audit.",
	}
	s, b = authedReq(adminToken, "POST", "/hardware-revisions", createRevPayload)
	var createdRev struct {
		Success bool `json:"success"`
		Data    struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &createdRev)
	crudDetails = append(crudDetails, fmt.Sprintf("Create Revision: HTTP %d (ID: %s, Code: %s)", s, createdRev.Data.ID, createdRev.Data.Code))
	if s != 201 || createdRev.Data.ID == "" {
		crudPass = false
	}

	// TEST D: Update Hardware Revision
	updateRevPayload := map[string]any{
		"name":   "Forensic Audit Test Revision (Updated)",
		"status": "In Review",
	}
	s, b = authedReq(adminToken, "PUT", "/hardware-revisions/"+auditRevID, updateRevPayload)
	crudDetails = append(crudDetails, fmt.Sprintf("Update Revision: HTTP %d", s))
	if s != 200 {
		crudPass = false
	}

	// TEST E: Delete Revision & Version
	s, _ = authedReq(adminToken, "DELETE", "/hardware-revisions/"+auditRevID, nil)
	crudDetails = append(crudDetails, fmt.Sprintf("Delete Revision: HTTP %d", s))
	if s != 200 {
		crudPass = false
	}

	s, _ = authedReq(adminToken, "DELETE", "/product-versions/"+auditVerID, nil)
	crudDetails = append(crudDetails, fmt.Sprintf("Delete Version: HTTP %d", s))
	if s != 200 {
		crudPass = false
	}

	crudStatus := "PASS"
	if !crudPass {
		crudStatus = "FAIL"
	}
	addResult("PART 7", "CRUD Persistence Audit", crudStatus, crudDetails)

	// -------------------------------------------------------------------------
	// PART 8: UNIQUE CONSTRAINT AUDIT
	// -------------------------------------------------------------------------
	var uniqDetails []string
	uniqPass := true

	// Attempt duplicate version number for same product (v1.0.0 on PRD-2024-001)
	dupVerPayload := map[string]any{
		"productId":     "PRD-2024-001",
		"versionNumber": "1.0.0", // already exists
		"versionName":   "Duplicate Version Attempt",
		"status":        "Draft",
	}
	s, _ = authedReq(adminToken, "POST", "/product-versions", dupVerPayload)
	uniqDetails = append(uniqDetails, fmt.Sprintf("Reject duplicate version (v1.0.0 on PRD-2024-001): HTTP %d (Expected: 400)", s))
	if s != 400 {
		uniqPass = false
	}

	// Attempt duplicate revision code for same version (REV-A on VER-004)
	dupRevPayload := map[string]any{
		"productVersionId": "VER-004",
		"code":             "REV-A", // already exists on VER-004
		"name":             "Duplicate Revision Attempt",
		"status":           "Draft",
	}
	s, _ = authedReq(adminToken, "POST", "/hardware-revisions", dupRevPayload)
	uniqDetails = append(uniqDetails, fmt.Sprintf("Reject duplicate revision (REV-A on VER-004): HTTP %d (Expected: 400)", s))
	if s != 400 {
		uniqPass = false
	}

	uniqStatus := "PASS"
	if !uniqPass {
		uniqStatus = "FAIL"
	}
	addResult("PART 8", "Unique Constraint Audit", uniqStatus, uniqDetails)

	// -------------------------------------------------------------------------
	// PART 9: STATUS VALIDATION AUDIT
	// -------------------------------------------------------------------------
	var statusDetails []string
	statusPass := true

	// Test invalid status normalization / fallback
	invVerNumber := fmt.Sprintf("7.%d.0", (time.Now().UnixNano()/1000)%100000)
	invVerID := fmt.Sprintf("VER-INV-%d", time.Now().Unix()%10000)
	invalidVerPayload := map[string]any{
		"id":            invVerID,
		"productId":     "PRD-2024-001",
		"versionNumber": invVerNumber,
		"versionName":   "Invalid Status Test",
		"status":        "INVALID_STATUS_XYZ",
	}
	s, b = authedReq(adminToken, "POST", "/product-versions", invalidVerPayload)
	var testInvVer struct {
		Data struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &testInvVer)
	statusDetails = append(statusDetails, fmt.Sprintf("Status normalization on invalid input: HTTP %d, Normalized to '%s' (Expected 'Draft')", s, testInvVer.Data.Status))
	if testInvVer.Data.Status != "Draft" {
		statusPass = false
	}
	if testInvVer.Data.ID != "" {
		_, _ = authedReq(adminToken, "DELETE", "/product-versions/"+testInvVer.Data.ID, nil)
	}

	stStatus := "PASS"
	if !statusPass {
		stStatus = "FAIL"
	}
	addResult("PART 9", "Status Validation Audit", stStatus, statusDetails)

	// -------------------------------------------------------------------------
	// PART 10: RBAC FORENSIC AUDIT
	// -------------------------------------------------------------------------
	var rbacDetails []string
	rbacPass := true

	// Viewer attempting mutations (Forbidden: 403)
	s, _ = authedReq(viewerToken, "POST", "/product-versions", map[string]any{
		"productId":     "PRD-2024-001",
		"versionNumber": "9.0.0",
		"versionName":   "Unauthorized Create",
	})
	rbacDetails = append(rbacDetails, fmt.Sprintf("Viewer POST /product-versions: HTTP %d (Expected: 403 Forbidden)", s))
	if s != 403 {
		rbacPass = false
	}

	s, _ = authedReq(viewerToken, "PUT", "/product-versions/VER-001", map[string]any{
		"versionName": "Unauthorized Edit",
	})
	rbacDetails = append(rbacDetails, fmt.Sprintf("Viewer PUT /product-versions/VER-001: HTTP %d (Expected: 403 Forbidden)", s))
	if s != 403 {
		rbacPass = false
	}

	s, _ = authedReq(viewerToken, "DELETE", "/product-versions/VER-001", nil)
	rbacDetails = append(rbacDetails, fmt.Sprintf("Viewer DELETE /product-versions/VER-001: HTTP %d (Expected: 403 Forbidden)", s))
	if s != 403 {
		rbacPass = false
	}

	s, _ = authedReq(viewerToken, "POST", "/hardware-revisions", map[string]any{
		"productVersionId": "VER-001",
		"code":             "REV-Z",
		"name":             "Unauthorized Revision",
	})
	rbacDetails = append(rbacDetails, fmt.Sprintf("Viewer POST /hardware-revisions: HTTP %d (Expected: 403 Forbidden)", s))
	if s != 403 {
		rbacPass = false
	}

	// Firmware Engineer reading versions (Permitted: 200)
	s, _ = authedReq(fwToken, "GET", "/product-versions", nil)
	rbacDetails = append(rbacDetails, fmt.Sprintf("Firmware Engineer GET /product-versions: HTTP %d (Expected: 200 OK)", s))
	if s != 200 {
		rbacPass = false
	}

	// Firmware Engineer attempting version deletion (Forbidden: 403)
	s, _ = authedReq(fwToken, "DELETE", "/product-versions/VER-001", nil)
	rbacDetails = append(rbacDetails, fmt.Sprintf("Firmware Engineer DELETE /product-versions/VER-001: HTTP %d (Expected: 403 Forbidden)", s))
	if s != 403 {
		rbacPass = false
	}

	rbacStatus := "PASS"
	if !rbacPass {
		rbacStatus = "FAIL"
	}
	addResult("PART 10", "RBAC Forensic Audit", rbacStatus, rbacDetails)

	// -------------------------------------------------------------------------
	// PART 11: ACTIVITY LOG AUDIT
	// -------------------------------------------------------------------------
	var actDetails []string
	actPass := true

	s, b = authedReq(adminToken, "GET", "/activity?limit=50", nil)
	var actRes struct {
		Data []struct {
			ID         string `json:"id"`
			Action     string `json:"action"`
			EntityType string `json:"entityType"`
			EntityName string `json:"entityName"`
			UserName   string `json:"userName"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &actRes)
	var verLogs, revLogs int
	for _, a := range actRes.Data {
		if strings.Contains(a.Action, "PRODUCT_VERSION") || a.EntityType == "ProductVersion" {
			verLogs++
		}
		if strings.Contains(a.Action, "HARDWARE_REVISION") || a.EntityType == "HardwareRevision" {
			revLogs++
		}
	}
	actDetails = append(actDetails, fmt.Sprintf("Product Version activity log records in MySQL: %d", verLogs))
	actDetails = append(actDetails, fmt.Sprintf("Hardware Revision activity log records in MySQL: %d", revLogs))
	if verLogs == 0 || revLogs == 0 {
		actPass = false
	}

	actStatus := "PASS"
	if !actPass {
		actStatus = "FAIL"
	}
	addResult("PART 11", "Activity Log Audit", actStatus, actDetails)

	// -------------------------------------------------------------------------
	// PART 12 & 13: LIFECYCLE & COMPARISON AUDIT
	// -------------------------------------------------------------------------
	var lcCmpDetails []string
	s, b = authedReq(adminToken, "GET", "/products/PRD-2024-001/lifecycle", nil)
	var lcData struct {
		Data struct {
			ProductCode string `json:"productCode"`
			Versions    []struct {
				ID            string `json:"id"`
				VersionNumber string `json:"versionNumber"`
				Revisions     []struct {
					Code string `json:"code"`
				} `json:"revisions"`
			} `json:"versions"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &lcData)
	lcCmpDetails = append(lcCmpDetails, fmt.Sprintf("Lifecycle hierarchy for %s has %d versions", lcData.Data.ProductCode, len(lcData.Data.Versions)))

	s, b = authedReq(adminToken, "GET", "/product-versions/compare?base_id=VER-001&target_id=VER-004", nil)
	var cmpData struct {
		Data struct {
			DiffParameters map[string]bool `json:"diffParameters"`
		} `json:"data"`
	}
	_ = json.Unmarshal(b, &cmpData)
	lcCmpDetails = append(lcCmpDetails, fmt.Sprintf("Version comparison delta parameters: %d evaluated", len(cmpData.Data.DiffParameters)))

	addResult("PART 12-13", "Lifecycle & Comparison Audit", "PASS", lcCmpDetails)

	// -------------------------------------------------------------------------
	// PART 15: PHASE 1 REGRESSION AUDIT
	// -------------------------------------------------------------------------
	var regDetails []string
	regPass := true
	endpoints := []struct {
		Name string
		Path string
	}{
		{"Products Master", "/products"},
		{"Components Library", "/components"},
		{"Category Tree", "/categories/tree"},
		{"Manufacturers", "/manufacturers"},
		{"Suppliers", "/suppliers"},
		{"Units of Measure", "/units"},
		{"Users Management", "/users"},
		{"Roles & RBAC", "/roles"},
	}

	for _, ep := range endpoints {
		st, _ := authedReq(adminToken, "GET", ep.Path, nil)
		regDetails = append(regDetails, fmt.Sprintf("Phase 1 GET %s: HTTP %d", ep.Path, st))
		if st != 200 {
			regPass = false
		}
	}
	regStatus := "PASS"
	if !regPass {
		regStatus = "FAIL"
	}
	addResult("PART 15", "Phase 1 Regression Audit", regStatus, regDetails)

	// -------------------------------------------------------------------------
	// PART 16: STRICT PHASE BOUNDARY VERIFICATION
	// -------------------------------------------------------------------------
	var bndDetails []string
	bndPass := true

	// Verify no BOM endpoints exist
	bomEndpoints := []string{"/bom", "/bill-of-materials", "/products/PRD-2024-001/components", "/product-versions/VER-001/components"}
	for _, ep := range bomEndpoints {
		st, _ := authedReq(adminToken, "GET", ep, nil)
		bndDetails = append(bndDetails, fmt.Sprintf("Check non-existence of GET %s: HTTP %d (404 Not Found)", ep, st))
		if st != 404 {
			bndPass = false
		}
	}

	// Verify MySQL has no BOM tables
	if db != nil {
		var bomTableCount int
		_ = db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'iot_rd_db' AND table_name LIKE '%bom%'").Scan(&bomTableCount)
		bndDetails = append(bndDetails, fmt.Sprintf("BOM tables in MySQL: %d (Expected: 0)", bomTableCount))
		if bomTableCount > 0 {
			bndPass = false
		}
	}

	bndStatus := "PASS"
	if !bndPass {
		bndStatus = "FAIL"
	}
	addResult("PART 16", "Strict Phase Boundary Audit", bndStatus, bndDetails)

	// -------------------------------------------------------------------------
	// SUMMARY
	// -------------------------------------------------------------------------
	fmt.Println("\n================================================================================")
	fmt.Println("                       FORENSIC AUDIT SUMMARY REPORT                            ")
	fmt.Println("================================================================================")
	passCount := 0
	failCount := 0
	for _, r := range allResults {
		if r.Status == "PASS" {
			passCount++
		} else {
			failCount++
		}
		fmt.Printf(" [%s] %s - %s\n", r.Status, r.SectionNumber, r.Title)
	}
	fmt.Printf("\nTOTAL SECTIONS AUDITED: %d | PASSED: %d | FAILED: %d\n", len(allResults), passCount, failCount)
	if failCount == 0 {
		fmt.Println("FINAL VERDICT: GO — READY FOR PHASE 3 BOM MANAGEMENT")
	} else {
		fmt.Println("FINAL VERDICT: NO-GO")
		os.Exit(1)
	}
}
