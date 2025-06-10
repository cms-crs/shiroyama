package repository

import (
	"authservice/src/model"
	"context"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) (*AuthRepository, error) {
	err := db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	return &AuthRepository{
		db: db,
	}, nil
}

func (repo *AuthRepository) CreateUser(ctx context.Context, user model.User) (uint, error) {
	// create user
	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (repo *AuthRepository) GetUser(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := repo.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
