package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"gorm.io/gorm"
	"net/url"
	"time"
	"unicode/utf8"
)

func ValidateMovie(db *gorm.DB, req requests.MovieCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	if len(req.GenreIDS) == 0 {
		errors["genres_ids"] = "Movie must have at least one id"
	} else {
		var countValidIds int64
		uniqueIDs := make(map[int]bool)
		for _, id := range req.GenreIDS {
			if id <= 0 {
				errors["genre_ids"] = "Genre ID must be positive"
				break
			}
			uniqueIDs[id] = true
		}
		db.Model(&models.Genre{}).Where("id IN ?", uniqueIDs).Count(&countValidIds)

		if countValidIds != int64(len(uniqueIDs)) {
			errors["genres_ids"] = "Something genre IDs not find in db , please input correct data"
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
