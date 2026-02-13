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

func ValidateUpdatePremiere(db *gorm.DB, req requests.PremiereUpdateRequest) (map[string]string, map[string]interface{}, bool) {
	errors := make(map[string]string)
	updates := make(map[string]interface{})
	var countPremiere int64
	errPremiere := db.Find(&models.Premiere{}).Where("id = ?", req.Id).Count(&countPremiere).Error
	if errPremiere != nil {
		errors["premiere_id"] = "Premiere not found, correct premier id"
	}
	if req.MovieID != nil {
		var countMovie int64
		db.Find(&models.Movie{}).Where("id = ?", *req.MovieID).Count(&countMovie)
		if countMovie == 0 {
			errors["movie_id"] = "Movie not found, correct movie id"
		} else {
			updates["movie_id"] = *req.MovieID
		}
	}
	if req.Hall != nil {
		if *req.Hall == "" {
			errors["hall"] = "Hall is required"
		} else {
			updates["hall"] = *req.Hall
		}
	}
	if req.Price != nil {
		if *req.Price <= 100 {
			errors["price"] = "Price must be greater than 100"
		} else {
			updates["price"] = *req.Price
		}
	}
	if req.StartTime != nil {
		if req.EndTime.Before(*req.StartTime) || req.EndTime.Sub(*req.StartTime) < 50 {
			errors["start_time"] = "End_time must be greater than StartTime more 50 minutes"
		} else {
			updates["start_time"] = *req.StartTime
		}
	}
	if req.EndTime != nil {
		if req.EndTime.Before(*req.StartTime) || req.EndTime.Sub(*req.StartTime) < 50 {
			errors["end_time"] = "End_time must be greater than StartTime more 50 minutes"
		} else {
			updates["end_time"] = *req.EndTime
		}
	}
	return errors, updates, len(errors) == 0
}
