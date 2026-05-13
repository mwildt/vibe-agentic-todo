package auth

import (
	"net/http"
	"vibe-agentic/middleware"

	"encoding/json"
	"github.com/go-playground/validator/v10"
	"vibe-agentic/auth/user"
)

var (
	sessionStore   SessionStore
	userRepo       user.UserRepository
	rateLimiter    *RateLimiter
	securityLogger *middleware.SecurityLogger
	validate       = validator.New()
)

type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50,alphanum"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// LoginResponse is no longer needed as session is managed via cookies
type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func RegisterHandlers(store SessionStore, userRepository user.UserRepository) {
	sessionStore = store
	userRepo = userRepository
	rateLimiter = NewRateLimiter()
	securityLogger = middleware.NewSecurityLogger()

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}

		// Rate limiting - maximum 5 attempts per minute per IP
		ip := r.RemoteAddr
		if !rateLimiter.AllowRequest(ip) {
			securityLogger.LogSecurityEvent("rate_limit_triggered", "", ip, false, map[string]interface{}{
				"endpoint": "/login",
				"limit":    "5 attempts per minute",
			})
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Validate Content-Type header
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			middleware.SanitizeError(w, err, http.StatusBadRequest)
			return
		}

		// Validate request structure and content
		if err := validate.Struct(req); err != nil {
			middleware.SanitizeError(w, err, http.StatusBadRequest)
			return
		}

		// Validate credentials (simple validation for now)
		if !isValidUser(req.Username, req.Password) {
			securityLogger.LogSecurityEvent("login_attempt", req.Username, ip, false, map[string]interface{}{
				"reason": "invalid_credentials",
			})
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Log successful login
		securityLogger.LogSecurityEvent("login_attempt", req.Username, ip, true, map[string]interface{}{
			"reason": "successful_authentication",
		})

		// Create session in store
		session, err := sessionStore.CreateSession(req.Username)
		if err != nil {
			securityLogger.LogSecurityEvent("session_creation", req.Username, ip, false, map[string]interface{}{
				"error": err.Error(),
			})
			middleware.SanitizeError(w, err, http.StatusInternalServerError)
			return
		}

		// Log successful session creation
		securityLogger.LogSecurityEvent("session_creation", req.Username, ip, true, map[string]interface{}{
			"session_id": session.ID,
		})

		// Set session ID as HTTP cookie instead of returning in JSON
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session.ID,
			Path:     "/",
			HttpOnly: true,  // Prevent JavaScript access
			Secure:   false, // Set to true in production with HTTPS
			SameSite: http.SameSiteLaxMode,
			MaxAge:   24 * 60 * 60, // 24 hours
		})

		// Return success response without session ID in body
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "Login successful"})
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
