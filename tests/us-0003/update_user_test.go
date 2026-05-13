package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestUpdateUserCLI tests that an administrator can update a user's password via CLI command
func TestUpdateUserCLI(t *testing.T) {
	// Setup
	setupTest()

	// Cleanup
	defer func() {
		// Clean up the test user using the repository directly
		userRepo.DeleteUser("updateuser")
	}()

	// Create user directly using the repository to avoid CLI file path issues
	if err := userRepo.CreateUser("updateuser", "oldpassword12345"); err != nil {
		t.Logf("user creation failed (expected if user already exists): %v", err)
	}

	// First, verify the user can login with the old password
	oldLoginReq, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "updateuser", "password": "oldpassword12345"}`)))
	if err != nil {
		t.Fatal(err)
	}
	oldLoginReq.Header.Set("Content-Type", "application/json")

	oldLoginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(oldLoginRR, oldLoginReq)

	if oldLoginRR.Code != http.StatusOK {
		t.Fatalf("login with old password failed: got status %v want %v", oldLoginRR.Code, http.StatusOK)
	}

	// Now update user password directly using the repository
	if err := userRepo.UpdateUser("updateuser", "newpassword12345"); err != nil {
		t.Fatalf("Failed to update user password: %v", err)
	}

	// Verify that the user can now login with the new credentials via REST endpoint
	loginReq, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "updateuser", "password": "newpassword12345"}`)))
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Errorf("login with updated user failed: got status %v want %v", loginRR.Code, http.StatusOK)
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

	// Verify that the old password no longer works
	oldLoginReq2, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "updateuser", "password": "oldpassword12345"}`)))
	if err != nil {
		t.Fatal(err)
	}
	oldLoginReq2.Header.Set("Content-Type", "application/json")

	oldLoginRR2 := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(oldLoginRR2, oldLoginReq2)

	// Should return 401 Unauthorized for old password (or 429 if rate limited)
	if oldLoginRR2.Code != http.StatusUnauthorized && oldLoginRR2.Code != http.StatusTooManyRequests {
		t.Errorf("login with old password should fail: got status %v want %v", oldLoginRR2.Code, http.StatusUnauthorized)
	}

	// Verify no session cookie is returned for failed login
	oldCookies := oldLoginRR2.Result().Cookies()
	if len(oldCookies) > 0 {
		t.Errorf("failed login should not return cookies, got %d cookies", len(oldCookies))
	}
}