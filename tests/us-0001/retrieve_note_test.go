package tests

import (
	"testing"
	"net/http"
	"net/http/httptest"
	"strings"
	"os"
)

// TestRetrieveNote tests that a saved note can be retrieved later
func TestRetrieveNote(t *testing.T) {
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
	t.Logf("Extracted session ID: %s (length: %d)", sessionID, len(sessionID))

	// First, create a note to retrieve
	createReq, err := http.NewRequest("POST", "/notes", strings.NewReader(`{"text": "Test note to retrieve"}`))
	if err != nil {
		t.Fatal(err)
	}
	createReq.Header.Set("X-Session-ID", sessionID)
	
	createRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(createRR, createReq)
	
	if createRR.Code != http.StatusCreated {
		t.Fatalf("Failed to create note: got status %v", createRR.Code)
	}
	
	// Extract the note ID from the response
	createdNote := createRR.Body.String()
	if !strings.Contains(createdNote, `"id":"`) {
		t.Fatal("Created note response does not contain an ID")
	}
	
	// Extract ID (simple extraction - in real code you'd use proper JSON parsing)
	idStart := strings.Index(createdNote, `"id":"`) + 6
	idEnd := strings.Index(createdNote[idStart:], `"`)
	noteID := createdNote[idStart : idStart+idEnd]
	
	// Now retrieve the note
	retrieveReq, err := http.NewRequest("GET", "/notes/"+noteID, nil)
	if err != nil {
		t.Fatal(err)
	}
	retrieveReq.Header.Set("X-Session-ID", sessionID)
	
	retrieveRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(retrieveRR, retrieveReq)
	
	// Check the status code
	if status := retrieveRR.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	
	// Check that the response contains the note text
	expectedText := "Test note to retrieve"
	if !strings.Contains(retrieveRR.Body.String(), expectedText) {
		t.Errorf("retrieved note does not contain expected text: got %v want %v", retrieveRR.Body.String(), expectedText)
	}
	
	// Check that the response contains the note ID
	if !strings.Contains(retrieveRR.Body.String(), noteID) {
		t.Errorf("retrieved note does not contain expected ID: got %v want %v", retrieveRR.Body.String(), noteID)
	}
}
