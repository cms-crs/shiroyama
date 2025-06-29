package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"taskservice/internal/config"
)

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConnections)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConnections)
	db.SetConnMaxLifetime(cfg.DB.MaxLifetime)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
