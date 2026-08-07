package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type DealerOrderService struct {
	repo repository.DealerOrder
}

func NewDealerOrderService(repo repository.DealerOrder) *DealerOrderService {
	return &DealerOrderService{repo: repo}
}

func (s *DealerOrderService) Save(DealerId int, list []models.CreateOrder) error {
	return s.repo.Save(DealerId, list)
}

func (s *DealerOrderService) Update(DealerOrderId int, list []models.CreateOrder) error {
	return s.repo.Update(DealerOrderId, list)
}

func (s *DealerOrderService) CompleTheDealerOrder(dealerOrderId int) error {
	return s.repo.CompleTheDealerOrder(dealerOrderId)
}

func (s *DealerOrderService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *DealerOrderService) GetOrders(pageSize, pageNumber, DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) (models.DealerOrderPage, error) {
	return s.repo.GetOrders(pageSize, pageNumber, DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear, language)

}

func (s *DealerOrderService) GetProduct(DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear int, status bool, language string) ([]models.ReadProductForTradePoint, error) {
	return s.repo.GetProduct(DealerId, startDay, startMonth, startYear, endDay, endMonth, endYear, status, language)
}

func (s *DealerOrderService) GetOrderList(Id int, language string) (*[]models.ReadOrderList, error) {
	return s.repo.GetOrderList(Id, language)
}

func (s *DealerOrderService) DeleteDealerOrder(orderId int) error {
	return s.repo.DeleteDealerOrder(orderId)
}
