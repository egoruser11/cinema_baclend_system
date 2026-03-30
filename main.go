package main

import (
	"cinema_backend_system/internal/handlers"
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/routes"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	dsn := utils.DbConfig
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Ошибка подключения к БД", "error", err)
		os.Exit(1)
	}

	if err := models.SetupDatabase(db); err != nil {
		logger.Error("Ошибка настройки БД", "error", err)
		os.Exit(1)
	}
	logger.Info("База данных настроена")

	authService := services.NewAuthService(db)
	userService := services.NewUserService(db)
	adminMovieService := services.NewAdminMovieService(db)
	adminPremiereService := services.NewAdminPremiereService(db)
	userPaymentMethodService := services.NewUserPaymentMethodService(db)
	userOrderService := services.NewUserOrderService(db)
	userReviewService := services.NewUserReviewService(db)

	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	adminMovieHandler := handlers.NewAdminMovieHandler(adminMovieService)
	adminPremiereHandler := handlers.NewAdminPremiereHandler(adminPremiereService)
	userPaymentMethodHandler := handlers.NewUserPaymentMethodHandler(userPaymentMethodService)
	userOrderHandler := handlers.NewUserOrderHandler(userOrderService)
	userReviewHandler := handlers.NewUserReviewHandler(userReviewService)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-User-ID", "X-Device-Info"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	routes.SetupAuthRoutes(e, authHandler, db)
	routes.SetupUserRoutes(e, userHandler, db)
	routes.SetupAdminMovieRoutes(e, adminMovieHandler, db)
	routes.SetupAdminPremiereRoutes(e, adminPremiereHandler, db)
	routes.SetupUserPaymentMethodRoutes(e, userPaymentMethodHandler, db)
	routes.SetupUserOrderRoutes(e, userOrderHandler, db)
	routes.SetupUserReviewRoutes(e, userReviewHandler, db)

	port := "localhost:8080"
	logger.Info("Запуск сервера", "port", port)
	if err := e.Start(port); err != nil {
		logger.Error("Ошибка сервера", "error", err)
	}
}
