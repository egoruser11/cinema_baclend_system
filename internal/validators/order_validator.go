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

func ValidateCreateOrder(c echo.Context, db *gorm.DB, req requests.OrderCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	var premiere models.Premiere
	err := db.Model(&models.Premiere{}).Where("id = ?", req.PremiereID).Find(&premiere).Error
	if err != nil {
		errors["premiere"] = "Premiere not found"
	}
	userId := c.Get("user_data").(*models.User).ID
	err = db.Model(&models.PaymentMethod{}).Where("id = ? AND user_id = ?", req.PaymentMethodID, userId).Error
	if err != nil {
		errors["payment_method"] = "PaymentMethod not found, or its not your payment method , please don do this!"
	}

	bookedSeatsMap := map[uint]int{}
	err = json.Unmarshal(premiere.BookedSeats, &bookedSeatsMap)
	if err != nil {
		errors["db"] = "BookedSeats json unmarshal error"
	}
	for rowBooked, seatBooked := range bookedSeatsMap {
		for row, seat := range req.Seats {
			if row == rowBooked {
				if seatBooked == seat {
					errors["seats"] = fmt.Sprintf("Seat booked for %s is already booked", seatBooked)
				}
			}
		}
	}
	if req.CountMinutesBeforePay != nil {
		if time.Now().Add(time.Duration(*req.CountMinutesBeforePay) * time.Minute).After(premiere.StartTime) {
			errors["count_minutes_before_pay"] = "Count minutes before pay not valid"
		}
		if *req.CountMinutesBeforePay > 300 {
			errors["count_minutes_before_pay"] = "Count minutes before pay not valid"
		}
	}
	if req.Coins != nil {
		var coins uint64
		db.Model(&models.User{}).Where("id = ?", userId).Pluck("coins", &coins)
		if coins < *req.Coins {
			errors["coins"] = fmt.Sprintf("Coins %d < %d KINGSLAYEEER", coins, *req.Coins)
		}
	}
	return errors, true
}
