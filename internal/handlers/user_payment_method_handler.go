package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserPaymentMethodHandler struct {
	userPaymentMethodService *services.UserPaymentMethodService
}

func NewUserPaymentMethodHandler(userPaymentMethodService *services.UserPaymentMethodService) *UserPaymentMethodHandler {
	return &UserPaymentMethodHandler{userPaymentMethodService: userPaymentMethodService}
}

func (handler *UserPaymentMethodHandler) Create(c echo.Context) error {
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

func (handler *UserPaymentMethodHandler) Update(c echo.Context) error {
	var req requests.PaymentMethodUpdateRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	method, err := handler.userPaymentMethodService.Update(c, req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, method)
}

func (handler *UserPaymentMethodHandler) Index(c echo.Context) error {
	var req requests.PaymentMethodIndexRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	methods, err := handler.userPaymentMethodService.Index(c, req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, methods)
}

func (handler *UserPaymentMethodHandler) Delete(c echo.Context) error {
	var req requests.PaymentMethodDeleteRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	deletedIds, err := handler.userPaymentMethodService.Delete(c, req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, deletedIds)
}

func (handler *UserPaymentMethodHandler) Show(c echo.Context) error {
	var req requests.PaymentMethodIdRequest
	if err := c.Bind(&req); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	method, err := handler.userPaymentMethodService.Show(c, req)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, method)
}
