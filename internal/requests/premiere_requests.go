package requests

import (
	"time"
)

type PremiereCreateRequest struct {
	MovieID     uint      `json:"movie_id" binding:"required"`
	Hall        string    `json:"hall" binding:"required"`
	Price       float64   `json:"price" binding:"required,min=0"`
	Rows        uint      `json:"rows" binding:"required,min=1,max=20"`
	SeatsPerRow uint      `json:"seats_per_row" binding:"required,min=1,max=30"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
}

type PremiereUpdateRequest struct {
	Id        uint       `json:"id" binding:"required"`
	MovieID   *uint      `json:"movie_id" binding:""`
	Hall      *string    `json:"hall" binding:""`
	Price     *float64   `json:"price" binding:""`
	StartTime *time.Time `json:"start_time" binding:""`
	EndTime   *time.Time `json:"end_time" binding:""`
}
