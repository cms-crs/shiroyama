package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) GetBoardWithLists(ctx context.Context, boardID string) (*boardv1.BoardWithLists, error) {
	const op = "BoardService.GetBoardWithLists"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", boardID),
	)

	if boardID == "" {
		log.Warn("Board ID is required")
		return nil, status.Error(codes.InvalidArgument, "board id is required")
	}

	board, err := service.repo.GetBoardWithLists(ctx, boardID)
	if err != nil {
		log.Error("Failed to get board with lists", "error", err)
		return nil, status.Error(codes.NotFound, "board not found")
	}

	return board, nil
}
