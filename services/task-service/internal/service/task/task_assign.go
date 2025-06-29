package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (s *Service) AssignUser(ctx context.Context, req *taskv1.AssignUserRequest) (*taskv1.Task, error) {
	const op = "TaskService.AssignUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("task_id", req.TaskId),
		slog.String("user_id", req.UserId),
	)

	if req.TaskId == "" {
		log.Warn("Task ID is required")
		return nil, status.Error(codes.InvalidArgument, "task_id is required")
	}

	if req.UserId == "" {
		log.Warn("User ID is required")
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	_, err := s.taskRepo.GetTask(ctx, req.TaskId)
	if err != nil {
		log.Error("Task not found", "task_id", req.TaskId, "error", err)
		return nil, status.Error(codes.NotFound, "task not found")
	}

	_, err = s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.UserId})
	if err != nil {
		log.Error("User not found", "user_id", req.UserId, "error", err)
		return nil, status.Error(codes.InvalidArgument, "user not found")
	}

	log.Info("Assigning user to task")

	task, err := s.taskRepo.AssignUser(ctx, req.TaskId, req.UserId)
	if err != nil {
		log.Error("Failed to assign user", "error", err)
		return nil, status.Error(codes.Internal, "failed to assign user")
	}

	log.Info("User assigned successfully")

	return task, nil
}
