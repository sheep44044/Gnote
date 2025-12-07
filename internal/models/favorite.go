package models

import "time"

type Favorite struct {
	UserID uint `gorm:"primaryKey"`
	NoteID uint `gorm:"primaryKey"`

	CreatedAt time.Time
}
