package notes

type NoteRepository interface {
	Save(note Note) error
	FindByID(id string) (Note, error)
}
