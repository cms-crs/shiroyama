package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.Board, error) {
	const op = "BoardService.UpdateBoard"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", req.Id),
	)

	if req.Id == "" {
		log.Warn("Board ID is required")
		return nil, status.Error(codes.InvalidArgument, "board id is required")
	}

	_, err := service.repo.GetBoard(ctx, req.Id)
	if err != nil {
		log.Error("Board not found", "board_id", req.Id, "error", err)
		return nil, status.Error(codes.NotFound, "board not found")
	}

	log.Info("Updating board")

	board, err := service.repo.UpdateBoard(ctx, req)
	if err != nil {
		log.Error("Failed to update board", "error", err)
		return nil, status.Error(codes.Internal, "failed to update board")
	}

	log.Info("Board updated successfully")

	return board, nil
}
