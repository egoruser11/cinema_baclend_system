package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	err := service.db.Transaction(func(tx *gorm.DB) error {
		lockedReview, err := lockReviewForUpdate(tx, req.ReviewID)
		if err != nil {
			return err
		}
		lockedMovie, err := lockMovieForUpdate(tx, lockedReview.MovieID)
		if err != nil {
			return err
		}
		lockedReview.Movie = *lockedMovie

		if req.Status == models.ReviewStatusApproved {
			var movieReviewsRatings []float64
			err = tx.Model(&models.Review{}).Where("movie_id = ? AND status = ?", lockedReview.MovieID,
				models.ReviewStatusApproved).Pluck("rating", &movieReviewsRatings).Error
			if err != nil {
				return err
			}
			if lockedReview.Status != models.ReviewStatusApproved {
				if err := RecalculateRatingMovie(tx, lockedReview.Movie, lockedReview.Rating, movieReviewsRatings, false); err != nil {
					return err
				}
			}
			if err := tx.Model(lockedReview).Update("status", models.ReviewStatusApproved).Error; err != nil {
				return err
			}
			review = *lockedReview
			review.Status = models.ReviewStatusApproved
			return nil
		}

		reviewUpdates := map[string]interface{}{
			"status": req.Status,
		}
		if req.ReasonReject != nil {
			reviewUpdates["reject_reason"] = *req.ReasonReject
		}
		if err := tx.Model(lockedReview).Updates(reviewUpdates).Error; err != nil {
			return err
		}
		review = *lockedReview
		review.Status = req.Status
		if req.ReasonReject != nil {
			review.RejectReason = *req.ReasonReject
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func lockReviewForUpdate(tx *gorm.DB, reviewID uint) (*models.Review, error) {
	var review models.Review
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", reviewID).
		First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func lockMovieForUpdate(tx *gorm.DB, movieID uint) (*models.Movie, error) {
	var movie models.Movie
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", movieID).
		First(&movie).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}
