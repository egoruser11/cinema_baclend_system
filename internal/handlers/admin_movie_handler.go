package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type AdminMovieHandler struct {
	adminMovieService *services.AdminMovieService
}

func NewAdminMovieHandler(adminMovieService *services.AdminMovieService) *AdminMovieHandler {
	return &AdminMovieHandler{adminMovieService: adminMovieService}
}

func (handler *AdminMovieHandler) Create(c echo.Context) error {
	var req requests.MovieCreateRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	movie, err := handler.adminMovieService.Create(req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.OK(c, *movie)
}
