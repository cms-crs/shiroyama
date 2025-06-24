package service

import (
	"context"
	"taskservice/internal/dto"
)

func (service *TaskService) UpdateTask(ctx context.Context, request *dto.UpdateTaskRequest) (*dto.UpdateTaskResponse, error) {
	task, err := service.taskRepository.UpdateTask(ctx, request)
	if err != nil {
		return nil, err
	}

	return dto.NewUpdateTaskResponse(task), nil
}
