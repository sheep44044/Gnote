package models

import "time"

type Tag struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_tag_name"`
	Name      string    `json:"name" gorm:"size:64;uniqueIndex:idx_user_tag_name"`
	Color     string    `json:"color" binding:"required"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Notes []Note `gorm:"many2many:note_tags;"`
}
