package adaptor

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DIContainer struct {
	database *sql.DB
	log      *zap.Logger
}

func NewDIContainer() (*DIContainer, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err) // Не удалось создать логгер
	}
	defer logger.Sync()

	connStr := "user=postgres password=postgres dbname=notesapp sslmode=disable host=127.0.0.1"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	return &DIContainer{
		database: db,
		log:      logger,
	}, nil
}

func (di *DIContainer) GetNoteAdaptor(ctx context.Context) *NoteStore {
	return NewNoteStore(di.database, di.log)
}

func (di *DIContainer) GetLogger() *zap.Logger {
	return di.log
}

func (di *DIContainer) Close() {
	di.database.Close()
}
