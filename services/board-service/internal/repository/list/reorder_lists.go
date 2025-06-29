package list

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
)

func (r *Repository) ReorderLists(ctx context.Context, req *boardv1.ReorderListsRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	for _, pos := range req.Positions {
		_, err = tx.ExecContext(ctx,
			"UPDATE lists SET position = $2 WHERE id = $1 AND board_id = $3",
			pos.ListId, pos.Position, req.BoardId)
		if err != nil {
			return fmt.Errorf("failed to update list position: %w", err)
		}
	}

	return tx.Commit()
}
