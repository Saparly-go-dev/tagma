package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type FurnitureService struct {
	repo repository.Furniture
}

func NewFurnitureService(repo repository.Furniture) *FurnitureService {
	return &FurnitureService{repo: repo}
}

func (s *FurnitureService) GetMobile(trade_point_id int, language string) ([]models.ReadFurnitureForMobile, error) {
	return s.repo.GetMobile(trade_point_id, language)
}

func (s *FurnitureService) Get(name string, pageSize, pageNUmber int) (models.FurniturePage, error) {
	return s.repo.Get(name, pageSize, pageNUmber)
}

func (s *FurnitureService) Create(data models.CreateFurniture) error {
	return s.repo.Create(data)
}

func (s *FurnitureService) Update(id int, data models.CreateFurniture) error {
	return s.repo.Update(id, data)
}

func (s *FurnitureService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *FurnitureService) GetFurnitureList(name, language string) ([]models.SelectObject, error) {
	return s.repo.GetFurnitureList(name, language)
}

func (s *FurnitureService) GetTradePointFurnitures(trade_point_id int, language string) ([]models.TradePointFurnitureRead, error) {
	return s.repo.GetTradePointFurnitures(trade_point_id, language)
}

func (s *FurnitureService) CreateTradePointFurniture(trade_point_id, furniture_id, count int) error {
	return s.repo.CreateTradePointFurniture(trade_point_id, furniture_id, count)
}

func (s *FurnitureService) UpdateTradePointFurniture(id, furniture_id, count int) error {
	return s.repo.UpdateTradePointFurniture(id, furniture_id, count)
}

func (s *FurnitureService) DeleteTradePointFurniture(id int) error {
	return s.repo.DeleteTradePointFurniture(id)
}
