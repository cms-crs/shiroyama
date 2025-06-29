package service

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"time"
)

func (s *Service) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.Task, error) {
	const op = "TaskService.CreateTask"

	log := s.log.With(
		slog.String("op", op),
		slog.String("list_id", req.ListId),
		slog.String("title", req.Title),
		slog.String("created_by", req.CreatedBy),
	)

	if req.ListId == "" {
		log.Warn("List ID is required")
		return nil, status.Error(codes.InvalidArgument, "list_id is required")
	}

	if req.Title == "" {
		log.Warn("Task title is required")
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}

	if req.CreatedBy == "" {
		log.Warn("Created by user ID is required")
		return nil, status.Error(codes.InvalidArgument, "created_by is required")
	}

	_, err := s.userClient.GetUser(ctx, &userv1.GetUserRequest{Id: req.CreatedBy})
	if err != nil {
		log.Error("User not found", "user_id", req.CreatedBy, "error", err)
		return nil, status.Error(codes.InvalidArgument, "user not found")
	}

	if err := s.validateListExists(ctx, req.ListId); err != nil {
		log.Error("List validation failed", "list_id", req.ListId, "error", err)
		return nil, status.Error(codes.InvalidArgument, "list not found")
	}

	now := timestamppb.New(time.Now())
	task := &taskv1.Task{
		ListId:      req.ListId,
		Title:       req.Title,
		Description: req.Description,
		Position:    req.Position,
		DueDate:     req.DueDate,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	log.Info("Creating task")

	createdTask, err := s.taskRepo.CreateTask(ctx, task)
	if err != nil {
		log.Error("Failed to create task", "error", err)
		return nil, status.Error(codes.Internal, "failed to create task")
	}

	log.Info("Task created successfully", "task_id", createdTask.Id)

	return createdTask, nil
}
