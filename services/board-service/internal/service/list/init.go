package list

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"log/slog"
)

type Repository interface {
	CreateList(ctx context.Context, req *boardv1.CreateListRequest) (*boardv1.List, error)
	UpdateList(ctx context.Context, req *boardv1.UpdateListRequest) (*boardv1.List, error)
	GetList(ctx context.Context, listID string) (*boardv1.List, error)
	DeleteList(ctx context.Context, listID string) error
	ReorderLists(ctx context.Context, req *boardv1.ReorderListsRequest) error
}

type Service struct {
	log  *slog.Logger
	repo Repository
}

func NewListService(log *slog.Logger, repo Repository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}
