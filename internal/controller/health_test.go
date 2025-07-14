package controller

import (
	"net/http"
	"testing"
)

// 测试健康连接
func TestHealth(t *testing.T) {
	// Making a GET request to the health check endpoint
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Checking if the status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", resp.StatusCode)
	}
}
