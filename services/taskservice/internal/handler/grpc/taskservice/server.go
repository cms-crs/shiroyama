package taskservicegrpc

import (
	"context"
	"errors"
	"github.com/ShiroyamaY/protos/gen/go/taskservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"log/slog"
	"taskservice/internal/dto"
)

type TaskService interface {
	Create(ctx context.Context,
		request dto.CreateTaskRequest,
	) (*dto.CreateTaskResponse, error)
	GetTask(ctx context.Context,
		request *dto.GetTaskRequest,
	) (*dto.GetTaskResponse, error)
	UpdateTask(ctx context.Context,
		request *dto.UpdateTaskRequest,
	) (*dto.UpdateTaskResponse, error)
	DeleteTask(ctx context.Context,
		request *dto.DeleteTaskRequest,
	) (*dto.DeleteTaskResponse, error)
}

type serverAPI struct {
	taskservice.UnimplementedTaskServiceServer
	log         *slog.Logger
	taskService TaskService
}

func Register(gRPC *grpc.Server, taskService TaskService, log *slog.Logger) {
	taskservice.RegisterTaskServiceServer(gRPC, serverAPI{
		taskService: taskService,
		log:         log,
	})
}

func (server serverAPI) CreateTask(ctx context.Context, req *taskservice.CreateTaskRequest) (*taskservice.Task, error) {
	const op = "serverAPI.CreateTask"

	log := server.log.With(
		slog.String("op", op),
	)

	if req == nil {
		log.Debug("create task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	createTaskRequest := dto.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := createTaskRequest.Validate(); err != nil {
		log.Debug("validate create task request failed", "error", err)
		return nil, err
	}

	task, err := server.taskService.Create(ctx, createTaskRequest)

	if err != nil {
		log.Error(op, "create task failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &taskservice.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
	}, nil
}

func (server serverAPI) GetTask(ctx context.Context, req *taskservice.GetTaskRequest) (*taskservice.Task, error) {
	const op = "serverAPI.GetTask"

	log := server.log.With(
		slog.String("op", op),
	)

	if req == nil {
		log.Debug("create task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	request := &dto.GetTaskRequest{
		ID: req.Id,
	}

	err := request.Validate()
	if err != nil {
		log.Debug("validate get task request failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	taskResponse, err := server.taskService.GetTask(ctx, request)
	if err != nil {
		log.Error("get task failed", "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error: "+err.Error())
	}

	return &taskservice.Task{
		Id:          taskResponse.ID,
		Title:       taskResponse.Title,
		Description: taskResponse.Description,
	}, nil
}

func (server serverAPI) UpdateTask(ctx context.Context, req *taskservice.UpdateTaskRequest) (*taskservice.Task, error) {
	const op = "serverAPI.GetTask"

	log := server.log.With(
		slog.String("op", op),
	)

	if req == nil {
		log.Debug("update task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	request := &dto.UpdateTaskRequest{
		ID:          req.Id,
		Title:       req.Title,
		Description: req.Description,
	}

	err := request.Validate()
	if err != nil {
		log.Debug("validate update task request failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	taskResponse, err := server.taskService.UpdateTask(ctx, request)
	if err != nil {
		log.Error("update task failed", "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskservice.Task{
		Id:          taskResponse.ID,
		Title:       taskResponse.Title,
		Description: taskResponse.Description,
	}, nil
}

func (server serverAPI) DeleteTask(ctx context.Context, req *taskservice.DeleteTaskRequest) (*taskservice.DeleteTaskResponse, error) {
	const op = "serverAPI.DeleteTask"

	log := server.log.With(
		slog.String("op", op),
	)

	if req == nil {
		log.Debug("delete task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	request := &dto.DeleteTaskRequest{
		ID: req.Id,
	}

	err := request.Validate()
	if err != nil {
		log.Debug("validate delete task request failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	response, err := server.taskService.DeleteTask(ctx, request)
	if err != nil {
		log.Error("delete task failed", "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskservice.DeleteTaskResponse{
		Id: response.ID,
	}, nil
}
