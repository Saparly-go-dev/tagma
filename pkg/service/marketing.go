package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type MarketingService struct {
	repo repository.Marketing
}

func NewMarketingService(repo repository.Marketing) *MarketingService {
	return &MarketingService{repo: repo}
}

func (s *MarketingService) GetMarketingMainPage(cityId, districtId, tradePointId, agentId, brand int, volume float64, taste, language string) ([]models.MarketingMainPage, error) {
	return s.repo.GetMarketingMainPage(cityId, districtId, tradePointId, agentId, brand, volume, taste, language)
}

func (s *MarketingService) MarketingMainPageExport(cityId, districtId, tradePointId, agentId, brand int, volume float64, taste, language string) ([]byte, error) {
	return s.repo.MarketingMainPageExport(cityId, districtId, tradePointId, agentId, brand, volume, taste, language)
}

func (s *MarketingService) GetDistribution(tradePointId int, language string) ([]models.ReadMarketingById, error) {
	return s.repo.GetDistribution(tradePointId, language)
}

func (s *MarketingService) DistributionExport(tradePointId int, language string) ([]byte, error) {
	return s.repo.DistributionExport(tradePointId, language)
}

func (s *MarketingService) GetRasprodano(tradePointId int, language string) ([]models.ReadMarketingById, error) {
	return s.repo.GetRasprodano(tradePointId, language)
}

func (s *MarketingService) RasprodanoExport(tradePointId int, language string) ([]byte, error) {
	return s.repo.RasprodanoExport(tradePointId, language)
}

func (s *MarketingService) GetWystawlaymost(tradePointId int, language string) ([]models.ReadMarketingById, error) {
	return s.repo.GetWystawlaymost(tradePointId, language)
}

func (s *MarketingService) WystawlaymostExport(tradePointId int, language string) ([]byte, error) {
	return s.repo.WystawlaymostExport(tradePointId, language)
}
