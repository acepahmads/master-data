package model

import (
	"time"

	"gorm.io/gorm"
)

type APIToken struct {
	ID           string         `gorm:"primaryKey;size:64" json:"id"`
	Name         string         `gorm:"size:128;not null" json:"name"`
	TokenHash    string         `gorm:"size:255;uniqueIndex;not null" json:"-"`
	TokenPreview string         `gorm:"size:64;not null" json:"tokenPreview"` // e.g. "iot_live_9a8f...4e1b"
	Scope        string         `gorm:"size:64;default:'ALL_READ'" json:"scope"` // ALL_READ, FULL_ACCESS
	ExpiresAt    *time.Time     `json:"expiresAt"` // nil means Unlimited / Never Expire
	LastUsedAt   *time.Time     `json:"lastUsedAt"`
	CreatedBy    string         `gorm:"size:64;not null" json:"createdBy"`
	Status       string         `gorm:"size:32;default:'Active'" json:"status"` // Active, Revoked
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
