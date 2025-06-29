package list

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) GetList(ctx context.Context, listID string) (*boardv1.List, error) {
	const op = "listService.GetList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("list_id", listID),
	)

	log.Info("Getting list")

	if listID == "" {
		log.Warn("List ID is empty")
		return nil, fmt.Errorf("list ID is required")
	}

	list, err := service.repo.GetList(ctx, listID)
	if err != nil {
		log.Error("Failed to get list", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("List retrieved successfully")

	return list, nil
}
