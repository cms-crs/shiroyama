package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) GetUserBoards(ctx context.Context, userID string) ([]*boardv1.Board, error) {
	const op = "BoardService.GetUserBoards"

	log := service.log.With(
		slog.String("op", op),
		slog.String("user_id", userID),
	)

	if userID == "" {
		log.Warn("User ID is required")
		return nil, status.Error(codes.InvalidArgument, "user id is required")
	}

	boards, err := service.repo.GetUserBoards(ctx, userID)
	if err != nil {
		log.Error("Failed to get user boards", "error", err)
		return nil, status.Error(codes.Internal, "failed to get user boards")
	}

	return boards, nil
}
