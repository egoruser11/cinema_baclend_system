package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"gorm.io/gorm"
)

func ValidateCreatePaymentMethod(req requests.PaymentMethodCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	if req.Details == "" {
		errors["details"] = "details is required"
	}
	if req.Type == "" || req.Type != models.PaymentCard && req.Type != models.PaymentWallet {
		errors["type"] = "Type is required"
	}

	return errors, len(errors) == 0
}

func ValidateUpdatePaymentMethod(db *gorm.DB, req requests.PaymentMethodUpdateRequest) (map[string]string, map[string]interface{}, bool) {
	errors := make(map[string]string)
	updates := make(map[string]interface{})
	if req.Id == nil {
		errors["id"] = "id is required"
	}
	var count int64
	err := db.Model(&models.PaymentMethod{}).Where("id = ?", *req.Id).Count(&count).Error
	if err != nil {
		errors["id"] = "Payment method with is di not found"
	}
	if req.Details != nil {
		if *req.Details != "" {
			updates["details"] = *req.Details
		} else {
			errors["details"] = "details is required if your input something"
		}
	}
	if req.Type != nil {
		if *req.Type == "" || *req.Type != models.PaymentCard && *req.Type != models.PaymentWallet {
			errors["type"] = "Type is required"
		} else {
			updates["type"] = *req.Type
		}
	}
	if req.IsDefault != nil {
		updates["is_default"] = *req.IsDefault
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	return errors, updates, len(errors) == 0
}
