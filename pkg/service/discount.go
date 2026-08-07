package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type DiscountService struct {
	repo repository.Discount
}

func NewDiscountService(repo repository.Discount) *DiscountService {
	return &DiscountService{repo: repo}
}

func (s *DiscountService) Create(discount models.CreateDiscount) error {
	return s.repo.Create(discount)
}

func (s *DiscountService) GetAll(pageSize, pageNumber int, name, language string) (*models.DiscountPage, error) {
	return s.repo.GetAll(pageSize, pageNumber, name, language)
}

func (s *DiscountService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *DiscountService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *DiscountService) Update(Id int, discount models.CreateDiscount) error {
	return s.repo.Update(Id, discount)
}

func (s *DiscountService) UpdateNameAndDate(id int, data models.UpdateDiscountInfo) error {
	return s.repo.UpdateNameAndDate(id, data)
}

func (s *DiscountService) GetTradePointByTradeAgent(tradeAgentId int) ([]models.SelectObject, error) {
	return s.repo.GetTradePointByTradeAgent(tradeAgentId)
}
