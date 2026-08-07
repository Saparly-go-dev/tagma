package service

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type CityService struct {
	repo repository.City
}

func NewCityService(repo repository.City) *CityService {
	return &CityService{repo: repo}
}

func (s *CityService) Create(city tagma.CreateCity) (int, error) {
	return s.repo.Create(city)
}

func (s *CityService) GetAll(pageSize, pageNumber int, name string) (*tagma.CityPage, error) {
	return s.repo.GetAll(pageSize, pageNumber, name)
}

func (s *CityService) ChangeStatus(cityId int) error {
	return s.repo.ChangeStatus(cityId)
}

func (s *CityService) GetById(cityId int) (tagma.City, error) {
	return s.repo.GetById(cityId)
}

func (s *CityService) Update(cityId int, city tagma.CreateCity) error {
	return s.repo.Update(cityId, city)
}
func (s *CityService) Delete(cityId int) error {
	return s.repo.Delete(cityId)
}
