package dto

import "github.com/go-playground/validator"

type TaskRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

type TaskResponse struct {
	Title       string
	Description string
}

func (r *TaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
