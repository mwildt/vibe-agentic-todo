package middleware

import (
	"log"
	"net/http"
)

// ErrorHandler handles errors consistently and sanitizes error messages
func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a response writer wrapper to capture status codes
		type responseWriter struct {
			http.ResponseWriter
			statusCode int
		}
		
		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		
		// Call the next handler
		next.ServeHTTP(wrapper, r)
		
		// Log errors for server-side issues
		if wrapper.statusCode >= 500 {
			log.Printf("ERROR: %s %s - Status: %d", r.Method, r.URL.Path, wrapper.statusCode)
		}
	})
}

// SanitizeError returns a generic error message to prevent information leakage
func SanitizeError(w http.ResponseWriter, err error, statusCode int) {
	log.Printf("ERROR: %v", err) // Log the actual error internally
	http.Error(w, "An error occurred", statusCode) // Return generic message to client
}