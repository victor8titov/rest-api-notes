package notes

import (
	"context"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

type CreateAction struct {
	store Store
	log   *zap.Logger
}

func NewCreateAction(
	store Store,
	log *zap.Logger,
) *CreateAction {
	return &CreateAction{
		store: store,
		log:   log,
	}
}

type CreateArgs struct {
	Label string   `json:"label"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

func (a *CreateAction) Do(ctx context.Context, args CreateArgs) (note.Note, error) {
	newNote := note.Note{
		ID:        uuid.UUID(ulid.Make()),
		Label:     args.Label,
		Body:      args.Body,
		Tags:      args.Tags,
		CreatedAt: time.Now(),
	}

	err := newNote.Validate()
	if err != nil {
		return note.Note{}, errors.WithMessage(err, "Failed validation, during saving note.")
	}

	a.log.Debug("Create note and validate it.", zap.Any("note", newNote))

	err = a.store.Create(ctx, newNote)
	if err != nil {
		return note.Note{}, errors.WithMessage(err, "Failed during save to store new Note")
	}

	return newNote, nil
}
