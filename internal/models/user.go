package models

import (
	"gorm.io/gorm"
	"time"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type UserStatus string

const (
	active   UserStatus = "active"
	blocked  UserStatus = "blocked"
	inactive UserStatus = "inactive"
)

type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:100;not null;unique" json:"username"`
	Email        string     `gorm:"size:255;not null;unique" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"` // не отдаем в JSON
	Role         UserRole   `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	Age          uint       `gorm:"check:age >= 10 AND age <= 99" json:"age"`
	MoneyBalance float64    `gorm:"type:decimal(10,2);default:0.00" json:"balance"`
	Status       UserStatus `gorm:"not null" json:"status"`
	CoinBalance  uint64     `gorm:"default:0" json:"coin_balance"`
	// Связи
	PaymentMethods []PaymentMethod `gorm:"foreignKey:UserID" json:"payment_methods,omitempty"`
	Tokens         []Token         `gorm:"foreignKey:UserID" json:"tokens,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
