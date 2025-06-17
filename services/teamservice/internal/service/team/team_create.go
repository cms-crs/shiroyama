package team

import (
	"context"
	"userservice/internal/dto"
	"userservice/internal/entity"
)

func (service *Service) CreateTeam(
	ctx context.Context,
	req *dto.CreateTeamRequest,
) (*dto.CreateTeamResponse, error) {
	// todo: implement check for createdBy

	team := &entity.Team{
		Name:        req.Name,
		Description: req.Description,
	}

	team, err := service.teamRepository.CreateTeam(ctx, team)

	if err != nil {
		return nil, err
	}

	return &dto.CreateTeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}, nil
}
