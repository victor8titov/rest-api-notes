package adaptor

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/victor8titov/rest-api-notes/internal/action/notes"
	"github.com/victor8titov/rest-api-notes/internal/entity/note"
	"go.uber.org/zap"
)

const NoteTable = "notes"

type NoteStore struct {
	db  *sql.DB
	log *zap.Logger
}

func NewNoteStore(db *sql.DB, logger *zap.Logger) *NoteStore {
	return &NoteStore{
		db:  db,
		log: logger,
	}
}

func (s *NoteStore) CreateTable(ctx context.Context) error {
	s.log.Debug("creating table", zap.Any("table", NoteTable))

	query := fmt.Sprintf(
		`CREATE TABLE IF NOT EXISTS %v (
			id UUID PRIMARY KEY NOT NULL,
			label TEXT NOT NULL,
			body TEXT,
			tags TEXT[],
			created_at timestamptz NOT NULL
		);`,
		NoteTable,
	)
	result, err := s.db.Exec(query)
	if err != nil {
		s.log.Error("create table", zap.Error(err))
		return errors.Wrapf(err, "create table %v", NoteTable)
	}

	n, err := result.RowsAffected()
	if err != nil {
		errors.WithMessage(err, "Failed during getting rows affected")
	}
	s.log.Debug("created table", zap.Any("rows", n))
	return nil

}

func (s *NoteStore) Create(ctx context.Context, note note.Note) error {
	s.log.Debug("saving note", zap.Any("note", note))

	query := fmt.Sprintf(
		`INSERT INTO %v (id, label, body, tags, created_at) VALUES ($1, $2, $3, $4, $5)`,
		NoteTable,
	)

	data, err := NoteFromEntity(note)
	if err != nil {
		return errors.WithMessage(err, "convert note to data for database")
	}

	result, err := s.db.Exec(
		query,
		data.ID,
		data.Label,
		data.Body,
		pq.Array(data.Tags),
		data.CreatedAt,
	)
	if err != nil {
		s.log.Debug("failed save to new note to db", zap.Any("err", err))
		return errors.Wrap(err, "save note to database")
	}

	n, err := result.RowsAffected()
	if err != nil {
		errors.WithMessage(err, "Failed during getting rows affected")
	}
	s.log.Debug("created", zap.Any("rows", n))

	return nil
}

func (s *NoteStore) Update(ctx context.Context, args notes.UpdateArgs) error {
	s.log.Debug("updating note", zap.Any("args", args))

	query := fmt.Sprintf(
		`UPDATE %v SET
			label = $1,
			body = $2,
			tags = $3
			WHERE id = $4
		`,
		NoteTable,
	)

	result, err := s.db.Exec(
		query,
		args.Label,
		args.Body,
		pq.Array(args.Tags),
		args.ID,
	)
	if err != nil {
		s.log.Debug("failed update note to db", zap.Any("err", err))
		return errors.Wrap(err, "update note to database")
	}

	count, err := result.RowsAffected()
	if err != nil {
		errors.WithMessage(err, "Failed during getting rows affected")
	}
	s.log.Debug("updated", zap.Any("count", count))

	return nil
}

func (s *NoteStore) Delete(ctx context.Context, noteIDs []uuid.UUID) error {
	s.log.Debug("deleting note by ids", zap.Any("note ids", noteIDs))

	idString := make([]string, len(noteIDs))
	for key, value := range noteIDs {
		idString[key] = value.String()
	}

	query := fmt.Sprintf(
		`DELETE FROM %v WHERE id=ANY($1::uuid[])`,
		NoteTable,
	)

	result, err := s.db.Exec(query, pq.Array(idString))
	if err != nil {
		s.log.Debug("Failed delete by notes id", zap.Any("error", err))
		return errors.Wrap(err, "delete notes by ids")
	}

	count, err := result.RowsAffected()
	if err != nil {
		errors.WithMessage(err, "Failed during getting rows affected")
	}

	s.log.Debug("deleted", zap.Any("count", count))

	return nil
}

func (s *NoteStore) GetByID(ctx context.Context, id uuid.UUID) (note.Note, error) {
	s.log.Debug("getting note by ID", zap.Any("noteID", id))

	query := fmt.Sprintf(
		`SELECT * FROM %v WHERE id = $1::uuid`,
		NoteTable,
	)

	rows, err := s.db.Query(query, id)
	if err != nil {
		return note.Note{}, errors.Wrap(err, "failed during get note by ID")
	}
	defer rows.Close()

	listNotes := []Note{}

	for rows.Next() {
		note := Note{}
		err = rows.Scan(&note.ID, &note.Label, &note.Body, pq.Array(&note.Tags), &note.CreatedAt)
		if err != nil {
			s.log.Debug("scan line", zap.Error(err))
			errors.WithMessage(err, "Failed during Scan rows to dest")
			continue
		}
		listNotes = append(listNotes, note)
	}

	if len(listNotes) == 0 {
		return note.Note{}, errors.New("Not Found")
	}

	firstNote := listNotes[0]
	result, err := NoteToEntity(firstNote)
	if err != nil {
		return note.Note{}, errors.WithMessage(err, "failed during convert note to entity")
	}

	return result, nil
}

func (s *NoteStore) Query(ctx context.Context, args notes.ListArgs) ([]note.Note, error) {
	s.log.Debug("getting notes with pagination and order", zap.Any("args", args))

	order, limit, offset := s.getPaginationParams(args)

	query := fmt.Sprintf(
		`SELECT * FROM %v ORDER BY %v LIMIT %v OFFSET %v`,
		NoteTable, order, limit, offset,
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return []note.Note{}, errors.Wrap(err, "failed during get notes")
	}
	defer rows.Close()

	notes := []note.Note{}

	for rows.Next() {
		note := Note{}
		err = rows.Scan(&note.ID, &note.Label, &note.Body, pq.Array(&note.Tags), &note.CreatedAt)
		if err != nil {
			s.log.Debug("scan line", zap.Error(err))
			errors.WithMessage(err, "Failed during Scan rows to dest")
			continue
		}
		n, err := NoteToEntity(note)
		if err != nil {
			s.log.Debug("convert to entity", zap.Error(err))
			continue
		}
		notes = append(notes, n)
	}

	if len(notes) == 0 {
		return []note.Note{}, errors.New("Not Found")
	}

	return notes, nil
}

func (s *NoteStore) getPaginationParams(args notes.ListArgs) (order, limit string, offset uint) {
	orderBy := "label"
	if args.SortBy == notes.SortFieldDate {
		orderBy = "created_at"
	}

	direction := "ASC"
	if args.SortDirection == notes.SortDirectionDesc {
		direction = "DESC"
	}
	order = fmt.Sprintf("%v %v", orderBy, direction)

	limit = "ALL"
	if args.Limit > 0 {
		limit = fmt.Sprintf("%v", args.Limit)
	}

	offset = args.Offset

	return order, limit, offset
}

func (s *NoteStore) Count(ctx context.Context) (uint, error) {
	s.log.Debug("counting notes")

	query := fmt.Sprintf(
		`SELECT COUNT(*) FROM %v`,
		NoteTable,
	)

	count := new(uint)
	err := s.db.QueryRow(query).Scan(count)
	if err != nil {
		return 0, errors.WithMessage(err, "count notes")
	}

	return *count, nil
}
