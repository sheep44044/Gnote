package models

import "time"

type UserFollow struct {
	ID         uint `gorm:"primaryKey"`
	FollowerID uint `gorm:"uniqueIndex:idx_follow_relation"` // 联合唯一索引
	FollowedID uint `gorm:"uniqueIndex:idx_follow_relation"`
	CreatedAt  time.Time
}

// TableName 对应 GORM 的一些 Hook 或者方法
func (UserFollow) TableName() string {
	return "user_follows"
}
