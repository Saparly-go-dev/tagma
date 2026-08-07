package service

import (
	"github.com/Saparly-go-dev/tagma/models"
	"github.com/Saparly-go-dev/tagma/pkg/repository"
)

type Trade_ChannelService struct {
	repo repository.Trade_Channel
}

func NewTradeChannelService(repo repository.Trade_Channel) *Trade_ChannelService {
	return &Trade_ChannelService{repo: repo}
}
func (s Trade_ChannelService) CreateType(tip models.ChannelObject) error {
	return s.repo.CreateType(tip)
}
func (s Trade_ChannelService) GetAllTypes() ([]models.ChannelType, error) {
	return s.repo.GetAllTypes()
}
func (s Trade_ChannelService) DeleteTypes(Id int) error {
	return s.repo.DeleteTypes(Id)
}
func (s Trade_ChannelService) UpdateTypes(Id int, tip models.ChannelObject) error {
	return s.repo.UpdateTypes(Id, tip)
}
func (s Trade_ChannelService) CreateStructure(structure models.ChannelObject) error {
	return s.repo.CreateStructure(structure)
}
func (s Trade_ChannelService) GetAllStructure() ([]models.ChannelStructure, error) {
	return s.repo.GetAllStructure()
}
func (s Trade_ChannelService) DeleteStructure(Id int) error {
	return s.repo.DeleteStructure(Id)
}
func (s Trade_ChannelService) UpdateStructure(Id int, structure models.ChannelObject) error {
	return s.repo.UpdateStructure(Id, structure)
}
func (s Trade_ChannelService) CreateSize(size models.ChannelObject) error {
	return s.repo.CreateSize(size)
}
func (s Trade_ChannelService) GetAllSizes() ([]models.ChannelSize, error) {
	return s.repo.GetAllSizes()
}
func (s Trade_ChannelService) DeleteSize(Id int) error {
	return s.repo.DeleteSize(Id)
}
func (s Trade_ChannelService) UpdateSize(Id int, size models.ChannelObject) error {
	return s.repo.UpdateSize(Id, size)
}
func (s Trade_ChannelService) CreateManagement(management models.ChannelObject) error {
	return s.repo.CreateManagement(management)
}
func (s Trade_ChannelService) GetAllManagements() ([]models.ChannelManagement, error) {
	return s.repo.GetAllManagements()
}
func (s Trade_ChannelService) DeleteManagement(Id int) error {
	return s.repo.DeleteManagement(Id)
}
func (s Trade_ChannelService) UpdateManagement(Id int, management models.ChannelObject) error {
	return s.repo.UpdateManagement(Id, management)
}
func (s Trade_ChannelService) Create(channel models.CreateChannel) error {
	return s.repo.Create(channel)
}
func (s Trade_ChannelService) GetChannelPage(pageSize, pageNumber, typeId int, language string) (*models.TradeChannelPage, error) {
	return s.repo.GetChannelPage(pageSize, pageNumber, typeId, language)
}
func (s Trade_ChannelService) ChangeStatus(Id int) error {
	return s.repo.ChangeStatus(Id)
}
func (s Trade_ChannelService) Delete(Id int) error {
	return s.repo.Delete(Id)
}
func (s Trade_ChannelService) Update(Id int, channel models.CreateChannel) error {
	return s.repo.Update(Id, channel)
}
func (s Trade_ChannelService) GetTypeForChannel(language string) (*[]models.SelectObject, error) {
	return s.repo.GetTypeForChannel(language)
}
func (s Trade_ChannelService) GetStructureForChannel(language string) (*[]models.SelectObject, error) {
	return s.repo.GetStructureForChannel(language)
}
func (s Trade_ChannelService) GetSizeForChannel(language string) (*[]models.SelectObject, error) {
	return s.repo.GetSizeForChannel(language)
}
func (s Trade_ChannelService) GetManagementForChannel(language string) (*[]models.SelectObject, error) {
	return s.repo.GetManagementForChannel(language)
}
