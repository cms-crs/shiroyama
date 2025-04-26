package dto

import (
	"github.com/go-playground/validator"
)

type DeleteTaskRequest struct {
	ID uint64 `json:"id" validate:"required"`
}

type DeleteTaskResponse struct {
	ID uint64
}

func (r *DeleteTaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
