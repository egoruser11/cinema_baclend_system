package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/utils"
	"gorm.io/gorm"
)

func RecalculateRatingMovie(db *gorm.DB, movie models.Movie, rating int, movieReviewsRatings []float64) error {
	sumOfRatings := utils.GetSLiceSum(movieReviewsRatings) + float64(rating)
	newMovieRatingCount := movie.RatingCount + 1
	newMovieRating := sumOfRatings / float64(newMovieRatingCount)
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
