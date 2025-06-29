package list

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) CreateList(ctx context.Context, req *boardv1.CreateListRequest) (*boardv1.List, error) {
	const op = "listService.CreateList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", req.BoardId),
		slog.String("name", req.Name),
	)

	log.Info("Creating list")

	if req.BoardId == "" {
		log.Warn("Board ID is empty")
		return nil, fmt.Errorf("board ID is required")
	}
	if req.Name == "" {
		log.Warn("List name is empty")
		return nil, fmt.Errorf("list name is required")
	}

	list, err := service.repo.CreateList(ctx, req)
	if err != nil {
		log.Error("Failed to create list", "error", err)
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("List created successfully", "list_id", list.Id)

	return list, nil
}
