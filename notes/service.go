package notes

import "github.com/google/uuid"

type NoteService struct {
	repo NoteRepository
}

func NewNoteService(repo NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(text string) (Note, error) {
	note := Note{
		ID:   uuid.New().String(),
		Text: text,
	}
	
	if err := s.repo.Save(note); err != nil {
		return Note{}, err
	}
	
	return note, nil
}

func (s *NoteService) GetNote(id string) (Note, error) {
	return s.repo.FindByID(id)
}
