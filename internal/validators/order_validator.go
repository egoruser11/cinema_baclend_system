package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ValidateCreateOrder(c echo.Context, db *gorm.DB, req requests.OrderCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	var premiere models.Premiere
	err := db.Model(&models.Premiere{}).Where("id = ?", req.PremiereID).Find(&premiere).Error
	if err != nil {
		errors["premiere"] = "Premiere not found"
	}
	userId := c.Get("user_data").(*models.User).ID
	bookedSeatsMap := map[int][]int{}
	err = json.Unmarshal(premiere.BookedSeats, &bookedSeatsMap)
	if err != nil {
		errors["db"] = "BookedSeats json unmarshal error"
	}
	for rowBooked, seatsBooked := range bookedSeatsMap {
		for _, seatBooked := range seatsBooked {
			for row, seats := range req.Seats {
				for _, seat := range seats {
					if int(row) == rowBooked {
						if seatBooked == int(seat) {
							errors["seats"] = fmt.Sprintf("Seat booked for %s is already booked", seatBooked)
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
				errors["seats"] = fmt.Sprintf("Seats cannot dublicate!", seat)
				break
			}
			seenSeats[key] = true
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
