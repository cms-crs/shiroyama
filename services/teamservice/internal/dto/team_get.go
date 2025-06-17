package dto

import "time"

type GetTeamResponse struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Role        string
}
