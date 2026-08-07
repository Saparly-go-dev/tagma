package service

import (
	"mime/multipart"

	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type TradePointService struct {
	repo repository.TradePoint
}

func NewTradePointService(repo repository.TradePoint) *TradePointService {

	return &TradePointService{repo: repo}
}

func (s *TradePointService) Create(point models.CreateTradePoint, image1, image2 *multipart.FileHeader) (int, error) {
	return s.repo.Create(point, image1, image2)
}

func (s *TradePointService) GetExcel(cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId int, name, language string) ([]byte, error) {
	return s.repo.GetExcel(cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId, name, language)
}

func (s *TradePointService) GetAll(pageSize, pageNumber, cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId, status int, name, code, contact, language string) (*models.TradePointPage, error) {
	return s.repo.GetAll(pageSize, pageNumber, cityId, districtId, agentId, dayId, ekspeditorId, typeId, structureId, sizeId, managementId, status, name, code, contact, language)
}

func (s *TradePointService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}

func (s *TradePointService) GetById(Id int) (models.TradePoint, error) {
	return s.repo.GetById(Id)
}

func (s *TradePointService) Update(Id int, point models.CreateTradePoint) error {
	return s.repo.Update(Id, point)
}

func (s *TradePointService) ChangeTradePointAgent(agent_id, ekspeditor_id, day_id int, ids []int) error {
	return s.repo.ChangeTradePointAgent(agent_id, ekspeditor_id, day_id, ids)
}

func (s *TradePointService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *TradePointService) GetGeneralInformation(Id int, language string) (*models.GeneralTradePoint, error) {
	return s.repo.GetGeneralInformation(Id, language)
}

func (s *TradePointService) DeleteImage(Id int) error {
	return s.repo.DeleteImage(Id)
}

func (s *TradePointService) GetCities(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetCities(name, language)
}
func (s *TradePointService) GetDistrictsByCity(cityId int, language string) (*[]models.SelectObject, error) {
	return s.repo.GetDistrictsByCity(cityId, language)
}

func (s *TradePointService) GetAgents(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetAgents(name, language)
}

func (s *TradePointService) GetChannels(name, language string) (*[]models.SelectObject, error) {
	return s.repo.GetChannels(name, language)
}

func (s *TradePointService) ReSaveImage(PointId int, file *multipart.FileHeader) error {
	return s.repo.ReSaveImage(PointId, file)
}

func (s *TradePointService) GetPointImages(Id int) (*[]models.ReadImage, error) {
	return s.repo.GetPointImages(Id)
}

func (s *TradePointService) GetInformationAboutOrder(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) ([]models.ReadOrder, error) {
	return s.repo.GetInformationAboutOrder(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear, language)
}
func (s *TradePointService) GetInformantionAboutProducts(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear int, language string) ([]models.ReadProductForTradePoint, error) {
	return s.repo.GetInformantionAboutProducts(tradePointId, startDay, startMonth, startYear, endDay, endMonth, endYear, language)
}
func (s *TradePointService) GetTradePointSaleHistoryForMonth(tradePointId int, language string) ([]models.SaleHistoryTradePoint, error) {
	return s.repo.GetTradePointSaleHistoryForMonth(tradePointId, language)
}
