package migrations

import (
	"context"

	"github.com/pkg/errors"
)

type Migration01 struct {
	migrator Migrator
}

func NewMigration01(ctx context.Context, migrator Migrator) *Migration01 {
	return &Migration01{
		migrator: migrator,
	}
}

func (m *Migration01) Up(ctx context.Context) error {
	err := m.migrator.CreateTable(ctx)
	if err != nil {
		return errors.WithMessage(err, "create notes table")
	}

	return nil
}
