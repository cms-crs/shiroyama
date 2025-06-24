package service

import (
	"context"
	"log/slog"
	"taskservice/internal/dto"
	"taskservice/internal/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetTask(ctx context.Context, taskID uint) (*entity.Task, error)
	UpdateTask(ctx context.Context, task *dto.UpdateTaskRequest) (*entity.Task, error)
	DeleteTask(ctx context.Context, taskID uint) (uint, error)
}

type TaskService struct {
	log            *slog.Logger
	taskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}
