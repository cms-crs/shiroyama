package team

import (
	"context"
	"taskservice/internal/kafka"
)

func (service *Service) RestoreUserTeams(ctx context.Context, userID string, data *kafka.TeamDeletionData) error {
	const op = "TeamService.RestoreUserTeams"

	service.log.Info("Restoring user teams", "user_id", userID, "teams_count", len(data.Teams), "op", op)

	if err := service.teamRepository.RestoreUserTeams(ctx, userID, data); err != nil {
		service.log.Error("Failed to restore user teams", "user_id", userID, "error", err, "op", op)
		return err
	}

	service.log.Info("Successfully restored user teams", "user_id", userID, "teams_count", len(data.Teams), "op", op)
	return nil
}
