package dto

import "userservice/internal/entity"

type TeamMember struct {
	UserId string
	TeamId string
	Role   string
	User   entity.User
}
