package dto

import (
	"time"
	"userservice/internal/entity"
)

type TeamMember struct {
	UserID    string
	TeamID    string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      entity.User
}
