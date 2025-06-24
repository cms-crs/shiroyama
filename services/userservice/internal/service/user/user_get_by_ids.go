package user

import (
	"context"
	"userservice/internal/dto"
)

func (service *Service) GetUsersByIds(ctx context.Context, IDs []string) (*dto.GetUsersByIdsResponse, error) {
	const op = "userservice.GetUsersByIds"

	users, err := service.userRepository.GetUsersByIds(ctx, IDs)

	if err != nil {
		return nil, err
	}

	var getUserResponses []*dto.GetUserResponse
	for _, user := range users {
		getUserResponses = append(getUserResponses, &dto.GetUserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return &dto.GetUsersByIdsResponse{Users: getUserResponses}, nil
}
