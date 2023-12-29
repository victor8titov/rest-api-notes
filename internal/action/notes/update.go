package notes

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

type UpdateArgs struct {
	ID    uuid.UUID `json:"id"`
	Label string    `json:"label"`
	Body  string    `json:"body"`
	Tags  []string  `json:"tags"`
}

type UpdateAction struct {
	store Store
	log   *zap.Logger
}

func NewUpdateAction(store Store, log *zap.Logger) *UpdateAction {
	return &UpdateAction{
		store: store,
		log:   log,
	}
}

func (a *UpdateAction) Do(ctx context.Context, args UpdateArgs) (note.Note, error) {

	err := a.store.Update(ctx, args)
	if err != nil {
		return note.Note{}, errors.WithMessage(err, "Failed action update")
	}

	getByIdAction := NewGetByIDAction(a.store, a.log)

	updatedNote, err := getByIdAction.Do(ctx, args.ID)
	if err != nil {
		return note.Note{}, errors.WithMessage(err, "failed during getting updated note")
	}

	a.log.Debug("Updated notes", zap.Any("note", updatedNote))

	return updatedNote, nil
}
