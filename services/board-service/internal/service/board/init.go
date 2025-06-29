package board

import (
	"context"
	"log/slog"

	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	teamv1 "github.com/cms-crs/protos/gen/go/team_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
)

type Repository interface {
	CreateBoard(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.Board, error)
	GetBoard(ctx context.Context, boardID string) (*boardv1.Board, error)
	GetBoardWithLists(ctx context.Context, boardID string) (*boardv1.BoardWithLists, error)
	UpdateBoard(ctx context.Context, req *boardv1.UpdateBoardRequest) (*boardv1.Board, error)
	DeleteBoard(ctx context.Context, boardID string) error
	GetUserBoards(ctx context.Context, userID string) ([]*boardv1.Board, error)
	GetTeamBoards(ctx context.Context, teamID string) ([]*boardv1.Board, error)
}

type UserServiceClient interface {
	GetUser(ctx context.Context, req *userv1.GetUserRequest) (*userv1.User, error)
}

type TeamServiceClient interface {
	GetTeam(ctx context.Context, req *teamv1.GetTeamRequest) (*teamv1.Team, error)
}

type Service struct {
	log        *slog.Logger
	repo       Repository
	userClient userv1.UserServiceClient
	teamClient teamv1.TeamServiceClient
}

func NewBoardService(log *slog.Logger, repo Repository, userClient userv1.UserServiceClient, teamClient teamv1.TeamServiceClient) *Service {
	return &Service{
		log:        log,
		repo:       repo,
		userClient: userClient,
		teamClient: teamClient,
	}
}
