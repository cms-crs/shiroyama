package user

import (
	"context"
	"userservice/internal/dto"
)

func (service *Service) GetUser(ctx context.Context, request *dto.GetUserRequest) (*dto.GetUserResponse, error) {
	user, err := service.userRepository.GetUser(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &dto.GetUserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
