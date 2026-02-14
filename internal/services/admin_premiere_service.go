package services

import (
	"cinema_backend_system/internal/models"
	"cinema_backend_system/internal/requests"
	"cinema_backend_system/internal/utils"
	"cinema_backend_system/internal/validators"
	"errors"
	"gorm.io/gorm"
	"time"
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
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
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
	err = service.db.Preload("Movie").Find(premiere).Error
	if err != nil {
		return nil, errors.New("Can not create premiere , please try again")
	}
	return premiere, nil
}

func (service *AdminPremiereService) Update(req requests.PremiereUpdateRequest) (*models.Premiere, error) {
	errorsValid, updates, ok := validators.ValidateUpdatePremiere(service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
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

func (service *AdminPremiereService) Index(req requests.PremiereIndexRequest) ([]*models.Premiere, error) {
	errorsValid, filter, ok := validators.ValidateIndexPremiers(service.db, req)
	if !ok {
		return nil, errors.New(utils.InputErrorsValid(errorsValid))
	}
	premieres, err := models.GetAvailablePremieres(service.db, req.MovieID)
	if len(premieres) == 0 {
		return []*models.Premiere{}, nil
	}
	premiereIds := []uint{}
	for _, premiere := range premieres {
		premiereIds = append(premiereIds, premiere.ID)
	}

	if err != nil {
		return nil, errors.New("Can not find premieres , please try again")
	}
	query := service.db.Model(&models.Premiere{}).Where("id in (?)", premiereIds).Preload("Movie")

	if len(filter) > 0 {
		dayPremiere, existsDayPremiere := filter["day_premiere"]
		if existsDayPremiere {
			date := dayPremiere.(time.Time)
			startDayPremiere := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
			endDayPremiere := startDayPremiere.AddDate(0, 0, 1)
			query = service.db.Where("start_time BETWEEN ? AND ? ", startDayPremiere, endDayPremiere)
		}

		sort, existsSort := filter["sort"]
		if existsSort {
			query = query.Order(sort.(string) + filter["order_type"].(string))
		}

		hourFromVal, exsistsHourFrom := filter["hour_from"]
		if exsistsHourFrom {
			hourFrom := hourFromVal.(time.Time)
			query = query.Where("CAST(start_time AS TIME) >= CAST(? AS TIME)", hourFrom.Format("15:04"))
		}
		hourToVal, exsistsHourTo := filter["hour_to"]
		if exsistsHourTo {
			hourTo := hourToVal.(time.Time)
			query = query.Where("CAST(start_time AS TIME) <= CAST(? AS TIME)", hourTo.Format("15:04"))
		}

		weekDayVal, exsistsWeekDay := filter["week_day"]
		if exsistsWeekDay {
			weekDay := weekDayVal.(int)
			query = query.Where("EXTRACT(DOW FROM start_time) = ?", weekDay)
		}

		maxPriceVal, exsistsMaxPrice := filter["max_price"]
		if exsistsMaxPrice {
			maxPrice := maxPriceVal.(int)
			query = query.Where("price <= ?", maxPrice)
		}
		minPriceVal, exsistsMinPrice := filter["min_price"]
		if exsistsMinPrice {
			minPrice := minPriceVal.(int)
			query = query.Where("price >= ?", minPrice)
		}

		offsetVal := filter["offset"]
		offset := offsetVal.(int)
		limitVal := filter["limit"]
		limit := limitVal.(int)
		query = query.Limit(limit).Offset(offset)

	}
	var result []*models.Premiere
	err = query.Find(&result).Error
	if err != nil {
		return nil, errors.New("Failed to filter premieres")
	}
	return result, err
}
