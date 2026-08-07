package repository

import (
	"errors"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type TradePointDebtPostgres struct {
	db *gorm.DB
}

func NewTradePointDebtPostgres(db *gorm.DB) *TradePointDebtPostgres {
	return &TradePointDebtPostgres{db: db}
}

func (r *TradePointDebtPostgres) Create(tradePointId int, currency float64) error {
	var tradePoint models.TradePoint

	if currency < 1 {
		return errors.New("Rugsat berilmeýär")
	}

	get_trade_point_request := r.db.Model(&models.TradePoint{}).Where("id = ?", tradePointId).Find(&tradePoint)

	if get_trade_point_request.Error != nil || tradePoint.Id == 0 {
		return errors.New("Söwda nokady tapylamdy")
	}

	newOrder := models.Order{
		TradePointId:  tradePoint.Id,
		TradeAgentId:  tradePoint.TradeAgentId,
		Sum:           currency,
		SaleSum:       currency,
		CreatedAt:     time.Now(),
		PaymentTypeId: 1,
		IsCredit:      true,
		Status:        true,
		IsClosed:      false,
	}

	save_order := r.db.Model(&models.Order{}).Create(&newOrder)
	return save_order.Error
}

func (r *TradePointDebtPostgres) GetDebtForMobile(trade_point_id int, language string) ([]models.TradePointDebtResponseForMobile, error) {
	result := []models.TradePointDebtResponseForMobile{}

	var orders []models.Order

	get_orders_not_closed := r.db.Model(&models.Order{}).Where("trade_point_id = ?", trade_point_id).
		Where("is_closed = ?", false).Preload("PaymentType").Find(&orders)

	if get_orders_not_closed.Error != nil {
		return result, get_orders_not_closed.Error
	}

	for _, order := range orders {
		var payments []models.Payment
		payments_sum := float64(0)

		get_order_payments := r.db.Model(&models.Payment{}).Where("order_id = ?", order.Id).Find(&payments)
		if get_order_payments.Error != nil {

		}

		for _, payment := range payments {
			payments_sum = payments_sum + payment.Currency
		}
		if payments_sum < order.SaleSum {
			newData := models.TradePointDebtResponseForMobile{
				TradePointId: trade_point_id,
				Currency:     order.SaleSum - payments_sum,
				CreatedAt:    order.CreatedAt.Format("02.01.2006"),
			}

			if language == "ru" {
				newData.PaymentType = order.PaymentType.NameRu
			} else {
				newData.PaymentType = order.PaymentType.NameTm
			}

			result = append(result, newData)
		}
	}

	return result, nil
}
