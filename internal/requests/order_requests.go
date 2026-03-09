package requests

type OrderCreateRequest struct {
	PremiereID uint            `gorm:"nit null" json:"premiere_ids" binding:"required"`
	UserId     uint            `gorm:"not null;index" json:"user_id"`
	Seats      map[uint][]uint `gorm:"type:text;not null" json:"seats"`
	Coins      *uint64         `json:"coins"`
}
