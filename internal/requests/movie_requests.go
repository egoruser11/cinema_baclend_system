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

type MovieUpdateRequest struct {
	Id          int               `json:"id"`
	Title       *string           `json:"title,omitempty"`
	Description *string           `json:"description,omitempty"`
	Duration    *int              `json:"duration,omitempty"`
	AgeRating   *models.AgeRating `json:"age_rating,omitempty"`
	PosterURL   *string           `json:"poster_url,omitempty"`
	TrailerURL  *string           `json:"trailer_url,omitempty"`
	ReleaseDate *time.Time        `json:"release_date,omitempty"`
	GenreIDS    []int             `json:"genre_ids,omitempty"`
}

type MovieIdRequest struct {
	Id uint `query:"id"`
}

type MovieIndexRequest struct {
	Sort   *string `query:"sort"` //rating , release date , duration , age_rating
	IsDesc *bool   `query:"is_desc"`
	Search *string `query:"search"` //title , description
	Offset *uint   `query:"offset"`
	Limit  *uint   `query:"limit"`
}
