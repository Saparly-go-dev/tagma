package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type TradePointDebtService struct {
	repo repository.TradePointDebt
}

func NewTradePointDebtService(repo repository.TradePointDebt) *TradePointDebtService {
	return &TradePointDebtService{repo: repo}
}

func (s *TradePointDebtService) Create(tradePointId int, currency float64) error {
	return s.repo.Create(tradePointId, currency)
}

func (s *TradePointDebtService) GetDebtForMobile(trade_point_id int, language string) ([]models.TradePointDebtResponseForMobile, error) {
	return s.repo.GetDebtForMobile(trade_point_id, language)
}
