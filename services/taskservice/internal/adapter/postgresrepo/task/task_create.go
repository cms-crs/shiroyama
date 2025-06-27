package postgresrepo

import (
	"context"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	return repository.db.WithContext(ctx).Create(task).Error
}
