package postgresrepo

import (
	"context"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) UpdateTask(ctx context.Context, task *entity.Task) (*entity.Task, error) {
	err := repository.db.Model(task).WithContext(ctx).Save(task).Error
	if err != nil {
		return nil, err
	}

	return task, nil
}
