package team

import (
	"database/sql"
	"log/slog"
)

type Repository struct {
	log *slog.Logger
	db  *sql.DB
}

func NewTeamRepository(log *slog.Logger, db *sql.DB) *Repository {
	return &Repository{log: log, db: db}
}
