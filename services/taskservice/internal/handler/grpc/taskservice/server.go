package taskservicegrpc

import (
	"context"
	"errors"
	"github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
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
	taskv1.UnimplementedTaskServiceServer
	log         *slog.Logger
	taskService TaskService
}

func Register(gRPC *grpc.Server, taskService TaskService, log *slog.Logger) {
	taskv1.RegisterTaskServiceServer(gRPC, serverAPI{
		taskService: taskService,
		log:         log,
	})
}

func (server serverAPI) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.Task, error) {
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

	return &taskv1.Task{
		Id:          strconv.Itoa(int(task.ID)),
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   timestamppb.New(task.CreatedAt),
		UpdatedAt:   timestamppb.New(task.UpdatedAt),
	}, nil
}

func (server serverAPI) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.Task, error) {
	const op = "serverAPI.GetTask"

	log := server.log.With(
		slog.String("op", op),
	)

	if req == nil {
		log.Debug("create task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	id, err := strconv.Atoi(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	request := &dto.GetTaskRequest{
		ID: uint(id),
	}

	err = request.Validate()
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

	return &taskv1.Task{
		Id:          strconv.Itoa(int(taskResponse.ID)),
		Title:       taskResponse.Title,
		Description: taskResponse.Description,
		CreatedAt:   timestamppb.New(taskResponse.CreatedAt),
		UpdatedAt:   timestamppb.New(taskResponse.UpdatedAt),
	}, nil
}

func (server serverAPI) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.Task, error) {
	const op = "serverAPI.GetTask"

	log := server.log.With(
		slog.String("op", op),
	)

	if req == nil {
		log.Debug("update task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	id, err := strconv.Atoi(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	request := &dto.UpdateTaskRequest{
		ID:          uint(id),
		Title:       req.Title,
		Description: req.Description,
	}

	err = request.Validate()
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

	return &taskv1.Task{
		Id:          strconv.Itoa(int(taskResponse.ID)),
		Title:       taskResponse.Title,
		Description: taskResponse.Description,
		CreatedAt:   timestamppb.New(taskResponse.CreatedAt),
		UpdatedAt:   timestamppb.New(taskResponse.UpdatedAt),
	}, nil
}

func (server serverAPI) DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*emptypb.Empty, error) {
	const op = "serverAPI.DeleteTask"

	log := server.log.With(
		slog.String("op", op),
	)

	if req == nil {
		log.Debug("delete task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	id, err := strconv.Atoi(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	request := &dto.DeleteTaskRequest{
		ID: uint(id),
	}

	err = request.Validate()
	if err != nil {
		log.Debug("validate delete task request failed", "error", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err = server.taskService.DeleteTask(ctx, request)
	if err != nil {
		log.Error("delete task failed", "error", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}
