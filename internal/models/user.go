package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null;size:50"`
	Password string `gorm:"not null;size:255"`
	Avatar   string `gorm:"size:255" json:"avatar,omitempty"`
	Bio      string `gorm:"type:text" json:"bio,omitempty"`

	FollowCount int `gorm:"default:0" json:"follow_count,omitempty"`
	FanCount    int `gorm:"default:0" json:"fan_count,omitempty"`
}
