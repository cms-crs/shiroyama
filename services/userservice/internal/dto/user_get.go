package dto

import (
	"github.com/go-playground/validator"
	"time"
)

type GetUserRequest struct {
	ID string `json:"id" validate:"required"`
}

type GetUserResponse struct {
	ID        string
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (req *GetUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}
