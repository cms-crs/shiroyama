package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) GetBoard(ctx context.Context, boardID string) (*boardv1.Board, error) {
	const op = "BoardService.GetBoard"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", boardID),
	)

	if boardID == "" {
		log.Warn("Board ID is required")
		return nil, status.Error(codes.InvalidArgument, "board id is required")
	}

	board, err := service.repo.GetBoard(ctx, boardID)
	if err != nil {
		log.Error("Failed to get board", "error", err)
		return nil, status.Error(codes.NotFound, "board not found")
	}

	return board, nil
}
