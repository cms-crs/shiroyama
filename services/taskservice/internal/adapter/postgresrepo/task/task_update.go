package postgresrepo

import (
	"context"
	"taskservice/internal/dto"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) UpdateTask(ctx context.Context, taskRequest *dto.UpdateTaskRequest) (*entity.Task, error) {
	var task entity.Task

	if err := repository.db.WithContext(ctx).First(&task, "id = ?", taskRequest.ID).Error; err != nil {
		return nil, err
	}

	err := repository.db.Model(task).WithContext(ctx).Save(task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}
