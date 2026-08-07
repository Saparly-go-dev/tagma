package repository

import (
	"fmt"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type TimeCheckerPostgres struct {
	db *gorm.DB
}

func NewTimeCheckerPostgres(db *gorm.DB) *TimeCheckerPostgres {
	return &TimeCheckerPostgres{db: db}
}

func (r *TimeCheckerPostgres) Orders() {
	var orders []models.Order
	fmt.Println("Timechecker Order start")

	get_orders_reques := r.db.Model(&models.Order{}).Where("is_closed = ?", false).Find(&orders)

	if get_orders_reques.Error != nil {
		fmt.Println("Timechecker Order error:", get_orders_reques.Error)
	}

	for _, order := range orders {
		var payments []models.Payment

		get_payments_for_order := r.db.Model(&models.Payment{}).Where("order_id = ?", order.Id).Find(&payments)

		if get_payments_for_order.Error != nil {
			fmt.Println("Timechecker Order error:", get_payments_for_order.Error)
		}

		var sum float64 = 0
		for _, payment := range payments {
			if payment.Status == true {
				sum += payment.Currency
			}
		}

		if sum == order.SaleSum {
			updateData := map[string]interface{}{
				"is_closed": true,
			}

			update_order_request := r.db.Model(&models.Order{}).Where("id = ?", order.Id).Updates(&updateData)
			if update_order_request.Error != nil {
				fmt.Println("Timechecker Order error:", update_order_request.Error)
			}
		}
	}

	fmt.Println("Timechecker Order end")
}
