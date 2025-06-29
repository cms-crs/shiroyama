package list

import (
	"context"
	"log/slog"

	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
)

type Repository interface {
	CreateList(ctx context.Context, req *boardv1.CreateListRequest) (*boardv1.List, error)
	UpdateList(ctx context.Context, req *boardv1.UpdateListRequest) (*boardv1.List, error)
	GetList(ctx context.Context, listID string) (*boardv1.List, error)
	DeleteList(ctx context.Context, listID string) error
	ReorderLists(ctx context.Context, req *boardv1.ReorderListsRequest) error
}

type BoardRepository interface {
	GetBoard(ctx context.Context, boardID string) (*boardv1.Board, error)
}

type Service struct {
	log       *slog.Logger
	listRepo  Repository
	boardRepo BoardRepository
}

func NewListService(log *slog.Logger, listRepo Repository, boardRepo BoardRepository) *Service {
	return &Service{
		log:       log,
		listRepo:  listRepo,
		boardRepo: boardRepo,
	}
}
