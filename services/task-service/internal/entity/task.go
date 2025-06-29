package entity

import (
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	ListID      string           `gorm:"size:255;not null;index"`
	Title       string           `gorm:"size:255;not null"`
	Description string           `gorm:"type:text"`
	Position    int32            `gorm:"not null;default:0"`
	DueDate     *time.Time       `gorm:"null"`
	CreatedBy   string           `gorm:"size:255;not null;index"`
	Assignments []TaskAssignment `gorm:"foreignKey:TaskID;constraint:OnDelete:CASCADE"`
	Labels      []Label          `gorm:"many2many:task_labels;constraint:OnDelete:CASCADE"`
}
