package dto

import (
	"github.com/go-playground/validator"
)

type DeleteTaskRequest struct {
	ID uint `json:"id" validate:"required"`
}

type DeleteTaskResponse struct {
	ID uint
}

func (r *DeleteTaskRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
