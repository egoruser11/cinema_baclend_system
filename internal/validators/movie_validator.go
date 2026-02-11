package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"fmt"
	"gorm.io/gorm"
	"net/url"
	"time"
	"unicode/utf8"
)

func ValidateCreateMovie(db *gorm.DB, req requests.MovieCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	if len(req.GenreIDS) == 0 {
		errors["genres_ids"] = "Movie must have at least one id"
	} else {
		var countValidIds int64
		db.Model(&models.Genre{}).Where("id IN ?", req.GenreIDS).Count(&countValidIds)
		uniqueSet := make(map[int]bool)
		for _, id := range req.GenreIDS {
			if id <= 0 {
				errors["genre_ids"] = "Genre ID must be positive"
				break
			}
			uniqueSet[id] = true
		}
		if countValidIds != int64(len(uniqueSet)) {
			errors["genres_ids"] = "Some genre IDs not found in database"
		}

	}

	if req.Title == "" {
		errors["title"] = "Title is required"
	} else if utf8.RuneCountInString(req.Title) < 2 {
		errors["title"] = "Title must be at least 2 characters"
	} else if utf8.RuneCountInString(req.Title) > 255 {
		errors["title"] = "Title must be less than 255 characters"
	}

	if req.Description == "" {
		errors["description"] = "Description is required"
	} else if utf8.RuneCountInString(req.Description) < 10 {
		errors["description"] = "Description must be at least 10 characters"
	}

	if req.Duration < 60 {
		errors["duration"] = "Duration must be at least 60 minutes"
	} else if req.Duration > 200 {
		errors["duration"] = "Duration must be less than 200 minutes"
	}

	validAgeRatings := map[models.AgeRating]bool{
		models.AgeRatingG:    true,
		models.AgeRatingPG:   true,
		models.AgeRatingPG13: true,
		models.AgeRatingR:    true,
		models.AgeRatingNC17: true,
	}
	if req.AgeRating == "" {
		errors["age_rating"] = "Age rating is required"
	} else if !validAgeRatings[req.AgeRating] {
		errors["age_rating"] = "Invalid age rating"
	}

	if req.PosterURL != "" && !isValidURL(req.PosterURL) {
		errors["poster_url"] = "Invalid poster URL"
	}
	if req.TrailerURL != "" && !isValidURL(req.TrailerURL) {
		errors["trailer_url"] = "Invalid trailer URL"
	}

	if req.ReleaseDate.IsZero() {
		errors["release_date"] = "Release date is required"
	} else if req.ReleaseDate.After(time.Now().AddDate(10, 0, 0)) {
		errors["release_date"] = "Release date cannot be more than 10 years in the future"
	} else if req.ReleaseDate.Before(time.Now().AddDate(-100, 0, 0)) {
		errors["release_date"] = "Release date cannot be more than 100 years in the past"
	}

	return errors, len(errors) == 0
}

func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func ValidateUpdateMovie(db *gorm.DB, req requests.MovieUpdateRequest) (map[string]string, map[string]interface{},
	[]int, bool) {
	errors := make(map[string]string)
	updates := make(map[string]interface{})
	var idGenres []int
	if len(req.GenreIDS) != 0 {
		var countValidIds int64
		db.Model(&models.Genre{}).Where("id IN ?", req.GenreIDS).Count(&countValidIds)
		uniqueSet := make(map[int]bool)
		for _, id := range req.GenreIDS {
			if id <= 0 {
				errors["genre_ids"] = "Genre ID must be positive"
				break
			}
			uniqueSet[id] = true
		}
		if countValidIds != int64(len(uniqueSet)) {
			errors["genres_ids"] = "Some genre IDs not found in database"
		} else {
			idGenres = req.GenreIDS
		}

	}

	if req.Title != nil {
		if *req.Title == "" {
			errors["title"] = "Title cannot be empty"
		} else if utf8.RuneCountInString(*req.Title) < 2 {
			errors["title"] = "Title must be at least 2 characters"
		} else if utf8.RuneCountInString(*req.Title) > 255 {
			errors["title"] = "Title must be less than 255 characters"
		} else {
			updates["title"] = *req.Title
		}
	}
	if req.Description != nil {
		if *req.Description == "" {
			errors["description"] = "Description is required"
		} else if utf8.RuneCountInString(*req.Description) < 10 {
			errors["description"] = "Description must be at least 10 characters"
		} else {
			updates["description"] = *req.Description
		}
	}

	if req.Duration != nil {
		if *req.Duration < 60 {
			errors["duration"] = "Duration must be at least 60 minutes"
		} else if *req.Duration > 200 {
			errors["duration"] = "Duration must be less than 200 minutes"
		} else {
			updates["duration"] = *req.Duration
		}
	}

	if req.AgeRating != nil {
		validAgeRatings := map[models.AgeRating]bool{
			models.AgeRatingG:    true,
			models.AgeRatingPG:   true,
			models.AgeRatingPG13: true,
			models.AgeRatingR:    true,
			models.AgeRatingNC17: true,
		}
		if !validAgeRatings[*req.AgeRating] {
			errors["age_rating"] = "Invalid age rating"
		} else {
			updates["age_rating"] = *req.AgeRating
		}
	}

	if req.ReleaseDate != nil {
		if req.ReleaseDate.After(time.Now().AddDate(10, 0, 0)) {
			errors["release_date"] = "Release date cannot be more than 10 years in the future"
		} else if req.ReleaseDate.Before(time.Now().AddDate(-100, 0, 0)) {
			errors["release_date"] = "Release date cannot be more than 100 years in the past"
		} else {
			updates["release_date"] = *req.ReleaseDate
		}
	}

	return errors, updates, idGenres, len(errors) == 0
}

func ValidateIndexMovie(req requests.MovieIndexRequest) (map[string]string, map[string]string) {
	filters := make(map[string]string)
	errors := make(map[string]string)
	if req.Search != nil {
		if *req.Search != "" {
			filters["search"] = *req.Search
		}
	}
	if req.Sort == nil && req.IsDesc != nil || req.Sort != nil && req.IsDesc == nil {
		errors["sort"] = "input second field , please"
	}
	if req.Sort != nil {
		if *req.Sort == "" && req.IsDesc == nil {
			errors["sort"] = "input second field , please"
		}
		filters["sort"] = *req.Sort
		filters["order_type"] = GetOrderType(*req.IsDesc)
	}
	if req.Offset != nil {
		filters["offset"] = fmt.Sprintf("%d", *req.Offset)
	} else {
		filters["offset"] = "0"
	}
	if req.Limit != nil {
		filters["limit"] = fmt.Sprintf("%d", *req.Limit)
	} else {
		filters["limit"] = "10"
	}

	return errors, filters
}

func GetOrderType(isDesc bool) string {
	if isDesc {
		return "DESC"
	} else {
		return "ASC"
	}
}
