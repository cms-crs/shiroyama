package database

import (
	"boardservice/internal/config"
	"database/sql"
)

func MustLoad(cfg *config.Config) *sql.DB {

	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err.Error())
	}

	return db
}
