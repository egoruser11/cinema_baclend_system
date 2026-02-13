package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/validators"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type AdminPremiereService struct {
	db *gorm.DB
}

func NewAdminPremiereService(db *gorm.DB) *AdminPremiereService {
	return &AdminPremiereService{db: db}
}

func (service *AdminPremiereService) Create(req requests.PremiereCreateRequest) (*models.Premiere, error) {
	errorsValid, ok := validators.ValidateCreatePremiere(service.db, req)
	if !ok {
		errorsRes := []string{}
		for field, err := range errorsValid {
			errorsRes = append(errorsRes, field+" :"+err)
		}
		return nil, errors.New(strings.Join(errorsRes, "\n"))
	}
	premiere := &models.Premiere{
		Hall:        req.Hall,
		MovieID:     req.MovieID,
		Price:       req.Price,
		Rows:        req.Rows,
		SeatsPerRow: req.SeatsPerRow,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		TotalSeats:  int(req.Rows * req.SeatsPerRow),
	}
	err := service.db.Create(premiere).Error
	if err != nil {
		return nil, errors.New("Can not create premiere , please try again")
	}

	return premiere, nil
}

func (service *AdminPremiereService) Update(req requests.PremiereUpdateRequest) (*models.Premiere, error) {
	errorsValid, updates, ok := validators.ValidateUpdatePremiere(service.db, req)
	if !ok {
		errorsRes := []string{}
		for field, err := range errorsValid {
			errorsRes = append(errorsRes, field+" :"+err)
		}
		return nil, errors.New(strings.Join(errorsRes, "\n"))
	}
	var premiere models.Premiere
	err := service.db.Preload("Movie").Where("id = ?", req.Id).First(&premiere).Error
	if err != nil {
		return nil, errors.New("Can not find premiere , please try again")
	}
	if len(updates) > 0 {
		service.db.Model(&premiere).Updates(updates)
	}
	return &premiere, nil
}
