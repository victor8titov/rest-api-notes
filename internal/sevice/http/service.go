package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/victor8titov/rest-api-notes/internal/action/notes"
	"github.com/victor8titov/rest-api-notes/internal/adaptor"
	"github.com/victor8titov/rest-api-notes/internal/migrations"
)

type Service struct {
	di    *adaptor.DIContainer
	route *chi.Mux
}

func NewService(di *adaptor.DIContainer) *Service {
	httpService := &Service{di: di}
	httpService.newRouter()

	return httpService
}

func (hs *Service) newRouter() {
	root := chi.NewRouter()

	root.Use(middleware.Logger)

	root.Route("/api/v1/migration", func(router chi.Router) {
		router.Get("/01", hs.handleMigration01)
	})

	root.Route("/api/v1/note", func(router chi.Router) {
		router.Post("/", hs.handleCreateNote)
		router.Get("/", hs.handleGetListNotes)
		router.Get("/{noteID}", hs.handleGetNoteByID)
		router.Delete("/", hs.handleDeleteNote)
		router.Put("/{noteID}", hs.handleUpdateNote)
	})

	hs.route = root
}

func (hs *Service) ListenAndServe(port int) error {
	err := http.ListenAndServe(":"+strconv.Itoa(port), hs.route)
	return errors.WithMessage(err, "Failed listen and serve http service")
}

func (hs *Service) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewCreateAction(store, log)
	handler := NewCreateNoteHandler(action, log)

	handler.Handle(w, r)
}

func (hs *Service) handleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewGetByIDAction(store, log)
	handler := NewGetByIDHandler(action, log)

	handler.Handle(w, r)
}

func (hs *Service) handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewUpdateAction(store, log)
	handler := NewUpdateNoteHandler(action, log)

	handler.Handle(w, r)
}

func (hs *Service) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewDeleteAction(store, log)
	handler := NewDeleteNoteByIDHandler(action, log)

	handler.Handle(w, r)
}

func (hs *Service) handleGetListNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewListAction(store, log)
	handler := NewListNotesHandler(action, log)

	handler.Handle(w, r)
}

func (hs *Service) handleMigration01(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()

	migration := migrations.NewMigration01(ctx, hs.di.GetNoteAdaptor(ctx))

	handler := NewMigration01Handler(migration, log)

	handler.Handle(w, r)
}
