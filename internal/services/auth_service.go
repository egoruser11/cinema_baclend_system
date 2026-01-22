package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"gorm.io/gorm"
	"strings"
	"time"
)

const lengthToken = 32

type AuthService struct {
	db *gorm.DB
}

type RegisterResult struct {
	User  *models.User
	Token string
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) Login(username, password, email string, deviceInfo string) (*models.User, string, error) {
	var user models.User

	if username != "" {
		if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
			return nil, "", errors.New("Invalid credentials")
		}
	} else {
		if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
			return nil, "", errors.New("Invalid credentials")
		}
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", errors.New("Invalid credentials")
	}

	if user.Status != models.Active {
		return nil, "", errors.New("User is not active")
	}

	tokenString, err := s.CreateToken(user.ID, deviceInfo)

	if err != nil {
		return nil, "", err
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

func (s *AuthService) Register(req validators.RegisterRequest) (*RegisterResult, error) {

	if errorsMsg, ok := validators.ValidateRegister(s.db, req); !ok {
		var errorMsgs []string
		for field, msg := range errorsMsg {
			errorMsgs = append(errorMsgs, field+": "+msg)
		}
		return nil, errors.New(strings.Join(errorMsgs, "\n"))
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("Failed to hash password")
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Age:          req.Age,
		Role:         models.RoleUser,
		Status:       models.Active,
		MoneyBalance: 0.00,
		CoinBalance:  0,
	}
	err = s.db.Create(user).Error

	if err != nil {
		return nil, errors.New("Failed to create user")
	}

	tokenString, err := s.CreateToken(user.ID, req.DeviceInfo)
	if err != nil {
		return &RegisterResult{User: user}, err
	}
	return &RegisterResult{
		User:  user,
		Token: tokenString,
	}, nil
}

func (s *AuthService) CreateToken(UserId uint, deviceInfo string) (string, error) {
	tokenString := utils.GenerateToken(lengthToken)
	token := models.Token{
		UserID:     UserId,
		Token:      tokenString,
		ExpiresAt:  time.Now().Add(24 * time.Hour),
		CreatedAt:  time.Now(),
		DeviceInfo: deviceInfo,
	}
	err := s.db.Create(&token).Error
	if err != nil {
		return "", errors.New("Failed to create token")
	}
	return tokenString, nil
}
