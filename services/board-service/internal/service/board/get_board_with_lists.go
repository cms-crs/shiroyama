package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) GetBoardWithLists(ctx context.Context, boardID string) (*boardv1.BoardWithLists, error) {
	const op = "boardService.GetBoardWithLists"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", boardID),
	)

	log.Info("Getting board with lists")

	if boardID == "" {
		log.Warn("Board ID is empty")
		return nil, fmt.Errorf("board ID is required")
	}

	board, err := service.repo.GetBoardWithLists(ctx, boardID)
	if err != nil {
		log.Error("Failed to get board with lists", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Board with lists retrieved successfully")

	return board, nil
}
