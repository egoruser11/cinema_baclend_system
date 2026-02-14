package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
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
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
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
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var movie models.Movie
	err1 := service.db.
		Preload("Genres").
		Preload("Reviews").
		Preload("Premieres").
		Where("id = ?", req.Id).
		First(&movie)

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

func (service *AdminMovieService) Delete(id uint) error {
	err := service.db.Delete(&models.Movie{}, id)
	if err.Error != nil {
		return errors.New("Failed to delete movie")
	}
	return nil
}

func (service *AdminMovieService) Show(id uint) (*models.Movie, error) {
	var movie models.Movie
	err := service.db.
		Preload("Genres").
		Preload("Reviews").
		Preload("Premieres").
		Where("id = ?", id).
		First(&movie).Error
	if err != nil {
		return nil, errors.New("Failed to find movie , try again")
	}
	return &movie, nil
}

//index

func (service *AdminMovieService) Index(req requests.MovieIndexRequest) ([]*models.Movie, error) {
	errorsVaild, filters := validators.ValidateIndexMovie(req)
	if len(errorsVaild) > 0 {
		returnedErrors := []string{}
		for field, err := range errorsVaild {
			returnedErrors = append(returnedErrors, field+": "+err)
		}
		return nil, errors.New(strings.Join(returnedErrors, "\n"))
	}
	fmt.Println(filters)
	var movies []*models.Movie
	query := service.db.Model(&models.Movie{}).Preload("Genres", "Reviews")
	if len(filters) > 0 {
		sort, existsSort := filters["sort"]
		if existsSort {
			orderType := filters["order_type"]
			query.Order(sort + " " + orderType)
		}
		search, existsSearch := filters["search"]
		if existsSearch {
			resSearch := "%" + strings.ToLower(search) + "%"
			query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", resSearch, resSearch)
		}
	}
	limit, _ := strconv.Atoi(filters["limit"])
	offset, _ := strconv.Atoi(filters["offset"])
	query.Limit(limit).Offset(offset).Find(&movies)
	return movies, nil
}
