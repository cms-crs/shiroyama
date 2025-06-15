package team

import (
	"database/sql"
	"log/slog"
	"userservice/internal/dto"
)

type Repository struct {
	Log *slog.Logger
	db  *sql.DB
}

func NewTeamRepository(log *slog.Logger, db *sql.DB) *Repository {
	return &Repository{Log: log, db: db}
}
