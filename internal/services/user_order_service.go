package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"cinema_backend_system/responses"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	var order models.Order
	err := service.db.Transaction(func(tx *gorm.DB) error {
		premiere, err := lockPremiereForUpdate(tx, req.PremiereID)
		if err != nil {
			return err
		}
		if err := ensureSeatsAvailable(premiere, seatsInOrder); err != nil {
			return err
		}

		totalAmount := premiere.Price * float64(countSeatsInOrder)
		order = models.Order{
			UserID:      userId,
			PremiereID:  req.PremiereID,
			Seats:       stringFormatSeats,
			TotalAmount: totalAmount,
			Status:      models.OrderPending,
		}
		if err := tx.Create(&order).Error; err != nil {
			return err
		}
		if err := ReserveSeats(tx, countSeatsInOrder, *premiere, nil, seatsInOrder, false, nil); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := service.db.Preload("Premiere.Movie").First(&order, order.ID).Error; err != nil {
		return nil, err
	}
	return &order, nil
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
	authUser := c.Get("user_data").(*models.User)
	var order models.Order
	err := service.db.Transaction(func(tx *gorm.DB) error {
		lockedUser, err := lockUserForUpdate(tx, authUser.ID)
		if err != nil {
			return err
		}

		lockedOrder, err := lockOrderForUpdate(tx, req.OrderID)
		if err != nil {
			return err
		}
		if lockedOrder.UserID != lockedUser.ID {
			return errors.New("access denied: order does not belong to this user")
		}
		if lockedOrder.Status != models.OrderPending {
			return errors.New("order is not pending")
		}

		var sumToAddCoins float64
		sumToPay := lockedOrder.TotalAmount
		var newUserCoinBalance float64
		if req.Coins != nil {
			sumToPay -= *req.Coins
			newUserCoinBalance = lockedUser.CoinBalance - *req.Coins
		} else {
			sumToAddCoins = lockedOrder.TotalAmount * 0.1
			newUserCoinBalance = lockedUser.CoinBalance + sumToAddCoins
		}

		newUserMoneyBalance := lockedUser.MoneyBalance - sumToPay
		updatesUser := map[string]interface{}{
			"money_balance": newUserMoneyBalance,
			"coin_balance":  newUserCoinBalance,
		}

		var coinsToPlus *float64
		if sumToAddCoins > 0 {
			coinsToPlus = &sumToAddCoins
		}

		updatesOrder := map[string]interface{}{
			"coins":         req.Coins,
			"status":        models.OrderPaid,
			"coins_to_plus": coinsToPlus,
		}
		if err := tx.Model(lockedOrder).Updates(updatesOrder).Error; err != nil {
			return err
		}
		if err := tx.Model(lockedUser).Updates(updatesUser).Error; err != nil {
			return err
		}

		operation := &models.Operation{
			UserID:  lockedUser.ID,
			OrderID: lockedOrder.ID,
			Amount:  sumToPay,
			Type:    models.Purchase,
			Status:  models.OperationStatusPaid,
		}
		if err := tx.Create(operation).Error; err != nil {
			return err
		}

		order = *lockedOrder
		order.Status = models.OrderPaid
		order.Coins = req.Coins
		order.CoinsToPlus = coinsToPlus
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := service.db.Preload("Premiere.Movie").First(&order, order.ID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (service *UserOrderService) Refund(c echo.Context, req requests.OrderRefundedRequest) (*models.Order, error) {
	errorsValid, ok := validators.ValidateOrderRefund(c, service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var order models.Order
	authUser := c.Get("user_data").(*models.User)
	err := service.db.Transaction(func(tx *gorm.DB) error {
		lockedUser, err := lockUserForUpdate(tx, authUser.ID)
		if err != nil {
			return err
		}
		lockedOrder, err := lockOrderForUpdate(tx, req.ID)
		if err != nil {
			return err
		}
		if lockedOrder.UserID != lockedUser.ID {
			return errors.New("access denied: order does not belong to this user")
		}
		if lockedOrder.Status != models.OrderPaid {
			return errors.New("order is not paid")
		}

		lockedPremiere, err := lockPremiereForUpdate(tx, lockedOrder.PremiereID)
		if err != nil {
			return err
		}
		lockedOrder.Premiere = *lockedPremiere

		var newUserCoinBalance float64
		var newUserMoneyBalance float64
		if lockedOrder.Coins == nil {
			if time.Now().Add(2 * time.Hour).After(lockedOrder.Premiere.StartTime) {
				newUserMoneyBalance = lockedUser.MoneyBalance + (lockedOrder.TotalAmount / 2)
				newUserCoinBalance = lockedUser.CoinBalance + (*lockedOrder.CoinsToPlus / 2)
			} else {
				newUserMoneyBalance = lockedUser.MoneyBalance + lockedOrder.TotalAmount
				newUserCoinBalance = lockedUser.CoinBalance + *lockedOrder.CoinsToPlus
			}
		} else {
			if lockedOrder.CoinsToPlus != nil && lockedUser.CoinBalance < *lockedOrder.CoinsToPlus {
				return errors.New("You not refund order")
			}
			newUserCoinBalance = lockedUser.CoinBalance - *lockedOrder.Coins
			if time.Now().Add(2 * time.Hour).After(lockedOrder.Premiere.StartTime) {
				newUserMoneyBalance = lockedUser.MoneyBalance + (lockedOrder.TotalAmount / 2)
			} else {
				newUserMoneyBalance = lockedUser.MoneyBalance + lockedOrder.TotalAmount
			}
		}

		updatesUser := map[string]interface{}{
			"money_balance": newUserMoneyBalance,
			"coin_balance":  newUserCoinBalance,
		}
		if err := UnReserveSeats(tx, &lockedOrder.Premiere, lockedOrder, nil, models.OrderRefunded); err != nil {
			return err
		}
		if err := tx.Model(lockedUser).Updates(updatesUser).Error; err != nil {
			return err
		}

		order = *lockedOrder
		order.Status = models.OrderRefunded
		order.Premiere.Movie = models.Movie{}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := service.db.Preload("Premiere.Movie").First(&order, order.ID).Error; err != nil {
		return nil, err
	}
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
		lockedPremiere, err := lockPremiereForUpdate(tx, order.Premiere.ID)
		if err != nil {
			return err
		}

		order.Premiere = *lockedPremiere

		if err := UnReserveForOneOrder(tx, order, order.Premiere, models.OrderPending); err != nil {
			return err
		}

		var rewritePremiere models.Premiere
		if err := tx.Model(&models.Premiere{}).Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", order.Premiere.ID).First(&rewritePremiere).Error; err != nil {
			return err
		}
		if err := ensureSeatsAvailable(&rewritePremiere, newBookedSeats); err != nil {
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

func lockPremiereForUpdate(tx *gorm.DB, premiereID uint) (*models.Premiere, error) {
	var premiere models.Premiere
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", premiereID).
		First(&premiere).Error
	if err != nil {
		return nil, err
	}
	return &premiere, nil
}

func lockOrderForUpdate(tx *gorm.DB, orderID uint) (*models.Order, error) {
	var order models.Order
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Preload("Premiere").
		Where("id = ?", orderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func lockUserForUpdate(tx *gorm.DB, userID uint) (*models.User, error) {
	var user models.User
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func ensureSeatsAvailable(premiere *models.Premiere, seats map[uint][]uint) error {
	if len(seats) == 0 {
		return errors.New("seats are required")
	}

	bookedSeats := map[int][]int{}
	if len(premiere.BookedSeats) > 0 {
		if err := json.Unmarshal(premiere.BookedSeats, &bookedSeats); err != nil {
			return errors.New("failed to read booked seats")
		}
	}

	seenSeats := make(map[string]struct{})
	for row, seatsInRow := range seats {
		if row == 0 || row > premiere.Rows {
			return fmt.Errorf("row %d is invalid", row)
		}
		for _, seat := range seatsInRow {
			if seat == 0 || seat > premiere.SeatsPerRow {
				return fmt.Errorf("seat %d in row %d is invalid", seat, row)
			}

			key := fmt.Sprintf("%d-%d", row, seat)
			if _, exists := seenSeats[key]; exists {
				return fmt.Errorf("seat %d in row %d is duplicated", seat, row)
			}
			seenSeats[key] = struct{}{}

			for _, bookedSeat := range bookedSeats[int(row)] {
				if bookedSeat == int(seat) {
					return fmt.Errorf("seat %d in row %d is already booked", seat, row)
				}
			}
		}
	}

	return nil
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

func (service *UserOrderService) OrderSummary(c echo.Context, req requests.OrderSummaryRequest) (responses.SummaryResponse, error) {
	var response responses.SummaryResponse
	var bonusesIncrease float64
	var bonusesDecrease float64
	var expiredAt *time.Time
	errorsValid, ok := validators.ValidateOrderSummary(c, service.db, req)
	if !ok {
		return response, errors.New(utils.InputErrorsValid(errorsValid))
	}
	var order models.Order
	err := service.db.Preload("Premiere").Model(&models.Order{}).Where("id = ?", req.ID).First(&order).Error
	if err != nil {
		return response, err
	}
	seats := utils.ParseSeats(order.Seats)
	if order.Status == models.OrderPaid {
		t := order.Premiere.StartTime
		expiredAt = &t
	} else if order.Status == models.OrderPending {
		t := order.CreatedAt.Add(30 * time.Minute)
		expiredAt = &t
	}
	if order.CoinsToPlus != nil {
		bonusesIncrease = *order.CoinsToPlus
	}
	if order.Coins != nil {
		bonusesDecrease = *order.Coins
	}
	response.BonusesIncrease = bonusesIncrease
	response.BonusesDecrease = bonusesDecrease
	response.ExpiredAt = expiredAt
	response.Seats = seats
	response.PremiereID = order.Premiere.ID
	response.MovieID = order.Premiere.MovieID
	response.Status = order.Status
	return response, nil
}
