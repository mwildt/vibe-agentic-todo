package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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

	sessionID := loginResp.SessionID

	// Verify session ID length
	if len(sessionID) != 64 {
		t.Errorf("Session ID has wrong length: got %d, want 64", len(sessionID))
	}

	// First, create a note to retrieve
	noteText := "Test note to retrieve"
	noteRequest := NoteRequest{Text: noteText}
	reqBody, err := json.Marshal(noteRequest)
	if err != nil {
		t.Fatalf("Failed to marshal note request: %v", err)
	}

	createReq, err := http.NewRequest("POST", "/notes", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	createReq.Header.Set("X-Session-ID", sessionID)
	createReq.Header.Set("Content-Type", "application/json")

	createRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(createRR, createReq)

	if createRR.Code != http.StatusCreated {
		t.Fatalf("Failed to create note: got status %v", createRR.Code)
	}

	// Parse the created note response
	var createdNote NoteResponse
	if err := json.NewDecoder(createRR.Body).Decode(&createdNote); err != nil {
		t.Fatalf("Failed to parse created note response: %v", err)
	}

	// Verify the created note has an ID
	if createdNote.ID == "" {
		t.Fatal("Created note response does not contain an ID")
	}

	// Now retrieve the note
	retrieveReq, err := http.NewRequest("GET", "/notes/"+createdNote.ID, nil)
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

	// Parse the retrieved note response
	var retrievedNote NoteResponse
	if err := json.NewDecoder(retrieveRR.Body).Decode(&retrievedNote); err != nil {
		t.Fatalf("Failed to parse retrieved note response: %v", err)
	}

	// Check that the response contains the note text
	if retrievedNote.Text != noteText {
		t.Errorf("retrieved note does not contain expected text: got %q want %q", retrievedNote.Text, noteText)
	}

	// Check that the response contains the note ID
	if retrievedNote.ID != createdNote.ID {
		t.Errorf("retrieved note has different ID: got %q want %q", retrievedNote.ID, createdNote.ID)
	}
}