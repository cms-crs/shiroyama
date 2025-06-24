package postgresrepo

import (
	"context"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) GetTask(ctx context.Context, ID uint) (*entity.Task, error) {
	var task entity.Task

	err := repository.db.WithContext(ctx).First(&task).Error
	if err != nil {
		return nil, err
	}

	return &task, nil
}
