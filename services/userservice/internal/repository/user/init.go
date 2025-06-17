package user

import (
	"database/sql"
	"log/slog"
)

type Repository struct {
	Log *slog.Logger
	db  *sql.DB
}

func NewUserRepository(log *slog.Logger, db *sql.DB) *Repository {
	return &Repository{Log: log, db: db}
}
