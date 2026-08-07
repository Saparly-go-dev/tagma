package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type PromoService struct {
	repo repository.Promo
}

func NewPromoService(repo repository.Promo) *PromoService {
	return &PromoService{repo: repo}
}

func (s *PromoService) Create(promo models.UpdatePromo) error {
	return s.repo.Create(promo)
}

func (s *PromoService) Update(promo_id int, promo models.UpdatePromo) error {
	return s.repo.Update(promo_id, promo)
}

func (s *PromoService) GetPage(pageSize, pageNumber int, name, language string) (models.PromoPage, error) {
	return s.repo.GetPage(pageSize, pageNumber, name, language)
}

func (s *PromoService) ChangeStatus(id int) error {
	return s.repo.ChangeStatus(id)
}

func (s *PromoService) Delete(id int) error {
	return s.repo.Delete(id)
}
