package service

import (
	"context"
	"taskservice/internal/dto"
)

func (service *TaskService) DeleteTask(ctx context.Context, request *dto.DeleteTaskRequest) (*dto.DeleteTaskResponse, error) {
	id, err := service.taskRepository.DeleteTask(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return &dto.DeleteTaskResponse{
		ID: id,
	}, nil
}
