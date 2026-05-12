package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"os"
)

// TestAuthenticationRequired tests that requests without valid session are rejected
func TestAuthenticationRequired(t *testing.T) {
	// Setup
	setupTest()
	
	// Cleanup
	defer func() {
		os.RemoveAll("./test_data")
	}()

	// Test POST without session header
	req, err := http.NewRequest("POST", "/notes", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Intentionally NOT setting X-Session-ID header
	
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	
	// Should return 401 Unauthorized
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
	
	// Test GET without session header
	getReq, err := http.NewRequest("GET", "/notes/test-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Intentionally NOT setting X-Session-ID header
	
	getRR := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(getRR, getReq)
	
	// Should return 401 Unauthorized
	if status := getRR.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}

// TestInvalidSession tests that requests with invalid session are rejected
func TestInvalidSession(t *testing.T) {
	// Setup
	setupTest()
	
	// Cleanup
	defer func() {
		os.RemoveAll("./test_data")
	}()

	// Test with invalid session ID
	req, err := http.NewRequest("POST", "/notes", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("X-Session-ID", "invalid-session-id")
	
	rr := httptest.NewRecorder()
	http.DefaultServeMux.ServeHTTP(rr, req)
	
	// Should return 401 Unauthorized
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}
}
