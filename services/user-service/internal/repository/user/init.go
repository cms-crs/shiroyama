package user

import (
	"context"
	"database/sql"
	"log/slog"
)

type Repository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewUserRepository(log *slog.Logger, db *sql.DB) *Repository {
	return &Repository{log: log, db: db}
}

func (repository *Repository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return repository.db.BeginTx(ctx, nil)
}
