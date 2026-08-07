package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type ProductService struct {
	repo repository.Product
}

func NewProductService(repo repository.Product) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(product models.CreateProduct) (int, error) {
	return s.repo.Create(product)
}

func (s *ProductService) GetAll(pageSize, pageNumber int, name, language string) (*models.ProductPage, error) {
	return s.repo.GetAll(pageSize, pageNumber, name, language)
}

func (s *ProductService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *ProductService) GetById(productId int) (models.CreateProduct, error) {
	return s.repo.GetById(productId)
}

func (s *ProductService) GetTypes(language string) ([]models.SelectObject, error) {
	return s.repo.GetTypes(language)
}

func (s *ProductService) Update(productId int, product models.CreateProduct) error {
	return s.repo.Update(productId, product)
}
func (s *ProductService) Delete(productId int) error {
	return s.repo.Delete(productId)
}

func (s *ProductService) GetProductBrands() ([]models.SelectObject, error) {
	return s.repo.GetProductBrands()
}

func (s *ProductService) GetInformantionAboutSales(brandId, productTypeId, agentId, cityId, districtId, tradePointId, channelTypeId, startDay, startMonth, startYear, endDay, endMonth, endYear int, volume float64, product_taste, language string) ([]models.ReadProductForTradePoint, error) {
	return s.repo.GetInformantionAboutSales(brandId, productTypeId, agentId, cityId, districtId, tradePointId, channelTypeId, startDay, startMonth, startYear, endDay, endMonth, endYear, volume, product_taste, language)
}

func (s *ProductService) GetDailySalesProductsInExcel(agentId, day, month, year int, language string) ([]byte, error) {
	return s.repo.GetDailySalesProductsInExcel(agentId, day, month, year, language)
}
