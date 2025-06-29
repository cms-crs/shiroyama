package list

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) GetList(ctx context.Context, listID string) (*boardv1.List, error) {
	const op = "ListService.GetList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("list_id", listID),
	)

	if listID == "" {
		log.Warn("List ID is required")
		return nil, status.Error(codes.InvalidArgument, "list id is required")
	}

	list, err := service.listRepo.GetList(ctx, listID)
	if err != nil {
		log.Error("Failed to get list", "error", err)
		return nil, status.Error(codes.NotFound, "list not found")
	}

	return list, nil
}
