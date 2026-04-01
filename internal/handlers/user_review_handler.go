package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

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

func (handler *UserReviewHandler) Update(c echo.Context) error {
	var request requests.ReviewUpdateRequest
	if err := c.Bind(&request); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	if err := validator.New().Struct(request); err != nil {
		return utils.BadRequest(c, err.Error())
	}
	review, err := handler.userReviewService.Update(c, request)
	if err != nil {
		return utils.BadRequest(c, err.Error())
	}
	return utils.OK(c, review)
}
