package team

import (
	"context"
	"taskservice/internal/dto"
)

func (service *Service) RemoveUserFromTeam(
	ctx context.Context,
	req *dto.RemoveUserFromTeamRequest,
) error {
	err := service.teamRepository.RemoveUserFromTeam(ctx, req)

	if err != nil {
		return err
	}

	return nil
}
