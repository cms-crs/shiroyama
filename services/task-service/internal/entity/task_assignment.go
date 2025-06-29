package entity

import (
	"gorm.io/gorm"
	"time"
)

type TaskAssignment struct {
	gorm.Model
	TaskID     uint      `gorm:"not null;index"`
	UserID     string    `gorm:"size:255;not null;index"`
	AssignedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Task       Task      `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
}
