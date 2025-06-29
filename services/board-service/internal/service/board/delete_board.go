package board

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) DeleteBoard(ctx context.Context, boardID string) error {
	const op = "BoardService.DeleteBoard"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", boardID),
	)

	if boardID == "" {
		log.Warn("Board ID is required")
		return status.Error(codes.InvalidArgument, "board id is required")
	}

	_, err := service.repo.GetBoard(ctx, boardID)
	if err != nil {
		log.Error("Board not found", "board_id", boardID, "error", err)
		return status.Error(codes.NotFound, "board not found")
	}

	log.Info("Deleting board")

	err = service.repo.DeleteBoard(ctx, boardID)
	if err != nil {
		log.Error("Failed to delete board", "error", err)
		return status.Error(codes.Internal, "failed to delete board")
	}

	log.Info("Board deleted successfully")

	return nil
}
