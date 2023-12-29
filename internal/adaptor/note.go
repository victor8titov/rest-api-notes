package adaptor

import (
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
)

type Note struct {
	ID        uuid.UUID `db:"id"`
	Label     string    `db:"label"`
	Body      string    `db:"body"`
	Tags      []string  `db:"tags"`
	CreatedAt time.Time `db:"created_at"`
}

func NoteFromEntity(entity note.Note) (Note, error) {
	return Note{
		ID:        entity.ID,
		Label:     entity.Label,
		Body:      entity.Body,
		Tags:      entity.Tags,
		CreatedAt: entity.CreatedAt,
	}, nil
}

func NoteToEntity(noteAdapter Note) (note.Note, error) {
	return note.Note{
		ID:        noteAdapter.ID,
		Label:     noteAdapter.Label,
		Body:      noteAdapter.Body,
		Tags:      noteAdapter.Tags,
		CreatedAt: noteAdapter.CreatedAt,
	}, nil
}
