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
	err := service.db.Transaction(func(tx *gorm.DB) error {
		lockedReview, err := lockReviewForUpdate(tx, req.ReviewID)
		if err != nil {
			return err
		}

		if lockedReview.Status == models.ReviewStatusApproved {
			lockedMovie, err := lockMovieForUpdate(tx, lockedReview.MovieID)
			if err != nil {
				return err
			}
			lockedReview.Movie = *lockedMovie

			var movieReviewsRatings []float64
			err = tx.Model(&models.Review{}).Where("movie_id = ? AND status = ?", lockedReview.MovieID,
				models.ReviewStatusApproved).Pluck("rating", &movieReviewsRatings).Error
			if err != nil {
				return err
			}
			if err := RecalculateRatingMovie(tx, lockedReview.Movie, lockedReview.Rating, movieReviewsRatings, true); err != nil {
				return err
			}
			updates["status"] = models.ReviewStatusUnderConsideration
		}

		if err := tx.Model(lockedReview).Updates(updates).Error; err != nil {
			return err
		}
		review = *lockedReview
		if status, exists := updates["status"]; exists {
			review.Status = status.(models.ReviewStatus)
		}
		if rating, exists := updates["rating"]; exists {
			review.Rating = rating.(int)
		}
		if comment, exists := updates["comment"]; exists {
			review.Comment = comment.(string)
		}
		return nil
	})
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
	return service.db.Transaction(func(tx *gorm.DB) error {
		lockedReview, err := lockReviewForUpdate(tx, req.ReviewID)
		if err != nil {
			return err
		}
		if lockedReview.Status == models.ReviewStatusApproved {
			lockedMovie, err := lockMovieForUpdate(tx, lockedReview.MovieID)
			if err != nil {
				return err
			}
			lockedReview.Movie = *lockedMovie

			var movieReviewsRatings []float64
			err = tx.Model(&models.Review{}).Where("movie_id = ? AND status = ?", lockedReview.MovieID,
				models.ReviewStatusApproved).Pluck("rating", &movieReviewsRatings).Error
			if err != nil {
				return err
			}
			if err := RecalculateRatingMovie(tx, lockedReview.Movie, lockedReview.Rating, movieReviewsRatings, true); err != nil {
				return err
			}
		}
		return tx.Delete(lockedReview).Error
	})
}
