package models

import (
	"time"
)

type Note struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `gorm:"index"`
	Title      string    `json:"title" binding:"required"`
	Content    string    `json:"content" binding:"required"`
	IsPrivate  bool      `gorm:"default:false" json:"is_private"`
	IsPinned   bool      `gorm:"default:false;index"` // 是否置顶
	IsFavorite bool      `gorm:"default:false;index"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	Tags       []Tag     `gorm:"many2many:note_tags;"`
}
