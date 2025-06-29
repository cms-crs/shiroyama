package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (s *Service) MoveTask(ctx context.Context, req *taskv1.MoveTaskRequest) (*taskv1.Task, error) {
	const op = "TaskService.MoveTask"

	log := s.log.With(
		slog.String("op", op),
		slog.String("task_id", req.TaskId),
		slog.String("to_list_id", req.ToListId),
	)

	if req.TaskId == "" {
		log.Warn("Task ID is required")
		return nil, status.Error(codes.InvalidArgument, "task_id is required")
	}

	if req.ToListId == "" {
		log.Warn("Target list ID is required")
		return nil, status.Error(codes.InvalidArgument, "to_list_id is required")
	}

	_, err := s.taskRepo.GetTask(ctx, req.TaskId)
	if err != nil {
		log.Error("Task not found", "task_id", req.TaskId, "error", err)
		return nil, status.Error(codes.NotFound, "task not found")
	}

	if err := s.validateListExists(ctx, req.ToListId); err != nil {
		log.Error("Target list validation failed", "list_id", req.ToListId, "error", err)
		return nil, status.Error(codes.InvalidArgument, "target list not found")
	}

	log.Info("Moving task")

	task, err := s.taskRepo.MoveTask(ctx, req)
	if err != nil {
		log.Error("Failed to move task", "error", err)
		return nil, status.Error(codes.Internal, "failed to move task")
	}

	log.Info("Task moved successfully")

	return task, nil
}
