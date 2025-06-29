package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (s *Service) GetTasksForUser(ctx context.Context, req *taskv1.GetTasksForUserRequest) (*taskv1.GetTasksForUserResponse, error) {
	const op = "TaskService.GetTasksForUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserId),
	)

	if req.UserId == "" {
		log.Warn("User ID is required")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	tasks, err := s.taskRepo.GetTasksForUser(ctx, req.UserId)
	if err != nil {
		log.Error("Failed to get tasks for user", "error", err)
		return nil, status.Error(codes.Internal, "failed to get tasks")
	}

	return &taskv1.GetTasksForUserResponse{
		Tasks: tasks,
	}, nil
}
