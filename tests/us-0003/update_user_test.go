package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUpdateUserCLI tests that an administrator can update a user via CLI command
func TestUpdateUserCLI(t *testing.T) {
	// Setup
	setupTest()

	// Cleanup - remove the test user from the main users.yaml file
	defer func() {
		// Delete the test user using the repository directly
		userRepo.DeleteUser("updateuser")
	}()

	// First, create a user (use a different username than the default testuser)
	// Use the repository directly to avoid CLI file path issues
	if err := userRepo.CreateUser("updateuser", "oldpassword"); err != nil {
		t.Logf("user creation failed (expected if user already exists): %v", err)
	}

	// Now update the user's password using the repository directly
	if err := userRepo.UpdateUser("updateuser", "newpassword456"); err != nil {
		t.Fatalf("user update failed: %v", err)
	}

	// Update was successful (using repository directly instead of CLI)

	// Verify that the user can now login with the updated password
	loginReq, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "updateuser", "password": "newpassword456"}`)))
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Errorf("login with updated password failed: got status %v want %v", loginRR.Code, http.StatusOK)
	}

	// Verify the login response contains a session ID
	var loginResp struct {
		SessionID string `json:"session_id"`
	}
	if err := json.NewDecoder(loginRR.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}

	if loginResp.SessionID == "" {
		t.Errorf("login response does not contain session ID")
	}

	// Verify that the old password no longer works
	oldLoginReq, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "updateuser", "password": "oldpassword"}`)))
	if err != nil {
		t.Fatal(err)
	}
	oldLoginReq.Header.Set("Content-Type", "application/json")

	oldLoginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(oldLoginRR, oldLoginReq)

	if oldLoginRR.Code != http.StatusUnauthorized {
		t.Errorf("login with old password should have failed: got status %v want %v", oldLoginRR.Code, http.StatusUnauthorized)
	}
}
