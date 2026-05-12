package middleware

import (
	"net/http"
	"vibe-agentic/auth"
)

var sessionStore auth.SessionStore

// SetSessionStore sets the session store for the auth middleware
func SetSessionStore(store auth.SessionStore) {
	sessionStore = store
}

// AuthMiddleware is a middleware that checks for valid session ID
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session ID in header
		sessionID := r.Header.Get("X-Session-ID")
		if sessionID == "" {
			http.Error(w, "Unauthorized: Session ID required", http.StatusUnauthorized)
			return
		}
		
		// Validate session against store
		if sessionStore == nil {
			http.Error(w, "Unauthorized: Session store not configured", http.StatusUnauthorized)
			return
		}
		
		_, valid := sessionStore.GetSession(sessionID)
		if !valid {
			http.Error(w, "Unauthorized: Invalid session", http.StatusUnauthorized)
			return
		}
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
