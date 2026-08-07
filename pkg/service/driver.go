package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type DriverService struct {
	repo repository.Driver
}

func NewDriverService(repo repository.Driver) *DriverService {
	return &DriverService{repo: repo}
}

func (s *DriverService) Create(driver models.CreateDriver) (int, error) {
	return s.repo.Create(driver)
}

func (s *DriverService) GetAll(pageNumber, pageSize int, name, language string) (*models.DriverPage, error) {
	return s.repo.GetAll(pageNumber, pageSize, name, language)
}

func (s *DriverService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *DriverService) GetById(Id int) (models.Driver, error) {
	return s.repo.GetById(Id)
}

func (s *DriverService) Update(Id int, driver models.CreateDriver) error {
	return s.repo.Update(Id, driver)
}

func (s *DriverService) Delete(Id int) error {
	return s.repo.Delete(Id)
}
