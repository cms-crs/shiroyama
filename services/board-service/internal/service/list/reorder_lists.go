package list

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

func (service *Service) ReorderLists(ctx context.Context, req *boardv1.ReorderListsRequest) error {
	const op = "listService.ReorderLists"

	log := service.log.With(
		slog.String("op", op),
		slog.String("board_id", req.BoardId),
		slog.Int("positions_count", len(req.Positions)),
	)

	log.Info("Reordering lists")

	if req.BoardId == "" {
		log.Warn("Board ID is empty")
		return fmt.Errorf("board ID is required")
	}
	if len(req.Positions) == 0 {
		log.Warn("Positions are empty")
		return fmt.Errorf("positions are required")
	}

	err := service.repo.ReorderLists(ctx, req)
	if err != nil {
		log.Error("Failed to reorder lists", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Lists reordered successfully")

	return nil
}
