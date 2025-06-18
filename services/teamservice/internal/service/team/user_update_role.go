package team

import (
	"context"
	"fmt"
	"userservice/internal/dto"
)

func (service *Service) UpdateUserRole(
	ctx context.Context,
	req *dto.UpdateUserRoleRequest,
) (*dto.TeamMember, error) {
	_, err := service.userClient.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	err = service.teamRepository.UpdateUserRole(ctx, req)

	if err != nil {
		return nil, err
	}

	return &dto.TeamMember{
		UserID: req.UserID,
		TeamID: req.TeamID,
		Role:   req.Role,
	}, nil
}
