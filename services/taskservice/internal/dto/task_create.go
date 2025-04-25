package dto

import (
	"fmt"
	"github.com/go-playground/validator"
	"taskservice/internal/entity"
)

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type CreateTaskResponse struct {
	ID          uint64
	Title       string
	Description string
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
	}
}
