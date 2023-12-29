package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/victor8titov/rest-api-notes/internal/action/notes"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

type CreateAction interface {
	Do(ctx context.Context, args notes.CreateArgs) (note.Note, error)
}

type CreateNoteHandler struct {
	action CreateAction
	log    *zap.Logger
}

func NewCreateNoteHandler(action CreateAction, log *zap.Logger) *CreateNoteHandler {
	return &CreateNoteHandler{action: action, log: log}
}

func (h *CreateNoteHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	var requestParams notes.CreateArgs
	err = json.Unmarshal(body, &requestParams)
	if err != nil {
		h.log.Debug("failed unmarshal request body", zap.Any("body", body), zap.Error(err))
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	newNote, err := h.action.Do(ctx, requestParams)
	if err != nil {
		h.log.Debug("failed create action", zap.Error(err))
		http.Error(w, "failed during creating", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(newNote)
	if err != nil {
		h.log.Debug("failed marshal new note", zap.Any("new note", newNote), zap.Error(err))
		http.Error(w, "failed during creating", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(res)
	if err != nil {
		h.log.Debug("failed during write response", zap.Error(err))
		http.Error(w, "failed during creating", http.StatusInternalServerError)
		return
	}

	h.log.Debug("Handled create note", zap.Any("newNote", newNote))
}
