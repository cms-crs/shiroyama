package dto

import "time"

type UpdateTeamRequest struct {
	ID          string
	Name        string
	Description string
}

type UpdateTeamResponse struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
