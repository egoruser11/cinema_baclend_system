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

func (service *UserReviewService) Update(c echo.Context, req requests.ReviewUpdateRequest) (*models.Review, error) {
	errorsValid, updates, ok := validators.ValidateUpdateReview(c, service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var review models.Review
	err := service.db.Preload("Movie").Model(&models.Review{}).Where("id = ?", req.ReviewID).First(&review).Error
	if err != nil {
		return nil, err
	}
	if review.Status == models.ReviewStatusApproved {
		var movieReviewsRatings []float64
		err = service.db.Model(&models.Review{}).Where("movie_id = ? AND status = ?", review.MovieID,
			models.ReviewStatusApproved).Pluck("rating", &movieReviewsRatings).Error
		if err != nil {
			return nil, err
		}
		err = RecalculateRatingMovie(service.db, review.Movie, review.Rating, movieReviewsRatings, true)
		if err != nil {
			return nil, err
		}
		updates["status"] = models.ReviewStatusUnderConsideration

	}
	err = service.db.Model(&review).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (service *UserReviewService) Delete(c echo.Context, req requests.ReviewDeleteRequest) error {
	errorsValid, ok := validators.ValidateDeleteReview(c, service.db, req)
	if !ok {
		return errors.New(utils.InputErrorsValid(errorsValid))
	}
	var review models.Review
	err := service.db.Preload("Movie").Model(&review).Where("id = ?", req.ReviewID).First(&review).Error
	if err != nil {
		return err
	}
	var movieReviewsRatings []float64
	err = service.db.Model(&models.Review{}).Where("movie_id = ? AND status = ?", review.MovieID,
		models.ReviewStatusApproved).Pluck("rating", &movieReviewsRatings).Error
	if err != nil {
		return nil
	}
	err = RecalculateRatingMovie(service.db, review.Movie, review.Rating, movieReviewsRatings, true)
	if err != nil {
		return err
	}
	err = service.db.Delete(&review).Error
	if err != nil {
		return err
	}
	return nil
}
