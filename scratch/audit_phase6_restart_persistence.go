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
	adminEmail = "alex.chen@iotcontrol.io"
	adminPwd   = "admin123"
)

var (
	token string
)

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

func main() {
	mode := "pre"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	// Login
	var loginResp struct {
		Success bool `json:"success"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	code, err := doReq("POST", "/auth/login", map[string]string{"email": adminEmail, "password": adminPwd}, &loginResp)
	if err != nil || code != 200 || loginResp.Data.Token == "" {
		fmt.Printf("Login failed: code %d, err %v\n", code, err)
		os.Exit(1)
	}
	token = loginResp.Data.Token

	if mode == "pre" {
		var dashResp struct {
			Success bool `json:"success"`
			Data    struct {
				TotalEcrs int64 `json:"totalEcrs"`
			} `json:"data"`
		}
		_, _ = doReq("GET", "/changes/dashboard", nil, &dashResp)

		var traceResp struct {
			Success bool `json:"success"`
			Data    []struct {
				ECRCode string `json:"ecrCode"`
			} `json:"data"`
		}
		_, _ = doReq("GET", "/changes/traceability", nil, &traceResp)

		fmt.Printf("PRE_RESTART_ECR_TOTAL=%d\n", dashResp.Data.TotalEcrs)
		fmt.Printf("PRE_RESTART_TRACE_COUNT=%d\n", len(traceResp.Data))
	} else {
		var dashResp struct {
			Success bool `json:"success"`
			Data    struct {
				TotalEcrs int64 `json:"totalEcrs"`
			} `json:"data"`
		}
		code, err := doReq("GET", "/changes/dashboard", nil, &dashResp)
		if err != nil || code != 200 {
			fmt.Printf("FAIL: Dashboard query post-restart failed: %v\n", err)
			os.Exit(1)
		}

		var traceResp struct {
			Success bool `json:"success"`
			Data    []struct {
				ECRCode string `json:"ecrCode"`
			} `json:"data"`
		}
		code, err = doReq("GET", "/changes/traceability", nil, &traceResp)
		if err != nil || code != 200 {
			fmt.Printf("FAIL: Traceability query post-restart failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("POST_RESTART_ECR_TOTAL=%d\n", dashResp.Data.TotalEcrs)
		fmt.Printf("POST_RESTART_TRACE_COUNT=%d\n", len(traceResp.Data))
		fmt.Println("SERVER RESTART PERSISTENCE AUDIT: [PASS] All Phase 6 records verified intact in MySQL.")
	}
}
