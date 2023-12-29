package note

import (
	"errors"
	"time"

	uuid "github.com/satori/go.uuid"
)

// Note Сущность заметка
type Note struct {
	ID        uuid.UUID `json:"id"`
	Label     string    `json:"label"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}

func (n Note) Validate() error {
	if n.ID == uuid.Nil {
		return errors.New("note ID is required")
	}
	if len(n.Label) == 0 {
		return errors.New("note label is required")
	}

	return nil
}

type ListNotes struct {
	Notes []Note `json:"notes"`
	Total uint   `json:"total"`
}
