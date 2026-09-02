package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080/api/v1"

type ApiResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
}

type User struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Department string `json:"department"`
	Role       *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"role"`
}

type LoginData struct {
	Token       string   `json:"token"`
	User        User     `json:"user"`
	Permissions []string `json:"permissions"`
}

var (
	totalPassed = 0
	totalFailed = 0
)

func record(testID, description string, passed bool, details string) {
	if passed {
		totalPassed++
		fmt.Printf("  [PASS] [%s] %s\n", testID, description)
	} else {
		totalFailed++
		fmt.Printf(" ![FAIL] [%s] %s: %s\n", testID, description, details)
	}
}

func doRequest(method, path string, token string, body interface{}) (int, []byte, error) {
	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, baseURL+path, bodyReader)
	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	return resp.StatusCode, respBytes, err
}

func main() {
	fmt.Println("================================================================================")
	fmt.Println(" AUTHENTICATION UX & JWT SESSION VERIFICATION SUITE")
	fmt.Println(" Target: IoT Product R&D Control Center")
	fmt.Println("================================================================================")

	// 1. Valid Admin Login
	code, respBytes, err := doRequest("POST", "/auth/login", "", map[string]string{
		"email":    "alex.chen@iotcontrol.io",
		"password": "admin123",
	})
	var adminResp struct {
		Success bool      `json:"success"`
		Data    LoginData `json:"data"`
		Message string    `json:"message"`
	}
	_ = json.Unmarshal(respBytes, &adminResp)
	adminToken := adminResp.Data.Token
	record("AUTH-01", "Valid Admin Login (alex.chen@iotcontrol.io)", err == nil && code == 200 && adminToken != "" && adminResp.Data.User.Name == "Alex Chen", fmt.Sprintf("Code: %d, err: %v", code, err))

	// 2. Valid HW Lead Login
	code, respBytes, err = doRequest("POST", "/auth/login", "", map[string]string{
		"email":    "marcus.vance@iotcontrol.io",
		"password": "admin123",
	})
	var hwResp struct {
		Success bool      `json:"success"`
		Data    LoginData `json:"data"`
	}
	_ = json.Unmarshal(respBytes, &hwResp)
	record("AUTH-02", "Valid HW Lead Login (marcus.vance@iotcontrol.io)", err == nil && code == 200 && hwResp.Data.Token != "" && hwResp.Data.User.Name == "Marcus Vance", fmt.Sprintf("Code: %d", code))

	// 3. Valid QA / R&D Lead Login
	code, respBytes, err = doRequest("POST", "/auth/login", "", map[string]string{
		"email":    "sarah.jenkins@iotcontrol.io",
		"password": "admin123",
	})
	var qaResp struct {
		Success bool      `json:"success"`
		Data    LoginData `json:"data"`
	}
	_ = json.Unmarshal(respBytes, &qaResp)
	record("AUTH-03", "Valid R&D Manager Login (sarah.jenkins@iotcontrol.io)", err == nil && code == 200 && qaResp.Data.Token != "" && (qaResp.Data.User.Name == "Sarah Jenkins" || qaResp.Data.User.Name == "Dr. Sarah Jenkins"), fmt.Sprintf("Code: %d, Name: %s", code, qaResp.Data.User.Name))

	// 4. Valid Viewer Login
	code, respBytes, err = doRequest("POST", "/auth/login", "", map[string]string{
		"email":    "maya.lin@iotcontrol.io",
		"password": "admin123",
	})
	var viewResp struct {
		Success bool      `json:"success"`
		Data    LoginData `json:"data"`
	}
	_ = json.Unmarshal(respBytes, &viewResp)
	viewerToken := viewResp.Data.Token
	record("AUTH-04", "Valid Viewer Login (maya.lin@iotcontrol.io)", err == nil && code == 200 && viewerToken != "" && viewResp.Data.User.Name == "Maya Lin", fmt.Sprintf("Code: %d", code))

	// 5. Invalid Password Error Handling
	code, _, _ = doRequest("POST", "/auth/login", "", map[string]string{
		"email":    "alex.chen@iotcontrol.io",
		"password": "incorrect_password",
	})
	record("AUTH-05", "Invalid Password Returns 400 Bad Request", code == 400, fmt.Sprintf("Status: %d", code))

	// 6. Non-existent User Error Handling
	code, _, _ = doRequest("POST", "/auth/login", "", map[string]string{
		"email":    "nonexistent.user@iotcontrol.io",
		"password": "admin123",
	})
	record("AUTH-06", "Non-existent User Returns 400 Bad Request", code == 400, fmt.Sprintf("Status: %d", code))

	// 7. Get Current Authenticated User (/auth/me) with Admin Token
	code, respBytes, err = doRequest("GET", "/auth/me", adminToken, nil)
	var meResp struct {
		Success bool `json:"success"`
		Data    struct {
			User        User     `json:"user"`
			Permissions []string `json:"permissions"`
		} `json:"data"`
	}
	_ = json.Unmarshal(respBytes, &meResp)
	record("AUTH-07", "Load Current User (GET /auth/me) with JWT", err == nil && code == 200 && meResp.Data.User.Email == "alex.chen@iotcontrol.io", fmt.Sprintf("User: %s, code: %d", meResp.Data.User.Email, code))

	// 8. Unauthenticated Access to /auth/me returns 401
	code, _, _ = doRequest("GET", "/auth/me", "", nil)
	record("AUTH-08", "Unauthenticated Request to /auth/me Returns 401 Unauthorized", code == 401, fmt.Sprintf("Code: %d", code))

	// 9. Unauthenticated Access to Protected API returns 401
	code, _, _ = doRequest("GET", "/products", "", nil)
	record("AUTH-09", "Unauthenticated Request to Protected Route Returns 401 Unauthorized", code == 401, fmt.Sprintf("Code: %d", code))

	// 10. Role Verification & RBAC Permissions Check
	code, _, _ = doRequest("POST", "/eco/ECO-2026-001/implement", viewerToken, nil)
	record("AUTH-10", "Viewer Token Mutation Request Blocked by RBAC (403 Forbidden)", code == 403, fmt.Sprintf("Code: %d", code))

	fmt.Println("================================================================================")
	fmt.Printf(" AUTHENTICATION SUITE RESULTS: %d PASSED, %d FAILED\n", totalPassed, totalFailed)
	fmt.Println("================================================================================")

	if totalFailed > 0 {
		os.Exit(1)
	}
}
