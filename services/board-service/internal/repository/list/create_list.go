package list

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (r *Repository) CreateList(ctx context.Context, req *boardv1.CreateListRequest) (*boardv1.List, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO lists (id, board_id, name, position, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query, id, req.BoardId, req.Name, req.Position, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	return &boardv1.List{
		Id:        id,
		BoardId:   req.BoardId,
		Name:      req.Name,
		Position:  req.Position,
		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}, nil
}
