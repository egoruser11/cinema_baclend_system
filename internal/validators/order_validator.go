package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"time"
)

func ValidateCreateOrder(db *gorm.DB, req requests.OrderCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	var premiere models.Premiere
	err := db.Model(&models.Premiere{}).Where("id = ?", req.PremiereID).Find(&premiere).Error
	if err != nil {
		errors["premiere"] = "Premiere not found"
	}
	bookedSeatsMap := map[int][]int{}
	if premiere.BookedSeats != nil {
		err = json.Unmarshal(premiere.BookedSeats, &bookedSeatsMap)
		if err != nil {
			errors["db"] = "BookedSeats json unmarshal error"
		}
		for rowBooked, seatsBooked := range bookedSeatsMap {
			for _, seatBooked := range seatsBooked {
				for row, seats := range req.Seats {
					if row > premiere.Rows {
						errors["rows"] = fmt.Sprintf("Row %d is uncorrect", row)
					}
					for _, seat := range seats {
						if int(row) == rowBooked {
							if seatBooked == int(seat) {
								errors["seats"] = fmt.Sprintf("Seat booked for %d is already booked", seatBooked)
							}
						}
					}
				}
			}
		}
	}
	seenSeats := map[string]bool{}
	for row, seats := range req.Seats {
		for _, seat := range seats {
			key := fmt.Sprintf("%d-%d", row, seat)
			if seenSeats[key] {
				errors["seats"] = fmt.Sprintf("Seats cannot dublicate! %d", seat)
				break
			}
			seenSeats[key] = true
		}
	}
	return errors, len(errors) == 0
}

func ValidatePaidOrder(c echo.Context, db *gorm.DB, req requests.OrderPaidRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	//если с момента создания заказа прошло более 30 мин, не даем его оплатить , удаляем заказ , удаляем места из брони.
	//если не хватает денег , выдаем ошибку
	//если введено много коинов , выдаем ошибку, если введено коинов столько , что превышает стоимость заказа, то вернем  ошибку.
	user := c.Get("user_data").(*models.User)
	var order models.Order
	err := db.Model(&models.Order{}).Where("id = ?", req.OrderID).Find(&order).Error
	if order.Status == models.OrderPaid {
		errors["orderPaid"] = "Order is paid!"
		return errors, false
	}
	if err != nil {
		errors["order"] = "Order not found"
	}
	if order.UserID != user.ID {
		errors["order"] = "User id does not match , dont do this"
	}
	fmt.Println(order.CreatedAt)
	if order.CreatedAt.Before(time.Now().Add(-30 * time.Minute)) {
		errors["orderExpired"] = "Order is too old, create new"
		//надо разбронировать места
		return errors, false
	}
	if req.Coins != nil {
		var coins float64
		db.Model(&models.User{}).Where("id = ?", user.ID).Pluck("coins", &coins)
		if coins < *req.Coins {
			errors["coins"] = fmt.Sprintf("Coins %d < %d KINGSLAYEEER", coins, *req.Coins)
		}
		if float64(*req.Coins) > order.TotalAmount {
			errors["coins"] = "Ypu input more then total amount of this order"
		}
	}
	_, exists := errors["coins"]
	if !exists {
		userBalance := user.MoneyBalance
		sumToPay := order.TotalAmount
		if req.Coins != nil {
			sumToPay -= float64(*req.Coins)
		}
		if sumToPay > userBalance {
			errors["money"] = "Money is not enough"
		}
	}
	return errors, len(errors) == 0
}

func ValidateOrderRefund(c echo.Context, db *gorm.DB, request requests.OrderRefundedRequest) (map[string]string, bool) {
	var order models.Order
	user := c.Get("user_data").(*models.User)
	errors := make(map[string]string)
	err := db.Preload("Premiere").Model(&models.Order{}).Where("id = ?", request.ID).Find(&order).Error
	if err != nil {
		errors["order"] = "Order not found"
		return errors, false
	}
	if user.ID != order.UserID {
		errors["order"] = "User id does not match , dont do this"
	}
	if order.Status != models.OrderPaid {
		errors["orderPaid"] = "Order is not paid!"
	}
	var premiere models.Premiere
	premiere = order.Premiere
	if time.Now().After(premiere.StartTime) {
		errors["premiere"] = "Premiere is start , ypu can not refund money , sorry!"
	}
	return errors, len(errors) == 0
}

func ValidateOrderIndex(c echo.Context, db *gorm.DB, req requests.OrderIndexRequest) (map[string]string, map[string]interface{}, bool) {
	user := c.Get("user_data").(*models.User)
	errors := make(map[string]string)
	filters := make(map[string]interface{})

	filters["user_id"] = user.ID

	if req.Status != nil {
		validStatuses := map[models.OrderStatus]bool{
			models.OrderPending:  true,
			models.OrderPaid:     true,
			models.OrderRefunded: true,
			models.OrderCanceled: true,
			models.OrderDeleted:  true,
		}
		if !validStatuses[*req.Status] {
			errors["status"] = "invalid order status"
		} else {
			filters["status"] = *req.Status
		}
	}

	if req.PremiereID != nil {
		var premiere models.Premiere
		if err := db.Where("id = ?", *req.PremiereID).First(&premiere).Error; err != nil {
			errors["premiere_id"] = "premiere not found"
		} else {
			filters["premiere_id"] = *req.PremiereID
		}
	}

	if req.DateFrom != nil && req.DateTo != nil {
		if req.DateFrom.After(*req.DateTo) {
			errors["date_range"] = "date_from must be before date_to"
		} else {
			filters["date_from"] = *req.DateFrom
			filters["date_to"] = *req.DateTo
		}
	} else if req.DateFrom != nil {
		filters["date_from"] = *req.DateFrom
	} else if req.DateTo != nil {
		filters["date_to"] = *req.DateTo
	}

	// Пагинация
	if req.Limit != nil {
		if *req.Limit < 1 || *req.Limit > 100 {
			errors["limit"] = "limit must be between 1 and 100"
		} else {
			filters["limit"] = *req.Limit
		}
	} else {
		filters["limit"] = 10
	}

	if req.Offset != nil {
		if *req.Offset < 0 {
			errors["offset"] = "offset must be >= 0"
		} else {
			filters["offset"] = *req.Offset
		}
	} else {
		filters["offset"] = 0
	}

	// Сортировка
	if req.Sort != nil {
		validSortFields := map[string]bool{
			"created_at":   true,
			"total_amount": true,
			"status":       true,
		}
		if !validSortFields[*req.Sort] {
			errors["sort"] = "invalid sort field, allowed: created_at, total_amount, status"
		} else {
			filters["sort"] = *req.Sort
			if req.Order != nil && (*req.Order == "ASC" || *req.Order == "DESC") {
				filters["order"] = *req.Order
			} else {
				filters["order"] = "DESC"
			}
		}
	}

	return errors, filters, len(errors) == 0
}

func ValidateOrderShow(c echo.Context, db *gorm.DB, req requests.OrderShowRequest) (map[string]string, *models.Order) {
	user := c.Get("user_data").(*models.User)
	errors := make(map[string]string)

	if req.ID == nil {
		errors["id"] = "order id is required"
		return errors, nil
	}

	var order models.Order
	err := db.Preload("Premiere.Movie").Where("id = ?", *req.ID).First(&order).Error
	if err != nil {
		errors["order"] = "order not found"
		return errors, nil
	}

	if order.UserID != user.ID {
		errors["order"] = "access denied: order does not belong to this user"
		return errors, nil
	}

	return errors, &order
}

func ValidateOrderDelete(c echo.Context, db *gorm.DB, req requests.OrderDeleteRequest) (map[string]string, *models.Order, bool) {
	user := c.Get("user_data").(*models.User)
	errors := make(map[string]string)

	if req.ID == nil {
		errors["id"] = "order id is required"
		return errors, nil, false
	}

	var order models.Order
	err := db.Where("id = ?", *req.ID).First(&order).Error
	if err != nil {
		errors["order"] = "order not found"
		return errors, nil, false
	}

	if order.UserID != user.ID {
		errors["order"] = "access denied: order does not belong to this user"
		return errors, nil, false
	}

	if order.Status != models.OrderPending && order.Status != models.OrderCanceled {
		errors["status"] = fmt.Sprintf("cannot delete order with status: %s", order.Status)
		return errors, nil, false
	}

	return errors, &order, true
}

func ValidateOrderUpdate(c echo.Context, db *gorm.DB, req requests.OrderUpdateRequest) (map[string]string, bool, map[uint][]uint) {
	user := c.Get("user_data").(*models.User)
	errors := make(map[string]string)
	var newSeatsBooked map[uint][]uint
	if req.ID == nil {
		errors["id"] = "order id is required"
		return errors, false, newSeatsBooked
	}

	var order models.Order
	err := db.Preload("Premiere").Where("id = ?", *req.ID).First(&order).Error
	if err != nil {
		errors["order"] = "order not found"
		return errors, false, newSeatsBooked
	}

	if order.UserID != user.ID {
		errors["order"] = "access denied: order does not belong to this user"
		return errors, false, newSeatsBooked
	}

	if order.Status != models.OrderPending {
		errors["status"] = "cannot update order: only pending orders can be updated"
		return errors, false, newSeatsBooked
	}

	if order.CreatedAt.Before(time.Now().Add(-30 * time.Minute)) {
		errors["orderExpired"] = "Order is too old, create new"
		return errors, false, newSeatsBooked
	}
	if req.Seats != nil {
		// Получаем текущие места заказа
		currentSeats := utils.ParseSeats(order.Seats)
		newSeatsBooked = make(map[uint][]uint)
		// Получаем уже забронированные места в премьере
		var bookedSeatsMap map[int][]int
		if order.Premiere.BookedSeats != nil {
			json.Unmarshal(order.Premiere.BookedSeats, &bookedSeatsMap)
		}
		for row, seatsRow := range req.Seats {
			if row > order.Premiere.Rows {
				errors["row"] = fmt.Sprintf("row %d exceeds maximum rows", row)
			}
			for _, seat := range seatsRow {
				isBookedByOther := false
				for _, booked := range bookedSeatsMap[int(row)] {
					if booked == int(seat) {
						isInCurrentOrder := false
						for _, currentSeat := range currentSeats[row] {
							if currentSeat == seat {
								isInCurrentOrder = true
								break
							}
						}
						if !isInCurrentOrder {
							isBookedByOther = true
						}
						break
					}
				}
				if isBookedByOther {
					errors["seats"] = fmt.Sprintf("seat %d in row %d is already booked by another user", seat, row)
				} else {
					newSeatsBooked[row] = append(newSeatsBooked[row], seat)
				}
			}
		}
		seenSeats := map[string]bool{}
		for row, seatsRow := range req.Seats {
			for _, seat := range seatsRow {
				key := fmt.Sprintf("%d-%d", row, seat)
				if seenSeats[key] {
					errors["seats"] = fmt.Sprintf("duplicate seat %d in row %d", seat, row)
				}
				seenSeats[key] = true
			}
		}
	}
	fmt.Println(newSeatsBooked)
	if len(newSeatsBooked) == 0 {
		errors["seats"] = "no new seats"
		return errors, false, newSeatsBooked
	}
	return errors, len(errors) == 0, newSeatsBooked
}
