package services

import (
	"gorm.io/gorm"
)

type UserOrderService struct {
	db *gorm.DB
}

func NewUserOrderService(db *gorm.DB) *UserOrderService {
	return &UserOrderService{db: db}
}

//func (service *UserOrderService) Create(req requests.UserOrderCreateRequest) (*models.UserOrder, error) {
//
//}
