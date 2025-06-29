package dto

import (
	"github.com/go-playground/validator"
	"time"
)

type GetTaskRequest struct {
	ID uint `json:"id" validate:"required"`
}

type GetTaskResponse struct {
	ID          uint `json:"id" validate:"required"`
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (r *GetTaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
