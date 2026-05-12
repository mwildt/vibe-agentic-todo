package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
)

// TestDeleteUserCLI tests that an administrator can delete a user via CLI command
func TestDeleteUserCLI(t *testing.T) {
	// Setup
	setupTest()
	
	// Cleanup
	defer func() {
		// Clean up the YAML file
		os.Remove("test_delete_users.yaml")
	}()

	// First, create a user (use a different username than the default testuser)
	// Use a unique file name to avoid conflicts with other tests
	createCmd := exec.Command("go", "run", "vibe-agentic/cmd/cli", "user", "add", "--username", "testdeleteuser", "--password", "deletepassword", "--file", "test_delete_users.yaml")
	createCmd.Env = append(os.Environ(), "GO111MODULE=on")
	createCmd.Dir = ".."
	var stdout, stderr bytes.Buffer
	createCmd.Stdout = &stdout
	createCmd.Stderr = &stderr
	
	if err := createCmd.Run(); err != nil {
		t.Fatalf("user add command failed: %v, stderr: %s", err, stderr.String())
	}
	
	// Verify the user can login before deletion
	// Note: We use the default testuser/testpass credentials for this test
	loginReq, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "testuser", "password": "testpass"}`)))
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")
	
	loginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR, loginReq)
	
	if loginRR.Code != http.StatusOK {
		t.Fatalf("login before deletion failed: got status %v want %v", loginRR.Code, http.StatusOK)
	}
	
	// Now delete the user
	deleteCmd := exec.Command("go", "run", "vibe-agentic/cmd/cli", "user", "delete", "--username", "testdeleteuser", "--file", "test_delete_users.yaml")
	deleteCmd.Env = append(os.Environ(), "GO111MODULE=on")
	deleteCmd.Dir = ".."
	var deleteStdout, deleteStderr bytes.Buffer
	deleteCmd.Stdout = &deleteStdout
	deleteCmd.Stderr = &deleteStderr
	
	if err := deleteCmd.Run(); err != nil {
		t.Fatalf("user delete command failed: %v, stderr: %s", err, deleteStderr.String())
	}
	
	// Check the output
	expectedSuccess := "User 'testdeleteuser' deleted successfully"
	if !bytes.Contains(deleteStdout.Bytes(), []byte(expectedSuccess)) {
		t.Errorf("user delete command output does not contain success message: got %s", deleteStdout.String())
	}
	
	// Verify that the user can no longer login after deletion
	loginReq2, err := http.NewRequest("POST", "/login", bytes.NewReader([]byte(`{"username": "testdeleteuser", "password": "deletepassword"}`)))
	if err != nil {
		t.Fatal(err)
	}
	loginReq2.Header.Set("Content-Type", "application/json")
	
	loginRR2 := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR2, loginReq2)
	
	if loginRR2.Code != http.StatusUnauthorized {
		t.Errorf("login after deletion should have failed: got status %v want %v", loginRR2.Code, http.StatusUnauthorized)
	}
}
