package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedAdmin(db *gorm.DB) error {
	// проверяем, существует ли админ
	var count int64
	db.Model(&User{}).Where("role = ?", RoleAdmin).Count(&count)

	if count > 0 {
		return nil // Админ уже есть
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash admin password")
	}

	admin := &User{
		Username:     "admin",
		Email:        "admin@cinema.com",
		PasswordHash: string(hashedPassword),
		Age:          30,
		Role:         RoleAdmin,
		Status:       Active,
		MoneyBalance: 10000.00,
		CoinBalance:  1000,
	}

	if err := db.Create(admin).Error; err != nil {
		return errors.New("failed to create admin user")
	}

	return nil
}

func SeedGenres(db *gorm.DB) error {
	genres := []string{
		"Драма", "Комедия", "Боевик", "Фантастика", "Ужасы",
		"Мелодрама", "Триллер", "Детектив", "Фэнтези", "Приключения",
		"Анимация", "Документальный", "Исторический", "Мюзикл",
	}

	for _, name := range genres {
		genre := &Genre{
			Name: name,
		}

		// Создаем если не существует
		db.Where(Genre{Name: name}).FirstOrCreate(genre)
	}

	return nil
}
