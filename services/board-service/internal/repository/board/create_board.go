package board

import (
	"context"
	"fmt"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (r *Repository) CreateBoard(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.Board, error) {
	id := uuid.New().String()
	now := time.Now()

	query := `
		INSERT INTO boards (id, name, description, team_id, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query, id, req.Name, req.Description, req.TeamId, req.CreatedBy, now, now)
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	return &boardv1.Board{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
		TeamId:      req.TeamId,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   timestamppb.New(now),
		UpdatedAt:   timestamppb.New(now),
	}, nil
}
