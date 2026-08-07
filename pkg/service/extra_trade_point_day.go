package service

import "github.com/Saparly-go-dev/tagma/pkg/repository"

type ExtraTradePointDayService struct {
	repo repository.ExtraTradePointDay
}

func NewExtraTradePointDayService(repo repository.ExtraTradePointDay) *ExtraTradePointDayService {
	return &ExtraTradePointDayService{repo: repo}
}

func (s *ExtraTradePointDayService) Update(trade_point_id int, list []int) error {
	return s.repo.Update(trade_point_id, list)
}

func (s *ExtraTradePointDayService) Get(trade_point_id int) []int {
	return s.repo.Get(trade_point_id)
}
