package board

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (r *Repository) GetBoard(ctx context.Context, boardID string) (*boardv1.Board, error) {
	query := `
		SELECT id, name, description, team_id, created_by, created_at, updated_at
		FROM boards
		WHERE id = $1
	`

	var board boardv1.Board
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, boardID).Scan(
		&board.Id,
		&board.Name,
		&board.Description,
		&board.TeamId,
		&board.CreatedBy,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("board not found")
		}
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	board.CreatedAt = timestamppb.New(createdAt)
	board.UpdatedAt = timestamppb.New(updatedAt)

	return &board, nil
}
