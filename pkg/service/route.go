package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type RouteService struct {
	repo repository.Route
}

func NewRouteService(repo repository.Route) *RouteService {
	return &RouteService{repo: repo}
}

func (s *RouteService) Create(route models.CreateRoute) (int, error) {
	return s.repo.Create(route)
}

func (s *RouteService) GetById(Id int, language string) (*models.ReadRoute, error) {
	return s.repo.GetById(Id, language)
}

func (s *RouteService) GetAll(pageSize, pageNumber, agentId, cityId, dayId int, language string) (*models.RoutePage, error) {
	return s.repo.GetAll(pageSize, pageNumber, agentId, cityId, dayId, language)
}

func (s *RouteService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *RouteService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *RouteService) Update(Id int, route models.CreateRoute) error {
	return s.repo.Update(Id, route)
}
func (s *RouteService) GetDays(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetDays(name, language)
}

func (s *RouteService) GetPoints(cityId, districtId int, name, language string) (*[]models.PointSelect, error) {
	return s.repo.GetPoints(cityId, districtId, name, language)
}
func (s *RouteService) GetEkspeditors(TradeAgentId int, language string) (*[]models.SelectObject, error) {
	return s.repo.GetEkspeditors(TradeAgentId, language)
}

func (s *RouteService) GetRouteOrder(agent_id, day_id int, language string) ([]models.SelectObject, error) {
	return s.repo.GetRouteOrder(agent_id, day_id, language)
}

func (s *RouteService) PostRouteOrder(agent_id, day_id int, list []models.SelectObject) error {
	return s.repo.PostRouteOrder(agent_id, day_id, list)
}

func (s *RouteService) ExcelRouteOrder(agent_id, day_id int, language string) ([]byte, error) {
	return s.repo.ExcelRouteOrder(agent_id, day_id, language)
}
