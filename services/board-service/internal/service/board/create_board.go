package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	teamv1 "github.com/cms-crs/protos/gen/go/team_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) CreateBoard(ctx context.Context, req *boardv1.CreateBoardRequest) (*boardv1.Board, error) {
	const op = "BoardService.CreateBoard"

	log := service.log.With(
		slog.String("op", op),
		slog.String("name", req.Name),
		slog.String("created_by", req.CreatedBy),
		slog.String("team_id", req.TeamId),
	)

	if req.Name == "" {
		log.Warn("Board name is required")
		return nil, status.Error(codes.InvalidArgument, "board name is required")
	}

	if req.CreatedBy == "" {
		log.Warn("Created by user ID is required")
		return nil, status.Error(codes.InvalidArgument, "created_by is required")
	}

	_, err := service.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.CreatedBy})
	if err != nil {
		log.Error("User not found", "user_id", req.CreatedBy, "error", err)
		return nil, status.Error(codes.InvalidArgument, "user not found")
	}

	if req.TeamId != "" {
		_, err := service.teamClient.GetTeam(ctx, &teamv1.GetTeamRequest{Id: req.TeamId})
		if err != nil {
			log.Error("Team not found", "team_id", req.TeamId, "error", err)
			return nil, status.Error(codes.InvalidArgument, "team not found")
		}
	}

	log.Info("Creating board")

	board, err := service.repo.CreateBoard(ctx, req)
	if err != nil {
		log.Error("Failed to create board", "error", err)
		return nil, status.Error(codes.Internal, "failed to create board")
	}

	log.Info("Board created successfully", "board_id", board.Id)

	return board, nil
}
