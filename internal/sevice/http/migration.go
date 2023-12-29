package http

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Migration interface {
	Up(ctx context.Context) error
}

type Migration01Handler struct {
	migration Migration
	log       *zap.Logger
}

func NewMigration01Handler(migration Migration, log *zap.Logger) *Migration01Handler {
	return &Migration01Handler{
		migration: migration,
		log:       log,
	}
}

func (h *Migration01Handler) Handle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	err := h.migration.Up(ctx)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("Error migration: %v", err)))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Success migration 01")))
}
