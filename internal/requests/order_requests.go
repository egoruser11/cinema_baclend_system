package requests

import (
	"time"
)

type OrderCreateRequest struct {
	MovieID     uint      `json:"movie_id" binding:"required"`
	Hall        string    `json:"hall" binding:"required"`
	Price       float64   `json:"price" binding:"required,min=0"`
	Rows        uint      `json:"rows" binding:"required,min=1,max=20"`
	SeatsPerRow uint      `json:"seats_per_row" binding:"required,min=1,max=30"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
}
