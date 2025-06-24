package team

import (
	"context"
	"fmt"
	"userservice/internal/dto"
	"userservice/internal/entity"
)

func (service *Service) CreateTeam(
	ctx context.Context,
	req *dto.CreateTeamRequest,
) (*dto.CreateTeamResponse, error) {
	_, err := service.userClient.GetUser(ctx, req.CreatedBy)
	if err != nil {
		return nil, fmt.Errorf("creator not found: %w", err)
	}

	team := &entity.Team{
		Name:        req.Name,
		Description: req.Description,
	}

	team, err = service.teamRepository.CreateTeam(ctx, team)

	if err != nil {
		return nil, err
	}

	teamMember := &entity.TeamMember{
		TeamID: team.ID,
		UserID: req.CreatedBy,
		Role:   "admin",
	}

	err = service.teamRepository.AddUserToTeam(ctx, teamMember)
	if err != nil {
		return nil, fmt.Errorf("failed to add creator to team: %w", err)
	}

	return &dto.CreateTeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}, nil
}
