package board

import (
	"context"
	boardv1 "github.com/cms-crs/protos/gen/go/board_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (service *Service) GetTeamBoards(ctx context.Context, teamID string) ([]*boardv1.Board, error) {
	const op = "BoardService.GetTeamBoards"

	log := service.log.With(
		slog.String("op", op),
		slog.String("team_id", teamID),
	)

	if teamID == "" {
		log.Warn("Team ID is required")
		return nil, status.Error(codes.InvalidArgument, "team id is required")
	}

	boards, err := service.repo.GetTeamBoards(ctx, teamID)
	if err != nil {
		log.Error("Failed to get team boards", "error", err)
		return nil, status.Error(codes.Internal, "failed to get team boards")
	}

	return boards, nil
}
