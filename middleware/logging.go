package middleware

import (
	"log"
	"net/http"
	"time"
)

// SecurityLogger logs security-relevant events
type SecurityLogger struct {
	logger *log.Logger
}

// NewSecurityLogger creates a new security logger
func NewSecurityLogger() *SecurityLogger {
	return &SecurityLogger{
		logger: log.Default(),
	}
}

// LogSecurityEvent logs a security event with details
func (sl *SecurityLogger) LogSecurityEvent(eventType, username, ipAddress string, success bool, details map[string]interface{}) {
	sl.logger.Printf("SECURITY_EVENT: type=%s username=%s ip=%s success=%t details=%v",
		eventType, username, ipAddress, success, details)
}

// LoggingMiddleware adds security logging to requests
func LoggingMiddleware(logger *SecurityLogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Create a response writer wrapper to capture status codes
		type responseWriter struct {
			http.ResponseWriter
			statusCode int
		}

		wrapper := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call the next handler
		next.ServeHTTP(wrapper, r)

		// Log the request
		logger.logger.Printf("REQUEST: %s %s %s - Status: %d, Duration: %v",
			r.Method, r.URL.Path, r.RemoteAddr, wrapper.statusCode, time.Since(startTime))
	})
}
