package migrations

import "context"

type Migrator interface {
	CreateTable(ctx context.Context) error
}

type DIContainer interface {
	GetNoteAdaptor(ctx context.Context) Migrator
}
