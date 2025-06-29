package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) GetUserBoards(ctx context.Context, userID string) ([]*boardv1.Board, error) {
	const op = "boardService.GetUserBoards"

	log := service.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	log.Info("Getting user boards")

	if userID == "" {
		log.Warn("User ID is empty")
		return nil, fmt.Errorf("user ID is required")
	}

	boards, err := service.repo.GetUserBoards(ctx, userID)
	if err != nil {
		log.Error("Failed to get user boards", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User boards retrieved successfully", "count", len(boards))

	return boards, nil
}
