package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type ContactService struct {
	repo repository.Contact
}

func NewContactService(repo repository.Contact) *ContactService {
	return &ContactService{repo: repo}
}

func (s *ContactService) Create(contact models.CreateContact) (int, error) {
	return s.repo.Create(contact)
}

func (s *ContactService) GetAll(PointId int, language string) (*[]models.ReadContact, error) {
	return s.repo.GetAll(PointId, language)
}

func (s *ContactService) GetById(Id int) (*models.Contact, error) {
	return s.repo.GetById(Id)
}

func (s *ContactService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *ContactService) Update(Id int, contact models.CreateContact) error {
	return s.repo.Update(Id, contact)
}

func (s *ContactService) GetKinds(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetKinds(name, language)
}

func (s *ContactService) GetPosts(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetPosts(name, language)
}

func (s *ContactService) GetUprs(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetUprs(name, language)
}
