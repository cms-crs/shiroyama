package repository

import (
	"authservice/src/config"
	"authservice/src/model"
	"context"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

type AuthRepository struct {
	db     *gorm.DB
	rdb    *redis.Client
	config *config.Config
	log    *slog.Logger
}

func NewAuthRepository(db *gorm.DB, rdb *redis.Client, config *config.Config, log *slog.Logger) (*AuthRepository, error) {
	err := createRoleType(db)
	if err != nil {
		log.Warn("some error happened while creating role type", err)
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, err
	}

	return &AuthRepository{
		db:     db,
		rdb:    rdb,
		config: config,
		log:    log,
	}, nil
}

func createRoleType(db *gorm.DB) error {
	err := db.Exec(`
		  DO $$
		  BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
			  CREATE TYPE user_role AS ENUM ('Admin', 'Regular');
			END IF;
		  END
		  $$;
		`).Error

	if err != nil {
		return err
	}

	return nil
}

func (repo *AuthRepository) CreateUser(ctx context.Context, user model.User) (uint, error) {
	if err := repo.db.WithContext(ctx).Create(&user).Error; err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (repo *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := repo.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *AuthRepository) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User

	if err := repo.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *AuthRepository) GetRefreshToken(ctx context.Context, userId string) (string, error) {
	if err := repo.rdb.Get(ctx, userId).Err(); err != nil {
		var user model.User
		if err := repo.db.WithContext(ctx).First(&user, "id = ?", userId).Error; err != nil {
			return "", err
		}
		return user.RefreshToken, nil
	}

	return repo.rdb.Get(ctx, userId).Result()
}

func (repo *AuthRepository) UpdateRefreshToken(ctx context.Context, user *model.User, token string) error {
	user.RefreshToken = token
	if err := repo.db.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}
	repo.rdb.Set(ctx, strconv.Itoa(int(user.ID)), token, repo.config.Redis.TTL)

	return nil
}

func (repo *AuthRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	tx := repo.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		repo.log.Error("Failed to begin transaction", "error", tx.Error)
		return nil, tx.Error
	}
	return tx, nil
}

func (repo *AuthRepository) SoftDeleteUserTx(ctx context.Context, tx *gorm.DB, userID string) error {
	result := tx.Model(&model.User{}).
		Where("user_id = ?", userID).
		Update("is_deleted", true)

	if result.Error != nil {
		repo.log.Error("Failed to soft delete user", "user_id", userID, "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Warn("No user found to soft delete", "user_id", userID)
	} else {
		repo.log.Info("User soft deleted successfully", "user_id", userID)
	}

	return nil
}

func (repo *AuthRepository) RestoreUserTx(ctx context.Context, tx *gorm.DB, userID string) error {
	result := tx.Model(&model.User{}).
		Where("user_id = ?", userID).
		Update("is_deleted", false)

	if result.Error != nil {
		repo.log.Error("Failed to restore user", "user_id", userID, "error", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		repo.log.Warn("No user found to restore", "user_id", userID)
	} else {
		repo.log.Info("User restored successfully", "user_id", userID)
	}

	return nil
}
