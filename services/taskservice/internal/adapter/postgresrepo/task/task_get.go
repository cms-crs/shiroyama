package postgresrepo

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"taskservice/internal/entity"
)

func (repository *TaskRepository) GetTask(ctx context.Context, ID uint64) (*entity.Task, error) {
	task := &entity.Task{
		ID: ID,
	}

	err := repository.db.WithContext(ctx).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("task with ID %d not found", ID)
		}
		return nil, err
	}

	return task, nil
}
