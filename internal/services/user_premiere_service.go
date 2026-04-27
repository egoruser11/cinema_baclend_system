package services

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"cinema_backend_system/responses"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
)

type UserPremiereService struct {
	db *gorm.DB
}

func NewUserPremiereService(db *gorm.DB) *UserPremiereService {
	return &UserPremiereService{db: db}
}

func (service *UserPremiereService) SeatMap(req requests.PremiereSeatMapRequest) (responses.SeatMapResponse, error) {
	errorsValid, premiere, ok := validators.ValidateSeatMap(service.db, req)
	var response responses.SeatMapResponse
	if !ok {
		return response, errors.New(utils.InputErrorsValid(errorsValid))
	}

	bookedSeatsMap := map[uint][]uint{}
	freeSeatsMap := map[uint][]uint{}
	if len(premiere.BookedSeats) > 0 {
		err := json.Unmarshal(premiere.BookedSeats, &bookedSeatsMap)
		if err != nil {
			return response, err
		}
	}

	rows := int(premiere.Rows)
	seatsPerRows := int(premiere.SeatsPerRow)
	if len(bookedSeatsMap) == 0 {
		for i := 1; i <= rows; i++ {
			for j := 1; j <= seatsPerRows; j++ {
				freeSeatsMap[uint(i)] = append(freeSeatsMap[uint(i)], uint(j))
			}
		}
		response.BusySeats = bookedSeatsMap
		response.FreeSeats = freeSeatsMap
		return response, nil
	}

	for i := 1; i <= rows; i++ {
		for j := 1; j <= seatsPerRows; j++ {
			if isSeatFree(bookedSeatsMap, uint(i), uint(j)) {
				freeSeatsMap[uint(i)] = append(freeSeatsMap[uint(i)], uint(j))
			}
		}
	}

	response.BusySeats = bookedSeatsMap
	response.FreeSeats = freeSeatsMap
	return response, nil
}

func isSeatFree(bookedSeatsMap map[uint][]uint, rowCheck, seatCheck uint) bool {
	for row, seats := range bookedSeatsMap {
		for _, seat := range seats {
			if row == rowCheck && seat == seatCheck {
				return false
			}
		}
	}
	return true
}
