package services

import (
	"cinema_backend_system/internal/models"
	"errors"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (service *UserService) Profile(userId uint) (*models.User, error) {
	var user models.User

	err := service.db.Where("id = ?", userId).First(&user)
	if err != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}
