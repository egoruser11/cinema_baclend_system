package requests

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
