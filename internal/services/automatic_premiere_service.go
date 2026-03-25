package services

import (
	"cinema_backend_system/internal/models"
	"encoding/json"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

func ReserveSeats(db *gorm.DB, bookedCount int, premiere models.Premiere, seats map[uint][]uint) error {
	var existingSeats map[int][]int
	json.Unmarshal(premiere.BookedSeats, &existingSeats)

	if existingSeats == nil {
		existingSeats = make(map[int][]int)
	}
	for row, seatsInRow := range seats {
		rowInt := int(row)
		for _, num := range seatsInRow {
			numInt := int(num)
			found := false
			for _, existing := range existingSeats[rowInt] {
				if existing == numInt {
					found = true
					break
				}
			}
			if !found {
				existingSeats[rowInt] = append(existingSeats[rowInt], numInt)
			}
		}
	}

	for row := range existingSeats {
		sort.Ints(existingSeats[row])
	}

	data, _ := json.Marshal(existingSeats)
	return db.Model(&premiere).Updates(map[string]interface{}{
		"booked_seats": data,
		"booked_count": bookedCount,
	}).Error
}

func UnReserveSeats(db *gorm.DB, premiere *models.Premiere, order *models.Order, orders []models.Order, status models.OrderStatus) error {
	if orders != nil {
		for _, orderCurr := range orders {
			err := UnReserveForOneOrder(db, orderCurr, orderCurr.Premiere, status)
			if err != nil {
				return err
			}
		}
	}
	err := UnReserveForOneOrder(db, *order, *premiere, status)
	return err
}

func findInRowsCopiesSeats(seatsInRowBooked []int, seatsInRowInOrder []int) []int {
	result := []int{}
	for i := 0; i < len(seatsInRowBooked); i++ {
		for j := 0; j < len(seatsInRowInOrder); j++ {
			if seatsInRowBooked[i] == seatsInRowInOrder[j] {
				result = append(result, seatsInRowBooked[i])
			}
		}
	}
	return result
}

func UnReserveForOneOrder(db *gorm.DB, order models.Order, premiere models.Premiere, status models.OrderStatus) error {
	seats := strings.TrimSuffix(order.Seats, ",")
	mapInOrderSeats := make(map[int][]int)
	newPremiereBookedSeats := make(map[int][]int)

	seatsArr := strings.Split(seats, ",")
	for _, seat := range seatsArr {
		parts := strings.Split(seat, " - ")

		row, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		seatNum, _ := strconv.Atoi(strings.TrimSpace(parts[1]))

		mapInOrderSeats[row] = append(mapInOrderSeats[row], seatNum)
	}
	var bookedSeats map[int][]int
	bookedCountNew := premiere.BookedCount
	json.Unmarshal(premiere.BookedSeats, &bookedSeats)
	for rowBooked, seatsInRowBooked := range bookedSeats {
		arrCopiesSeats := []int{}
		for rowInOrder, seatsInRowInOrder := range mapInOrderSeats {
			if rowInOrder == rowBooked {
				arrCopiesSeats = findInRowsCopiesSeats(seatsInRowBooked, seatsInRowInOrder)
			}
		}
		if len(arrCopiesSeats) > 0 {
			for _, seat := range seatsInRowBooked {
				isAdd := false
				for i := 0; i < len(arrCopiesSeats); i++ {
					if seat == arrCopiesSeats[i] {
						bookedCountNew--
						isAdd = true
					}
				}
				if !isAdd {
					newPremiereBookedSeats[rowBooked] = append(newPremiereBookedSeats[rowBooked], seat)
				}
			}
		} else {
			for _, seat := range seatsInRowBooked {
				newPremiereBookedSeats[rowBooked] = append(newPremiereBookedSeats[rowBooked], seat)
			}
		}
	}
	newBookedSeatsJson, err := json.Marshal(newPremiereBookedSeats)
	if err != nil {
		return err
	}
	db.Model(&premiere).Updates(map[string]interface{}{
		"booked_seats": newBookedSeatsJson,
		"booked_count": bookedCountNew,
	})
	switch status {
	case models.OrderDeleted:
		db.Delete(&order)
	default:
		db.Model(&order).Update("status", status)
	}
	return nil
}
