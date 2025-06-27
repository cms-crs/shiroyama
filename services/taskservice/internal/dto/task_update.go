package dto

import (
	"github.com/go-playground/validator"
	"taskservice/internal/entity"
	"time"
)

type UpdateTaskRequest struct {
	ID          uint   `json:"id" validate:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskResponse struct {
	ID          uint
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *UpdateTaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func NewUpdateTaskResponse(task *entity.Task) *UpdateTaskResponse {
	return &UpdateTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
