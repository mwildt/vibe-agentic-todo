package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"os"
)

// TestLoginSuccess tests that a user can login with valid credentials and receives a session cookie
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

	// Parse the response using proper JSON unmarshaling
	var loginResp LoginResponse
	if err := json.NewDecoder(rr.Body).Decode(&loginResp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	// Verify login was successful
	if loginResp.Status != "success" {
		t.Errorf("login response status is %s, want success", loginResp.Status)
	}

	// Verify session cookie is present
	cookies := rr.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No cookies returned in login response")
	}

	sessionCookie := cookies[0]
	if sessionCookie.Name != "session_id" {
		t.Errorf("expected session_id cookie, got %s", sessionCookie.Name)
	}

	// Verify session ID length (64 characters = 32 bytes in hex)
	if len(sessionCookie.Value) != 64 {
		t.Errorf("session ID has wrong length: got %d characters, want 64", len(sessionCookie.Value))
	}

	// Verify cookie properties
	if !sessionCookie.HttpOnly {
		t.Errorf("session cookie should be HttpOnly")
	}

	if sessionCookie.Path != "/" {
		t.Errorf("session cookie path should be /, got %s", sessionCookie.Path)
	}
}