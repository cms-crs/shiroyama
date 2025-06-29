package entity

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"size:255"`
}
