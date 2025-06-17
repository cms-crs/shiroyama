package team

import (
	"context"
	"userservice/internal/dto"
)

func (service *Service) GetTeam(
	ctx context.Context,
	ID string,
) (*dto.GetTeamResponse, error) {
	team, err := service.teamRepository.GetTeam(ctx, ID)

	if err != nil {
		return nil, err
	}

	return &dto.GetTeamResponse{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}, nil
}
