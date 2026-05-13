package tests

// LoginResponse represents the response from a successful login
// This is shared between test files to avoid duplicate declarations
type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NoteRequest represents a request to create a note
type NoteRequest struct {
	Text string `json:"text"`
}

// NoteResponse represents a note in the response
type NoteResponse struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
