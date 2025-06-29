package user

import (
	"context"
	"userservice/internal/dto"
	"userservice/internal/entity"
)

func (service *Service) CreateUser(ctx context.Context, request *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	user := &entity.User{
		Email:    request.Email,
		Username: request.Username,
	}

	user, err := service.userRepository.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.CreateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
