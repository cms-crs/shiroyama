package dto

import "time"

type UpdateUserRequest struct {
	ID       string
	Username string
	Password string
}

type UpdateUserResponse struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
