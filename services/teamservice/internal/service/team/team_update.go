package team

import (
	"context"
	"userservice/internal/dto"
	"userservice/internal/entity"
)

func (service *Service) UpdateTeam(
	ctx context.Context,
	req *dto.UpdateTeamRequest,
) (*dto.UpdateTeamResponse, error) {
	team := &entity.Team{
		ID:          req.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	team, err := service.teamRepository.UpdateTeam(ctx, team)

	if err != nil {
		return nil, err
	}

	return &dto.UpdateTeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}, nil
}
