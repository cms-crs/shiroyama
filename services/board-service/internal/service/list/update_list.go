package list

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) UpdateList(ctx context.Context, req *boardv1.UpdateListRequest) (*boardv1.List, error) {
	const op = "listService.UpdateList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("list_id", req.Id),
		slog.String("name", req.Name),
	)

	log.Info("Updating list")

	if req.Id == "" {
		log.Warn("List ID is empty")
		return nil, fmt.Errorf("list ID is required")
	}
	if req.Name == "" {
		log.Warn("List name is empty")
		return nil, fmt.Errorf("list name is required")
	}

	list, err := service.repo.UpdateList(ctx, req)
	if err != nil {
		log.Error("Failed to update list", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("List updated successfully")

	return list, nil
}
