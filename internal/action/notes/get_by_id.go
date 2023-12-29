package notes

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

type GetByIDAction struct {
	store Store
	log   *zap.Logger
}

func NewGetByIDAction(
	store Store,
	log *zap.Logger,
) *GetByIDAction {
	return &GetByIDAction{
		store: store,
		log:   log,
	}
}

func (a *GetByIDAction) Do(ctx context.Context, noteID uuid.UUID) (note.Note, error) {
	n, err := a.store.GetByID(ctx, noteID)
	switch {
	case err != nil && (errors.Is(err, NotFound) || err.Error() == NotFound.Error()):
		return note.Note{}, err
	case err != nil:
		return note.Note{}, errors.WithMessage(err, "Failed during getting from store")
	}

	a.log.Debug("Getting note from store.", zap.Any("note", n))

	return n, nil
}
