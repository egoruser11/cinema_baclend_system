package requests

type OrderCreateRequest struct {
	PremiereID            uint         `gorm:"nit null" json:"premiere_ids" binding:"required"`
	UserId                uint         `gorm:"not null;index" json:"user_id"`
	PaymentMethodID       uint         `gorm:"not null;index" json:"payment_method_id" binding:"required"`
	Seats                 map[uint]int `gorm:"type:text;not null" json:"seats"`
	CountMinutesBeforePay *uint        `json:"count_minutes_before_pay"`
	Coins                 *uint64      `json:"coins"`
}
