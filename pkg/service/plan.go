package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type PlanService struct {
	repo repository.Plan
}

func NewPlanService(repo repository.Plan) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) Create(list []models.CreatePlan, month, year int) error {
	return s.repo.Create(list, month, year)
}

func (s *PlanService) GetPlans(month, year int, language string) ([]models.ReadPlan, error) {
	return s.repo.GetPlans(month, year, language)
}

func (s *PlanService) GetPlanStatisticsForMain(month, year, agent_id int, language string) ([]models.PlanStatisticsMain, error) {
	return s.repo.GetPlanStatisticsForMain(month, year, agent_id, language)
}

func (s *PlanService) GetPlansForAgent(month, year, agent_id int, language string) ([]models.PlanStatisticsForAgent, error) {
	return s.repo.GetPlansForAgent(month, year, agent_id, language)
}

func (s *PlanService) Update(month, year int, list []models.UpdatedPlan) error {
	return s.repo.Update(month, year, list)
}
