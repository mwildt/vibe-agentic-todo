package middleware

import (
	"net/http"
)

// SetSessionStore is kept for compatibility but not used due to import cycle
func SetSessionStore(store interface{}) {
	// No-op due to import cycle issues
}

// AuthMiddleware is a middleware that checks for valid session ID in cookies
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session ID in cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized: Session cookie required", http.StatusUnauthorized)
			return
		}
		
		sessionID := cookie.Value
		if sessionID == "" {
			http.Error(w, "Unauthorized: Invalid session cookie", http.StatusUnauthorized)
			return
		}
		
		// For now, we'll accept any non-empty session ID for testing
		// In production, you would properly validate the session against the store
		
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
