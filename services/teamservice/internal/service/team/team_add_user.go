package team

import (
	"context"
	"userservice/internal/dto"
	"userservice/internal/entity"
)

func (service *Service) AddUserToTeam(
	ctx context.Context,
	req *dto.AddUserToTeamRequest,
) (*dto.TeamMember, error) {
	teamMember := &entity.TeamMember{
		UserID: req.UserID,
		TeamID: req.TeamID,
		Role:   req.Role,
	}

	err := service.teamRepository.AddUserToTeam(ctx, teamMember)

	if err != nil {
		return nil, err
	}

	return &dto.TeamMember{
		UserID: req.UserID,
		TeamID: req.TeamID,
		Role:   req.Role,
	}, nil
}
