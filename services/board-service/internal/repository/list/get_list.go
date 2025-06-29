package list

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (r *Repository) GetList(ctx context.Context, listID string) (*boardv1.List, error) {
	query := `
		SELECT id, board_id, name, position, created_at, updated_at
		FROM lists
		WHERE id = $1
	`

	var list boardv1.List
	var createdAt, updatedAt time.Time

	err := r.db.QueryRowContext(ctx, query, listID).Scan(
		&list.Id,
		&list.BoardId,
		&list.Name,
		&list.Position,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("list not found")
		}
		return nil, fmt.Errorf("failed to get list: %w", err)
	}

	list.CreatedAt = timestamppb.New(createdAt)
	list.UpdatedAt = timestamppb.New(updatedAt)

	return &list, nil
}
