package dto

import (
	"github.com/go-playground/validator"
	"taskservice/internal/entity"
)

type UpdateTaskRequest struct {
	ID          uint64 `json:"id" validate:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTaskResponse struct {
	ID          uint64
	Title       string
	Description string
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
	}
}
