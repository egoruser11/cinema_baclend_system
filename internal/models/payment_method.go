package models

import (
	"gorm.io/gorm"
	"time"
)

type PaymentMethodType string

const (
	PaymentCard   PaymentMethodType = "card"
	PaymentWallet PaymentMethodType = "wallet"
	PaymentCoins  PaymentMethodType = "coins"
)

type PaymentMethod struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	UserID    uint              `gorm:"not null" json:"user_id"`
	User      User              `gorm:"foreignKey:UserID" json:"-"`
	Type      PaymentMethodType `gorm:"type:varchar(20);not null" json:"type"`
	Details   string            `gorm:"type:text" json:"details,omitempty"` // зашифрованные данные карты
	IsDefault bool              `gorm:"default:false" json:"is_default"`
	IsActive  bool              `gorm:"default:true;" json:"is_active"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
