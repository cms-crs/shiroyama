package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) GetTeamBoards(ctx context.Context, teamID string) ([]*boardv1.Board, error) {
	const op = "boardService.GetTeamBoards"

	log := service.log.With(
		slog.String("op", op),
		slog.String("team_id", teamID),
	)

	log.Info("Getting team boards")

	if teamID == "" {
		log.Warn("Team ID is empty")
		return nil, fmt.Errorf("team ID is required")
	}

	boards, err := service.repo.GetTeamBoards(ctx, teamID)
	if err != nil {
		log.Error("Failed to get team boards", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Team boards retrieved successfully", "count", len(boards))

	return boards, nil
}
