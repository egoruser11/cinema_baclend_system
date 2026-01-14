package models

import (
	"encoding/json"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

type Seat struct {
	Row    int  `json:"row"`
	Number int  `json:"number"`
	Booked bool `json:"booked"`
}

type Premiere struct {
	ID      uint    `gorm:"primaryKey" json:"id"`
	MovieID uint    `gorm:"not null;index" json:"movie_id"`
	Movie   Movie   `gorm:"foreignKey:MovieID" json:"-"`
	Hall    string  `gorm:"size:100;not null" json:"hall"` // "Зал 1", "IMAX"
	Price   float64 `gorm:"type:decimal(10,2);not null" json:"price"`

	// Конфигурация зала
	Rows        int `gorm:"not null" json:"rows"`
	SeatsPerRow int `gorm:"not null" json:"seats_per_row"`

	// Занятые места - храним в JSON
	BookedSeats datatypes.JSON `gorm:"type:jsonb" json:"booked_seats"` // массив Seat

	// Статистика
	TotalSeats  int `gorm:"-" json:"total_seats"`
	BookedCount int `gorm:"-" json:"booked_count"`

	// Связи
	Orders []Order `gorm:"foreignKey:PremiereID" json:"orders,omitempty"`

	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   time.Time `gorm:"not null" json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GORM Hooks для вычисляемых полей
func (p *Premiere) AfterFind(tx *gorm.DB) (err error) {
	p.TotalSeats = p.Rows * p.SeatsPerRow

	// Считаем занятые места из JSON
	var seats []Seat
	if err := json.Unmarshal(p.BookedSeats, &seats); err == nil {
		for _, seat := range seats {
			if seat.Booked {
				p.BookedCount++
			}
		}
	}
	return
}
