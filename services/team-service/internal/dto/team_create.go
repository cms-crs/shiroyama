package dto

import "time"

type CreateTeamRequest struct {
	Name        string
	Description string
	CreatedBy   string
}

type CreateTeamResponse struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
