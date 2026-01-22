package handlers

import (
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username   string `json:"username,omitempty"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,required"`
	DeviceInfo string `json:"device_info,required"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if (req.Password == "" || req.DeviceInfo == "") || ((req.Username == "" && req.Email == "") ||
		(req.Username != "" && req.Email != "")) {
		return utils.BadRequest(c, "invalid credentials")
	}
	user, token, err := h.authService.Login(req.Username, req.Email, req.Password, req.DeviceInfo)
	if err != nil {
		switch err.Error() {
		case "Invalid credentials":
			return utils.Unauthorized(c, "invalid credentials")
		case "User is not active":
			return utils.Unauthorized(c, "Account is not active")
		default:
			return utils.InternalServerError(c, "login failed")
		}
	}
	c.Response().Header().Set("Authorization", "Bearer "+token)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":          user.ID,
			"username":    user.Username,
			"email":       user.Email,
			"role":        user.Role,
			"status":      user.Status,
			"device_info": req.DeviceInfo,
		},
		"token": token,
	}
	//после логина перекидывает на Profile , отдаем device_info , user_id, токен создаем, не отдаем
	return utils.OK(c, response)
}
