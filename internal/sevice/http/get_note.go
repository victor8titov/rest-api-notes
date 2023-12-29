package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/action/notes"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

type GetByIDAction interface {
	Do(ctx context.Context, noteID uuid.UUID) (note.Note, error)
}

type GetByIDHandler struct {
	action GetByIDAction
	log    *zap.Logger
}

func NewGetByIDHandler(action GetByIDAction, log *zap.Logger) *GetByIDHandler {
	return &GetByIDHandler{
		action: action,
		log:    log,
	}
}

func (h *GetByIDHandler) Handle(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")

	h.log.Debug("hadle get note by id", zap.Any("note id", noteID))

	id, err := uuid.FromString(noteID)
	if err != nil {
		h.log.Debug("failed to covert type string to uuid")
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	note, err := h.action.Do(ctx, id)
	if err != nil && (errors.Is(err, notes.NotFound) || err.Error() == notes.NotFound.Error()) {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		h.log.Debug("failed during action doing", zap.Error(err))
		http.Error(w, "failed during inner process", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(note)
	if err != nil {
		h.log.Debug("failed marshal note", zap.Any("note", note), zap.Error(err))
		http.Error(w, "failed during inner process", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		h.log.Debug("failed during write response", zap.Error(err))
		http.Error(w, "failed during creating", http.StatusInternalServerError)
		return
	}
}
