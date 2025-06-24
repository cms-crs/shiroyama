package database

import (
	"gorm.io/gorm"
	"taskservice/internal/config"
	"taskservice/internal/entity"
)

func MustLoad(cfg *config.Config) *gorm.DB {

	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err.Error())
	}

	err = db.AutoMigrate(&entity.Task{})
	if err != nil {
		panic(err.Error())
	}

	return db
}
