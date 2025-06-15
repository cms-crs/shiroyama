package dto

import (
	userv1 "github.com/cms-crs/protos/gen/go/user_service"
	"github.com/go-playground/validator"
	"time"
)

type CreateUserRequest struct {
	Username string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type CreateUserResponse struct {
	ID        string
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (req *CreateUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}

func NewCreateUserRequest(req *userv1.CreateUserRequest) *CreateUserRequest {
	return &CreateUserRequest{
		req.Username,
		req.Email,
	}
}
