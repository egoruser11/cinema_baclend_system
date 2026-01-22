package middleware

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strings"
	"time"
)

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

			var token models.Token

			err := db.
				Preload("User").
				Where("token = ?", tokenString).
				First(&token).Error

			if err != nil {
				return utils.Unauthorized(c, "Invalid token")
			}
			if token.ExpiresAt.Before(time.Now()) {
				db.Delete(&token)
				return utils.Unauthorized(c, "Token is expired")
			}
			if token.User.Status != models.Active {
				return utils.Unauthorized(c, "User is not active")
			}
			c.Set("user_id", token.UserID)
			c.Set("user_role", token.User.Role)
			c.Set("token", tokenString)
			c.Set("device_info", token.DeviceInfo)
			c.Set("user", &token.User)

			// 7. Передаем управление следующему обработчику
			return next(c)
		}

	}
}
