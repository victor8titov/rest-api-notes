package notes

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
)

type Store interface {
	Create(ctx context.Context, note note.Note) error
	Update(ctx context.Context, args UpdateArgs) error
	Delete(ctx context.Context, ids []uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (note.Note, error)
	Query(ctx context.Context, args ListArgs) ([]note.Note, error)
	Count(ctx context.Context) (uint, error)
}

var NotFound = errors.New("Not Found")
