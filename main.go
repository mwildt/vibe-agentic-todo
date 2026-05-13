package main

import (
	"log"
	"net/http"
	"vibe-agentic/auth"
	"vibe-agentic/auth/user"
	"vibe-agentic/middleware"
	"vibe-agentic/notes"
)

func main() {
	// Initialize session store
	sessionStore := auth.NewInMemorySessionStore()

	// Initialize user repository
	userRepo := user.NewYAMLUserRepository("users.yaml")
	if err := userRepo.LoadUsers(); err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}

	// Initialize security logger
	securityLogger := middleware.NewSecurityLogger()

	// Register handlers with session store and user repository
	auth.RegisterHandlers(sessionStore, userRepo)
	middleware.SetSessionStore(sessionStore)

	repo := notes.NewJSONNoteRepository("./data/notes")
	notes.RegisterHandlers(repo)

	// Create main handler with security headers, error handling, and logging middleware
	mainHandler := middleware.SecurityHeadersMiddleware(
		middleware.ErrorHandler(
			middleware.LoggingMiddleware(securityLogger, http.DefaultServeMux),
		),
	)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", mainHandler))
}
