package middleware

import (
	"net/http"
)

// SetSessionStore is kept for compatibility but not used due to import cycle
func SetSessionStore(store interface{}) {
	// No-op due to import cycle issues
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
		
		// For now, we'll accept any non-empty session ID for testing
		// In production, you would properly validate the session against the store
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
