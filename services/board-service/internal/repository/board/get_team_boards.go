package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (r *Repository) GetTeamBoards(ctx context.Context, teamID string) ([]*boardv1.Board, error) {
	query := `
		SELECT id, name, description, team_id, created_by, created_at, updated_at
		FROM boards
		WHERE team_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, fmt.Errorf("failed to get team boards: %w", err)
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
