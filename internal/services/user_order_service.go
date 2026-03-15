package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserOrderService struct {
	db *gorm.DB
}

func NewUserOrderService(db *gorm.DB) *UserOrderService {
	return &UserOrderService{db: db}
}

func (service *UserOrderService) Create(c echo.Context, req requests.OrderCreateRequest) (*models.Order, error) {
	errorsValid, ok := validators.ValidateCreateOrder(service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	userId := c.Get("user_data").(*models.User).ID
	seatsInOrder := req.Seats
	stringFormatSeats := ""
	var countSeatsInOrder int
	for row, seats := range seatsInOrder {
		for _, seat := range seats {
			stringFormatSeats += fmt.Sprintf("%d - %d,", row, seat)

		}
		countSeatsInOrder += len(seats)
	}
	var premiere models.Premiere
	totalAmount := float64(int(premiere.Price) * countSeatsInOrder)
	service.db.Model(&models.Premiere{}).Where("id = ?", req.PremiereID).First(&premiere)
	//Создать заказ со статусом  ожидание , забронировать места на премьере.
	order := &models.Order{
		UserID:      userId,
		PremiereID:  req.PremiereID,
		Seats:       stringFormatSeats,
		TotalAmount: totalAmount,
		Status:      models.OrderPending,
	}
	if err := service.db.Create(order).Error; err != nil {
		return nil, err
	}
	if err := service.db.Preload("Premiere.Movie").First(order, order.ID).Error; err != nil {
		return nil, err
	}
	err := ReserveSeats(service.db, countSeatsInOrder, premiere, seatsInOrder)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (service *UserOrderService) Paid(c echo.Context, req requests.OrderPaidRequest) (*models.Order, error) {
	errorsValid, ok := validators.ValidatePaidOrder(c, service.db, req)
	if !ok {
		_, isOrderExpired := errorsValid["orderExpired"]
		if isOrderExpired {
			err := UnReserveSeats(service.db, models.Premiere{}, models.Order{})
			if err != nil {
				return nil, err
			}
			return nil, errors.New(utils.InputErrorsValid(errorsValid))
		}
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var order *models.Order
	service.db.Model(&models.Order{}).Where("id = ?", req.OrderID).First(&order)
	user := c.Get("user_data").(*models.User)
	sumToPay := order.TotalAmount
	var newUserCoinBalance uint64
	if req.Coins != nil {
		sumToPay -= float64(*req.Coins)
		newUserCoinBalance = user.CoinBalance - *req.Coins
	} else {
		sumToAddCoins := order.TotalAmount * 0.1
		newUserCoinBalance = user.CoinBalance + uint64(sumToAddCoins)
	}
	newUserMoneyBalance := user.MoneyBalance - sumToPay
	updatesUser := map[string]interface{}{
		"money_balance": newUserMoneyBalance,
		"coin_balance":  newUserCoinBalance,
	}
	updatesOrder := map[string]interface{}{
		"coins": func() interface{} {
			if req.Coins != nil {
				return *req.Coins
			}
			return nil
		}(),
		"status": models.OrderPaid,
	}
	service.db.Model(&models.Order{}).Where("id = ?", req.OrderID).Updates(updatesOrder)
	service.db.Model(&models.User{}).Where("id = ?", user.ID).Updates(updatesUser)
	operation := &models.Operation{
		UserID:  user.ID,
		OrderID: order.ID,
		Amount:  sumToPay,
		Type:    models.Purchase,
		Status:  models.OperationStatusPaid,
	}
	err := service.db.Create(operation).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}
