package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type UserReviewHandler struct {
	userReviewService *services.UserReviewService
}

func NewUserReviewHandler(userReviewService *services.UserReviewService) *UserReviewHandler {
	return &UserReviewHandler{userReviewService: userReviewService}
}

func (handler *UserReviewHandler) Create(c echo.Context) error {
	var request requests.ReviewCreateRequest
	if err := c.Bind(&request); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	review, err := handler.userReviewService.Create(c, request)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, review)
}
