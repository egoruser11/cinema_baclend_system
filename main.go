package main

import (
	"cinema_backend_system/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"os"
)

const dbConfig = "host=localhost user=user password=password dbname=mydb port=5432 sslmode=disable"

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	dsn := dbConfig
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
	logger.Info("Сидеры и миграции запустились")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		logger.Info("Получен запрос на /")
		return c.String(http.StatusOK, "Привет от Echo и GORM!")
	})

	port := "localhost:8080"
	logger.Info("Запуск сервера", "port", port)
	if err := e.Start(port); err != nil {
		logger.Error("Ошибка сервера", "error", err)
	}
}
