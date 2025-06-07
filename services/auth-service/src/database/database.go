package database

import (
	"authservice/src/config"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func MustConnect(cfg *config.Config) *gorm.DB {
	db, err := Connect(cfg)
	if err != nil {
		panic(err)
	}
	return db
}

func Connect(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DB.Host,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
