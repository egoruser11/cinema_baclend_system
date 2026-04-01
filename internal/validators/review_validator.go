package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ValidateCreateReview(c echo.Context, db *gorm.DB, req requests.ReviewCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	if req.Comment != "" {
		if !utils.ValidateCommentInReview(req.Comment) {
			errors["comment"] = "comment is invalid"
			return errors, false
		}
	} else {
		errors["comment"] = "input normal comment , please"
	}

	user := c.Get("user_data").(*models.User)
	var countOrdersUserInMovie int64
	var countReviewsUserInMovie int64
	var premieresMovieIds []int
	err := db.Model(&models.Premiere{}).Where("movie_id = ?", req.MovieID).Pluck("id", &premieresMovieIds).Error
	if err != nil {
		errors["premiere"] = err.Error()
		return errors, false
	}
	err = db.Model(&models.Order{}).Where("premiere_id IN ? AND user_id = ? AND status = ?", premieresMovieIds, user.ID, models.OrderPaid).Count(&countOrdersUserInMovie).Error
	if err != nil {
		errors["order"] = err.Error()
		return errors, false
	}
	if countOrdersUserInMovie == 0 {
		errors["order"] = "you don have orders with this movie"
	}
	err = db.Model(&models.Review{}).Where("movie_id = ? AND user_id = ?", req.MovieID, user.ID).Count(&countReviewsUserInMovie).Error
	if err != nil {
		errors["review"] = err.Error()
		return errors, false
	}
	if countReviewsUserInMovie != 0 {
		errors["review"] = "you don have reviews with this movie, you can update it , but not write new"
		return errors, false
	}
	return errors, len(errors) == 0
}

func ValidateApproveReview(db *gorm.DB, req requests.ReviewApproveRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	err := db.Model(&models.Review{}).Where("id = ?", req.ReviewID).Error
	if err != nil {
		errors["review"] = err.Error()
		return errors, false
	}
	if req.ReasonReject != nil {
		if *req.ReasonReject != "" {
			if req.Status != models.ReviewStatusApproved {
				errors["reject"] = "status approve and reason cannot be together , please correct data"
				return errors, false
			}
		}
		errors["reject"] = "input normal reject , please"
	}
	return errors, len(errors) == 0
}

func ValidateUpdateReview(c echo.Context, db *gorm.DB, req requests.ReviewUpdateRequest) (map[string]string, map[string]interface{}, bool) {
	errors := make(map[string]string)
	updates := map[string]interface{}{}
	user := c.Get("user_data").(*models.User)
	var review models.Review
	err := db.Model(&models.Review{}).Where("id = ?", req.ReviewID).First(&review).Error
	if err != nil {
		errors["review"] = err.Error()
		return errors, updates, false
	}
	if review.UserID != user.ID {
		errors["review"] = "it is not your review!"
		return errors, updates, false
	}
	if req.Rating == nil && req.Comment == nil {
		errors["data"] = "input normal data , please"
		return errors, updates, false
	}
	if req.Comment != nil {
		if *req.Comment != "" {
			if !utils.ValidateCommentInReview(*req.Comment) {
				errors["comment"] = "comment is invalid"
				return errors, updates, false
			}
			updates["comment"] = *req.Comment
		} else {
			errors["comment"] = "input normal comment , please"
		}
	}
	if req.Rating != nil {
		if *req.Rating == review.Rating && req.Comment == nil {
			updates["data"] = "input normal data , please "
		}
		updates["rating"] = *req.Rating
	}
	return errors, updates, len(errors) == 0
}
