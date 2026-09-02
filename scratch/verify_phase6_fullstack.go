package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	baseURL  = "http://localhost:8080/api/v1"
	email    = "alex.chen@iotcontrol.io"
	password = "admin123"
)

var (
	token   string
	passCt  int
	failCt  int
)

func assert(condition bool, testName string, details string) {
	if condition {
		passCt++
		fmt.Printf(" [PASS] %s\n", testName)
	} else {
		failCt++
		fmt.Printf("![FAIL] %s: %s\n", testName, details)
	}
}

func doReq(method, path string, body interface{}, target interface{}) (int, error) {
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, bodyReader)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, err
	}

	if target != nil && len(respBytes) > 0 {
		_ = json.Unmarshal(respBytes, target)
	}

	return resp.StatusCode, nil
}

type APIEnvelope struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
	Errors  interface{}     `json:"errors"`
}

func main() {
	fmt.Println("================================================================================")
	fmt.Println(" PHASE 6: ENGINEERING CHANGE MANAGEMENT (ECR / ECO) FULL-STACK VERIFICATION")
	fmt.Println("================================================================================")

	// Step 1: Login & RBAC Token Acquisition
	var loginResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token       string   `json:"token"`
			Permissions []string `json:"permissions"`
			User        struct {
				Username string `json:"username"`
				Role     string `json:"role"`
			} `json:"user"`
		} `json:"data"`
	}

	loginPayload := map[string]string{
		"email":    email,
		"password": password,
	}

	code, err := doReq("POST", "/auth/login", loginPayload, &loginResp)
	assert(err == nil && code == 200 && loginResp.Data.Token != "", "1. Authentication & JWT Acquisition", fmt.Sprintf("Code: %d", code))
	token = loginResp.Data.Token

	// Verify Phase 6 RBAC permissions
	var permMap = make(map[string]bool)
	for _, p := range loginResp.Data.Permissions {
		permMap[p] = true
	}

	expectedPerms := []string{
		"changes.view", "ecr.create", "ecr.edit", "ecr.review", "ecr.approve",
		"eco.create", "eco.approve", "eco.implement", "eco.verify",
	}
	allPermsPresent := true
	for _, ep := range expectedPerms {
		if !permMap[ep] {
			allPermsPresent = false
			fmt.Printf("   Missing expected permission: %s\n", ep)
		}
	}
	assert(allPermsPresent, "2. Phase 6 RBAC Permissions Registered in Master Seeder", "All 9 permissions verified in admin role")

	// Step 2: Dashboard Metrics Endpoint
	var dashResp struct {
		Success bool `json:"success"`
		Data    struct {
			TotalEcrs                  int64 `json:"totalEcrs"`
			DraftEcrs                  int64 `json:"draftEcrs"`
			UnderReviewEcrs            int64 `json:"underReviewEcrs"`
			ApprovedEcrs               int64 `json:"approvedEcrs"`
			TotalEcos                  int64 `json:"totalEcos"`
			DraftEcos                  int64 `json:"draftEcos"`
			ApprovedEcos               int64 `json:"approvedEcos"`
			ImplementedEcos            int64 `json:"implementedEcos"`
			VerifiedEcos               int64 `json:"verifiedEcos"`
			ClosedEcos                 int64 `json:"closedEcos"`
			CriticalHighEcrs           int64 `json:"criticalHighEcrs"`
			AvgImplementationCycleDays int   `json:"avgImplementationCycleDays"`
		} `json:"data"`
	}

	code, err = doReq("GET", "/changes/dashboard", nil, &dashResp)
	assert(err == nil && code == 200 && dashResp.Success && dashResp.Data.TotalEcrs >= 1, "3. Changes Dashboard KPI Aggregations (GET /changes/dashboard)", fmt.Sprintf("Total ECRs: %d, Total ECOs: %d", dashResp.Data.TotalEcrs, dashResp.Data.TotalEcos))

	// Step 3: ECR Registry & Seeded Findings Ingestion Verification
	var ecrListResp struct {
		Success bool `json:"success"`
		Data    struct {
			Items []struct {
				ID                       string `json:"id"`
				Code                     string `json:"code"`
				Title                    string `json:"title"`
				Status                   string `json:"status"`
				Priority                 string `json:"priority"`
				Origin                   string `json:"origin"`
				SourceHardwareRevisionID string `json:"sourceHardwareRevisionId"`
			} `json:"items"`
			Total int64 `json:"total"`
		} `json:"data"`
	}

	code, err = doReq("GET", "/ecr", nil, &ecrListResp)
	assert(err == nil && code == 200 && ecrListResp.Data.Total >= 2, "4. ECR Registry Query (GET /ecr)", fmt.Sprintf("Found %d ECR records", ecrListResp.Data.Total))

	// Step 4: Create a New ECR (Dynamic Workflow Test)
	var createEcrResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID                       string `json:"id"`
			Code                     string `json:"code"`
			Title                    string `json:"title"`
			Status                   string `json:"status"`
			Priority                 string `json:"priority"`
			SourceHardwareRevisionID string `json:"sourceHardwareRevisionId"`
		} `json:"data"`
	}

	newEcrPayload := map[string]interface{}{
		"title":                    "Antenna Matching Network Impedance Optimization",
		"origin":                   "TEST_FINDING",
		"sourceFindingId":          "FND-2024-001",
		"productId":                "PRD-2024-001",
		"sourceHardwareRevisionId": "REV-001",
		"changeType":               "HARDWARE_DESIGN",
		"priority":                 "CRITICAL",
		"description":              "Optimize Pi-matching network components (L1, C12, C14) to eliminate RF return loss degradation at 2.4GHz ISM band.",
		"businessReason":           "Enhance transmission range and pass FCC certification.",
		"technicalJustification":   "Smith chart analysis indicates 12dB return loss improvement.",
		"assignedLead":             "Senior RF Engineer",
	}

	code, err = doReq("POST", "/ecr", newEcrPayload, &createEcrResp)
	assert(err == nil && code == 201 && createEcrResp.Data.ID != "", "5. Create New ECR Draft (POST /ecr)", fmt.Sprintf("ECR Code: %s", createEcrResp.Data.Code))
	testEcrID := createEcrResp.Data.ID
	testEcrCode := createEcrResp.Data.Code

	// Step 5: Add Change Items to ECR
	var itemResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID                string `json:"id"`
			ChangeType        string `json:"changeType"`
			AffectedReference string `json:"affectedReference"`
		} `json:"data"`
	}

	changeItemPayload := map[string]interface{}{
		"changeType":        "COMPONENT_VALUE",
		"affectedReference": "C12",
		"title":             "Tune Matching Capacitor C12 from 2.2pF to 1.8pF",
		"description":       "Impedance tuning for 50 Ohm RF trace match",
		"currentValue":      "2.2pF 0402",
		"proposedValue":     "1.8pF 0402 high-Q RF cap",
	}

	code, err = doReq("POST", fmt.Sprintf("/ecr/%s/items", testEcrID), changeItemPayload, &itemResp)
	assert(err == nil && code == 201 && itemResp.Data.ID != "", "6. Add Line-Item Change to ECR (POST /ecr/:id/items)", fmt.Sprintf("Added Item for %s", itemResp.Data.AffectedReference))

	// Step 6: Save 5-Domain Impact Analysis
	var impactResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID          string  `json:"id"`
			Domain      string  `json:"domain"`
			ImpactLevel string  `json:"impactLevel"`
			CostDelta   float64 `json:"costDelta"`
		} `json:"data"`
	}

	costImpactPayload := map[string]interface{}{
		"domain":              "COST",
		"impactLevel":         "LOW",
		"description":         "Minimal unit cost delta of $0.02 for RF grade dielectric ceramic capacitor.",
		"costDelta":           0.02,
		"scheduleImpactWeeks": 1,
		"scrapDisposition":    "REWORK",
	}

	code, err = doReq("POST", fmt.Sprintf("/ecr/%s/impact", testEcrID), costImpactPayload, &impactResp)
	assert(err == nil && (code == 200 || code == 201) && impactResp.Data.Domain == "COST", "7. Record 5-Domain Cross-Functional Impact (POST /ecr/:id/impact)", fmt.Sprintf("Domain: %s, Level: %s", impactResp.Data.Domain, impactResp.Data.ImpactLevel))

	// Step 7: Submit ECR for Review
	var submitResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}

	code, err = doReq("POST", fmt.Sprintf("/ecr/%s/submit", testEcrID), nil, &submitResp)
	assert(err == nil && code == 200 && submitResp.Data.Status == "UNDER_REVIEW", "8. Submit ECR for Technical Review (POST /ecr/:id/submit)", fmt.Sprintf("Status: %s", submitResp.Data.Status))

	// Step 8: Multi-Disciplinary Reviews
	var reviewResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID         string `json:"id"`
			Discipline string `json:"discipline"`
			Decision   string `json:"decision"`
		} `json:"data"`
	}

	reviewPayload := map[string]interface{}{
		"discipline": "HARDWARE",
		"decision":   "APPROVE",
		"comments":   "RF simulation validated; schematic match verified.",
	}

	code, err = doReq("POST", fmt.Sprintf("/ecr/%s/reviews", testEcrID), reviewPayload, &reviewResp)
	assert(err == nil && code == 201 && reviewResp.Data.Decision == "APPROVE", "9. Submit Discipline Technical Review (POST /ecr/:id/reviews)", fmt.Sprintf("Discipline: %s", reviewResp.Data.Discipline))

	// Step 9: CCB Multi-Tier Approval Workflow
	var apprResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID       string `json:"id"`
			Stage    string `json:"stage"`
			Decision string `json:"decision"`
		} `json:"data"`
	}

	ccbStages := []string{"TIER_1_TECH_LEAD", "TIER_2_ENG_MANAGER", "TIER_3_QUALITY_DIRECTOR"}
	for idx, stage := range ccbStages {
		apprPayload := map[string]interface{}{
			"stage":         stage,
			"decision":      "APPROVED",
			"justification": fmt.Sprintf("CCB Stage %d sign-off authorized.", idx+1),
		}
		code, err = doReq("POST", fmt.Sprintf("/ecr/%s/approvals", testEcrID), apprPayload, &apprResp)
		assert(err == nil && code == 201 && apprResp.Data.Decision == "APPROVED", fmt.Sprintf("10.%d. CCB Approval Sign-off: %s (POST /ecr/:id/approvals)", idx+1, stage), fmt.Sprintf("Stage: %s", stage))
	}

	// Verify ECR transitioned to APPROVED
	var verifyEcrResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID     string `json:"id"`
			Status string `json:"status"`
			ECO    *struct {
				ID string `json:"id"`
			} `json:"eco"`
		} `json:"data"`
	}

	code, err = doReq("GET", fmt.Sprintf("/ecr/%s", testEcrID), nil, &verifyEcrResp)
	assert(err == nil && code == 200 && verifyEcrResp.Data.Status == "APPROVED", "11. ECR Full CCB Approval Status Transition", fmt.Sprintf("ECR Status: %s", verifyEcrResp.Data.Status))

	// CRITICAL GUARDRAIL #1: ECR must not automatically create an ECO
	assert(verifyEcrResp.Data.ECO == nil, "12. CRITICAL GUARDRAIL #1: ECR Approval does NOT auto-create ECO", "ECO is nil; awaiting explicit conversion action")

	// Step 10: Explicit Conversion to ECO
	var createEcoResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID                       string  `json:"id"`
			Code                     string  `json:"code"`
			Title                    string  `json:"title"`
			Status                   string  `json:"status"`
			SourceHardwareRevisionID string  `json:"sourceHardwareRevisionId"`
			TargetHardwareRevisionID *string `json:"targetHardwareRevisionId"`
		} `json:"data"`
	}

	ecoPayload := map[string]interface{}{
		"ecrId":               testEcrID,
		"title":               fmt.Sprintf("Implement %s: Antenna Impedance Tuning", testEcrCode),
		"implementationOwner": "Lead RF HW Engineer",
		"scrapDisposition":    "REWORK",
	}

	code, err = doReq("POST", "/eco", ecoPayload, &createEcoResp)
	testEcoID := createEcoResp.Data.ID
	testEcoCode := createEcoResp.Data.Code
	assert(err == nil && code == 201 && testEcoID != "", "13. Explicit ECO Creation from Approved ECR (POST /eco)", fmt.Sprintf("ECO Code: %s, Status: %s", testEcoCode, createEcoResp.Data.Status))

	// CRITICAL GUARDRAIL #2: ECO creation does NOT modify hardware revision
	assert(createEcoResp.Data.Status == "DRAFT" && createEcoResp.Data.TargetHardwareRevisionID == nil, "14. CRITICAL GUARDRAIL #2 (Part A): ECO Initialized in DRAFT, Target Revision NOT yet created", fmt.Sprintf("ECO %s has no target revision; historical REV-001 untouched", testEcoCode))

	// Step 11: Approve ECO for Execution
	var approveEcoResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}

	code, err = doReq("POST", fmt.Sprintf("/eco/%s/approve", testEcoID), nil, &approveEcoResp)
	assert(err == nil && code == 200 && approveEcoResp.Data.Status == "APPROVED", "15. Approve ECO for Execution Baseline Generation (POST /eco/:id/approve)", fmt.Sprintf("ECO Status: %s", approveEcoResp.Data.Status))

	// Step 12: BOM Change Preview Engine
	var previewResp struct {
		Success bool `json:"success"`
		Data    struct {
			SourceRevisionCode string  `json:"sourceRevisionCode"`
			TargetRevisionCode string  `json:"targetRevisionCode"`
			SourceBOMCost      float64 `json:"sourceBomCost"`
			TargetBOMCost      float64 `json:"targetBomCost"`
			NetCostDelta       float64 `json:"netCostDelta"`
			TotalModified      int     `json:"totalModified"`
			TotalUnchanged     int     `json:"totalUnchanged"`
			Items              []struct {
				ReferenceDesignator string  `json:"referenceDesignator"`
				ChangeStatus        string  `json:"changeStatus"`
				SourceTotalCost     float64 `json:"sourceTotalCost"`
				TargetTotalCost     float64 `json:"targetTotalCost"`
			} `json:"items"`
		} `json:"data"`
	}

	code, err = doReq("GET", fmt.Sprintf("/eco/%s/change-preview", testEcoID), nil, &previewResp)
	assert(err == nil && code == 200 && previewResp.Success && len(previewResp.Data.Items) > 0, "16. Automated BOM Delta Preview Diff Engine (GET /eco/:id/change-preview)", fmt.Sprintf("Source Cost: $%.4f, Target Cost: $%.4f, Delta: $%.4f", previewResp.Data.SourceBOMCost, previewResp.Data.TargetBOMCost, previewResp.Data.NetCostDelta))

	// Step 13: ATOMIC ECO IMPLEMENTATION ENGINE (Spawns Target Hardware Revision & Clones BOM)
	var implementEcoResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID                        string  `json:"id"`
			Code                      string  `json:"code"`
			Status                    string  `json:"status"`
			RequiresRegressionTesting bool    `json:"requiresRegressionTesting"`
			TargetHardwareRevisionID  *string `json:"targetHardwareRevisionId"`
			TargetHardwareRevision    *struct {
				ID            string `json:"id"`
				Code          string `json:"code"`
				Name          string `json:"name"`
				Status        string `json:"status"`
				ChangeSummary string `json:"changeSummary"`
			} `json:"targetHardwareRevision"`
			TargetBOM *struct {
				ID      string  `json:"id"`
				BOMCode string  `json:"bomCode"`
				TotalCost float64 `json:"totalCost"`
			} `json:"targetBOM"`
		} `json:"data"`
	}

	code, err = doReq("POST", fmt.Sprintf("/eco/%s/implement", testEcoID), nil, &implementEcoResp)
	assert(err == nil && code == 200 && implementEcoResp.Data.Status == "IMPLEMENTED", "17. ATOMIC ECO IMPLEMENTATION ENGINE (POST /eco/:id/implement)", fmt.Sprintf("Status: %s", implementEcoResp.Data.Status))

	// Verify Target Revision created in Draft status with distinct code
	assert(implementEcoResp.Data.TargetHardwareRevision != nil && implementEcoResp.Data.TargetHardwareRevision.Code != "" && implementEcoResp.Data.TargetHardwareRevision.Status == "Draft", "18. Baseline Spawning: New Target Hardware Revision Generated in Draft Status", fmt.Sprintf("Target Revision Code: %s, Status: %s", implementEcoResp.Data.TargetHardwareRevision.Code, implementEcoResp.Data.TargetHardwareRevision.Status))

	// Verify Target BOM cloned and created
	assert(implementEcoResp.Data.TargetBOM != nil && implementEcoResp.Data.TargetBOM.BOMCode != "", "19. BOM Cloning & Delta Application: Target BOM Generated with Line Deltas", fmt.Sprintf("Target BOM Code: %s, Cost: $%.4f", implementEcoResp.Data.TargetBOM.BOMCode, implementEcoResp.Data.TargetBOM.TotalCost))

	// Verify CRITICAL GUARDRAIL #3: Regression Testing Mandate Set
	assert(implementEcoResp.Data.RequiresRegressionTesting == true, "20. CRITICAL GUARDRAIL #3: Regression Testing Mandate Enforced on Implemented ECO", "requires_regression_testing = true")

	// Verify Historical Baseline Immutability (Source REV-001 unchanged)
	var sourceRevCheck struct {
		Success bool `json:"success"`
		Data    struct {
			ID   string `json:"id"`
			Code string `json:"code"`
		} `json:"data"`
	}
	code, err = doReq("GET", "/hardware-revisions/REV-001", nil, &sourceRevCheck)
	assert(err == nil && code == 200 && sourceRevCheck.Data.ID == "REV-001" && sourceRevCheck.Data.Code == "REV-A", "21. HISTORICAL BASELINE IMMUTABILITY: Source REV-001 (REV-A) remains intact and unmodified", fmt.Sprintf("Historical baseline %s (%s) preserved untouched", sourceRevCheck.Data.ID, sourceRevCheck.Data.Code))

	// Step 14: Verify ECO Regression Execution
	var verifyEcoActionResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}

	verifyPayload := map[string]interface{}{
		"notes": "Regression RF sweep test passed: return loss improved to 18.4dB at 2.45GHz. Zero thermal abnormalities.",
	}

	code, err = doReq("POST", fmt.Sprintf("/eco/%s/verify", testEcoID), verifyPayload, &verifyEcoActionResp)
	assert(err == nil && code == 200 && verifyEcoActionResp.Data.Status == "VERIFIED", "22. ECO Post-Implementation Regression Verification (POST /eco/:id/verify)", fmt.Sprintf("Status: %s", verifyEcoActionResp.Data.Status))

	// Step 15: Close ECO
	var closeEcoResp struct {
		Success bool `json:"success"`
		Data    struct {
			ID     string `json:"id"`
			Status string `json:"status"`
		} `json:"data"`
	}

	code, err = doReq("POST", fmt.Sprintf("/eco/%s/close", testEcoID), nil, &closeEcoResp)
	assert(err == nil && code == 200 && closeEcoResp.Data.Status == "CLOSED", "23. Formally Close ECO (POST /eco/:id/close)", fmt.Sprintf("Status: %s", closeEcoResp.Data.Status))

	// Step 16: End-to-End Closed-Loop Traceability Matrix (DAG)
	var traceResp struct {
		Success bool `json:"success"`
		Data    []struct {
			ECRID                string  `json:"ecrId"`
			ECRCode              string  `json:"ecrCode"`
			ECRTitle             string  `json:"ecrTitle"`
			ECRStatus            string  `json:"ecrStatus"`
			SourceFindingID      *string `json:"sourceFindingId"`
			SourceFindingCode    *string `json:"sourceFindingCode"`
			SourceRevisionCode   string  `json:"sourceRevisionCode"`
			ECOID                *string `json:"ecoId"`
			ECOCode              *string `json:"ecoCode"`
			ECOStatus            *string `json:"ecoStatus"`
			TargetRevisionCode   *string `json:"targetRevisionCode"`
			TargetBOMCode        *string `json:"targetBomCode"`
			ClosedLoopVerified   bool    `json:"closedLoopVerified"`
			CCBApprovalCount     int     `json:"ccbApprovalCount"`
		} `json:"data"`
	}

	code, err = doReq("GET", "/changes/traceability", nil, &traceResp)
	assert(err == nil && code == 200 && traceResp.Success && len(traceResp.Data) >= 2, "24. Closed-Loop Traceability Directed Acyclic Graph (GET /changes/traceability)", fmt.Sprintf("Verified %d complete traceability lineage chains", len(traceResp.Data)))

	// Check that our newly completed change is reflected as closed loop verified
	foundCompletedChain := false
	for _, c := range traceResp.Data {
		if c.ECRCode == testEcrCode && c.ClosedLoopVerified {
			foundCompletedChain = true
			fmt.Printf("   Trace Chain Verified: Finding[%v] -> %s -> CCB[%d] -> %s[%s] -> TargetRev[%s] -> Verified[%v]\n",
				c.SourceFindingCode != nil, c.ECRCode, c.CCBApprovalCount, *c.ECOCode, *c.ECOStatus, *c.TargetRevisionCode, c.ClosedLoopVerified)
		}
	}
	assert(foundCompletedChain, "25. Closed-Loop Verification Integrity in Traceability Matrix", "Defect -> ECR -> CCB -> ECO -> Target Rev -> Regressions -> Closed Loop fully validated")

	// Final Summary Report
	fmt.Println("================================================================================")
	fmt.Printf(" PHASE 6 VERIFICATION COMPLETE: %d PASSED, %d FAILED\n", passCt, failCt)
	fmt.Println("================================================================================")

	if failCt > 0 {
		os.Exit(1)
	}
}
