package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type GiftService struct {
	repo repository.Gift
}

func NewGiftService(repo repository.Gift) *GiftService {
	return &GiftService{repo: repo}
}

func (s *GiftService) Create(gift models.CreateGift) error {
	return s.repo.Create(gift)
}

func (s *GiftService) GetPage(pageSize, pageNumber int, name, language string) (*models.GiftPage, error) {
	return s.repo.GetPage(pageSize, pageNumber, name, language)
}

func (s *GiftService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *GiftService) GetById(Id int) (*models.Gift, error) {
	return s.repo.GetById(Id)
}

func (s *GiftService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *GiftService) Update(Id int, gift models.CreateGift) error {
	return s.repo.Update(Id, gift)
}

func (s *GiftService) GetGiftTypes(language string) (*[]models.SelectObject, error) {
	return s.repo.GetGiftTypes(language)
}
