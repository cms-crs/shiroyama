package entity

import "gorm.io/gorm"

type Label struct {
	gorm.Model
	BoardID string `gorm:"size:255;not null;index"`
	Name    string `gorm:"size:255;not null"`
	Color   string `gorm:"size:7;not null;default:'#808080'"`
	Tasks   []Task `gorm:"many2many:task_labels;constraint:OnDelete:CASCADE"`
}
