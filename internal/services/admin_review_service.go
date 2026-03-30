package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"gorm.io/gorm"
)

type AdminReviewService struct {
	db *gorm.DB
}

func NewAdminReviewService(db *gorm.DB) *AdminReviewService {
	return &AdminReviewService{db: db}
}

func (service *AdminReviewService) Approve(req requests.ReviewApproveRequest) (*models.Review, error) {
	errorsValid, ok := validators.ValidateApproveReview(service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var review models.Review
	err := service.db.Preload("Movie").Model(&review).Where("id = ?", req.ReviewID).First(&review).Error
	if err != nil {
		return nil, err
	}
	if req.Status == models.ReviewStatusApproved {
		var movieReviewsRatings []float64
		err = service.db.Model(&models.Review{}).Where("movie_id = ? AND status = ?", review.MovieID,
			models.ReviewStatusApproved).Pluck("rating", &movieReviewsRatings).Error
		if err != nil {
			return nil, err
		}
		err := RecalculateRatingMovie(service.db, review.Movie, review.Rating, movieReviewsRatings)
		if err != nil {
			return nil, err
		}
		err = service.db.Model(&review).Update("status", models.ReviewStatusApproved).Error
		if err != nil {
			return nil, err
		}
		return &review, nil
	}
	reviewUpdates := map[string]interface{}{
		"status":        req.Status,
		"reject_reason": *req.ReasonReject,
	}
	err = service.db.Model(&review).Updates(reviewUpdates).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}
