package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) GetBoard(ctx context.Context, boardID string) (*boardv1.Board, error) {
	const op = "boardService.GetBoard"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", boardID),
	)

	log.Info("Getting board")

	if boardID == "" {
		log.Warn("Board ID is empty")
		return nil, fmt.Errorf("board ID is required")
	}

	board, err := service.repo.GetBoard(ctx, boardID)
	if err != nil {
		log.Error("Failed to get board", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Board retrieved successfully")

	return board, nil
}
