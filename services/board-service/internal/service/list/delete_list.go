package list

import (
	"context"
	"fmt"
	"log/slog"
)

func (service *Service) DeleteList(ctx context.Context, listID string) error {
	const op = "listService.DeleteList"

	log := service.log.With(
		slog.String("op", op),
		slog.String("list_id", listID),
	)

	log.Info("Deleting list")

	if listID == "" {
		log.Warn("List ID is empty")
		return fmt.Errorf("list ID is required")
	}

	err := service.repo.DeleteList(ctx, listID)
	if err != nil {
		log.Error("Failed to delete list", "error", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("List deleted successfully")

	return nil
}
