package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"gorm.io/gorm"
	"time"
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

func ValidateIndexPremiers(db *gorm.DB, req requests.PremiereIndexRequest) (map[string]string, map[string]interface{}, bool) {
	errors := make(map[string]string)
	filter := make(map[string]interface{})
	if req.Sort != nil {
		if (*req.Sort != "price" && *req.Sort != "booked_count" && *req.Sort != "total_seats" && *req.Sort != "start_time") || *req.Sort == "" {
			errors["sort"] = "Sort param is not correct"
		} else {
			filter["sort"] = req.Sort
		}
	}
	if req.IsDesc != nil {
		if *req.IsDesc {
			filter["order_type"] = "DESC"
		} else {
			filter["order_type"] = "ASC"
		}
	}

	if req.DayPremiere != nil {
		date, err := time.Parse("2006-01-02", *req.DayPremiere)
		if err != nil {
			errors["day_premiere"] = "DayPremiere format is invalid, use YYYY-MM-DD (e.g., 2026-01-26)"
		} else {
			filter["day_premiere"] = date
		}
	}

	if req.HourFrom != nil {
		date, err := time.Parse("15:04", *req.HourFrom)
		if err != nil {
			errors["hour_from"] = "HourFrom format is invalid, use HH:MM (e.g., 14:30)"
		} else {
			filter["hour_from"] = date
		}
	}

	if req.HourTo != nil {
		date, err := time.Parse("15:04", *req.HourTo)
		if err != nil {
			errors["hour_to"] = "HourTo format is invalid, use HH:MM (e.g., 22:00)"
		} else {
			filter["hour_to"] = date
		}
	}

	hourTo, isExistsHourTo := filter["hour_to"]
	hourFrom, isExistsHourFrom := filter["hour_from"]
	if isExistsHourFrom && isExistsHourTo {
		fromTime, okFrom := hourFrom.(time.Time)
		toTime, okTo := hourTo.(time.Time)
		if okFrom && okTo {
			if toTime.Before(fromTime) {
				errors["hours"] = "INVALID hours given"
				delete(filter, "hour_from")
				delete(filter, "hour_to")
			}
		}
	}

	if req.WeekDay != nil {
		if !(*req.WeekDay >= 0 && *req.WeekDay < 7) {
			errors["week_day"] = "WeekDay param is invalid , correct him"
		} else {
			filter["week_day"] = *req.WeekDay
		}
	}

	if req.PriceMin != nil {
		if *req.PriceMin <= 100 {
			errors["price_min"] = "Price min param is invalid , correct min"
		} else {
			filter["price_min"] = *req.PriceMin
		}
	}

	if req.PriceMax != nil {
		if *req.PriceMax <= 100 {
			errors["price_max"] = "Price max param is invalid , correct max"
		} else {
			filter["price_max"] = *req.PriceMax
		}
	}

	priceMin, isExistspriceMin := filter["price_min"]
	priceMax, isExistspriceMax := filter["price_max"]
	if isExistspriceMin && isExistspriceMax {
		priceMinInt, okFrom := priceMin.(float64)
		priceMaxInt, okTo := priceMax.(float64)
		if okFrom && okTo {
			if priceMinInt > priceMaxInt {
				errors["prices"] = "Price min param is invalid , correct min"
				delete(filter, "price_min")
				delete(filter, "price_max")
			}
		}
	}

	if req.Offset != nil {
		filter["offset"] = int(*req.Offset)
	} else {
		filter["offset"] = 0
	}
	if req.Limit != nil {
		filter["limit"] = int(*req.Limit)
	} else {
		filter["limit"] = 10
	}

	return errors, filter, len(errors) == 0
}
