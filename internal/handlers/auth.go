package handlers

import (
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"github.com/labstack/echo/v4"
	"strings"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	DeviceInfo string `json:"device_info"`
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
	user, token, err := h.authService.Login(req.Username, req.Password, req.Email, req.DeviceInfo)
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

func (h *AuthHandler) Register(c echo.Context) error {
	var req validators.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, "invalid request body")
	}

	registerResult, err := h.authService.Register(req)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "Username is required"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Username must be at least"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Username must be less than"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Username can only contain"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Email is required"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Invalid email"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Password is required"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Password must be at least"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Password must be less than"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Age must be at least"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Age must be less than"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Device info is required"):
			return utils.BadRequest(c, errMsg)
		case strings.Contains(errMsg, "Username or email already exists"):
			return utils.BadRequest(c, errMsg)
		case errMsg == "Failed to hash password":
			return utils.BadRequest(c, "Failed to hash password")
		case errMsg == "Failed to create user":
			return utils.BadRequest(c, "Failed to create user")
		default:
			return utils.InternalServerError(c, "registration failed: "+errMsg)
		}
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
