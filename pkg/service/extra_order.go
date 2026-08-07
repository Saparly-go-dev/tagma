package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type ExtraOrderService struct {
	repo repository.ExtraOrder
}

func NewExtraOrderService(repo repository.ExtraOrder) *ExtraOrderService {
	return &ExtraOrderService{repo: repo}
}

func (s *ExtraOrderService) GetTradePoint(day_id, trade_agent_id int, language string) (*[]models.TradePointMobile, error) {
	return s.repo.GetTradePoint(day_id, trade_agent_id, language)
}

func (s *ExtraOrderService) GetOrder(day, month, year, tradePointId int, language string) (models.GetOrderObject, error) {
	return s.repo.GetOrder(day, month, year, tradePointId, language)
}
