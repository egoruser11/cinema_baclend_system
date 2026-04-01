package requests

import (
	"cinema_backend_system/internal/models"
	"time"
)

type OrderCreateRequest struct {
	PremiereID uint            `gorm:"nit null" json:"premiere_id" binding:"required"`
	Seats      map[uint][]uint `gorm:"type:text;not null" json:"seats"`
}

type OrderPaidRequest struct {
	OrderID uint     `gorm:"not null" json:"order_id"`
	Coins   *float64 `json:"coins"`
}
type OrderRefundedRequest struct {
	ID uint `gorm:"not null" json:"id"`
}
type OrderIndexRequest struct {
	Status     *models.OrderStatus `query:"status"`
	PremiereID *uint               `query:"premiere_id"`
	DateFrom   *time.Time          `query:"date_from"`
	DateTo     *time.Time          `query:"date_to"`
	Limit      *int                `query:"limit"`
	Offset     *int                `query:"offset"`
	Sort       *string             `query:"sort"`
	Order      *string             `query:"order"`
}

type OrderShowRequest struct {
	ID *uint `json:"id" query:"id"`
}

type OrderDeleteRequest struct {
	ID *uint `query:"id"`
}

type OrderUpdateRequest struct {
	ID    *uint           `json:"id"`
	Seats map[uint][]uint `json:"seats"`
}
