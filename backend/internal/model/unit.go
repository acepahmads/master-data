package model

import (
	"time"

	"gorm.io/gorm"
)

type Unit struct {
	ID          string         `gorm:"primaryKey;size:64" json:"id"`
	Code        string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Symbol      string         `gorm:"size:32" json:"symbol"`
	Category    string         `gorm:"size:64" json:"category"`
	Description string         `gorm:"size:255" json:"description"`
	Status      string         `gorm:"size:32;default:'Active'" json:"status"`
	IsDefault   bool           `gorm:"default:false" json:"isDefault"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
