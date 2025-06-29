package list

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) CreateList(ctx context.Context, req *boardv1.CreateListRequest) (*boardv1.List, error) {
	const op = "ListService.CreateList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", req.BoardId),
		slog.String("name", req.Name),
	)

	if req.BoardId == "" {
		log.Warn("Board ID is required")
		return nil, status.Error(codes.InvalidArgument, "board_id is required")
	}

	if req.Name == "" {
		log.Warn("List name is required")
		return nil, status.Error(codes.InvalidArgument, "list name is required")
	}

	_, err := service.boardRepo.GetBoard(ctx, req.BoardId)
	if err != nil {
		log.Error("Board not found", "board_id", req.BoardId, "error", err)
		return nil, status.Error(codes.InvalidArgument, "board not found")
	}

	log.Info("Creating list")

	list, err := service.listRepo.CreateList(ctx, req)
	if err != nil {
		log.Error("Failed to create list", "error", err)
		return nil, status.Error(codes.Internal, "failed to create list")
	}

	log.Info("List created successfully", "list_id", list.Id)

	return list, nil
}
