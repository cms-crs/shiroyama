package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

func (s *Service) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.Task, error) {
	const op = "TaskService.UpdateTask"

	log := s.log.With(
		slog.String("op", op),
		slog.String("task_id", req.Id),
	)

	if req.Id == "" {
		log.Warn("Task ID is required")
		return nil, status.Error(codes.InvalidArgument, "task id is required")
	}

	_, err := s.taskRepo.GetTask(ctx, req.Id)
	if err != nil {
		log.Error("Task not found", "task_id", req.Id, "error", err)
		return nil, status.Error(codes.NotFound, "task not found")
	}

	log.Info("Updating task")

	task, err := s.taskRepo.UpdateTask(ctx, req)
	if err != nil {
		log.Error("Failed to update task", "error", err)
		return nil, status.Error(codes.Internal, "failed to update task")
	}

	log.Info("Task updated successfully")

	return task, nil
}
