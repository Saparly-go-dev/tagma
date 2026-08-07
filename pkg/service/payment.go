package service

import (
	"errors"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type PaymentService struct {
	repo repository.Payment
}

func NewPaymentService(repo repository.Payment) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) GetTradePointsForPayment(day, month, year, agentId int, language string) (*[]models.TradePointForPayment, error) {
	return s.repo.GetTradePointsForPayment(day, month, year, agentId, language)
}

func (s *PaymentService) SavePayment(agentId int, data models.CreatePayment) error {
	if agentId <= 0 || data.OrderId <= 0 || data.PaymentTypeId <= 0 {
		return errors.New("agent, order and payment type are required")
	}
	if data.Currency <= 0 {
		return errors.New("payment currency must be positive")
	}
	return s.repo.SavePayment(agentId, data)
}

func (s *PaymentService) GetPaymentTypes(language string) ([]models.SelectObject, error) {
	return s.repo.GetPaymentTypes(language)
}

func (s *PaymentService) GetPaymentForOrder(orderId int, language string) ([]models.ReadPayment, error) {
	return s.repo.GetPaymentForOrder(orderId, language)
}

func (s *PaymentService) GetAgentPayments(agent_id, payment_type_id int, status bool, point_code, language string) ([]models.ReadPayment, error) {
	return s.repo.GetAgentPayments(agent_id, payment_type_id, status, point_code, language)
}

func (s *PaymentService) ConfirmPayments(agent_id, price1, price2 int, ids []int) error {
	return s.repo.ConfirmPayments(agent_id, price1, price2, ids)
}

func (s *PaymentService) UpdatePayment(id, price1, price2 int) error {
	return s.repo.UpdatePayment(id, price1, price2)
}

func (s *PaymentService) UpdateDebt(agent_id, price1, price2 int) error {
	return s.repo.UpdateDebt(agent_id, price1, price2)
}

func (s *PaymentService) AddSelfDebt(agent_id int, currency float64) error {
	return s.repo.AddSelfDebt(agent_id, currency)
}

func (s *PaymentService) GetTradePointAllDebtSum(tradePointId int) float64 {
	return s.repo.GetTradePointAllDebtSum(tradePointId)
}

func (s *PaymentService) PostPayment(currency float64, agent_id, tradepointId, PaymentTypeId int, is_agent bool) error {
	if currency <= 0 {
		return errors.New("payment currency must be positive")
	}
	if agent_id <= 0 || tradepointId <= 0 || PaymentTypeId <= 0 {
		return errors.New("agent, trade point and payment type are required")
	}
	return s.repo.PostPayment(currency, agent_id, tradepointId, PaymentTypeId, is_agent)
}

func (s *PaymentService) GetTradePointPaymentsHistory(agent_id, tradePointId int, language string) ([]models.ReadPayment, error) {
	return s.repo.GetTradePointPaymentsHistory(agent_id, tradePointId, language)
}

func (s *PaymentService) GetDailyPaymentsInformations(agent_id, day, month, year int, language string) (models.ObjectWithIdNameAndCurrency, error) {
	return s.repo.GetDailyPaymentsInformations(agent_id, day, month, year, language)
}

func (s *PaymentService) GetTradePointsDebtList(agentid int, language string) ([]models.TradePointMobile, error) {
	return s.repo.GetTradePointsDebtList(agentid, language)
}
