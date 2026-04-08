package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/validators"
	"errors"
	"github.com/labstack/echo/v4"
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
	result := service.db.Where("id = ?", userId).First(&user)
	if result.Error != nil {
		return nil, errors.New("User not found")
	}
	return &user, nil
}

func (service *UserService) Update(c echo.Context, req requests.UpdateUserRequest) (*models.User, error) {
	errorsValid, updates, ok := validators.ValidateUpdateUser(c, req)
	if !ok {
		return nil, errors.New(errorsValid["error"])
	}
	var user models.User
	err := service.db.Model(&models.User{}).Where("id = ?", updates["id"]).Find(&user).Error
	if err != nil {
		return nil, errors.New("User not found")
	}
	if len(updates) == 1 {
		return &user, nil
	}
	err = service.db.Model(&user).Updates(updates).Error
	err = service.db.First(&user, user.ID).Error
	if err != nil {
		return nil, errors.New("Failed to load updated user")
	}
	return &user, nil
}
