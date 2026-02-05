package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/validators"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type AdminMovieService struct {
	db *gorm.DB
}

func NewAdminMovieService(db *gorm.DB) *AdminMovieService {
	return &AdminMovieService{db: db}
}

func (service *AdminMovieService) CreateMovie(req requests.MovieCreateRequest) (*models.Movie, error) {

	errorsValid, ok := validators.ValidateMovie(service.db, req)
	if !ok {
		var errorMsgs []string
		for field, err := range errorsValid {
			errorMsgs = append(errorMsgs, field+": "+err)
		}
		return nil, errors.New(strings.Join(errorMsgs, "\n"))
	}
	var genres []models.Genre
	service.db.Model(&models.Genre{}).Where("id IN ?", req.GenreIDS).Find(&genres)

	movie := &models.Movie{
		Title:       req.Title,
		Description: req.Description,
		Duration:    req.Duration,
		AgeRating:   req.AgeRating,
		PosterURL:   req.PosterURL,
		TrailerURL:  req.TrailerURL,
		ReleaseDate: req.ReleaseDate,
		Rating:      0.0,
		RatingCount: 0,
		Genres:      genres,
	}
	err := service.db.Create(movie).Error
	if err != nil {
		return nil, errors.New("Failed to create movie")
	}

	return movie, nil
}
