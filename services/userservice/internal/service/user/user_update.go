package user

import (
	"context"
	"userservice/internal/dto"
	"userservice/internal/entity"
)

func (service *Service) UpdateUser(ctx context.Context, request *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	user := &entity.User{
		ID:       request.ID,
		Username: request.Username,
	}

	user, err := service.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &dto.UpdateUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
