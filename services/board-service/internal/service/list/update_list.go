package list

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) UpdateList(ctx context.Context, req *boardv1.UpdateListRequest) (*boardv1.List, error) {
	const op = "ListService.UpdateList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("list_id", req.Id),
	)

	if req.Id == "" {
		log.Warn("List ID is required")
		return nil, status.Error(codes.InvalidArgument, "list id is required")
	}

	_, err := service.listRepo.GetList(ctx, req.Id)
	if err != nil {
		log.Error("List not found", "list_id", req.Id, "error", err)
		return nil, status.Error(codes.NotFound, "list not found")
	}

	log.Info("Updating list")

	list, err := service.listRepo.UpdateList(ctx, req)
	if err != nil {
		log.Error("Failed to update list", "error", err)
		return nil, status.Error(codes.Internal, "failed to update list")
	}

	log.Info("List updated successfully")

	return list, nil
}
