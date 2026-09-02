package model

import (
	"time"

	"gorm.io/gorm"
)

type ComponentCategory struct {
	ID              string               `gorm:"primaryKey;size:64" json:"id"`
	Name            string               `gorm:"size:255;not null" json:"name"`
	Code            string               `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Description     string               `gorm:"type:text" json:"description"`
	ParentID        *string              `gorm:"size:64;index" json:"parentId"`
	Parent          *ComponentCategory   `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children        []ComponentCategory  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Status          string               `gorm:"size:32;default:'Active'" json:"status"`
	TotalComponents int                  `gorm:"default:0" json:"totalComponents"`
	CreatedAt       time.Time            `json:"createdAt"`
	UpdatedAt       time.Time            `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt       `gorm:"index" json:"-"`
}
