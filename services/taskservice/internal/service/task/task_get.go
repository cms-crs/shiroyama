package service

import (
	"context"
	"fmt"
	"taskservice/internal/dto"
)

func (service *TaskService) GetTask(ctx context.Context, taskRequest *dto.GetTaskRequest) (*dto.GetTaskResponse, error) {
	task, err := service.taskRepository.GetTask(ctx, taskRequest.ID)

	if err != nil {
		return nil, fmt.Errorf("GetTask: %w", err)
	}

	return dto.NewGetTaskResponse(task), nil
}
