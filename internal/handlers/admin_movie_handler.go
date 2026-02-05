package handlers

import "cinema_backend_system/internal/services"

type AdminMovieHandler struct {
	adminMovieService services.AdminMovieService
}

func NewAdminMovieHandler(adminMovieService services.AdminMovieService) *AdminMovieHandler {
	return &AdminMovieHandler{adminMovieService: adminMovieService}
}
