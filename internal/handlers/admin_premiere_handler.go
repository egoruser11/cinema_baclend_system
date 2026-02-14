package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AdminPremiereHandler struct {
	adminPremiereService *services.AdminPremiereService
}

func NewAdminPremiereHandler(adminPremiereService *services.AdminPremiereService) *AdminPremiereHandler {
	return &AdminPremiereHandler{adminPremiereService: adminPremiereService}
}

func (handler *AdminPremiereHandler) Create(c echo.Context) error {
	var req requests.PremiereCreateRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	premiere, err := handler.adminPremiereService.Create(req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return c.JSON(http.StatusOK, premiere)
}

func (handler *AdminPremiereHandler) Update(c echo.Context) error {
	var req requests.PremiereUpdateRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	premiere, err := handler.adminPremiereService.Update(req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return c.JSON(http.StatusOK, premiere)
}

func (handler *AdminPremiereHandler) Index(c echo.Context) error {
	var req requests.PremiereIndexRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	premieres, err := handler.adminPremiereService.Index(req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return c.JSON(http.StatusOK, premieres)
}
