package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserPaymentMethodService struct {
	db *gorm.DB
}

func NewUserPaymentMethodService(db *gorm.DB) *UserPaymentMethodService {
	return &UserPaymentMethodService{db: db}
}

func (service *UserPaymentMethodService) Create(c echo.Context, req requests.PaymentMethodCreateRequest) (*models.PaymentMethod, error) {
	errorsValid, ok := validators.ValidateCreatePaymentMethod(req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	userId := c.Get("user_data").(*models.User).ID
	paymentMethod := &models.PaymentMethod{
		Type:      req.Type,
		Details:   req.Details,
		IsDefault: req.IsDefault,
		UserID:    userId,
	}

	if err := service.db.Create(paymentMethod).Error; err != nil {
		return nil, errors.New("Failed to create payment method")
	}

	return paymentMethod, nil

}

func (service *UserPaymentMethodService) Update(c echo.Context, req requests.PaymentMethodUpdateRequest) (*models.PaymentMethod, error) {
	errorsValid, updates, ok := validators.ValidateUpdatePaymentMethod(service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var paymentMethod models.PaymentMethod
	userId := c.Get("user_data").(*models.User).ID
	service.db.Model(&models.PaymentMethod{}).Where("id = ?", *req.Id).First(&paymentMethod)
	if paymentMethod.UserID != userId {
		return nil, errors.New("Payment method not found, please dont do this")
	}
	service.db.Model(&paymentMethod).Updates(updates)

	return &paymentMethod, nil
}
