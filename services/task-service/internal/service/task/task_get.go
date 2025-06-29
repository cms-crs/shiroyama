package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (s *Service) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.Task, error) {
	const op = "TaskService.GetTask"

	log := s.log.With(
		slog.String("op", op),
		slog.String("task_id", req.Id),
	)

	if req.Id == "" {
		log.Warn("Task ID is required")
		return nil, status.Error(codes.InvalidArgument, "task id is required")
	}

	task, err := s.taskRepo.GetTask(ctx, req.Id)
	if err != nil {
		log.Error("Failed to get task", "error", err)
		return nil, status.Error(codes.NotFound, "task not found")
	}

	return task, nil
}
