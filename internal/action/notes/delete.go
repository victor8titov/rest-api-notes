package notes

import (
	"context"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type DeleteAction struct {
	store Store
	log   *zap.Logger
}

func NewDeleteAction(store Store, log *zap.Logger) *DeleteAction {
	return &DeleteAction{
		store: store,
		log:   log,
	}
}

func (a *DeleteAction) Do(ctx context.Context, noteIDs []uuid.UUID) error {
	err := a.store.Delete(ctx, noteIDs)
	if err != nil {
		return errors.WithMessage(err, "Failed during action deleting")
	}

	a.log.Debug("Deleted notes", zap.Any("noteIDs", noteIDs))

	return nil
}
