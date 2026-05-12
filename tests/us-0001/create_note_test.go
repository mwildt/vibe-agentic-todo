package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"vibe-agentic/notes"
)

var (
	repo     notes.NoteRepository
	setupOnce sync.Once
)

func setupTest() {
	setupOnce.Do(func() {
		repo = notes.NewJSONNoteRepository("./test_data/notes")
		notes.RegisterHandlers(repo)
	})
}

// TestCreateNote tests that a user can create a new note
func TestCreateNote(t *testing.T) {
	// Setup
	setupTest()

	// Cleanup
	defer func() {
		os.RemoveAll("./test_data")
	}()

	// Create a request to the notes endpoint
	req, err := http.NewRequest("POST", "/notes", strings.NewReader(`{"text": "Test note"}`))
	if err != nil {
		t.Fatal(err)
	}

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
