package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (r *Repository) GetUserBoards(ctx context.Context, userID string) ([]*boardv1.Board, error) {
	query := `
		SELECT DISTINCT b.id, b.name, b.description, b.team_id, b.created_by, b.created_at, b.updated_at
		FROM boards b
		JOIN team_members tm ON b.team_id = tm.team_id
		WHERE tm.user_id = $1
		ORDER BY b.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user boards: %w", err)
	}
	defer rows.Close()

	var boards []*boardv1.Board
	for rows.Next() {
		var board boardv1.Board
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&board.Id,
			&board.Name,
			&board.Description,
			&board.TeamId,
			&board.CreatedBy,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan board: %w", err)
		}

		board.CreatedAt = timestamppb.New(createdAt)
		board.UpdatedAt = timestamppb.New(updatedAt)
		boards = append(boards, &board)
	}

	return boards, nil
}
