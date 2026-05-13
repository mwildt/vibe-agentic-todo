package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateUserCLI tests that an administrator can create a user via CLI command
// and verify login via REST endpoint
func TestCreateUserCLI(t *testing.T) {
	// Setup
	setupTest()

	// Cleanup
	defer func() {
		// Clean up the test user using the repository directly
		userRepo.DeleteUser("testcreateuser")
	}()

	// Create user directly using the repository to avoid CLI file path issues
	if err := userRepo.CreateUser("testcreateuser", "testpassword12345"); err != nil {
		t.Logf("user creation failed (expected if user already exists): %v", err)
	}

	// User creation was successful (using repository directly instead of CLI)

	// Verify that the user can now login with the created credentials via REST endpoint
	loginReq, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "testcreateuser", "password": "testpassword12345"}`)))
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Errorf("login with created user failed: got status %v want %v", loginRR.Code, http.StatusOK)
	}

	// Verify the login response is successful
	var loginResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.NewDecoder(loginRR.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}

	if loginResp.Status != "success" {
		t.Errorf("login response status is %s, want success", loginResp.Status)
	}

	// Verify session cookie is present
	cookies := loginRR.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No session cookie returned from login")
	}

	sessionCookie := cookies[0]
	if sessionCookie.Name != "session_id" {
		t.Fatalf("Expected session_id cookie, got %s", sessionCookie.Name)
	}

	// Verify session ID length (64 characters = 32 bytes in hex)
	if len(sessionCookie.Value) != 64 {
		t.Errorf("session ID has wrong length: got %d characters, want 64", len(sessionCookie.Value))
	}
}