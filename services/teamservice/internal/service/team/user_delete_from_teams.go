package team

import (
	"context"
	"taskservice/internal/kafka"
)

func (service *Service) DeleteUserFromAllTeams(ctx context.Context, userID string) (*kafka.TeamDeletionData, error) {
	const op = "TeamService.DeleteUserFromAllTeams"

	service.log.Info("Deleting user from all teams", "user_id", userID, "op", op)

	deletionData, err := service.teamRepository.DeleteUserFromAllTeams(ctx, userID)
	if err != nil {
		service.log.Error("Failed to delete user from all teams", "user_id", userID, "error", err, "op", op)
		return nil, err
	}

	service.log.Info("Successfully deleted user from all teams", "user_id", userID, "teams_count", len(deletionData.Teams), "op", op)
	return deletionData, nil
}
