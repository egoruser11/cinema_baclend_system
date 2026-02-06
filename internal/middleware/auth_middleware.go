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
				return utils.Unauthorized(c, "Uncorrect auth token header")
			}

			tokenString := parts[1]

			if tokenString == "" {
				return utils.Unauthorized(c, "empty token")
			}
			var token models.Token
			err := db.Preload("User").Where("token = ?", tokenString).First(&token).Error

			if err != nil {
				return utils.Unauthorized(c, "invalid token")
			}

			if token.ExpiresAt.Before(time.Now()) {
				db.Delete(&token)
				return utils.Unauthorized(c, "token is expired")
			}
			if token.User.Status != models.ActiveUserStatus {
				return utils.Unauthorized(c, "user is not active")
			}
			userIdStr := c.QueryParam("user_id")
			deviceInfo := c.QueryParam("device_info")
			if userIdStr == "" || deviceInfo == "" {
				return utils.Unauthorized(c, "invalid query params")
			}
			intUserId, _ := strconv.ParseUint(userIdStr, 10, 0)
			resUserId := uint(intUserId)
			if token.DeviceInfo != deviceInfo || token.User.ID != resUserId {
				return utils.Unauthorized(c, "!!!invalid query params!!!")
			}
			c.Set("user_data", &token.User)
			return next(c)
		}
	}
}
