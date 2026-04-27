package responses

type SeatMapResponse struct {
	FreeSeats map[uint][]uint `json:"free_seats"`
	BusySeats map[uint][]uint `json:"busy_seats"`
}
