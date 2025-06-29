package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"time"
)

func (r *Repository) UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.Board, error) {
	now := time.Now()
	query := `
		UPDATE boards 
		SET name = $2, description = $3, updated_at = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, req.Id, req.Name, req.Description, now)
	if err != nil {
		return nil, fmt.Errorf("failed to update board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("board not found")
	}

	return r.GetBoard(ctx, req.Id)
}
