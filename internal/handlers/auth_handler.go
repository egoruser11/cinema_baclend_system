package handlers

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
	"strings"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (handler *AuthHandler) Login(c echo.Context) error {
	var req requests.LoginRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	if (req.Password == "" || req.DeviceInfo == "") || ((req.Username == "" && req.Email == "") ||
		(req.Username != "" && req.Email != "")) {
		return utils.BadRequest(c, "invalid credentials")
	}
	user, token, err := handler.authService.Login(req.Username, req.Password, req.Email, req.DeviceInfo)
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
	return utils.OK(c, response)
}

func (handler *AuthHandler) Register(c echo.Context) error {
	var req requests.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	registerResult, err := handler.authService.Register(req)
	if err != nil {
		errMsg := err.Error()
		return utils.BadRequest(c, errMsg)
	}

	c.Response().Header().Set("Authorization", "Bearer "+registerResult.Token)

	response := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       registerResult.User.ID,
			"username": registerResult.User.Username,
			"email":    registerResult.User.Email,
			"role":     registerResult.User.Role,
			"status":   registerResult.User.Status,
		},
		"token": registerResult.Token,
	}

	return utils.OK(c, response)
}

func (handler *AuthHandler) Logout(c echo.Context) error {
	var req requests.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	token := c.Get("token").(string)

	err := handler.authService.Logout(token, req)
	if err != nil {
		return utils.InternalServerError(c, "logout failed")
	}
	return utils.OK(c, "Logout ok!")
}

func (handler *AuthHandler) ResetPassword(c echo.Context) error {
	var reqReset requests.ResetPasswordRequest
	if err := c.Bind(&reqReset); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}
	user := c.Get("user_data").(*models.User)
	authHeader := c.Request().Header.Get("Authorization")
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return utils.Unauthorized(c, "Uncorrect auth header")
	}
	currentToken := parts[1]
	err := handler.authService.ResetPassword(currentToken, user, reqReset.OldPassword, reqReset.NewPassword)
	if err != nil {
		return utils.InternalServerError(c, "reset password failed")
	}
	return utils.OK(c, "reset password ok!")
}
