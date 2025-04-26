package service

import (
	"context"
	"log/slog"
	"taskservice/internal/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task *entity.Task) error
	GetTask(ctx context.Context, taskID uint64) (*entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error)
	DeleteTask(ctx context.Context, taskID uint64) (uint64, error)
}

type TaskService struct {
	log            *slog.Logger
	taskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}
