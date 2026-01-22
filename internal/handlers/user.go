package handlers

import (
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	userService *services.UserService
}

type ProfileRequest struct {
	UserId     uint   `json:"user_id"`
	DeviceInfo string `json:"device_info"`
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (handler *UserHandler) Profile(c echo.Context) error {
	var req ProfileRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	user, err := handler.userService.Profile(req.UserId)
	if err != nil {
		return utils.BadRequest(c, "User not found")
	}
	data := map[string]interface{}{
		"userData": map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
		},
		"authData": map[string]interface{}{
			"deviceInfo": req.DeviceInfo,
			"userId":     user.ID,
		},
	}
	return utils.OK(c, data)
}
