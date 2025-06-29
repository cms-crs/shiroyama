package postgresrepo

import (
	"gorm.io/gorm"
	"taskservice/internal/entity"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {

	err := db.AutoMigrate(&entity.Task{})
	if err != nil {
		panic(err)
	}

	return &TaskRepository{db: db}
}
