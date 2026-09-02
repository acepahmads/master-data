package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func main() {
	baseURL := "http://localhost:8080/api/v1"

	// 1. Test Login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "alex.chen@iotcontrol.io",
		"password": "admin123",
	})
	resp, err := http.Post(baseURL+"/auth/login", "application/json", bytes.NewBuffer(loginBody))
	if err != nil {
		fmt.Printf("Login failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var loginRes struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&loginRes)

	if !loginRes.Success || loginRes.Data.Token == "" {
		fmt.Println("Cannot proceed without valid auth token.")
		return
	}

	token := loginRes.Data.Token
	client := &http.Client{}

	authedGet := func(path string) (int, string) {
		req, _ := http.NewRequest("GET", baseURL+path, nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r, e := client.Do(req)
		if e != nil {
			return 500, e.Error()
		}
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)
		return r.StatusCode, string(body)
	}

	authedPost := func(path string, payload any) (int, string) {
		b, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", baseURL+path, bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		r, e := client.Do(req)
		if e != nil {
			return 500, e.Error()
		}
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)
		return r.StatusCode, string(body)
	}

	authedPut := func(path string, payload any) (int, string) {
		b, _ := json.Marshal(payload)
		req, _ := http.NewRequest("PUT", baseURL+path, bytes.NewBuffer(b))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		r, e := client.Do(req)
		if e != nil {
			return 500, e.Error()
		}
		defer r.Body.Close()
		body, _ := io.ReadAll(r.Body)
		return r.StatusCode, string(body)
	}

	// Flow A: GET /product-versions
	status, body := authedGet("/product-versions")
	fmt.Printf("[A] GET /product-versions: HTTP %d (Length: %d bytes)\n", status, len(body))

	// Flow B: GET /hardware-revisions
	status, body = authedGet("/hardware-revisions")
	fmt.Printf("[B] GET /hardware-revisions: HTTP %d (Length: %d bytes)\n", status, len(body))

	// Flow C: GET /products/PRD-2024-001/lifecycle
	status, body = authedGet("/products/PRD-2024-001/lifecycle")
	fmt.Printf("[C] GET /products/PRD-2024-001/lifecycle: HTTP %d (Body: %s)\n", status, body)

	// Flow D: GET /product-versions/compare
	status, body = authedGet("/product-versions/compare?base_id=VER-001&target_id=VER-004")
	fmt.Printf("[D] GET /product-versions/compare: HTTP %d\n", status)

	// Flow E: POST /product-versions (Create Version)
	newVer := map[string]any{
		"productId":     "PRD-2024-001",
		"versionNumber": "3.0.0",
		"versionName":   "Edge AI Neural Compute Baseline",
		"status":        "Draft",
		"owner":         "Alex Chen",
		"description":   "Next-generation edge computer with on-board NPU.",
		"releaseNotes":  "1. NPU integration.\n2. Dual GbE.\n3. USB 3.2.",
	}
	status, body = authedPost("/product-versions", newVer)
	fmt.Printf("[E] POST /product-versions: HTTP %d (Body: %s)\n", status, body)

	// Flow F: PUT /product-versions/VER-004 (Update Version)
	updateVer := map[string]any{
		"versionName": "LTE-M / NB-IoT Cellular Edge Hub (Updated)",
		"status":      "Active",
	}
	status, body = authedPut("/product-versions/VER-004", updateVer)
	fmt.Printf("[F] PUT /product-versions/VER-004: HTTP %d\n", status)

	// Flow G: POST /hardware-revisions (Create Revision)
	newRev := map[string]any{
		"productVersionId": "VER-004",
		"code":             "REV-D",
		"name":             "Thermal Dissipation Optimization",
		"status":           "Draft",
		"engineer":         "Marcus Vance",
		"changeSummary":    "Enhanced ground plane thermal vias underneath MPU and power inductors.",
	}
	status, body = authedPost("/hardware-revisions", newRev)
	fmt.Printf("[G] POST /hardware-revisions: HTTP %d (Body: %s)\n", status, body)

	// Flow H: PUT /hardware-revisions/REV-002 (Update Revision)
	updateRev := map[string]any{
		"name":   "RF Power Filtering & RS485 Isolation Improvement (Certified)",
		"status": "Released",
	}
	status, body = authedPut("/hardware-revisions/REV-002", updateRev)
	fmt.Printf("[H] PUT /hardware-revisions/REV-002: HTTP %d\n", status)
}
