package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/validators"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
)

type AdminMovieService struct {
	db *gorm.DB
}

func NewAdminMovieService(db *gorm.DB) *AdminMovieService {
	return &AdminMovieService{db: db}
}

func (service *AdminMovieService) Create(req requests.MovieCreateRequest) (*models.Movie, error) {

	errorsValid, ok := validators.ValidateCreateMovie(service.db, req)
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

func (service *AdminMovieService) Update(req requests.MovieUpdateRequest) (*models.Movie, error) {
	errorsValid, updateData, genreIds, ok := validators.ValidateUpdateMovie(service.db, req)
	if !ok {
		var errorMsgs []string
		for field, err := range errorsValid {
			errorMsgs = append(errorMsgs, field+": "+err)
		}
		return nil, errors.New(strings.Join(errorMsgs, "\n"))
	}
	var movie models.Movie
	err1 := service.db.Preload("Genres").Where("id = ?", req.Id).First(&movie)

	if err1.Error != nil {
		return nil, errors.New("Failed to find movie , unncorrect id")
	}

	if len(updateData) > 0 {
		service.db.Model(&movie).Updates(updateData)
	}

	if len(genreIds) > 0 {
		var newGenres []models.Genre
		service.db.Where("id IN ?", genreIds).Find(&newGenres)
		fmt.Println(genreIds)
		if err := service.db.Model(&movie).Association("Genres").Replace(newGenres); err != nil {
			return nil, errors.New("Failed to update genres: ")
		}
	}

	return &movie, nil
}
