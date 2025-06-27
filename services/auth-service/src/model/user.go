package model

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
)

type role string

const (
	Admin   role = "Admin"
	Regular role = "Regular"
)

func (r *role) Scan(value interface{}) error {
	switch v := value.(type) {
	case string:
		*r = role(v)
	case []byte:
		*r = role(v)
	default:
		return fmt.Errorf("cannot scan role from %T", value)
	}
	return nil
}

func (r role) Value() (driver.Value, error) {
	return string(r), nil
}

func (r role) String() string {
	return string(r)
}

type User struct {
	gorm.Model
	Email        string `gorm:"unique"`
	Password     []byte
	RefreshToken string
	Role         role `gorm:"type:user_role;default:'Regular';not null"`
	UserID       string `gorm:"unique"`
}
