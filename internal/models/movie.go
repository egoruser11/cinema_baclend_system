package models

import (
	"gorm.io/gorm"
	"time"
)

type AgeRating string

const (
	AgeRatingG    AgeRating = "G"     // 0+
	AgeRatingPG   AgeRating = "PG"    // 10+
	AgeRatingPG13 AgeRating = "PG-13" // 13+
	AgeRatingR    AgeRating = "R"     // 16+
	AgeRatingNC17 AgeRating = "NC-17" // 18+
)

type Movie struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Duration    int       `gorm:"not null" json:"duration"` // в минутах
	AgeRating   AgeRating `gorm:"type:varchar(10);not null" json:"age_rating"`
	PosterURL   string    `gorm:"size:500" json:"poster_url"`
	TrailerURL  string    `gorm:"size:500" json:"trailer_url"`
	ReleaseDate time.Time `json:"release_date"`
	Rating      float64   `gorm:"type:decimal(3,1);default:0.0" json:"rating"`
	RatingCount int       `gorm:"default:0" json:"rating_count"`

	// Связи
	Genres    []Genre    `gorm:"many2many:movie_genres;" json:"genres,omitempty"`
	Premieres []Premiere `gorm:"foreignKey:MovieID" json:"premieres,omitempty"`
	Reviews   []Review   `gorm:"foreignKey:MovieID" json:"reviews,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
