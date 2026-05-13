package tests

// LoginResponse represents the response from a successful login
// This is shared between test files to avoid duplicate declarations
type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
