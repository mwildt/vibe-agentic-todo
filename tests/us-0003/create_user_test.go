package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestCreateUserCLI tests that an administrator can create a user via CLI command
// and verify login via REST endpoint
func TestCreateUserCLI(t *testing.T) {
	// Setup
	setupTest()

	// Cleanup
	defer func() {
		// Clean up the YAML file
		os.Remove("users.yaml")
	}()

	// No need to build separately, we'll use go run directly

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Construct absolute path to YAML file
	yamlPath := filepath.Join(cwd, "users.yaml")

	// Run the CLI command to create a user
	// We need to change to the project root directory first
	createCmd := exec.Command("go", "run", "vibe-agentic/cmd/cli", "user", "add", "--username", "testuser", "--password", "testpassword123", "--file", yamlPath)
	createCmd.Env = append(os.Environ(), "GO111MODULE=on")
	createCmd.Dir = ".."
	var stdout, stderr bytes.Buffer
	createCmd.Stdout = &stdout
	createCmd.Stderr = &stderr

	err = createCmd.Run()

	// Check that the command executed successfully
	if err != nil {
		t.Errorf("user add command failed: %v, stderr: %s", err, stderr.String())
	}

	// Check the output
	expectedSuccess := "User 'testuser' created successfully"
	if !bytes.Contains(stdout.Bytes(), []byte(expectedSuccess)) {
		t.Errorf("user add command output does not contain success message: got %s", stdout.String())
	}

	// Verify that the YAML file was created and contains the user
	if _, err := os.Stat("users.yaml"); os.IsNotExist(err) {
		t.Errorf("YAML file was not created")
	}

	// Verify that the user can now login with the created credentials via REST endpoint
	loginReq, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "testuser", "password": "testpassword123"}`)))
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Errorf("login with created user failed: got status %v want %v", loginRR.Code, http.StatusOK)
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

	// Verify session ID length (64 characters = 32 bytes in hex)
	if len(loginResp.SessionID) != 64 {
		t.Errorf("session ID has wrong length: got %d characters, want 64", len(loginResp.SessionID))
	}
}
