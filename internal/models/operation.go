package models

import (
	"time"
)

type OperationStatus string

type OperationType string

const (
	Purchase      OperationType = "purchase"
	Replenishment OperationType = "replenishment"
)

type Operation struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	UserID          uint            `gorm:"not null;index" json:"user_id"`
	User            User            `gorm:"foreignKey:UserID" json:"-"`
	OrderID         uint            `gorm:"nullable" json:"order_id"`
	PaymentMethodId uint            `gorm:"not null" json:"payment_method_id"`
	Order           Order           `gorm:"foreignKey:OrderID" json:"-"`
	PaymentMethod   PaymentMethod   `gorm:"foreignKey:PaymentMethodId" json:"-"`
	Amount          float64         `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status          OperationStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	Type            OperationType   `gorm:"type:varchar(20);not null;default:'pending'" json:"type"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
