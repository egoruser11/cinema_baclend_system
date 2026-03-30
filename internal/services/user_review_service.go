package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserReviewService struct {
	db *gorm.DB
}

func NewUserReviewService(db *gorm.DB) *UserReviewService {
	return &UserReviewService{db: db}
}

func (service *UserReviewService) Create(c echo.Context, req requests.ReviewCreateRequest) (*models.Review, error) {
	errorsValid, ok := validators.ValidateCreateReview(c, service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var movie models.Movie
	err := service.db.Model(&movie).Where("id = ?", req.MovieID).First(&movie).Error
	if err != nil {
		return nil, err
	}
	var movieReviewsRatings []float64
	err = service.db.Model(&models.Review{}).Where("movie_id = ?", req.MovieID).Pluck("rating", &movieReviewsRatings).Error
	if err != nil {
		return nil, err
	}
	sumOfRatings := utils.GetSLiceSum(movieReviewsRatings) + float64(req.Rating)
	newMovieRatingCount := movie.RatingCount + 1
	newMovieRating := sumOfRatings / float64(newMovieRatingCount)
	movieUpdateData := map[string]interface{}{
		"rating_count": newMovieRatingCount,
		"rating":       newMovieRating,
	}
	err = service.db.Model(&movie).Updates(movieUpdateData).Error
	if err != nil {
		return nil, err
	}
	review := &models.Review{
		MovieID:   movie.ID,
		UserID:    c.Get("user_data").(*models.User).ID,
		Comment:   req.Comment,
		Status:    models.ReviewStatusUnderConsideration,
		IsVisible: true,
		Rating:    req.Rating,
	}
	err = service.db.Create(&review).Error
	if err != nil {
		return nil, err
	}

	return review, nil
}
