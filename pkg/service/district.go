package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type DistrictService struct {
	repo repository.District
}

func NewDistrictService(repo repository.District) *DistrictService {
	return &DistrictService{repo: repo}
}

func (s *DistrictService) Create(district models.CreateDistrict) error {
	return s.repo.Create(district)
}

func (s *DistrictService) GetPage(pageSize, pageNumber int, name, language string) (*models.DistrictPage, error) {
	return s.repo.GetPage(pageSize, pageNumber, name, language)
}

func (s *DistrictService) GetById(Id int) (*models.District, error) {
	return s.repo.GetById(Id)
}

func (s *DistrictService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *DistrictService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *DistrictService) Update(Id int, district models.CreateDistrict) error {
	return s.repo.Update(Id, district)
}
