package entity

import "time"

type User struct {
	ID        string
	Email     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	IsDeleted bool
}
