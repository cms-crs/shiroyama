package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

func (s *Service) DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*emptypb.Empty, error) {
	const op = "TaskService.DeleteTask"

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

	log.Info("Deleting task")

	err = s.taskRepo.DeleteTask(ctx, req.Id)
	if err != nil {
		log.Error("Failed to delete task", "error", err)
		return nil, status.Error(codes.Internal, "failed to delete task")
	}

	log.Info("Task deleted successfully")

	return &emptypb.Empty{}, nil
}
