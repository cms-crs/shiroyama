package postgresrepo

import (
	"context"
	"taskservice/internal/entity"
)

func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}
