package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"primaryKey;size:64" json:"id"`
	Name         string         `gorm:"size:255;not null" json:"name"`
	Username     string         `gorm:"size:100;uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	RoleID       string         `gorm:"size:64;index" json:"roleId"`
	Role         *Role          `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Department   string         `gorm:"size:255" json:"department"`
	Avatar       string         `gorm:"size:512" json:"avatar"`
	Status       string         `gorm:"size:32;default:'Active'" json:"status"`
	TwoFactor    bool           `gorm:"default:false" json:"twoFactor"`
	LastLoginAt  *time.Time     `json:"lastLoginAt"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Role struct {
	ID          string         `gorm:"primaryKey;size:64" json:"id"`
	Name        string         `gorm:"size:100;uniqueIndex;not null" json:"name"`
	Code        string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Description string         `gorm:"size:512" json:"description"`
	IsSystem    bool           `gorm:"default:false" json:"isSystem"`
	Status      string         `gorm:"size:32;default:'Active'" json:"status"`
	UsersCount  int            `gorm:"-" json:"usersCount"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Permission struct {
	ID          string    `gorm:"primaryKey;size:64" json:"id"`
	Module      string    `gorm:"size:64;index;not null" json:"module"`
	Action      string    `gorm:"size:64;not null" json:"action"`
	Code        string    `gorm:"size:128;uniqueIndex;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RolePermission struct {
	RoleID       string `gorm:"primaryKey;size:64"`
	PermissionID string `gorm:"primaryKey;size:64"`
}
