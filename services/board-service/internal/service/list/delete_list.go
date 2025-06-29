package list

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) DeleteList(ctx context.Context, listID string) error {
	const op = "ListService.DeleteList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("list_id", listID),
	)

	if listID == "" {
		log.Warn("List ID is required")
		return status.Error(codes.InvalidArgument, "list id is required")
	}

	_, err := service.listRepo.GetList(ctx, listID)
	if err != nil {
		log.Error("List not found", "list_id", listID, "error", err)
		return status.Error(codes.NotFound, "list not found")
	}

	log.Info("Deleting list")

	err = service.listRepo.DeleteList(ctx, listID)
	if err != nil {
		log.Error("Failed to delete list", "error", err)
		return status.Error(codes.Internal, "failed to delete list")
	}

	log.Info("List deleted successfully")

	return nil
}
