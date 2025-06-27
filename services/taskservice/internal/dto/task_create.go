package dto

import (
	"fmt"
	"github.com/go-playground/validator"
	"taskservice/internal/entity"
	"time"
)

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	ID          uint
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *CreateTaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

func NewCreateTaskResponse(task *entity.Task) *CreateTaskResponse {
	fmt.Println(task)
	return &CreateTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}

func NewGetTaskResponse(task *entity.Task) *GetTaskResponse {
	return &GetTaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
