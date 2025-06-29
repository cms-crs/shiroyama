package team

import (
	"context"
	"taskservice/internal/dto"
)

func (service *Service) GetUserTeams(
	ctx context.Context,
	UserID string,
) ([]*dto.GetTeamResponse, error) {
	userTeams, err := service.teamRepository.GetUserTeams(ctx, UserID)
	if err != nil {
		return nil, err
	}

	var getUserTeamResponses []*dto.GetTeamResponse
	for _, userTeam := range userTeams {
		getUserTeamResponses = append(getUserTeamResponses, &dto.GetTeamResponse{
			ID:          userTeam.ID,
			Name:        userTeam.Name,
			Description: userTeam.Description,
			CreatedAt:   userTeam.CreatedAt,
			UpdatedAt:   userTeam.UpdatedAt,
		})
	}

	return getUserTeamResponses, nil
}
