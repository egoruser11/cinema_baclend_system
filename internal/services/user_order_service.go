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
	"time"
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
		_, isOrderPaid := errorsValid["orderPaid"]
		var order *models.Order
		service.db.Preload("Premiere").Model(&models.Order{}).Where("id = ?", req.OrderID).First(&order)
		if isOrderExpired && !isOrderPaid {
			err := UnReserveSeats(service.db, &order.Premiere, order, nil, models.OrderPaid)
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
	var sumToAddCoins float64
	sumToPay := order.TotalAmount
	var newUserCoinBalance uint64
	if req.Coins != nil {
		sumToPay -= float64(*req.Coins)
		newUserCoinBalance = user.CoinBalance - *req.Coins
	} else {
		sumToAddCoins = order.TotalAmount * 0.1
		newUserCoinBalance = user.CoinBalance + uint64(sumToAddCoins)
	}
	newUserMoneyBalance := user.MoneyBalance - sumToPay
	updatesUser := map[string]interface{}{
		"money_balance": newUserMoneyBalance,
		"coin_balance":  newUserCoinBalance,
	}
	updatesOrder := map[string]interface{}{
		"coins":  req.Coins,
		"status": models.OrderPaid,
		"coins_to_add": func() *float64 {
			if sumToAddCoins > 0 {
				return &sumToAddCoins
			}
			return nil
		},
	}
	service.db.Preload("Premiere.Movie").Model(&models.Order{}).Where("id = ?", req.OrderID).Updates(updatesOrder)
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

func (service *UserOrderService) Refund(c echo.Context, req requests.OrderRefundedRequest) (*models.Order, error) {
	errorsValid, ok := validators.ValidateOrderRefund(c, service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	user := c.Get("user_data").(*models.User)
	var order models.Order
	err := service.db.Preload("Premiere.Movie").Model(&models.Order{}).Where("id = ?", req.ID).Find(&order).Error
	if err != nil {
		return nil, err
	}
	var newUserCoinBalance uint64
	var newUserMoneyBalance float64
	if order.Coins == nil {
		if time.Now().Add(2 * time.Hour).After(order.Premiere.StartTime) {
			newUserMoneyBalance = user.MoneyBalance + (order.TotalAmount / 2)
			newUserCoinBalance = user.CoinBalance + (*order.CoinsToPlus / 2)
		} else {
			newUserMoneyBalance = user.MoneyBalance + order.TotalAmount
			newUserCoinBalance = user.CoinBalance + *order.CoinsToPlus
		}
	} else {
		if user.CoinBalance < *order.CoinsToPlus {
			var orders []models.Order
			err = service.db.Model(&models.Order{}).Where("user_id = ? AND coins != ? AND status = ? AND created_at > ?",
				user.ID, nil, models.OrderPaid, order.CreatedAt).Find(&orders).Error
			if err != nil {
				return nil, err
			}
			if len(orders) == 0 {
				return nil, errors.New("ypu cannot refund order")
			}
			ids := []uint{}
			sumCoinsInOtherOrders := 0
			for _, curOrder := range orders {
				if curOrder.Coins != nil {
					sumCoinsInOtherOrders += int(*curOrder.Coins)
					ids = append(ids, curOrder.ID)
				}
				if sumCoinsInOtherOrders+int(user.CoinBalance) >= int(*order.CoinsToPlus) {
					break
				}
			}
			if sumCoinsInOtherOrders+int(user.CoinBalance) < int(*order.CoinsToPlus) {
				return nil, errors.New("Can not refund order")
			}
			var ordersRefund []models.Order
			err = service.db.Model(&models.Order{}).Where("id in ?", ids).Find(&ordersRefund).Error
			if err != nil {
				return nil, err
			}
			err = UnReserveSeats(service.db, &order.Premiere, nil, ordersRefund, models.OrderCanceled)
			if err != nil {
				return nil, err
			}
		} else {
			newUserCoinBalance = user.CoinBalance - *order.Coins
		}
		if time.Now().Add(2 * time.Hour).After(order.Premiere.StartTime) {
			newUserMoneyBalance = user.MoneyBalance + (order.TotalAmount / 2)
		} else {
			newUserMoneyBalance = user.MoneyBalance + order.TotalAmount
		}
	}
	updatesUser := map[string]interface{}{
		"money_balance": newUserMoneyBalance,
		"coin_balance":  newUserCoinBalance,
	}
	err = UnReserveSeats(service.db, &order.Premiere, &order, nil, models.OrderRefunded)
	if err != nil {
		return nil, err
	}
	service.db.Model(&user).Updates(updatesUser)
	return &order, nil
}

//index , show , delete , update , refund
