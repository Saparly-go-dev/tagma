package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type ProductTypeService struct {
	repo repository.ProductType
}

func NewProductTypeService(repo repository.ProductType) *ProductTypeService {
	return &ProductTypeService{repo: repo}
}

func (s *ProductTypeService) Create(productType models.CreateProductType) (int, error) {
	return s.repo.Create(productType)
}

func (s *ProductTypeService) GetAll(pageSize, pageNumber int, name string) (*models.ProductTypePage, error) {
	return s.repo.GetAll(pageSize, pageNumber, name)
}

func (s *ProductTypeService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *ProductTypeService) GetById(productId int) (models.ProductType, error) {
	return s.repo.GetById(productId)
}

func (s *ProductTypeService) Update(productId int, productType models.CreateProductType) error {
	return s.repo.Update(productId, productType)
}
func (s *ProductTypeService) Delete(productId int) error {
	return s.repo.Delete(productId)
}

func (S *ProductTypeService) GetAllVolumes() (*[]models.Volume, error) {
	return S.repo.GetAllVolumes()
}
