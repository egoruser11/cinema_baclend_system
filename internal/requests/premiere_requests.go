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

type PremiereIndexRequest struct {
	MovieID *uint `query:"movie_id" binding:"required"`

	Sort   *string `query:"sort"` //price , booked_count , total_seats , start_time
	IsDesc *bool   `query:"is_desc"`

	DayPremiere *string `query:"day_premiere"` // "1999-12-22"

	HourFrom *string `query:"hour_from"` // "15-00"
	HourTo   *string `query:"hour_to"`   // "15-00"

	PriceMin *float64 `query:"price_min"`
	PriceMax *float64 `query:"price_max"`

	WeekDay *uint `query:"week_day"`

	Offset *uint `query:"offset"`
	Limit  *uint `query:"limit"`
}

type PremiereIdRequest struct {
	Id uint `query:"id"`
}
