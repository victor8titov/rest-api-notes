package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/action/notes"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

type UpdateAction interface {
	Do(ctx context.Context, args notes.UpdateArgs) (note.Note, error)
}

type RequestUpdateNote struct {
	Label string   `json:"label"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

type UpdateNoteHandler struct {
	action UpdateAction
	log    *zap.Logger
}

func NewUpdateNoteHandler(action UpdateAction, log *zap.Logger) *UpdateNoteHandler {
	return &UpdateNoteHandler{
		action: action,
		log:    log,
	}
}

func (h *UpdateNoteHandler) Handle(w http.ResponseWriter, r *http.Request) {
	noteID := chi.URLParam(r, "noteID")

	h.log.Debug("hading update note by id", zap.Any("note id", noteID))

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		h.log.Debug("invalid Content-Type header", zap.Any("contentType", contentType))
		http.Error(w, "invalid request header", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		h.log.Debug("invalid request body", zap.Any("body", body), zap.Error(err))
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	var requestParams RequestUpdateNote
	err = json.Unmarshal(body, &requestParams)
	if err != nil {
		h.log.Debug("failed unmarshal request body", zap.Any("body", body), zap.Error(err))
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	id, err := uuid.FromString(noteID)
	if err != nil {
		h.log.Debug("invalid note id", zap.Any("noteID", noteID), zap.Error(err))
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	args := notes.UpdateArgs{
		ID:    id,
		Label: requestParams.Label,
		Body:  requestParams.Body,
		Tags:  requestParams.Tags,
	}

	ctx := r.Context()
	updatedNote, err := h.action.Do(ctx, args)
	if err != nil {
		h.log.Debug("failed during action doing", zap.Error(err))
		http.Error(w, "failed during inner process", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(updatedNote)
	if err != nil {
		h.log.Debug("failed marshal note", zap.Any("note", updatedNote), zap.Error(err))
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
