package repository

import (
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type ExtraTradePointDayPostgres struct {
	db *gorm.DB
}

func NewExtraTradePointDayPostgres(db *gorm.DB) *ExtraTradePointDayPostgres {
	return &ExtraTradePointDayPostgres{db: db}
}

func (r *ExtraTradePointDayPostgres) Update(trade_point_id int, list []int) error {
	var old_data []models.ExtraTradePointDay

	fmt.Print(trade_point_id, list)

	get_old_data_request := r.db.Model(&models.ExtraTradePointDay{}).Where("trade_point_id = ?", trade_point_id).Find(&old_data)

	if get_old_data_request.Error != nil {
		return get_old_data_request.Error
	}

	if len(old_data) > 0 {

		delete_old_data_request := r.db.Model(&models.ExtraTradePointDay{}).Delete(&old_data)

		if delete_old_data_request.Error != nil {
			return delete_old_data_request.Error
		}
	}

	for _, item := range list {
		new_data := models.ExtraTradePointDay{
			TradePointId: trade_point_id,
			DayId:        item,
		}

		create_request := r.db.Model(&models.ExtraTradePointDay{}).Create(&new_data)

		if create_request.Error != nil {
			return create_request.Error
		}
	}

	return nil
}

func (r *ExtraTradePointDayPostgres) Get(trade_point_id int) []int {
	var result []int
	var list []models.ExtraTradePointDay

	get_data_request := r.db.Model(&models.ExtraTradePointDay{}).Where("trade_point_id = ?", trade_point_id).Find(&list)

	if get_data_request.Error != nil {
		return result
	}

	for _, item := range list {
		result = append(result, item.DayId)
	}
	return result
}
