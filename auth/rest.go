package auth

import (
	"encoding/json"
	"net/http"
)

var sessionStore SessionStore

type LoginRequest struct {
	Username string `json:"username"`
	Password  string `json:"password"`
}

type LoginResponse struct {
	SessionID string `json:"session_id"`
}

func RegisterHandlers(store SessionStore) {
	sessionStore = store
	
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		
		// Validate credentials (simple validation for now)
		if !isValidUser(req.Username, req.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		
		// Create session in store
		session, err := sessionStore.CreateSession(req.Username)
		if err != nil {
			http.Error(w, "Failed to create session", http.StatusInternalServerError)
			return
		}
		
		// Return session ID
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(LoginResponse{SessionID: session.ID})
	})
}

// isValidUser validates user credentials (simple implementation for now)
func isValidUser(username, password string) bool {
	// Simple validation: username and password should not be empty
	// In a real application, this would check against a user store
	return username != "" && password != "" && username == "testuser" && password == "testpass"
}
