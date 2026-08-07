package service

import (
	"errors"
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type OrderService struct {
	repo repository.Order
}

func NewOrderService(repo repository.Order) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) GetTradePoints(day, month, year, agentid int, language string) (*[]models.TradePointMobile, error) {
	return s.repo.GetTradePoints(day, month, year, agentid, language)
}

func (s *OrderService) GetTradePointsForMarketing(day, month, year, agentid int, language string) (*[]models.TradePointMobile, error) {
	return s.repo.GetTradePointsForMarketing(day, month, year, agentid, language)
}

func (s *OrderService) GetProducts(language string) (*[]models.ProductDistinct, error) {
	return s.repo.GetProducts(language)
}

func (s *OrderService) GetBrand(language string) (*[]string, error) {
	return s.repo.GetBrand(language)
}

func (s *OrderService) GetModel(brandId, ProductTypeId int, language string) (*[]string, error) {
	return s.repo.GetModel(brandId, ProductTypeId, language)
}

func (s *OrderService) GetTypes(language string) (*[]string, error) {
	return s.repo.GetTypes(language)
}

func (s *OrderService) GetVolume(brand, model, tip string) (*[]models.SelectObject, error) {
	return s.repo.GetVolume(brand, model, tip)
}

func (s *OrderService) GetSum(Id, Count int) (float64, error) {
	return s.repo.GetSum(Id, Count)
}

func (s *OrderService) SaveOrder(day, month, year, tradePointId, paymentTypeId int, first, is_credit bool, orders []models.CreateOrder) error {
	if tradePointId <= 0 || paymentTypeId <= 0 {
		return errors.New("trade point and payment type are required")
	}
	if _, err := validDate(day, month, year); err != nil {
		return err
	}
	if len(orders) == 0 {
		return errors.New("order must contain at least one product")
	}
	for _, order := range orders {
		if order.Id <= 0 || order.Count <= 0 {
			return errors.New("product id and count must be positive")
		}
	}
	return s.repo.SaveOrder(day, month, year, tradePointId, paymentTypeId, first, is_credit, orders)
}

func validDate(day, month, year int) (time.Time, error) {
	if year < 2000 || month < 1 || month > 12 || day < 1 || day > 31 {
		return time.Time{}, errors.New("invalid order date")
	}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if date.Day() != day || int(date.Month()) != month || date.Year() != year {
		return time.Time{}, errors.New("invalid order date")
	}
	return date, nil
}

func (s *OrderService) Marketing(day, month, year, tradePointId int, list []models.MarketingPost) error {
	return s.repo.Marketing(day, month, year, tradePointId, list)
}

func (s *OrderService) GetOrderPage(pageSize, pageNumber, tradeAgentId, cityId, districtId, day, month, year int, name, language, location, code, order string, status, is_status_active bool) (*models.OrderPage, error) {
	return s.repo.GetOrderPage(pageSize, pageNumber, tradeAgentId, cityId, districtId, day, month, year, name, language, location, code, order, status, is_status_active)
}

func (s *OrderService) OrderExcel(tradeAgentId, cityId, districtId int, name, language string) ([]byte, error) {
	return s.repo.OrderExcel(tradeAgentId, cityId, districtId, name, language)
}

func (s *OrderService) GetOrderList(OrderId int, language string) (*[]models.ReadOrderList, error) {
	return s.repo.GetOrderList(OrderId, language)
}

func (s *OrderService) GetOrders(day, month, year, tradePointId int, language string) (models.GetOrderObject, error) {
	return s.repo.GetOrders(day, month, year, tradePointId, language)
}

func (s *OrderService) CloseTheOrder(orderId int) error {
	return s.repo.CloseTheOrder(orderId)
}
