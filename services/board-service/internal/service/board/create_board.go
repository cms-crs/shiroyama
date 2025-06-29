package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) CreateBoard(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.Board, error) {
	const op = "boardService.CreateBoard"

	log := service.log.With(
		slog.String("op", op),
		slog.String("name", req.Name),
		slog.String("created_by", req.CreatedBy),
	)

	log.Info("Creating board")

	if req.Name == "" {
		log.Warn("Board name is empty")
		return nil, fmt.Errorf("board name is required")
	}
	if req.CreatedBy == "" {
		log.Warn("Created by is empty")
		return nil, fmt.Errorf("created_by is required")
	}

	board, err := service.repo.CreateBoard(ctx, req)
	if err != nil {
		log.Error("Failed to create board", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Board created successfully", "board_id", board.Id)

	return board, nil
}
