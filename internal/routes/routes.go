package routes

import (
	"cinema_backend_system/internal/handlers"
	"cinema_backend_system/internal/middleware"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func SetupUserRoutes(e *echo.Echo, userHandler *handlers.UserHandler, db *gorm.DB) {
	protected := e.Group("/api/v1/")
	protected.Use(middleware.AuthMiddleware(db))
	{
		protected.GET("profile", userHandler.Profile)
	}
}

func SetupAuthRoutes(e *echo.Echo, authHandler *handlers.AuthHandler) {
	e.POST("api/v1/auth/login", authHandler.Login)
	e.POST("api/v1/auth/register", authHandler.Register)
}

//добавить adminMiddleware
func SetupAdminRoutes(e *echo.Echo, userHandler *handlers.UserHandler, db *gorm.DB) {
	protected := e.Group("/api/v1/")
	protected.Use(middleware.AuthMiddleware(db))
	{
		protected.GET("movie/create", userHandler.Profile)
	}
}
