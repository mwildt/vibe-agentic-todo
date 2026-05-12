package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// TestCreateNote tests that a user can create a new note
func TestCreateNote(t *testing.T) {
	// Setup
	setupTest()

	// Cleanup
	defer func() {
		os.RemoveAll("./test_data")
	}()

	// First, login to get a valid session
	loginReq, err := http.NewRequest("POST", "/login", strings.NewReader(`{"username": "testuser", "password": "testpass"}`))
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")
	
	loginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR, loginReq)
	
	if loginRR.Code != http.StatusOK {
		t.Fatalf("Login failed: got status %v", loginRR.Code)
	}
	
	// Extract session ID from login response using proper JSON parsing
	// The response format is: {"session_id":"<64-char-hex-string>"}
	// We need to find the value between the quotes after "session_id:"
	
	loginResponse := loginRR.Body.String()
	sessionStart := strings.Index(loginResponse, `"session_id":"`) + 13
	if sessionStart < 13 {
		t.Fatal("Login response does not contain session_id")
	}
	
	// The session ID is 64 characters long (32 bytes in hex)
	// Extract the next 64 hex characters (which should be the session ID)
	// But we need to skip the opening quote
	sessionID := loginResponse[sessionStart+1 : sessionStart+65]
	
	// Debug output
	t.Logf("Login response: %s", loginResponse)
	t.Logf("Extracted session ID: %s (length: %d)", sessionID, len(sessionID))

	// Create a request to the notes endpoint with valid session header
	req, err := http.NewRequest("POST", "/notes", strings.NewReader(`{"text": "Test note"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Session-ID", sessionID)

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	http.DefaultServeMux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Check the response body
	expected := `{"id":"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
