package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
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
	}
	if err != nil {
		errors["order"] = "Order not found"
	}
	if order.UserID != user.ID {
		errors["order"] = "User id does not match , dont do this"
	}
	if order.CreatedAt.Before(time.Now().Add(-30 * time.Minute)) {
		errors["orderExpired"] = "Order is too old, create new"
		//надо разбронировать места
		return errors, false
	}
	if req.Coins != nil {
		var coins uint64
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
