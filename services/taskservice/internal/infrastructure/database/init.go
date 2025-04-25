package database

import (
	"gorm.io/gorm"
	"taskservice/internal/config"
)

func MustLoad(cfg *config.Config) *gorm.DB {

	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err.Error())
	}

	return db
}
