package models

import (
	"time"
)

type Note struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" binding:"required"`
	Content   string    `json:"content" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Tags      []Tag     `gorm:"many2many:note_tags;"`
}
