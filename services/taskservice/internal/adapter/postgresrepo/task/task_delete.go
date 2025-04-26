package postgresrepo

import (
	"context"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) DeleteTask(ctx context.Context, taskID uint64) (uint64, error) {
	task := &entity.Task{
		ID: taskID,
	}

	err := repository.db.WithContext(ctx).Delete(&task).Error
	if err != nil {
		return taskID, err
	}

	return taskID, nil
}
