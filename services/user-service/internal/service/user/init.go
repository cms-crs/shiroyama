package user

import (
	"context"
	"log/slog"
	"userservice/internal/entity"
)

type Repository interface {
	CreateUser(
		ctx context.Context,
		req *entity.User,
	) (*entity.User, error)
	GetUser(
		ctx context.Context,
		ID string,
	) (*entity.User, error)
	GetUsersByIds(
		ctx context.Context,
		IDs []string,
	) ([]*entity.User, error)
	UpdateUser(
		ctx context.Context,
		req *entity.User,
	) (*entity.User, error)
}

type Service struct {
	log            *slog.Logger
	userRepository Repository
}

func NewUserService(log *slog.Logger, userRepository Repository) *Service {
	return &Service{
		log:            log,
		userRepository: userRepository,
	}
}
