package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
	"github.com/victor8titov/rest-api-notes/internal/action/notes"
	"github.com/victor8titov/rest-api-notes/internal/adaptor"
	"github.com/victor8titov/rest-api-notes/internal/migrations"

	_ "github.com/victor8titov/rest-api-notes/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Service struct {
	di    *adaptor.DIContainer
	route *chi.Mux
}

// @title REST API Notes API
// @version 1.0
// @description Simple app notes.

// @contact.name Viktor
// @contact.email nulltomato@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /api/v1
func NewService(di *adaptor.DIContainer) *Service {
	httpService := &Service{di: di}
	httpService.newRouter()

	return httpService
}

func (hs *Service) newRouter() {
	root := chi.NewRouter()

	root.Use(middleware.Logger)
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	root.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:300/*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	root.Get("/api/v1/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/api/v1/swagger/doc.json"),
	))

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

// handleCreateNote
//
//	@Summary	Create note.
//	@Accept		json
//	@Produce	json
//	@Param		note	body	notes.CreateArgs	true	"fields for new note"
//	@Success	200	{object}	note.Note	"Ok"
//	@Failure		400		{string}	string	"invalid request params"
//	@Failure		404		{string}	string	"not found"
//	@Failure		500		{string}	string	"failed during inner process"
//	@Router			/note  [post]
func (hs *Service) handleCreateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewCreateAction(store, log)
	handler := NewCreateNoteHandler(action, log)

	handler.Handle(w, r)
}

// handleGetNoteByID
//
//	@Summary	Get note by ID.
//	@Produce	json
//	@Param		noteID	path	string	true	"ID of note that you want getting"
//	@Success	200	{object}	note.Note	"Ok"
//	@Failure		400		{string}	string	"invalid request params"
//	@Failure		404		{string}	string	"not found"
//	@Failure		500		{string}	string	"failed during inner process"
//	@Router			/note/{noteID}  [get]
func (hs *Service) handleGetNoteByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewGetByIDAction(store, log)
	handler := NewGetByIDHandler(action, log)

	handler.Handle(w, r)
}

// handleUpdateNote
//
//	@Summary	Update note.
//	@Produce	json
//	@Param		noteID	path	string	true	"ID of note that you want updating"
//	@Param		fields	body	RequestUpdateNote	true	"fields for updating note"
//	@Success	200	{object}	note.Note	"Updated note"
//	@Failure		400		{string}	string	"invalid request params"
//	@Failure		404		{string}	string	"not found"
//	@Failure		500		{string}	string	"failed during inner process"
//	@Router		/note/{noteID} [put]
func (hs *Service) handleUpdateNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewUpdateAction(store, log)
	handler := NewUpdateNoteHandler(action, log)

	handler.Handle(w, r)
}

// handleDeleteNote
//
//	@Summary	Delete note by ID.
//	@Param	noteID	path	string	true	"ID of note that you want to delete"
//	@Success	200	{string}	string	"Success deleting"
//	@Failure		400		{string}	string	"invalid request params"
//	@Failure		404		{string}	string	"not found"
//	@Failure		500		{string}	string	"failed during inner process"
//	@Router		/note/{noteID} [delete]
func (hs *Service) handleDeleteNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := hs.di.GetLogger()
	store := hs.di.GetNoteAdaptor(ctx)

	action := notes.NewDeleteAction(store, log)
	handler := NewDeleteNoteByIDHandler(action, log)

	handler.Handle(w, r)
}

// handleGetListNotes
//
//	@Summary		Getting list of notes.
//	@Description	Getting list with pagination.
//	@Accept			json
//	@Produce		json
//	@Param	pagination	body	RequestListNotes true	"params for pagination"
//	@Success		200		{object}	note.ListNotes			"ok"
//	@Failure		400		{string}	string	"invalid request params"
//	@Failure		404		{string}	string	"not found"
//	@Failure		500		{string}	string	"failed during inner process"
//	@Router			/note  [get]
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
