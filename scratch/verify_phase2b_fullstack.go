package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type VerificationReport struct {
	TotalTests  int
	PassedTests int
	FailedTests int
	Results     []TestResult
}

type TestResult struct {
	ID       string
	Category string
	Name     string
	Status   string // PASS, FAIL
	Details  string
}

func main() {
	baseURL := "http://localhost:8080/api/v1"
	client := &http.Client{}
	report := VerificationReport{}

	record := func(id, cat, name string, passed bool, details string) {
		report.TotalTests++
		status := "PASS"
		if passed {
			report.PassedTests++
		} else {
			report.FailedTests++
			status = "FAIL"
		}
		report.Results = append(report.Results, TestResult{
			ID:       id,
			Category: cat,
			Name:     name,
			Status:   status,
			Details:  details,
		})
		fmt.Printf("[%s] %s: %s - %s\n", status, id, name, details)
	}

	// 1. Authenticate as Super Admin (Alex Chen)
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "alex.chen@iotcontrol.io",
		"password": "admin123",
	})
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		fmt.Printf("Fatal: Login connection failed: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var loginRes struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&loginRes)
	token := loginRes.Data.Token

	record("TEST-AUTH-01", "Authentication", "Super Admin JWT Login", loginRes.Success && token != "", "Token acquired successfully")

	authedReq := func(method, path string, payload any) (int, []byte) {
		var req *http.Request
		if payload != nil {
			b, _ := json.Marshal(payload)
			req, _ = http.NewRequest(method, baseURL+path, bytes.NewBuffer(b))
			req.Header.Set("Content-Type", "application/json")
		} else {
			req, _ = http.NewRequest(method, baseURL+path, nil)
		}
		req.Header.Set("Authorization", "Bearer "+token)
		r, e := client.Do(req)
		if e != nil {
			return 500, []byte(e.Error())
		}
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)
		return r.StatusCode, body
	}

	// 2. Test GET /product-versions
	status, body := authedReq("GET", "/product-versions", nil)
	var verListRes struct {
		Success bool `json:"success"`
		Data    []struct {
			ID            string `json:"id"`
			VersionNumber string `json:"versionNumber"`
			ProductCode   string `json:"productCode"`
		} `json:"data"`
		Meta struct {
			Total int `json:"total"`
		} `json:"meta"`
	}
	_ = json.Unmarshal(body, &verListRes)
	record("TEST-VER-01", "Version API", "List Product Versions", status == 200 && verListRes.Success && len(verListRes.Data) > 0, fmt.Sprintf("Retrieved %d versions", len(verListRes.Data)))

	// 3. Test GET /hardware-revisions
	status, body = authedReq("GET", "/hardware-revisions", nil)
	var revListRes struct {
		Success bool `json:"success"`
		Data    []struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &revListRes)
	record("TEST-REV-01", "Revision API", "List Hardware Revisions", status == 200 && revListRes.Success && len(revListRes.Data) > 0, fmt.Sprintf("Retrieved %d hardware revisions", len(revListRes.Data)))

	// 4. Test GET /products/PRD-2024-001/lifecycle
	status, body = authedReq("GET", "/products/PRD-2024-001/lifecycle", nil)
	var lcRes struct {
		Success bool `json:"success"`
		Data    struct {
			ProductCode string `json:"productCode"`
			Versions    []struct {
				ID            string `json:"id"`
				VersionNumber string `json:"versionNumber"`
				Revisions     []struct {
					ID   string `json:"id"`
					Code string `json:"code"`
				} `json:"revisions"`
			} `json:"versions"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &lcRes)
	hasNestedRevs := false
	for _, v := range lcRes.Data.Versions {
		if len(v.Revisions) > 0 {
			hasNestedRevs = true
			break
		}
	}
	record("TEST-LC-01", "Lifecycle API", "Product Lifecycle 3-Tier Hierarchy", status == 200 && lcRes.Data.ProductCode == "PRD-2024-001" && len(lcRes.Data.Versions) > 0 && hasNestedRevs, fmt.Sprintf("Product %s has %d versions with nested revisions", lcRes.Data.ProductCode, len(lcRes.Data.Versions)))

	// 5. Test Version Comparison
	status, body = authedReq("GET", "/product-versions/compare?base_id=VER-001&target_id=VER-004", nil)
	var cmpRes struct {
		Success bool `json:"success"`
		Data    struct {
			BaseVersion   any `json:"baseVersion"`
			TargetVersion any `json:"targetVersion"`
			DiffParameters map[string]bool `json:"diffParameters"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &cmpRes)
	record("TEST-CMP-01", "Comparison API", "Compare Version Parameters", status == 200 && cmpRes.Success && len(cmpRes.Data.DiffParameters) > 0, fmt.Sprintf("%d parameters compared", len(cmpRes.Data.DiffParameters)))

	// 6. Test Create Version + Unique Constraint
	uniqueVerNum := "9.9.0"
	newVerPayload := map[string]any{
		"productId":     "PRD-2024-001",
		"versionNumber": uniqueVerNum,
		"versionName":   "Automated Full-Stack Verification Baseline",
		"status":        "Draft",
		"owner":         "Alex Chen",
		"description":   "Automated verification test baseline.",
		"releaseNotes":  "Automated verification test notes.",
	}
	status, body = authedReq("POST", "/product-versions", newVerPayload)
	var createdVerRes struct {
		Success bool `json:"success"`
		Data    struct {
			ID            string `json:"id"`
			VersionNumber string `json:"versionNumber"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &createdVerRes)
	record("TEST-VER-02", "Version API", "Create Product Version (201 Created)", status == 201 && createdVerRes.Success && createdVerRes.Data.ID != "", fmt.Sprintf("Created Version ID: %s (v%s)", createdVerRes.Data.ID, createdVerRes.Data.VersionNumber))

	// Test Duplicate Version Number rejection
	statusDup, _ := authedReq("POST", "/product-versions", newVerPayload)
	record("TEST-VER-03", "Version API", "Reject Duplicate Version Number (400 Bad Request)", statusDup == 400, "Duplicate version rejected by constraint")

	// 7. Test Create Hardware Revision for the new Version
	newRevPayload := map[string]any{
		"productVersionId": createdVerRes.Data.ID,
		"code":             "REV-A1",
		"name":             "High-Reliability Power Regulation",
		"status":           "In Review",
		"engineer":         "Alex Chen",
		"changeSummary":    "Switched inductor to shielded automotive grade ferrite core.",
	}
	status, body = authedReq("POST", "/hardware-revisions", newRevPayload)
	var createdRevRes struct {
		Success bool `json:"success"`
		Data    struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &createdRevRes)
	record("TEST-REV-02", "Revision API", "Create Hardware Revision (201 Created)", status == 201 && createdRevRes.Success && createdRevRes.Data.ID != "", fmt.Sprintf("Created Revision ID: %s (%s)", createdRevRes.Data.ID, createdRevRes.Data.Code))

	// Test Duplicate Revision Code rejection
	statusRevDup, _ := authedReq("POST", "/hardware-revisions", newRevPayload)
	record("TEST-REV-03", "Revision API", "Reject Duplicate Revision Code (400 Bad Request)", statusRevDup == 400, "Duplicate revision code rejected by constraint")

	// 8. Test Update Product Version
	updateVerPayload := map[string]any{
		"versionName": "Automated Full-Stack Verification Baseline (Updated)",
		"status":      "Active",
	}
	status, _ = authedReq("PUT", "/product-versions/"+createdVerRes.Data.ID, updateVerPayload)
	record("TEST-VER-04", "Version API", "Update Product Version (200 OK)", status == 200, fmt.Sprintf("Updated Version %s", createdVerRes.Data.ID))

	// 9. Test Update Hardware Revision
	updateRevPayload := map[string]any{
		"name":   "High-Reliability Power Regulation (Approved)",
		"status": "Approved",
	}
	status, _ = authedReq("PUT", "/hardware-revisions/"+createdRevRes.Data.ID, updateRevPayload)
	record("TEST-REV-04", "Revision API", "Update Hardware Revision (200 OK)", status == 200, fmt.Sprintf("Updated Revision %s", createdRevRes.Data.ID))

	// 10. Test Audit Log Persistence
	status, body = authedReq("GET", "/activity?limit=20", nil)
	var actRes struct {
		Success bool `json:"success"`
		Data    []struct {
			Action     string `json:"action"`
			EntityType string `json:"entityType"`
		} `json:"data"`
	}
	_ = json.Unmarshal(body, &actRes)
	hasVerLog := false
	hasRevLog := false
	for _, a := range actRes.Data {
		if strings.Contains(a.Action, "VERSION") || a.EntityType == "ProductVersion" {
			hasVerLog = true
		}
		if strings.Contains(a.Action, "REVISION") || a.EntityType == "HardwareRevision" {
			hasRevLog = true
		}
	}
	record("TEST-ACT-01", "Activity Audit", "Verify Product Version Audit Logging", hasVerLog, "Version activity recorded in activity_logs")
	record("TEST-ACT-02", "Activity Audit", "Verify Hardware Revision Audit Logging", hasRevLog, "Revision activity recorded in activity_logs")

	// 11. Test Delete Hardware Revision & Version
	status, _ = authedReq("DELETE", "/hardware-revisions/"+createdRevRes.Data.ID, nil)
	record("TEST-REV-05", "Revision API", "Delete Hardware Revision (200 OK)", status == 200, fmt.Sprintf("Deleted Revision %s", createdRevRes.Data.ID))

	status, _ = authedReq("DELETE", "/product-versions/"+createdVerRes.Data.ID, nil)
	record("TEST-VER-05", "Version API", "Delete Product Version (200 OK)", status == 200, fmt.Sprintf("Deleted Version %s", createdVerRes.Data.ID))

	// 12. Test Clean Boundary: Verify NO BOM Endpoints exist
	statusBOM1, _ := authedReq("GET", "/bom", nil)
	statusBOM2, _ := authedReq("GET", "/bill-of-materials", nil)
	statusBOM3, _ := authedReq("GET", "/products/PRD-2024-001/components", nil)
	bomExcluded := (statusBOM1 == 404) && (statusBOM2 == 404) && (statusBOM3 == 404)
	record("TEST-BND-01", "Boundary Integrity", "Zero BOM Contamination (Strict Boundary)", bomExcluded, "Confirmed BOM and Component Relationship endpoints are absent")

	fmt.Println("==================================================")
	fmt.Printf("FULL-STACK VERIFICATION RESULT: %d / %d TESTS PASSED (Failed: %d)\n", report.PassedTests, report.TotalTests, report.FailedTests)
	fmt.Println("==================================================")

	if report.FailedTests > 0 {
		os.Exit(1)
	}
}
