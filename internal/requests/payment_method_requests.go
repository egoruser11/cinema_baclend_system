package requests

import (
	"cinema_backend_system/internal/models"
)

type PaymentMethodCreateRequest struct {
	Type      models.PaymentMethodType `json:"type" binding:"required"`
	Details   string                   `json:"details" binding:"required"`
	IsDefault bool                     `json:"is_default" binding:"required,min=1,max=20"`
}

type PaymentMethodUpdateRequest struct {
	Id        *uint                     `json:"id" binding:"required"`
	Type      *models.PaymentMethodType `json:"type" binding:"required"`
	Details   *string                   `json:"details" binding:"required"`
	IsDefault *bool                     `json:"is_default" binding:"required"`
	IsActive  *bool                     `json:"is_active" binding:"required"`
}

type PaymentMethodIndexRequest struct {
	Sort   *string `query:"sort"`   //active, not_active
	Search *string `query:"search"` //details
	Offset *uint   `query:"offset"`
	Limit  *uint   `query:"limit"`
}

type PaymentMethodDeleteRequest struct {
	Ids []*uint `query:"ids" binding:"required"`
}

type PaymentMethodIdRequest struct {
	Id uint `query:"id" binding:"required"`
}
