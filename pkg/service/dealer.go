package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type DealerService struct {
	repo repository.Dealer
}

func NewDealerService(repo repository.Dealer) *DealerService {
	return &DealerService{repo: repo}
}

func (s *DealerService) Create(dealer models.CreateDealer) error {
	return s.repo.Create(dealer)
}

func (s *DealerService) GetPage(pageSize, pageNumber, cityId int, name, language string) (models.DealerPage, error) {
	return s.repo.GetPage(pageSize, pageNumber, cityId, name, language)
}

func (s *DealerService) Update(Id int, dealer models.CreateDealer) error {
	return s.repo.Update(Id, dealer)
}

func (s *DealerService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *DealerService) GetDealerForSelect(name, language string) ([]models.SelectObject, error) {
	return s.repo.GetDealerForSelect(name, language)
}

func (s *DealerService) TopUpDealerAccount(dealerId, price1, price2 int) error {
	return s.repo.TopUpDealerAccount(dealerId, price1, price2)
}
