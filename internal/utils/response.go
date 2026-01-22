package utils

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message,omitempty"`
}
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code"`
}

func OK(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "BAD_REQUEST",
	})
}

func Unauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "UNAUTHORIZED",
	})
}

func InternalServerError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Error:   message,
		Code:    "INTERNAL_SERVER_ERROR",
	})
}
