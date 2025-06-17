package team

import (
	"context"
	"userservice/internal/dto"
)

func (service *Service) UpdateUserRole(
	ctx context.Context,
	req *dto.UpdateUserRoleRequest,
) (*dto.TeamMember, error) {
	err := service.teamRepository.UpdateUserRole(ctx, req)

	if err != nil {
		return nil, err
	}

	return &dto.TeamMember{
		UserID: req.UserID,
		TeamID: req.TeamID,
		Role:   req.Role,
	}, nil
}
