package middleware

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

// AuthMiddleware проверяет токен, user_id и device_info
func AuthMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return utils.Unauthorized(c, "Missing Authorization header")
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return utils.Unauthorized(c, "Invalid Authorization header")
			}
			tokenString := parts[1]

			userIDStr := c.Request().Header.Get("X-User-ID")
			if userIDStr == "" {
				return utils.Unauthorized(c, "Missing X-User-ID header")
			}

			userID, err := strconv.ParseUint(userIDStr, 10, 32)
			if err != nil {
				return utils.BadRequest(c, "Invalid user_id format")
			}

			deviceInfo := c.Request().Header.Get("X-Device-Info")
			if deviceInfo == "" {
				return utils.Unauthorized(c, "Missing X-Device-Info header")
			}

			var token models.Token
			err = db.
				Preload("User").
				Where("token = ? AND user_id = ? AND device_info = ?",
					tokenString, uint(userID), deviceInfo).
				First(&token).Error

			if err != nil {
				return utils.Unauthorized(c, "Invalid authentication data")
			}

			if token.ExpiresAt.Before(time.Now()) {
				db.Delete(&token)
				return utils.Unauthorized(c, "Token is expired")
			}

			if token.User.Status != models.ActiveUserStatus {
				return utils.Unauthorized(c, "User is not active")
			}

			c.Set("user_id", token.UserID)
			c.Set("user_role", token.User.Role)
			c.Set("token", tokenString)
			c.Set("device_info", token.DeviceInfo)
			c.Set("user", &token.User)

			return next(c)
		}
	}
}
