package tests

// LoginResponse represents the response from a successful login
// This is shared between test files to avoid duplicate declarations
type LoginResponse struct {
	SessionID string `json:"session_id"`
}
