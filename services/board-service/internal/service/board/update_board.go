package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.Board, error) {
	const op = "BoardService.UpdateBoard"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", req.Id),
		slog.String("name", req.Name),
	)

	log.Info("Updating board")

	if req.Id == "" {
		log.Warn("Board ID is empty")
		return nil, fmt.Errorf("board ID is required")
	}
	if req.Name == "" {
		log.Warn("Board name is empty")
		return nil, fmt.Errorf("board name is required")
	}

	board, err := service.repo.UpdateBoard(ctx, req)
	if err != nil {
		log.Error("Failed to update board", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Board updated successfully")

	return board, nil
}
