package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/utils"
	"gorm.io/gorm"
)

func RecalculateRatingMovie(db *gorm.DB, movie models.Movie, rating int, movieReviewsRatings []float64, isDeleteData bool) error {
	var sumOfRatings float64
	var newMovieRatingCount int
	var newMovieRating float64
	if isDeleteData {
		sumOfRatings = utils.GetSLiceSum(movieReviewsRatings) - float64(rating)
	} else {
		sumOfRatings = utils.GetSLiceSum(movieReviewsRatings) + float64(rating)
	}
	if isDeleteData {
		newMovieRatingCount = movie.RatingCount - 1
		if newMovieRatingCount <= 0 {
			newMovieRating = 0
			newMovieRatingCount = 0
			goto update
		}
	} else {
		newMovieRatingCount = movie.RatingCount + 1
	}
	newMovieRating = sumOfRatings / float64(newMovieRatingCount)
update:
	movieUpdateData := map[string]interface{}{
		"rating_count": newMovieRatingCount,
		"rating":       newMovieRating,
	}
	err := db.Model(&movie).Updates(movieUpdateData).Error
	if err != nil {
		return err
	}
	return nil
}
