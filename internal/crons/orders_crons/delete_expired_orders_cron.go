package orders_crons

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/services"
	"gorm.io/gorm"
	"time"
)

func Cron(db *gorm.DB) error {
	var expiredOrders []models.Order
	cutoffTime := time.Now().Add(-30 * time.Minute)
	err := db.Preload("Premiere").Where("status = ? AND created_at < ?", models.OrderPending, cutoffTime).Find(&expiredOrders).Error
	if err != nil {
		return err
	}
	if len(expiredOrders) == 0 {
		return nil
	}
	err = services.UnReserveSeats(db, nil, nil, expiredOrders, models.OrderDeleted) // тут происходит удаление
	if err != nil {
		return err
	}
	return nil
}
