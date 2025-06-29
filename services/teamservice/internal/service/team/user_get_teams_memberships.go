package team

import (
	"context"
	"taskservice/internal/kafka"
)

func (service *Service) GetUserTeamMemberships(ctx context.Context, userID string) (*kafka.TeamDeletionData, error) {
	const op = "TeamService.GetUserTeamMemberships"

	service.log.Debug("Getting user team memberships", "user_id", userID, "op", op)

	deletionData, err := service.teamRepository.GetUserTeamMemberships(ctx, userID)
	if err != nil {
		service.log.Error("Failed to get user team memberships", "user_id", userID, "error", err, "op", op)
		return nil, err
	}

	return deletionData, nil
}
