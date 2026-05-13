package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
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
	loginReq, err := http.NewRequest("POST", "/login", bytes.NewBufferString(`{"username": "testuser", "password": "testpass"}`))
	if err != nil {
		t.Fatal(err)
	}
	loginReq.Header.Set("Content-Type", "application/json")

	loginRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(loginRR, loginReq)

	if loginRR.Code != http.StatusOK {
		t.Fatalf("Login failed: got status %v", loginRR.Code)
	}

	// Parse login response using proper JSON unmarshaling
	var loginResp LoginResponse
	if err := json.NewDecoder(loginRR.Body).Decode(&loginResp); err != nil {
		t.Fatalf("Failed to parse login response: %v", err)
	}

	// Verify login was successful
	if loginResp.Status != "success" {
		t.Fatalf("Login failed: status is %s, want success", loginResp.Status)
	}

	// Extract session cookie for subsequent requests
	cookies := loginRR.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No session cookie returned from login")
	}

	sessionCookie := cookies[0]
	if sessionCookie.Name != "session_id" {
		t.Fatalf("Expected session_id cookie, got %s", sessionCookie.Name)
	}

	// Verify session ID length
	if len(sessionCookie.Value) != 64 {
		t.Errorf("Session ID has wrong length: got %d, want 64", len(sessionCookie.Value))
	}

	// Create a request to the notes endpoint with valid session cookie
	noteRequest := NoteRequest{Text: "Test note"}
	reqBody, err := json.Marshal(noteRequest)
	if err != nil {
		t.Fatalf("Failed to marshal note request: %v", err)
	}

	req, err := http.NewRequest("POST", "/notes", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	// Add the session cookie to the request
	req.AddCookie(sessionCookie)
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	http.DefaultServeMux.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Parse the response body properly
	var createdNote NoteResponse
	if err := json.NewDecoder(rr.Body).Decode(&createdNote); err != nil {
		t.Fatalf("Failed to parse created note response: %v", err)
	}

	// Verify the note has an ID
	if createdNote.ID == "" {
		t.Errorf("Created note response does not contain an ID")
	}

	// Verify the note contains the expected text
	if createdNote.Text != "Test note" {
		t.Errorf("Created note does not contain expected text: got %q want %q", createdNote.Text, "Test note")
	}
}