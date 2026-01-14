package models

import (
	"time"
)

type OrderStatus string

const (
	OrderPending   OrderStatus = "pending"
	OrderPaid      OrderStatus = "paid"
	OrderCancelled OrderStatus = "cancelled"
	OrderRefunded  OrderStatus = "refunded"
)

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null;index" json:"user_id"`
	User        User        `gorm:"foreignKey:UserID" json:"-"`
	PremiereID  uint        `gorm:"not null;index" json:"premiere_id"`
	Premiere    Premiere    `gorm:"foreignKey:PremiereID" json:"premiere"`
	Seats       string      `gorm:"type:text;not null" json:"seats"` // "A1,A2,B3"
	TotalAmount float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      OrderStatus `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
