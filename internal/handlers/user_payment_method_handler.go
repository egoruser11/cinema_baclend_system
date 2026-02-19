package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"fmt"
	"github.com/labstack/echo/v4"
)

type UserPaymentMethodHandler struct {
	userPaymentMethodService *services.UserPaymentMethodService
}

func NewUserPaymentMethodHandler(userPaymentMethodService *services.UserPaymentMethodService) *UserPaymentMethodHandler {
	return &UserPaymentMethodHandler{userPaymentMethodService: userPaymentMethodService}
}

func (handler *UserPaymentMethodHandler) Create(c echo.Context) error {
	fmt.Sprintln("AAA")
	var req requests.PaymentMethodCreateRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	method, err := handler.userPaymentMethodService.Create(c, req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, method)
}
