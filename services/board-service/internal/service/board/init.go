package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

type Repository interface {
	CreateBoard(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.Board, error)
	GetBoard(ctx context.Context, boardID string) (*boardv1.Board, error)
	GetBoardWithLists(ctx context.Context, boardID string) (*boardv1.BoardWithLists, error)
	UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.Board, error)
	GetUserBoards(ctx context.Context, userID string) ([]*boardv1.Board, error)
	GetTeamBoards(ctx context.Context, teamID string) ([]*boardv1.Board, error)
}

type Service struct {
	log  *slog.Logger
	repo Repository
}

func NewBoardService(log *slog.Logger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}
