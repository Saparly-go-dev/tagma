package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type MerchandiserService struct {
	repo repository.Merchandiser
}

func NewMerchandiserService(repo repository.Merchandiser) *MerchandiserService {
	return &MerchandiserService{repo: repo}
}

func (s *MerchandiserService) Create(merchandiser models.CreateMerchandiser) (int, error) {
	return s.repo.Create(merchandiser)
}

func (s *MerchandiserService) GetAll(pageNumber, pageSize int, name, language string) (*models.MerchandiserPage, error) {
	return s.repo.GetAll(pageNumber, pageSize, name, language)
}

func (s *MerchandiserService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *MerchandiserService) GetById(Id int) (models.Merchandiser, error) {
	return s.repo.GetById(Id)
}

func (s *MerchandiserService) Update(Id int, merchandiser models.CreateMerchandiser) error {
	return s.repo.Update(Id, merchandiser)
}

func (s *MerchandiserService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *MerchandiserService) GetMerchandiserForSelect(language, name string) ([]models.SelectObject, error) {
	return s.repo.GetMerchandiserForSelect(language, name)
}

func (s *MerchandiserService) GetMarketingPageForMobile(day, month, year, userId int, language string) (*[]models.TradePointMobile, error) {
	return s.repo.GetMarketingPageForMobile(day, month, year, userId, language)
}
