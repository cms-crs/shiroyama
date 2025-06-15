package dto

import (
	"github.com/go-playground/validator"
)

type GetUsersByIdsRequest struct {
	IDs []string `json:"ids" validate:"required"`
}

type GetUsersByIdsResponse struct {
	Users []*GetUserResponse
}

func (req *GetUsersByIdsRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(req)
}
