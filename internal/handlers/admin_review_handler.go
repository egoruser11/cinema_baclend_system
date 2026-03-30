package handlers

import (
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/services"
	"cinema_backend_system/internal/utils"
	"github.com/labstack/echo/v4"
)

type AdminReviewHandler struct {
	adminReviewService *services.AdminReviewService
}

func NewAdminReviewHandler(adminReviewService *services.AdminReviewService) *AdminReviewHandler {
	return &AdminReviewHandler{adminReviewService: adminReviewService}
}

func (handler *AdminReviewHandler) Approve(c echo.Context) error {
	var request requests.ReviewApproveRequest
	if err := c.Bind(&request); err != nil {
		return utils.BadRequest(c, "Invalid request body")
	}
	review, err := handler.adminReviewService.Approve(request)
	if err != nil {
		return err
	}
	return utils.OK(c, review)
}
