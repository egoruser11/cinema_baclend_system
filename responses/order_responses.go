package responses

import (
	"cinema_backend_system/internal/models"
	"time"
)

type SummaryResponse struct {
	BonusesDecrease float64            `json:"bonuses_decrease"`
	BonusesIncrease float64            `json:"bonuses_increase"`
	ExpiredAt       *time.Time         `json:"expired_at"`
	MovieID         uint               `json:"movie_id"`
	PremiereID      uint               `json:"premiere_id"`
	Seats           map[uint][]uint    `json:"seats"`
	Status          models.OrderStatus `json:"status"`
}
