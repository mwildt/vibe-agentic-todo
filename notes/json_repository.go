package notes

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type JSONNoteRepository struct {
	baseDir string
}

func NewJSONNoteRepository(baseDir string) *JSONNoteRepository {
	return &JSONNoteRepository{baseDir: baseDir}
}

func (r *JSONNoteRepository) Save(note Note) error {
	// Create the notes directory if it doesn't exist
	if err := os.MkdirAll(r.baseDir, 0755); err != nil {
		return err
	}
	
	// Create the note file
	filePath := filepath.Join(r.baseDir, note.ID+".json")
	data, err := json.Marshal(note)
	if err != nil {
		return err
	}
	
	return os.WriteFile(filePath, data, 0644)
}

func (r *JSONNoteRepository) FindByID(id string) (Note, error) {
	filePath := filepath.Join(r.baseDir, id+".json")
	
	// Read the note file
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Note{}, errors.New("note not found")
		}
		return Note{}, err
	}
	
	// Parse the JSON data
	var note Note
	if err := json.Unmarshal(data, &note); err != nil {
		return Note{}, err
	}
	
	return note, nil
}
