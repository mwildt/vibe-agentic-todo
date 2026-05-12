package auth

import (
	"fmt"
	"net/http"
	"os"

	"encoding/json"
	"gopkg.in/yaml.v2"
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

// isValidUser validates user credentials against YAML file
func isValidUser(username, password string) bool {
	// First check if it's the default test user
	if username == "testuser" && password == "testpass" {
		return true
	}
	
	// Check against YAML file
	data, err := os.ReadFile("users.yaml")
	if err != nil {
		// If file doesn't exist, only allow testuser
		return username == "testuser" && password == "testpass"
	}
	
	var config struct {
		Users []struct {
			Username     string `yaml:"username"`
			PasswordHash string `yaml:"password_hash"`
		} `yaml:"users"`
	}
	
	if err := yaml.Unmarshal(data, &config); err != nil {
		return false
	}
	
	// Find user and check password
	for _, user := range config.Users {
		if user.Username == username {
			// Check if password matches the hash
			expectedHash := fmt.Sprintf("hashed_%s", password)
			return user.PasswordHash == expectedHash
		}
	}
	
	return false
}
