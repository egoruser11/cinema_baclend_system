package requests

type OrderCreateRequest struct {
	PremiereID uint            `gorm:"nit null" json:"premiere_id" binding:"required"`
	Seats      map[uint][]uint `gorm:"type:text;not null" json:"seats"`
	Coins      *uint64         `json:"coins"`
}
