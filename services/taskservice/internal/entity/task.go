package entity

type Task struct {
	ID          uint64 `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"size:255;not null"`
	Description string `gorm:"size:255"`
}
