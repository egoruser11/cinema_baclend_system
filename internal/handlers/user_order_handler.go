package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserOrderHandler struct {
	userOrderService *services.UserOrderService
}

func NewUserOrderHandler(userOrderService *services.UserOrderService) *UserOrderHandler {
	return &UserOrderHandler{userOrderService: userOrderService}
}

func (handler *UserOrderHandler) Create(c echo.Context) error {
	var req requests.OrderCreateRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	order, err := handler.userOrderService.Create(c, req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.OK(c, order)
}

func (handler *UserOrderHandler) Paid(c echo.Context) error {
	var req requests.OrderPaidRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	order, err := handler.userOrderService.Paid(c, req)
	if err != nil {
		return utils.InternalServerError(c, err.Error())
	}
	return utils.OK(c, order)
}
