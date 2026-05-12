package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"os"
)

// TestLoginSuccess tests that a user can login with valid credentials and receives a session ID
func TestLoginSuccess(t *testing.T) {
	// Setup
	setupTest()
	
	// Cleanup
	defer func() {
		os.RemoveAll("./test_data")
	}()

	// Create login request with valid credentials
	req, err := http.NewRequest("POST", "/login", strings.NewReader(`{"username": "testuser", "password": "testpass"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	
	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check that the response contains a session ID
	expected := `"session_id":"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("login response does not contain session ID: got %v", rr.Body.String())
	}
	
	// Check that the session ID is 64 characters (32 bytes in hex)
	sessionStart := strings.Index(rr.Body.String(), expected) + len(expected)
	sessionEnd := strings.Index(rr.Body.String()[sessionStart:], `"`)
	sessionID := rr.Body.String()[sessionStart : sessionStart+sessionEnd]
	
	t.Logf("Session ID length: %d", len(sessionID))
	
	if len(sessionID) != 64 {
		t.Errorf("session ID has wrong length: got %d characters, want 64", len(sessionID))
	}
}
