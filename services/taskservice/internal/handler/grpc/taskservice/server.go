package taskservicegrpc

import (
	"context"
	"github.com/ShiroyamaY/protos/gen/go/taskservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"taskservice/internal/dto"
)

type TaskService interface {
	Create(ctx context.Context,
		request dto.CreateTaskRequest,
	) (*dto.CreateTaskResponse, error)
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
		log.Info("create task request is nil")
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}

	createTaskRequest := dto.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := createTaskRequest.Validate(); err != nil {
		log.Info("validate create task request failed", "error", err)
		return nil, err
	}

	task, err := server.taskService.Create(ctx, createTaskRequest)

	if err != nil {
		log.Error(op, "create task failed", "error", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &taskservice.Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
	}, nil
}
func (server serverAPI) GetTask(context.Context, *taskservice.GetTaskRequest) (*taskservice.Task, error) {
	panic("implement me")
}
func (server serverAPI) UpdateTask(context.Context, *taskservice.UpdateTaskRequest) (*taskservice.Task, error) {
	panic("implement me")
}
func (server serverAPI) DeleteTask(context.Context, *taskservice.DeleteTaskRequest) (*taskservice.DeleteTaskResponse, error) {
	panic("implement me")
}
