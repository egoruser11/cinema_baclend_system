package models

import (
	"errors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Seat struct {
	Row    int `json:"row"`
	Number int `json:"number"`
}

type Premiere struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	MovieID uint    `gorm:"not null;index" json:"movie_id"`
	Movie   Movie   `gorm:"foreignKey:MovieID" json:""`
	Hall    string  `gorm:"size:100;not null" json:"hall"` // "Зал 1", "IMAX"
	Price   float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	Rows        uint `gorm:"not null" json:"rows"`
	SeatsPerRow uint `gorm:"not null" json:"seats_per_row"`

	BookedSeats datatypes.JSON `gorm:"type:jsonb" json:"booked_seats"` // массив Seat

	TotalSeats  int `gorm:"not null" json:"total_seats"`
	BookedCount int `gorm:"not null" json:"booked_count"`

	Orders []Order `gorm:"foreignKey:PremiereID" json:"orders,omitempty"`

	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAvailablePremieres(db *gorm.DB, movieId uint) ([]Premiere, error) {
	var premieres []Premiere

	sql := `
	   SELECT p.*
	   FROM premieres p
	   WHERE p.movie_id = ?
	     AND p.booked_count < (p.rows * p.seats_per_row)
	     AND p.start_time > NOW()
	   ORDER BY p.start_time ASC
	`

	err := db.Raw(sql, movieId).Scan(&premieres).Error
	if err != nil {
		return nil, errors.New("Can not get premieres , please try again")
	}
	return premieres, nil
}
