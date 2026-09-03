package model

import (
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        string         `gorm:"primaryKey;size:64" json:"id"`
	Code      string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name      string         `gorm:"size:255;index;not null" json:"name"`
	Company   string         `gorm:"size:255;index" json:"company"`
	Phone     string         `gorm:"size:64" json:"phone"`
	Email     string         `gorm:"size:255" json:"email"`
	Descr     string         `gorm:"type:text" json:"descr"`
	Status    string         `gorm:"size:32;default:'Active'" json:"status"`
	Source    string         `gorm:"size:64;default:'MANUAL'" json:"source"` // IMPORT_XLS, MANUAL
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
