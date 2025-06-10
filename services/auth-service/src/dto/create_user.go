package dto

import "github.com/go-playground/validator"

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (req *CreateUserRequest) Validate() error {
	return validator.New().Struct(req)
}
