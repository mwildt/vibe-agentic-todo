package notes

import (
	"encoding/json"
	"net/http"
	"vibe-agentic/middleware"
)

type Note struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func RegisterHandlers(repo NoteRepository) {
	service := NewNoteService(repo)
	
	// Create handlers with auth middleware
	notesHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			// Validate Content-Type header
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
				return
			}

			var requestBody struct {
				Text string `json:"text"`
			}
			
			if r.Body == nil {
				http.Error(w, "Request body is required", http.StatusBadRequest)
				return
			}
			
			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				middleware.SanitizeError(w, err, http.StatusBadRequest)
				return
			}

			// Validate text length
			const MaxNoteTextLength = 10000
			if len(requestBody.Text) > MaxNoteTextLength {
				http.Error(w, "Note text too long", http.StatusBadRequest)
				return
			}
			
			note, err := service.CreateNote(requestBody.Text)
			if err != nil {
				middleware.SanitizeError(w, err, http.StatusInternalServerError)
				return
			}
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(note)
		} else if r.Method == http.MethodGet {
			// Handle GET /notes (list all notes - not implemented yet)
			http.NotFound(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
	
	// Handle GET /notes/{id}
	notesIDHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		
		// Extract note ID from URL
		noteID := r.URL.Path[len("/notes/"):]
		if noteID == "" {
			http.Error(w, "note ID is required", http.StatusBadRequest)
			return
		}
		
		note, err := service.GetNote(noteID)
		if err != nil {
			if err.Error() == "note not found" {
				http.Error(w, "note not found", http.StatusNotFound)
			} else {
				middleware.SanitizeError(w, err, http.StatusInternalServerError)
			}
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note)
	})
	
	// Register handlers with auth middleware
	http.Handle("/notes", middleware.AuthMiddleware(notesHandler))
	http.Handle("/notes/", middleware.AuthMiddleware(notesIDHandler))
}
