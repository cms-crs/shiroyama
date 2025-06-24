package postgresrepo

import (
	"context"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) DeleteTask(ctx context.Context, taskID uint) (uint, error) {
	var task entity.Task

	if err := repository.db.WithContext(ctx).First(&task, "id = ?", taskID).Error; err != nil {
		return 0, err
	}

	err := repository.db.WithContext(ctx).Delete(&task).Error
	if err != nil {
		return taskID, err
	}

	return taskID, nil
}
