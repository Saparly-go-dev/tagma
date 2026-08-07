package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type PromotionService struct {
	repo repository.Promotion
}

func NewPromotionService(repo repository.Promotion) *PromotionService {
	return &PromotionService{repo: repo}
}

func (s *PromotionService) Create(promotion models.CreatePromotion) error {
	return s.repo.Create(promotion)
}

func (s *PromotionService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *PromotionService) GetAll(pageSize, pageNumber int, name, language string) (*models.PromotionPage, error) {
	return s.repo.GetAll(pageSize, pageNumber, name, language)
}

func (s *PromotionService) Update(id int, data models.CreatePromotion) error {
	return s.repo.Update(id, data)
}

func (s *PromotionService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *PromotionService) UpdateNameAndTime(id int, data models.UpdateDiscountInfo) error {
	return s.repo.UpdateNameAndTime(id, data)
}
