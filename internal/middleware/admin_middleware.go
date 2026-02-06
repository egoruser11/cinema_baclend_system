package middleware

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user_data").(*models.User)
			if user.Role != models.RoleAdmin {
				return utils.Unauthorized(c, "Only admin condition")
			}
			return next(c)
		}
	}
}
