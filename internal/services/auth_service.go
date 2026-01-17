package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/utils"
	"errors"
	"gorm.io/gorm"
	"time"
)

const lengthToken = 32

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Login(username, password string, deviceInfo string) (*models.User, string, error) {
	var user models.User

	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, "", errors.New("Invalid username or password")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", errors.New("Invalid username or password")
	}

	if user.Status != models.Active {
		return nil, "", errors.New("User is not active , contact support")
	}

	tokenString := utils.GenerateToken(lengthToken)

	token := models.Token{
		UserID:     user.ID,
		Token:      tokenString,
		ExpiresAt:  time.Now().Add(time.Hour * 24),
		CreatedAt:  time.Now(),
		DeviceInfo: deviceInfo,
	}

	if err := s.db.Create(&token).Error; err != nil {
		return nil, "", errors.New("Failed to create token")
	}

	return &user, tokenString, nil
}

func (s *AuthService) Logout(tokenString string, isFullLogout bool) error {
	var token models.Token
	result := s.db.Model(&token).Where("token = ?", tokenString).First(&token)

	if result.Error != nil {
		return errors.New("Failed to logout, token not exists")
	}

	if isFullLogout {
		result = s.db.Where("user_id = ?", token.UserID).Delete(&models.Token{})
		if result.Error != nil {
			return errors.New("Failed to logout of all devices")
		}
		return nil
	}

	result = s.db.Delete(&token)
	if result.Error != nil {
		return errors.New("Failed to logout")
	}

	return nil
}
