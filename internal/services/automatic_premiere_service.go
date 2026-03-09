package services

import (
	"cinema_backend_system/internal/models"
	"encoding/json"
	"gorm.io/gorm"
	"sort"
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
