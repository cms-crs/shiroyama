package postgresrepo

import (
	"context"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) GetTask(ctx context.Context, ID uint64) (*entity.Task, error) {
	task := &entity.Task{
		ID: ID,
	}

	err := repository.db.WithContext(ctx).First(&task).Error
	if err != nil {
		return nil, err
	}

	return task, nil
}
