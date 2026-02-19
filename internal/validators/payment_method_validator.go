package validators

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
)

func ValidateCreatePaymentMethod(req requests.PaymentMethodCreateRequest) (map[string]string, bool) {
	errors := make(map[string]string)
	if req.Details == "" {
		errors["details"] = "details is required"
	}
	if req.Type == "" {
		errors["type"] = "Type is required"
	}
	if req.Type != models.PaymentCard && req.Type != models.PaymentWallet {
		errors["type"] = "Type is incorrect"
	}

	return errors, len(errors) == 0
}
