package auth

import (
	"net/http"

	"encoding/json"
	"vibe-agentic/auth/user"
)

var (
	sessionStore SessionStore
	userRepo     user.UserRepository
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	SessionID string `json:"session_id"`
}

func RegisterHandlers(store SessionStore, userRepository user.UserRepository) {
	sessionStore = store
	userRepo = userRepository

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

// isValidUser validates user credentials using the user repository
func isValidUser(username, password string) bool {

	// Check against user repository
	if userRepo == nil {
		return false
	}

	storedUser, found := userRepo.GetUser(username)
	if !found {
		return false
	}

	// Check if password matches using the user's VerifyPassword method
	return storedUser.VerifyPassword(password)
}
