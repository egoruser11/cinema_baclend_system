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
	err := ReserveSeats(service.db, countSeatsInOrder, premiere, nil, seatsInOrder, false, nil)
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
		err := ReserveSeats(service.db, countSeatsInOrder, order.Premiere, nil, seatsInOrder, false, nil)
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

// Update обновляет места в заказе (только для неоплаченных заказов)
func (service *UserOrderService) Update(c echo.Context, req requests.OrderUpdateRequest) (*models.Order, error) {
	errorsValid, ok, newBookedSeats := validators.ValidateOrderUpdate(c, service.db, req)
	if !ok {
		if _, isExpired := errorsValid["orderExpired"]; isExpired {
			var order models.Order
			if err := service.db.Preload("Premiere").Where("id = ?", req.ID).First(&order).Error; err != nil {
				return nil, errors.New("order not found")
			}
			if err := UnReserveSeats(service.db, &order.Premiere, &order, nil, models.OrderDeleted); err != nil {
				return nil, err
			}
			return nil, errors.New("order expired, please create a new one")
		}
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}

	user := c.Get("user_data").(*models.User)

	var order models.Order
	if err := service.db.Preload("Premiere.Movie").Where("id = ?", *req.ID).First(&order).Error; err != nil {
		return nil, errors.New("order not found")
	}
	currentSeats := utils.ParseSeats(order.Seats)
	countSeatsOldOrder := 0
	for _, row := range currentSeats {
		countSeatsOldOrder += len(row)
	}

	if order.UserID != user.ID {
		return nil, errors.New("access denied: order does not belong to this user")
	}

	if req.Seats == nil {
		return &order, nil
	}

	seatsStr := ""
	countSeats := 0
	for row, seatsRow := range req.Seats {
		for _, seat := range seatsRow {
			seatsStr += fmt.Sprintf("%d - %d,", row, seat)
			countSeats++
		}
	}

	newTotalAmount := float64(order.Premiere.Price) * float64(countSeats)

	err := service.db.Transaction(func(tx *gorm.DB) error {
		if err := UnReserveForOneOrder(tx, order, order.Premiere, models.OrderPending); err != nil {
			return err
		}

		var rewritePremiere models.Premiere
		if err := tx.Model(&models.Premiere{}).Where("id = ?", order.Premiere.ID).First(&rewritePremiere).Error; err != nil {
			return err
		}
		if err := ReserveSeats(tx, countSeats, rewritePremiere, &order.Premiere, newBookedSeats, true, &countSeatsOldOrder); err != nil {
			return err
		}

		updates := map[string]interface{}{
			"seats":        seatsStr,
			"total_amount": newTotalAmount,
		}
		if err := tx.Model(&order).Updates(updates).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	if err = service.db.Preload("Premiere.Movie").First(&order, order.ID).Error; err != nil {
		return nil, err
	}

	return &order, nil
}

func (service *UserOrderService) Index(c echo.Context, req requests.OrderIndexRequest) ([]*models.Order, error) {
	errorsValid, filters, ok := validators.ValidateOrderIndex(c, service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}

	var orders []*models.Order

	query := service.db.Model(&models.Order{}).Preload("Premiere.Movie")

	if userID, exists := filters["user_id"]; exists {
		query = query.Where("user_id = ?", userID)
	}
	if status, exists := filters["status"]; exists {
		query = query.Where("status = ?", status)
	}
	if premiereID, exists := filters["premiere_id"]; exists {
		query = query.Where("premiere_id = ?", premiereID)
	}
	if dateFrom, exists := filters["date_from"]; exists {
		query = query.Where("created_at >= ?", dateFrom)
	}
	if dateTo, exists := filters["date_to"]; exists {
		query = query.Where("created_at <= ?", dateTo)
	}

	if limit, exists := filters["limit"]; exists {
		query = query.Limit(limit.(int))
	}
	if offset, exists := filters["offset"]; exists {
		query = query.Offset(offset.(int))
	}

	if sort, exists := filters["sort"]; exists {
		order := filters["order"].(string)
		query = query.Order(fmt.Sprintf("%s %s", sort, order))
	} else {
		query = query.Order("created_at DESC")
	}

	if err := query.Find(&orders).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (service *UserOrderService) Show(c echo.Context, req requests.OrderShowRequest) (*models.Order, error) {
	errorsValid, order := validators.ValidateOrderShow(c, service.db, req)
	if len(errorsValid) > 0 {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	return order, nil
}

func (service *UserOrderService) Delete(c echo.Context, req requests.OrderDeleteRequest) error {
	errorsValid, order, ok := validators.ValidateOrderDelete(c, service.db, req)
	if !ok {
		return errors.New(utils.InputErrorsValid(errorsValid))
	}

	if err := UnReserveSeats(service.db, &order.Premiere, order, nil, models.OrderDeleted); err != nil {
		return err
	}

	return nil
}

func getInfoFromOrder(order models.Order) (int, map[uint][]uint) {
	seatsInOrder := utils.ParseSeats(order.Seats)
	var countSeatsInOrder int
	for _, seats := range seatsInOrder {
		countSeatsInOrder += len(seats)
	}
	return countSeatsInOrder, seatsInOrder
}
