package list

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"time"
)

func (r *Repository) UpdateList(ctx context.Context, req *boardv1.UpdateListRequest) (*boardv1.List, error) {
	now := time.Now()
	query := `
		UPDATE lists 
		SET name = $2, updated_at = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, req.Id, req.Name, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("list not found")
	}

	return r.GetList(ctx, req.Id)
}
