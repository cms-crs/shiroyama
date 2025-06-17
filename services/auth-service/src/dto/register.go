package dto

import "github.com/go-playground/validator"

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	UserId uint `json:"user_id"`
}

func (req *RegisterRequest) Validate() error {
	return validator.New().Struct(req)
}
