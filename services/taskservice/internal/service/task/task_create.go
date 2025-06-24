package service

import (
	"context"
	"taskservice/internal/dto"
	"taskservice/internal/entity"
)

func (service *TaskService) Create(ctx context.Context, req dto.CreateTaskRequest) (*dto.CreateTaskResponse, error) {
	task := &entity.Task{
		Title:       req.Title,
		Description: req.Description,
	}

	if err := service.taskRepository.Create(ctx, task); err != nil {
		return nil, err
	}

	return dto.NewCreateTaskResponse(task), nil
}
