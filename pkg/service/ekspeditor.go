package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type EkspeditorService struct {
	repo repository.Ekspeditor
}

func NewEkspeditorService(repo repository.Ekspeditor) *EkspeditorService {
	return &EkspeditorService{repo: repo}
}

func (s *EkspeditorService) Create(ekspeditor models.CreateEkspeditor) (int, error) {
	return s.repo.Create(ekspeditor)
}

func (s *EkspeditorService) GetAll(pageNumber, pageSize int, name, language string) (*models.EkspeditorPage, error) {
	return s.repo.GetAll(pageNumber, pageSize, name, language)
}

func (s *EkspeditorService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *EkspeditorService) GetById(Id int) (models.Ekspeditor, error) {
	return s.repo.GetById(Id)
}

func (s *EkspeditorService) Update(Id int, ekspeditor models.CreateEkspeditor) error {
	return s.repo.Update(Id, ekspeditor)
}

func (s *EkspeditorService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *EkspeditorService) GetProductCount(day, month, year, ekspeditorId int) (*map[int]int, error) {
	return s.repo.GetProductCount(day, month, year, ekspeditorId)
}

func (s *EkspeditorService) GetPointProductList(day, month, year, tradePointId int, language string) (models.GetOrderObject, error) {
	return s.repo.GetPointProductList(day, month, year, tradePointId, language)
}

func (s *EkspeditorService) SaveProductList(day, month, year, ekspeditorId int, orders []models.CreateOrder) error {
	return s.repo.SaveProductList(day, month, year, ekspeditorId, orders)
}

func (s *EkspeditorService) GetTradePoints(day, month, year, ekspeditorId int, language string) (*[]models.TradePointMobile, error) {
	return s.repo.GetTradePoints(day, month, year, ekspeditorId, language)
}

func (s *EkspeditorService) OrderStatus(day, month, year, tradePointId, ekspeditorId int) error {
	return s.repo.OrderStatus(day, month, year, tradePointId, ekspeditorId)
}

func (s *EkspeditorService) EkspeditorProducts(day, month, year, ekspeditorId int) (*[]models.EkspData, error) {
	return s.repo.EkspeditorProducts(day, month, year, ekspeditorId)
}

func (s *EkspeditorService) GetGalyndy(ekspeditorId int) (*[]models.CreateOrder, error) {
	return s.repo.GetGalyndy(ekspeditorId)
}

func (s *EkspeditorService) GetEkspeditorGalyndy(ekspeditorId int, language string) ([]models.ReadProductForTradePoint, error) {
	return s.repo.GetEkspeditorGalyndy(ekspeditorId, language)
}

func (s *EkspeditorService) DeleEkspData() error {
	return s.repo.DeleEkspData()
}

func (s *EkspeditorService) GetProductCountWithMinus(day, month, year, ekspeditorId int) (*map[int]int, error) {
	return s.repo.GetProductCountWithMinus(day, month, year, ekspeditorId)
}

func (s *EkspeditorService) GetProductInTheCarWithDifference(ekspeditorId int, language string) ([]models.ReadProductForTradePoint, error) {
	return s.repo.GetProductInTheCarWithDifference(ekspeditorId, language)
}

func (s *EkspeditorService) GetProductInTheCarWithDifferenceInArray(language string) ([]models.GetGalyndyInArray, error) {
	return s.repo.GetProductInTheCarWithDifferenceInArray(language)
}

func (s *EkspeditorService) GetProductPriceInTheCarWithDifference(ekspeditorId int) (float64, error) {
	return s.repo.GetProductPriceInTheCarWithDifference(ekspeditorId)
}
