package service

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type Trade_CategoryService struct {
	repo repository.Trade_Category
}

func NewTradeCategoryService(repo repository.Trade_Category) *Trade_CategoryService {
	return &Trade_CategoryService{repo: repo}
}

func (s *Trade_CategoryService) Create(category tagma.CreateTrade_category) (int, error) {
	return s.repo.Create(category)
}

func (s *Trade_CategoryService) GetAll(pageNumber, pageSize int, name string) (*tagma.TradeCategoryPage, error) {
	return s.repo.GetAll(pageNumber, pageSize, name)
}

func (s *Trade_CategoryService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *Trade_CategoryService) GetById(Id int) (tagma.Trade_category, error) {
	return s.repo.GetById(Id)
}

func (s *Trade_CategoryService) Update(Id int, category tagma.CreateTrade_category) error {
	return s.repo.Update(Id, category)
}
func (s *Trade_CategoryService) Delete(Id int) error {
	return s.repo.Delete(Id)
}
