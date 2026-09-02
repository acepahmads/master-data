package model

import (
	"time"

	"gorm.io/gorm"
)

type Manufacturer struct {
	ID              string         `gorm:"primaryKey;size:64" json:"id"`
	Name            string         `gorm:"size:255;index;not null" json:"name"`
	Code            string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Country         string         `gorm:"size:100" json:"country"`
	CountryCode     string         `gorm:"size:8" json:"countryCode"`
	Website         string         `gorm:"size:255" json:"website"`
	Status          string         `gorm:"size:32;default:'Approved'" json:"status"`
	Rating          string         `gorm:"size:64;default:'Approved'" json:"rating"`
	ContactPerson   string         `gorm:"size:255" json:"contactPerson"`
	ContactEmail    string         `gorm:"size:255" json:"contactEmail"`
	Description     string         `gorm:"type:text" json:"description"`
	TotalComponents int            `gorm:"default:0" json:"totalComponents"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type Supplier struct {
	ID              string         `gorm:"primaryKey;size:64" json:"id"`
	Code            string         `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Name            string         `gorm:"size:255;index;not null" json:"name"`
	Country         string         `gorm:"size:100" json:"country"`
	CountryCode     string         `gorm:"size:8" json:"countryCode"`
	Address         string         `gorm:"size:512" json:"address"`
	ContactPerson   string         `gorm:"size:255" json:"contact"`
	Email           string         `gorm:"size:255" json:"email"`
	Phone           string         `gorm:"size:64" json:"phone"`
	Website         string         `gorm:"size:255" json:"website"`
	Tier            string         `gorm:"size:100" json:"tier"`
	Status          string         `gorm:"size:32;default:'Active'" json:"status"`
	Terms           string         `gorm:"size:100" json:"terms"`
	LeadTimeAvg     string         `gorm:"size:64" json:"leadTimeAvg"`
	Description     string         `gorm:"type:text" json:"description"`
	TotalComponents int            `gorm:"default:0" json:"totalComponents"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
