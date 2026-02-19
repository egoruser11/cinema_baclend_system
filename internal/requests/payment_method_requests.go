package requests

import (
	"cinema_backend_system/internal/models"
)

type PaymentMethodCreateRequest struct {
	Type      models.PaymentMethodType `json:"type" binding:"required"`
	Details   string                   `json:"details" binding:"required"`
	IsDefault bool                     `json:"is_default" binding:"required,min=1,max=20"`
}
