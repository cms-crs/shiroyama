package team

import (
	"context"
	"fmt"
	"taskservice/internal/dto"
	"taskservice/internal/entity"
)

func (service *Service) AddUserToTeam(
	ctx context.Context,
	req *dto.AddUserToTeamRequest,
) (*dto.TeamMember, error) {
	_, err := service.userClient.GetUser(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	teamMember := &entity.TeamMember{
		UserID: req.UserID,
		TeamID: req.TeamID,
		Role:   req.Role,
	}

	err = service.teamRepository.AddUserToTeam(ctx, teamMember)

	if err != nil {
		return nil, err
	}

	return &dto.TeamMember{
		UserID: req.UserID,
		TeamID: req.TeamID,
		Role:   req.Role,
	}, nil
}
