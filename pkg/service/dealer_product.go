package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type DealerProductService struct {
	repo repository.DealerProduct
}

func NewDealerProductService(repo repository.DealerProduct) *DealerProductService {
	return &DealerProductService{repo: repo}
}

func (s *DealerProductService) Create(product models.CreateDealerProduct) (int, error) {
	return s.repo.Create(product)
}

func (s *DealerProductService) GetAll(pageSize, pageNumber int, name, language string) (*models.DealerProductPage, error) {
	return s.repo.GetAll(pageSize, pageNumber, name, language)
}

func (s *DealerProductService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *DealerProductService) Delete(productId int) error {
	return s.repo.Delete(productId)
}

func (s *DealerProductService) Update(DealerProductId int, DealerProduct models.CreateDealerProduct) error {
	return s.repo.Update(DealerProductId, DealerProduct)
}
