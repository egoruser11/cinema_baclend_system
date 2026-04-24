package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"fmt"
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

func (handler *AdminMovieHandler) Update(c echo.Context) error {
	var req requests.MovieUpdateRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	movie, err := handler.adminMovieService.Update(req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.OK(c, *movie)
}

func (handler *AdminMovieHandler) Delete(c echo.Context) error {
	var req requests.MovieIdRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	err := handler.adminMovieService.Delete(req.Id)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.OK(c, fmt.Sprintf("Movie %d deleted", req.Id))
}

func (handler *AdminMovieHandler) Show(c echo.Context) error {
	var req requests.MovieIdRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	movie, err := handler.adminMovieService.Show(req.Id)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.OK(c, *movie)
}

func (handler *AdminMovieHandler) Index(c echo.Context) error {
	var req requests.MovieIndexRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	movies, err := handler.adminMovieService.Index(req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}

	return utils.OK(c, movies)
}
