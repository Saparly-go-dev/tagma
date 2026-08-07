package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
	"mime/multipart"
)

type ManualService struct {
	repo repository.Manual
}

func NewManualService(repo repository.Manual) *ManualService {
	return &ManualService{repo: repo}
}

func (s *ManualService) Create(data models.CreateManual) error {
	return s.repo.Create(data)
}

func (s *ManualService) Get(language string) (*models.ReadManual, error) {
	return s.repo.Get(language)
}

func (s *ManualService) UpdateManual(Id int, data models.CreateManual) error {
	return s.repo.UpdateManual(Id, data)
}

func (s *ManualService) DeleteImage(Id int) error {
	return s.repo.DeleteImage(Id)
}

func (s *ManualService) SaveImage(Id int, file *multipart.FileHeader) error {
	return s.repo.SaveImage(Id, file)
}

func (s *ManualService) GetForMobile(language string) (*models.ReadManualForMobile, error) {
	return s.repo.GetForMobile(language)
}
