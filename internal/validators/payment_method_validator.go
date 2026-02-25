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

func ValidateIndexPaymentMethod(req requests.PaymentMethodIndexRequest) (map[string]string, map[string]interface{}, bool) {
	errors := make(map[string]string)
	filters := make(map[string]interface{})
	if req.Search != nil {
		if *req.Search != "" {
			filters["search"] = *req.Search
		} else {
			errors["search"] = "search is required if your input"
		}
	}
	if req.Limit != nil {
		if *req.Limit != 0 {
			filters["limit"] = int(*req.Limit)
		}
	} else {
		filters["limit"] = 10
	}
	if req.Offset != nil {
		filters["offset"] = int(*req.Offset)
	} else {
		filters["offset"] = 0
	}

	if req.Sort != nil {
		if *req.Sort == "" || (*req.Sort != "active" && *req.Sort != "not_active") {
			errors["sort"] = "sort is incorrect"
		} else {
			filters["sort"] = *req.Sort
		}
	}
	return errors, filters, len(errors) == 0
}
