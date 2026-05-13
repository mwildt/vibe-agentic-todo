package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"os"
)

// TestLoginFailure tests that login with invalid credentials returns 401
func TestLoginFailure(t *testing.T) {
	// Setup
	setupTest()
	
	// Cleanup
	defer func() {
		os.RemoveAll("./test_data")
	}()

	// Test with invalid credentials
	req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"username": "wronguser", "password": "wrongpass"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	
	// Should return 401 Unauthorized for invalid credentials
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
	
	// Check that the response contains an error message
	expectedError := "Invalid credentials"
	if !strings.Contains(rr.Body.String(), expectedError) {
		t.Errorf("login failure response does not contain expected error message: got %v want %v", rr.Body.String(), expectedError)
	}
	
	// Verify that no session ID is returned by trying to parse the response
	var loginResp LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&loginResp); err == nil {
		// If parsing succeeds, session_id should be empty
		if loginResp.SessionID != "" {
			t.Errorf("login failure response should not contain session_id, but got: %s", loginResp.SessionID)
		}
	}
	// If parsing fails, that's also acceptable for an error response
}

// TestLoginFailureEmptyCredentials tests that login with empty credentials returns 401
func TestLoginFailureEmptyCredentials(t *testing.T) {
	// Setup
	setupTest()
	
	// Cleanup
	defer func() {
		os.RemoveAll("./test_data")
	}()

	// Test with empty credentials
	req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"username": "", "password": ""}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	
	// Should return 400 Bad Request due to validation failure
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
}
