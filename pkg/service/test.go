package service

import (
	"github.com/Saparly-go-dev/tagma"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type TestService struct {
	repo repository.Test
}

func NewTestService(repo repository.Test) *TestService {
	return &TestService{repo: repo}
}

func (s *TestService) Create(test tagma.Test) (int, error) {
	return s.repo.Create(test)
}

func (s *TestService) GetAll() ([]tagma.Test, error) {
	return s.repo.GetAll()
}
