package handlers

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler *UserHandler) Profile(c echo.Context) error {
	deviceInfo := c.QueryParam("device_info")

	user := c.Get("user_data").(*models.User)
	data := map[string]interface{}{
		"userData": map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		},
		"authData": map[string]interface{}{
			"deviceInfo": deviceInfo,
			"userId":     user.ID,
		},
	}
	return utils.OK(c, data)
}

func (handler *UserHandler) Update(c echo.Context) error {
	var req requests.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	user, err := handler.userService.Update(c, req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, user)
}
