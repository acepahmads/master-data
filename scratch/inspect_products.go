package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	// Login
	resp, _ := http.Post("http://localhost:8080/api/v1/auth/login", "application/json", nil)
	// We will query directly after server restarts
	fmt.Println("Inspecting products")
}
