package database

import (
	"database/sql"
	"userservice/internal/config"
)

func MustLoad(cfg *config.Config) *sql.DB {

	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err.Error())
	}

	return db
}
