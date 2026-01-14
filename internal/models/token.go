package models

import (
	"gorm.io/gorm"
	"time"
)

type Token struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	UserID     uint   `gorm:"not null;index" json:"user_id"`
	User       User   `gorm:"foreignKey:UserID" json:"-"`
	Token      string `gorm:"type:text;not null;unique" json:"token"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
	DeviceInfo string `gorm:"size:255" json:"device_info"`

	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
