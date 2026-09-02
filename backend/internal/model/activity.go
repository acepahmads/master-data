package model

import (
	"time"
)

type ActivityLog struct {
	ID          string    `gorm:"primaryKey;size:64" json:"id"`
	UserID      string    `gorm:"size:64;index" json:"userId"`
	UserName    string    `gorm:"size:255" json:"user"`
	UserAvatar  string    `gorm:"size:512" json:"userAvatar"`
	Module      string    `gorm:"size:64;index;not null" json:"module"`
	EntityType  string    `gorm:"size:64;index;not null" json:"targetType"`
	EntityID    string    `gorm:"size:64;index;not null" json:"targetId"`
	EntityName  string    `gorm:"size:255" json:"targetName"`
	Action      string    `gorm:"size:255;not null" json:"action"`
	Description string    `gorm:"type:text" json:"description"`
	BadgeColor  string    `gorm:"size:32;default:'blue'" json:"badgeColor"`
	Timestamp   string    `gorm:"-" json:"timestamp"`
	Date        string    `gorm:"-" json:"date"`
	CreatedAt   time.Time `json:"createdAt"`
}
