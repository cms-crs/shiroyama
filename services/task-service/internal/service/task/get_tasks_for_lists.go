package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (s *Service) GetTasksForLists(ctx context.Context, req *taskv1.GetTasksForListsRequest) (*taskv1.GetTasksForListsResponse, error) {
	const op = "TaskService.GetTasksForLists"

	log := s.log.With(
		slog.String("op", op),
		slog.Int("list_count", len(req.ListIds)),
	)

	if len(req.ListIds) == 0 {
		log.Warn("List IDs are required")
		return nil, status.Error(codes.InvalidArgument, "list_ids are required")
	}

	tasks, err := s.taskRepo.GetTasksForLists(ctx, req.ListIds)
	if err != nil {
		log.Error("Failed to get tasks for lists", "error", err)
		return nil, status.Error(codes.Internal, "failed to get tasks")
	}

	return &taskv1.GetTasksForListsResponse{
		Tasks: tasks,
	}, nil
}
