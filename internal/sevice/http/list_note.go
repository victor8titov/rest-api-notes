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

type ListAction interface {
	Do(ctx context.Context, args notes.ListArgs) (note.ListNotes, error)
}

type RequestListNotes struct {
	SortBy        string `json:"sortBy"`
	SortDirection uint   `json:"direction"`
	Offset        uint   `json:"offset"`
	Limit         uint   `json:"limit"`
}

type ListNotesHandler struct {
	action ListAction
	log    *zap.Logger
}

func NewListNotesHandler(action ListAction, log *zap.Logger) *ListNotesHandler {
	return &ListNotesHandler{action: action, log: log}
}

func (h *ListNotesHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	var requestParams RequestListNotes
	err = json.Unmarshal(body, &requestParams)
	if err != nil {
		h.log.Debug("failed unmarshal request body", zap.Any("body", body), zap.Error(err))
		http.Error(w, "invalid request params", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	args := notes.ListArgs{
		SortBy:        notes.SortField(requestParams.SortBy),
		SortDirection: notes.SortDirection(requestParams.SortDirection),
		Offset:        requestParams.Offset,
		Limit:         requestParams.Limit,
	}
	list, err := h.action.Do(ctx, args)
	if err != nil {
		h.log.Debug("failed during action doing", zap.Error(err))
		http.Error(w, "failed during inner process", http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(list)
	if err != nil {
		h.log.Debug("failed marshal list note", zap.Any("list", list), zap.Error(err))
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
