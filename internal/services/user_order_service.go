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
	"strconv"
	"strings"
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
	service.db.Model(&models.Premiere{}).Where("id = ?", req.PremiereID).First(&premiere)
	totalAmount := float64(int(premiere.Price) * countSeatsInOrder)
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
			err := UnReserveSeats(service.db, &order.Premiere, order, nil, models.OrderDeleted)
			if err != nil {
				return nil, err
			}
			return nil, errors.New(utils.InputErrorsValid(errorsValid))
		}
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var order *models.Order
	service.db.Preload("Premiere").Model(&models.Order{}).Where("id = ?", req.OrderID).First(&order)
	user := c.Get("user_data").(*models.User)
	if order.Status != models.OrderPending {
		countSeatsInOrder, seatsInOrder := getInfoFromOrder(*order)
		err := ReserveSeats(service.db, countSeatsInOrder, order.Premiere, seatsInOrder)
		if err != nil {
			return nil, err
		}
	}
	var sumToAddCoins float64
	sumToPay := order.TotalAmount
	var newUserCoinBalance float64
	if req.Coins != nil {
		sumToPay -= float64(*req.Coins)
		newUserCoinBalance = user.CoinBalance - *req.Coins
	} else {
		sumToAddCoins = order.TotalAmount * 0.1
		newUserCoinBalance = user.CoinBalance + sumToAddCoins
	}
	newUserMoneyBalance := user.MoneyBalance - sumToPay
	updatesUser := map[string]interface{}{
		"money_balance": newUserMoneyBalance,
		"coin_balance":  newUserCoinBalance,
	}
	var coinsToPlus *float64
	if sumToAddCoins > 0 {
		coinsToPlus = &sumToAddCoins
	} else {
		coinsToPlus = nil
	}

	updatesOrder := map[string]interface{}{
		"coins":         req.Coins,
		"status":        models.OrderPaid,
		"coins_to_plus": coinsToPlus,
	}
	err := service.db.Preload("Premiere.Movie").Model(&order).Where("id = ?", req.OrderID).Updates(updatesOrder).Error
	if err != nil {
		return nil, err
	}
	err = service.db.Model(&user).Where("id = ?", user.ID).Updates(updatesUser).Error
	if err != nil {
		return nil, err
	}
	operation := &models.Operation{
		UserID:  user.ID,
		OrderID: order.ID,
		Amount:  sumToPay,
		Type:    models.Purchase,
		Status:  models.OperationStatusPaid,
	}
	err = service.db.Create(operation).Error
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
	var newUserCoinBalance float64
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
			return nil, errors.New("You not refund order")
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
	err = service.db.Transaction(func(tx *gorm.DB) error {
		if err = UnReserveSeats(tx, &order.Premiere, &order, nil, models.OrderRefunded); err != nil {
			return err
		}
		if err = tx.Model(&user).Updates(updatesUser).Error; err != nil {
			return err
		}
		return nil
	})
	return &order, nil
}

func getInfoFromOrder(order models.Order) (int, map[uint][]uint) {
	seatsInOrder := parseSeats(order.Seats)
	var countSeatsInOrder int
	for _, seats := range seatsInOrder {
		countSeatsInOrder += len(seats)
	}
	return countSeatsInOrder, seatsInOrder
}

// index , show , delete , update , refund
func parseSeats(seatsStr string) map[uint][]uint {
	result := make(map[uint][]uint)
	seatsStr = strings.TrimSuffix(seatsStr, ",")
	if seatsStr == "" {
		return result
	}
	seatsArr := strings.Split(seatsStr, ",")

	for _, seat := range seatsArr {
		parts := strings.Split(strings.TrimSpace(seat), " - ")
		if len(parts) != 2 {
			continue
		}

		row, err := strconv.ParseUint(parts[0], 10, 0)
		if err != nil {
			continue
		}

		num, err := strconv.ParseUint(parts[1], 10, 0)
		if err != nil {
			continue
		}

		rowUint := uint(row)
		numUint := uint(num)

		exists := false
		for _, existing := range result[rowUint] {
			if existing == numUint {
				exists = true
				break
			}
		}

		if !exists {
			result[rowUint] = append(result[rowUint], numUint)
		}
	}

	return result
}
