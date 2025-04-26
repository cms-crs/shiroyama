package dto

import "github.com/go-playground/validator"

type GetTaskRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

type GetTaskResponse struct {
	ID          uint64 `json:"id" validate:"required"`
	Title       string
	Description string
}

func (r *GetTaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
