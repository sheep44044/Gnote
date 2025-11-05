package models

import "time"

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `gorm:"uniqueIndex"`
	Color     string    `json:"color" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Notes     []Note    `gorm:"many2many:note_tags;"`
}
