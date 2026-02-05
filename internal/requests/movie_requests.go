package requests

import (
	"cinema_backend_system/internal/models"
	"time"
)

type MovieCreateRequest struct {
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Duration    int              `json:"duration"`
	AgeRating   models.AgeRating `json:"age_rating"`
	PosterURL   string           `json:"poster_url"`
	TrailerURL  string           `json:"trailer_url"`
	ReleaseDate time.Time        `json:"release_date"`
	GenreIDS    []int            `json:"genre_ids"`
}
