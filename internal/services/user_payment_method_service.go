package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"strings"
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
		return nil, errors.New("Payment method not found, please dont check others mayments methods!")
	}
	service.db.Model(&paymentMethod).Updates(updates)

	return &paymentMethod, nil
}

func (service *UserPaymentMethodService) Index(c echo.Context, req requests.PaymentMethodIndexRequest) ([]*models.PaymentMethod, error) {
	errorsValid, filters, ok := validators.ValidateIndexPaymentMethod(req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	userId := c.Get("user_data").(*models.User).ID
	var paymentMethods []*models.PaymentMethod
	query := service.db.Model(&models.PaymentMethod{}).Where("user_id = ?", userId)
	var count int64
	query.Count(&count)
	if count == 0 {
		return nil, errors.New("Payment method not found")
	}
	if len(filters) > 0 {
		sort, existsSort := filters["sort"]
		if existsSort {
			if sort == "active" {
				query.Where("is_active = ?", true)
			} else {
				query.Where("is_active = ?", false)
			}
		}
		search, existsSearch := filters["search"]
		if existsSearch {
			searchString := search.(string)
			resSearch := "%" + strings.ToLower(searchString) + "%"
			query.Where("LOWER(details) LIKE ?", resSearch)
		}
	}
	limit := filters["limit"].(int)
	offset := filters["offset"].(int)
	query.Limit(limit).Offset(offset).Find(&paymentMethods)
	return paymentMethods, nil
}

func (service UserPaymentMethodService) Delete(c echo.Context, req requests.PaymentMethodDeleteRequest) ([]uint, error) {
	var userPaymentMethods []models.PaymentMethod
	userId := c.Get("user_data").(*models.User).ID
	if len(req.Ids) == 0 {
		return nil, errors.New("No payment method to delete , please enter same 1 value")
	}
	err := service.db.Model(&models.PaymentMethod{}).Where("user_id = ?", userId).Find(&userPaymentMethods)
	if err.Error != nil {
		return nil, errors.New("No , she is not yours to torment")
	}
	if len(userPaymentMethods) == 0 {
		return nil, errors.New("Payment methodx not found , you can not delete")
	}
	deleteIds := []uint{}
	for _, paymentMethod := range userPaymentMethods {
		for _, id := range req.Ids {
			if *id == paymentMethod.ID {
				deleteIds = append(deleteIds, *id)
			}
		}
	}
	if len(deleteIds) > 0 {
		return deleteIds, nil
	}
	return nil, errors.New("Your input not your payment method, dont do this!")
}

func (service UserPaymentMethodService) Show(c echo.Context, req requests.PaymentMethodIdRequest) (*models.PaymentMethod, error) {
	var paymentMethods []models.PaymentMethod
	userId := c.Get("user_data").(*models.User).ID
	err := service.db.Model(&models.PaymentMethod{}).Where("user_id = ?", userId).Find(&paymentMethods).Error
	if err != nil {
		return nil, errors.New("No , she is not yours to torment")
	}
	for i := 0; i < len(paymentMethods); i++ {
		if paymentMethods[i].ID == req.Id {
			return &paymentMethods[i], nil
		}
	}
	return nil, errors.New("Your input not your payment method, dont do this!")
}
