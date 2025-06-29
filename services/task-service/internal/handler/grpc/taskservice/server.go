package handler

import (
	"context"
	taskv1 "github.com/cms-crs/protos/gen/go/task_service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
)

type TaskService interface {
	CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.Task, error)
	GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.Task, error)
	UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.Task, error)
	DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*emptypb.Empty, error)
	MoveTask(ctx context.Context, req *taskv1.MoveTaskRequest) (*taskv1.Task, error)
	AssignUser(ctx context.Context, req *taskv1.AssignUserRequest) (*taskv1.Task, error)
	UnassignUser(ctx context.Context, req *taskv1.UnassignUserRequest) (*taskv1.Task, error)
	GetTasksForLists(ctx context.Context, req *taskv1.GetTasksForListsRequest) (*taskv1.GetTasksForListsResponse, error)
	GetTasksForUser(ctx context.Context, req *taskv1.GetTasksForUserRequest) (*taskv1.GetTasksForUserResponse, error)
}

type Handler struct {
	taskv1.UnimplementedTaskServiceServer
	log         *slog.Logger
	taskService TaskService
}

func NewHandler(log *slog.Logger, taskService TaskService) *Handler {
	return &Handler{
		log:         log,
		taskService: taskService,
	}
}

func (h *Handler) CreateTask(ctx context.Context, req *taskv1.CreateTaskRequest) (*taskv1.Task, error) {
	const op = "handler.CreateTask"

	log := h.log.With(
		slog.String("op", op),
		slog.String("list_id", req.ListId),
		slog.String("title", req.Title),
		slog.String("created_by", req.CreatedBy),
	)

	log.Info("Create task request received")

	task, err := h.taskService.CreateTask(ctx, req)
	if err != nil {
		log.Error("Failed to create task", "error", err)
		return nil, err
	}

	log.Info("Task created successfully", "task_id", task.Id)

	return task, nil
}

func (h *Handler) GetTask(ctx context.Context, req *taskv1.GetTaskRequest) (*taskv1.Task, error) {
	const op = "handler.GetTask"

	log := h.log.With(
		slog.String("op", op),
		slog.String("task_id", req.Id),
	)

	log.Info("Get task request received")

	task, err := h.taskService.GetTask(ctx, req)
	if err != nil {
		log.Error("Failed to get task", "error", err)
		return nil, err
	}

	log.Info("Task retrieved successfully")

	return task, nil
}

func (h *Handler) UpdateTask(ctx context.Context, req *taskv1.UpdateTaskRequest) (*taskv1.Task, error) {
	const op = "handler.UpdateTask"

	log := h.log.With(
		slog.String("op", op),
		slog.String("task_id", req.Id),
		slog.String("title", req.Title),
	)

	log.Info("Update task request received")

	task, err := h.taskService.UpdateTask(ctx, req)
	if err != nil {
		log.Error("Failed to update task", "error", err)
		return nil, err
	}

	log.Info("Task updated successfully")

	return task, nil
}

func (h *Handler) DeleteTask(ctx context.Context, req *taskv1.DeleteTaskRequest) (*emptypb.Empty, error) {
	const op = "handler.DeleteTask"

	log := h.log.With(
		slog.String("op", op),
		slog.String("task_id", req.Id),
	)

	log.Info("Delete task request received")

	result, err := h.taskService.DeleteTask(ctx, req)
	if err != nil {
		log.Error("Failed to delete task", "error", err)
		return nil, err
	}

	log.Info("Task deleted successfully")

	return result, nil
}

func (h *Handler) MoveTask(ctx context.Context, req *taskv1.MoveTaskRequest) (*taskv1.Task, error) {
	const op = "handler.MoveTask"

	log := h.log.With(
		slog.String("op", op),
		slog.String("task_id", req.TaskId),
		slog.String("to_list_id", req.ToListId),
		slog.Int("position", int(req.Position)),
	)

	log.Info("Move task request received")

	task, err := h.taskService.MoveTask(ctx, req)
	if err != nil {
		log.Error("Failed to move task", "error", err)
		return nil, err
	}

	log.Info("Task moved successfully")

	return task, nil
}

func (h *Handler) AssignUser(ctx context.Context, req *taskv1.AssignUserRequest) (*taskv1.Task, error) {
	const op = "handler.AssignUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("task_id", req.TaskId),
		slog.String("user_id", req.UserId),
	)

	log.Info("Assign user request received")

	task, err := h.taskService.AssignUser(ctx, req)
	if err != nil {
		log.Error("Failed to assign user", "error", err)
		return nil, err
	}

	log.Info("User assigned successfully")

	return task, nil
}

func (h *Handler) UnassignUser(ctx context.Context, req *taskv1.UnassignUserRequest) (*taskv1.Task, error) {
	const op = "handler.UnassignUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("task_id", req.TaskId),
		slog.String("user_id", req.UserId),
	)

	log.Info("Unassign user request received")

	task, err := h.taskService.UnassignUser(ctx, req)
	if err != nil {
		log.Error("Failed to unassign user", "error", err)
		return nil, err
	}

	log.Info("User unassigned successfully")

	return task, nil
}

func (h *Handler) GetTasksForLists(ctx context.Context, req *taskv1.GetTasksForListsRequest) (*taskv1.GetTasksForListsResponse, error) {
	const op = "handler.GetTasksForLists"

	log := h.log.With(
		slog.String("op", op),
		slog.Int("list_count", len(req.ListIds)),
	)

	log.Info("Get tasks for lists request received")

	response, err := h.taskService.GetTasksForLists(ctx, req)
	if err != nil {
		log.Error("Failed to get tasks for lists", "error", err)
		return nil, err
	}

	log.Info("Tasks for lists retrieved successfully", "task_count", len(response.Tasks))

	return response, nil
}

func (h *Handler) GetTasksForUser(ctx context.Context, req *taskv1.GetTasksForUserRequest) (*taskv1.GetTasksForUserResponse, error) {
	const op = "handler.GetTasksForUser"

	log := h.log.With(
		slog.String("op", op),
		slog.String("user_id", req.UserId),
	)

	log.Info("Get tasks for user request received")

	response, err := h.taskService.GetTasksForUser(ctx, req)
	if err != nil {
		log.Error("Failed to get tasks for user", "error", err)
		return nil, err
	}

	log.Info("Tasks for user retrieved successfully", "task_count", len(response.Tasks))

	return response, nil
}

func Register(gRPCServer interface{}, log *slog.Logger, taskService TaskService) {
	handler := NewHandler(log, taskService)
	taskv1.RegisterTaskServiceServer(gRPCServer.(*grpc.Server), handler)
}
