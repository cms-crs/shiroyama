package service

import (
	"context"
	"taskservice/internal/dto"
	"taskservice/internal/entity"
)

func (service *TaskService) UpdateTask(ctx context.Context, request *dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error) {
	task := &entity.Task{
		ID:          request.ID,
		Title:       request.Title,
		Description: request.Description,
	}

	task, err := service.taskRepository.UpdateTask(ctx, task)
	if err != nil {
		return nil, err
	}

	return dto.NewUpdateTaskResponse(task), nil
}
