package service

import (
	"fmt"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type TradeAgentService struct {
	repo repository.TradeAgent
}

func NewTradeAgentService(repo repository.TradeAgent) *TradeAgentService {
	fmt.Println("service geldi")
	return &TradeAgentService{repo: repo}
}

func (s *TradeAgentService) Create(agent models.CreateTradeAgent) (int, error) {
	return s.repo.Create(agent)
}

func (s *TradeAgentService) GetAll(pageSize, pageNumber int, name, language string) (*models.TradeAgentPage, error) {
	return s.repo.GetAll(pageSize, pageNumber, name, language)
}

func (s *TradeAgentService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *TradeAgentService) GetById(Id int) (models.TradeAgent, error) {
	return s.repo.GetById(Id)
}

func (s *TradeAgentService) Update(Id int, agent models.CreateTradeAgent) error {
	return s.repo.Update(Id, agent)
}
func (s *TradeAgentService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *TradeAgentService) GetDrivers(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetDrivers(name, language)
}

func (s *TradeAgentService) GetEkspeditors(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetEkspeditors(name, language)
}

func (s *TradeAgentService) AgentExcel(name, language string) ([]byte, error) {
	return s.repo.AgentExcel(name, language)
}

func (s *TradeAgentService) GetAgentDebt(AgentId int, language string) (models.AgentDebt, error) {
	return s.repo.GetAgentDebt(AgentId, language)
}

func (s *TradeAgentService) GetMerchendisingInformation(day, month, year, agentid int, language string) (models.PointsMerchendising_List, error) {
	return s.repo.GetMerchendisingInformation(day, month, year, agentid, language)
}

func (s *TradeAgentService) CreateCSVFileForFactory(day, month, year, agentid int, language string) ([]byte, error) {
	return s.repo.CreateCSVFileForFactory(day, month, year, agentid, language)
}

func (s *TradeAgentService) GetTradeAgentRevenue(month, year, agent_id int, language string) (float64, error) {
	return s.repo.GetTradeAgentRevenue(month, year, agent_id, language)
}
