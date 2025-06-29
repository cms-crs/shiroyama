package list

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) ReorderLists(ctx context.Context, req *boardv1.ReorderListsRequest) error {
	const op = "ListService.ReorderLists"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", req.BoardId),
	)

	if req.BoardId == "" {
		log.Warn("Board ID is required")
		return status.Error(codes.InvalidArgument, "board_id is required")
	}

	if len(req.Positions) == 0 {
		log.Warn("List positions are required")
		return status.Error(codes.InvalidArgument, "list positions are required")
	}

	_, err := service.boardRepo.GetBoard(ctx, req.BoardId)
	if err != nil {
		log.Error("Board not found", "board_id", req.BoardId, "error", err)
		return status.Error(codes.InvalidArgument, "board not found")
	}

	log.Info("Reordering lists")

	err = service.listRepo.ReorderLists(ctx, req)
	if err != nil {
		log.Error("Failed to reorder lists", "error", err)
		return status.Error(codes.Internal, "failed to reorder lists")
	}

	log.Info("Lists reordered successfully")

	return nil
}
