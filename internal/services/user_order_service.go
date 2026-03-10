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
)

type UserOrderService struct {
	db *gorm.DB
}

func NewUserOrderService(db *gorm.DB) *UserOrderService {
	return &UserOrderService{db: db}
}

func (service *UserOrderService) Create(c echo.Context, req requests.OrderCreateRequest) (*models.Order, error) {
	errorsValid, ok := validators.ValidateCreateOrder(c, service.db, req)
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
	totalAmount := float64(int(premiere.Price) * countSeatsInOrder)
	service.db.Model(&models.Premiere{}).Where("id = ?", req.PremiereID).First(&premiere)
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
	err := ReserveSeats(service.db, countSeatsInOrder, premiere, seatsInOrder)
	if err != nil {
		return nil, err
	}
	return order, nil
}
