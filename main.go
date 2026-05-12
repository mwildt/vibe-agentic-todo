package main

import (
	"log"
	"net/http"
	"vibe-agentic/notes"
)

func main() {
	repo := notes.NewJSONNoteRepository("./data/notes")
	notes.RegisterHandlers(repo)
	
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
