package repository

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

// Write entry points live here so transaction boundaries stay separate from
// query/report code in the larger repository implementations.
func (r *OrderPostgres) SaveOrder(
	day, month, year, tradePointID, paymentTypeID int,
	first, isCredit bool,
	orders []models.CreateOrder,
) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := *r
		txRepo.db = tx
		return txRepo.saveOrder(day, month, year, tradePointID, paymentTypeID, first, isCredit, orders)
	})
}

func (r *PaymentPostgres) SavePayment(agentID int, data models.CreatePayment) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := *r
		txRepo.db = tx
		return txRepo.savePayment(agentID, data)
	})
}

func (r *PaymentPostgres) PostPayment(currency float64, agentID, tradePointID, paymentTypeID int, isAgent bool) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := *r
		txRepo.db = tx
		return txRepo.postPayment(currency, agentID, tradePointID, paymentTypeID, isAgent)
	})
}

func (r *UserPostgres) Create(user tagma.CreateUser) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := *r
		txRepo.db = tx
		return txRepo.create(user)
	})
}
