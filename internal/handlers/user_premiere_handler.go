package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserPremiereHandler struct {
	userPremiereService *services.UserPremiereService
}

func NewUserPremiereHandler(userPremiereService *services.UserPremiereService) *UserPremiereHandler {
	return &UserPremiereHandler{userPremiereService: userPremiereService}
}

func (handler *UserPremiereHandler) SeatMap(c echo.Context) error {
	var req requests.PremiereSeatMapRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	seatMap, err := handler.userPremiereService.SeatMap(req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.OK(c, seatMap)
}
