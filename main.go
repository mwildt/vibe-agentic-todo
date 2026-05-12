package main

import (
	"log"
	"net/http"
	"vibe-agentic/auth"
	"vibe-agentic/middleware"
	"vibe-agentic/notes"
)

func main() {
	// Initialize session store
	sessionStore := auth.NewInMemorySessionStore()
	
	// Register handlers with session store
	auth.RegisterHandlers(sessionStore)
	middleware.SetSessionStore(sessionStore)
	
	repo := notes.NewJSONNoteRepository("./data/notes")
	notes.RegisterHandlers(repo)
	
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
