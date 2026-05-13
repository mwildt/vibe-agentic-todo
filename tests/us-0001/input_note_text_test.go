package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestInputNoteText tests that a user can input text into a note
func TestInputNoteText(t *testing.T) {
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

	// Extract session cookie for subsequent requests
	cookies := loginRR.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("No session cookie returned from login")
	}

	sessionCookie := cookies[0]
	if sessionCookie.Name != "session_id" {
		t.Fatalf("Expected session_id cookie, got %s", sessionCookie.Name)
	}

	// Verify session ID length (from cookie)
	if len(sessionCookie.Value) != 64 {
		t.Errorf("Session ID has wrong length: got %d, want 64", len(sessionCookie.Value))
	}

	// Create a note with specific text content
	noteText := "This is a test note with specific content that the user wants to input"
	noteRequest := NoteRequest{Text: noteText}
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

	// Parse the created note response
	var createdNote NoteResponse
	if err := json.NewDecoder(rr.Body).Decode(&createdNote); err != nil {
		t.Fatalf("Failed to parse created note response: %v", err)
	}

	// Verify that the created note contains the exact text that was input
	if createdNote.Text != noteText {
		t.Errorf("created note does not contain expected text: got %v want %v", createdNote.Text, noteText)
	}

	// Verify that the note has an ID
	if createdNote.ID == "" {
		t.Fatal("Created note response does not contain an ID")
	}

	// Verify that the note was created with the correct text by retrieving it
	retrieveReq, err := http.NewRequest("GET", "/notes/"+createdNote.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add the session cookie to the retrieve request
	retrieveReq.AddCookie(sessionCookie)

	retrieveRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(retrieveRR, retrieveReq)

	// Check the status code
	if status := retrieveRR.Code; status != http.StatusOK {
		t.Errorf("retrieve handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Parse the retrieved note response
	var retrievedNote NoteResponse
	if err := json.NewDecoder(retrieveRR.Body).Decode(&retrievedNote); err != nil {
		t.Fatalf("Failed to parse retrieved note response: %v", err)
	}

	// Verify that the retrieved note contains the exact text that was input
	if retrievedNote.Text != noteText {
		t.Errorf("retrieved note does not contain expected text: got %v want %v", retrievedNote.Text, noteText)
	}

	// Verify that the retrieved note has the same ID
	if retrievedNote.ID != createdNote.ID {
		t.Errorf("retrieved note has different ID: got %v want %v", retrievedNote.ID, createdNote.ID)
	}
}