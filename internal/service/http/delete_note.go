package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type DeleteAction interface {
	Do(ctx context.Context, noteIDs []uuid.UUID) error
}

type DeleteRequest struct {
	NoteID []uuid.UUID `json:"noteId"`
}

type DeleteNoteByIDHandler struct {
	action DeleteAction
	log    *zap.Logger
}

func NewDeleteNoteByIDHandler(action DeleteAction, log *zap.Logger) *DeleteNoteByIDHandler {
	return &DeleteNoteByIDHandler{
		action: action,
		log:    log,
	}
}

func (h *DeleteNoteByIDHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	var requestParams DeleteRequest
	err = json.Unmarshal(body, &requestParams)
	if err != nil {
		h.log.Debug("failed unmarshal request body", zap.Any("body", body), zap.Error(err))
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	if len(requestParams.NoteID) == 0 {
		h.log.Debug("failed unmarshal request body", zap.Any("body", body), zap.Error(err))
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	err = h.action.Do(ctx, requestParams.NoteID)
	if err != nil {
		h.log.Debug("failed deleting action", zap.Error(err))
		http.Error(w, "failed during deleting", http.StatusInternalServerError)
		return
	}
}
