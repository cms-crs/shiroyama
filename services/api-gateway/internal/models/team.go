package models

import "time"

type Team struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TeamMember struct {
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	User   User   `json:"user"`
}

type TeamWithRole struct {
	Team Team   `json:"team"`
	Role string `json:"role"`
}

type CreateTeamRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	CreatedBy   string `json:"created_by" binding:"required"`
}

type UpdateTeamRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
}

type AddUserToTeamRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role" binding:"required,oneof=admin member viewer"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin member viewer"`
}

type GetUserTeamsResponse struct {
	Teams []TeamWithRole `json:"teams"`
}

type GetTeamMembersResponse struct {
	Members []TeamMember `json:"members"`
}
