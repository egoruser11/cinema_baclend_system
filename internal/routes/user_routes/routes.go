package routes

import (
	"cinema_backend_system/internal/handlers"
	"cinema_backend_system/internal/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupRoutes(e *echo.Echo, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler, db *gorm.DB) {
	e.POST("api/v1/auth/login", authHandler.Login)
	e.POST("api/v1/auth/register", authHandler.Register)

	protected := e.Group("/api/v1/")
	protected.Use(middleware.AuthMiddleware(db))
	{
		protected.GET("profile", userHandler.Profile)
	}

}
