package dto

import (
	"taskservice/internal/entity"
	"time"
)

type TeamMember struct {
	UserID    string
	TeamID    string
	Role      string
	CreatedAt time.Time
	UpdatedAt time.Time
	User      entity.User
}
