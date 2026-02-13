package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"gorm.io/gorm"
)

func ValidateCreatePremiere(db *gorm.DB, req requests.PremiereCreateRequest) (map[string]string, bool) {

	errors := make(map[string]string)
	if req.Hall == "" {
		errors["hall"] = "Hall is required"
	}
	if req.Price <= 100 {
		errors["price"] = "Price must be greater than 100"
	}
	var movie models.Movie
	err := db.Where("id = ?", req.MovieID).First(&movie).Error
	if err != nil {
		errors["movie_id"] = "Movie not found, correct movie id"
	}
	if req.EndTime.Before(req.StartTime) || req.EndTime.Sub(req.StartTime) < 50 {
		errors["end_time"] = "End_time must be greater than StartTime more 50 minutes"
	}

	return errors, len(errors) == 0
}
