package notes

import (
	"context"

	"github.com/pkg/errors"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

type SortField string

const (
	SortFieldLabel SortField = "label"
	SortFieldDate  SortField = "created_at"
)

type SortDirection int

const (
	SortDirectionAsc SortDirection = iota
	SortDirectionDesc
)

type ListArgs struct {
	SortBy        SortField     `json:"sortBy"`
	SortDirection SortDirection `json:"direction"`
	Offset        uint          `json:"offset"`
	Limit         uint          `json:"limit"`
}

type ListAction struct {
	store Store
	log   *zap.Logger
}

func NewListAction(store Store, log *zap.Logger) *ListAction {
	return &ListAction{store: store, log: log}
}

func (a *ListAction) Do(ctx context.Context, args ListArgs) (note.ListNotes, error) {
	result, err := a.store.Query(ctx, args)
	switch {
	case errors.Cause(err) == NotFound || errors.Is(err, NotFound):
		return note.ListNotes{}, nil
	case err != nil:
		return note.ListNotes{}, errors.WithMessage(err, "list notes")
	}

	total, err := a.store.Count(ctx)
	if err != nil {
		return note.ListNotes{}, errors.WithMessage(err, "list notes")
	}

	return note.ListNotes{
		Notes: result,
		Total: total,
	}, nil

}
