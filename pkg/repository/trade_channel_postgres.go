package repository

import (
	"time"

	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
)

type TradeChannelPostgres struct {
	db *gorm.DB
}

// NewTradeChannelPostgres initializes a new TradeChannelPostgres instance
func NewTradeChannelPostgres(db *gorm.DB) *TradeChannelPostgres {
	return &TradeChannelPostgres{db: db}
}

func (r *TradeChannelPostgres) CreateType(tip models.ChannelObject) error {
	newType := models.ChannelType{
		NameRu: tip.NameRu,
		NameTm: tip.NameTm,
	}

	result := r.db.Model(&models.ChannelType{}).Create(&newType)

	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *TradeChannelPostgres) GetAllTypes() ([]models.ChannelType, error) {
	var types []models.ChannelType
	result := r.db.Model(&models.ChannelType{}).Order("id").Find(&types)
	if result.Error != nil {
		return nil, result.Error
	}
	return types, nil
}

func (r *TradeChannelPostgres) DeleteTypes(Id int) error {
	result := r.db.Delete(&models.ChannelType{}, Id)
	return result.Error
}

func (r *TradeChannelPostgres) UpdateTypes(Id int, tip models.ChannelObject) error {
	updateData := map[string]interface{}{
		"name_ru": tip.NameRu,
		"name_tm": tip.NameTm,
	}

	result := r.db.Model(&models.ChannelType{}).Where("id = ?", Id).Updates(updateData)

	return result.Error
}

func (r *TradeChannelPostgres) CreateStructure(structure models.ChannelObject) error {
	newStructure := models.ChannelStructure{
		NameRu: structure.NameRu,
		NameTm: structure.NameTm,
	}

	result := r.db.Model(&models.ChannelStructure{}).Create(&newStructure)

	return result.Error
}

func (r *TradeChannelPostgres) GetAllStructure() ([]models.ChannelStructure, error) {
	var structures []models.ChannelStructure

	result := r.db.Model(&models.ChannelStructure{}).Order("id").Find(&structures)

	return structures, result.Error
}

func (r *TradeChannelPostgres) DeleteStructure(Id int) error {
	result := r.db.Delete(&models.ChannelStructure{}, Id)

	return result.Error
}

func (r *TradeChannelPostgres) UpdateStructure(Id int, structure models.ChannelObject) error {
	updateData := map[string]interface{}{
		"name_ru": structure.NameRu,
		"name_tm": structure.NameTm,
	}

	result := r.db.Model(&models.ChannelStructure{}).Where("id = ?", Id).Updates(updateData)

	return result.Error
}

func (r *TradeChannelPostgres) CreateSize(size models.ChannelObject) error {
	newSize := models.ChannelSize{
		NameRu: size.NameRu,
		NameTm: size.NameTm,
	}

	result := r.db.Model(&models.ChannelSize{}).Create(&newSize)
	return result.Error
}

func (r *TradeChannelPostgres) GetAllSizes() ([]models.ChannelSize, error) {
	var sizes []models.ChannelSize
	result := r.db.Model(&models.ChannelSize{}).Order("id").Find(&sizes)

	return sizes, result.Error
}

func (r *TradeChannelPostgres) DeleteSize(Id int) error {
	result := r.db.Delete(&models.ChannelSize{}, Id)
	return result.Error
}

func (r *TradeChannelPostgres) UpdateSize(Id int, size models.ChannelObject) error {
	updateData := map[string]interface{}{
		"name_ru": size.NameRu,
		"name_tm": size.NameTm,
	}

	result := r.db.Model(&models.ChannelSize{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *TradeChannelPostgres) CreateManagement(management models.ChannelObject) error {
	newManagement := models.ChannelManagement{
		NameRu: management.NameRu,
		NameTm: management.NameTm,
	}

	result := r.db.Model(&models.ChannelManagement{}).Create(&newManagement)
	return result.Error
}

func (r *TradeChannelPostgres) GetAllManagements() ([]models.ChannelManagement, error) {
	var managements []models.ChannelManagement

	result := r.db.Model(&models.ChannelManagement{}).Order("id").Find(&managements)
	return managements, result.Error
}

func (r *TradeChannelPostgres) DeleteManagement(Id int) error {
	result := r.db.Delete(&models.ChannelManagement{}, Id)
	return result.Error
}

func (r *TradeChannelPostgres) UpdateManagement(Id int, management models.ChannelObject) error {
	updateData := map[string]interface{}{
		"name_ru": management.NameRu,
		"name_tm": management.NameTm,
	}

	result := r.db.Model(&models.ChannelManagement{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *TradeChannelPostgres) Create(channel models.CreateChannel) error {
	newChannel := models.TradeChannel{
		TypeId:       channel.TypeId,
		StructureId:  channel.StructureId,
		SizeId:       channel.SizeId,
		ManagementId: channel.ManagementId,
		Status:       channel.Status,
		UpdatedAt:    time.Now(),
	}

	result := r.db.Model(&models.TradeChannel{}).Create(&newChannel)
	return result.Error
}

func (r *TradeChannelPostgres) GetChannelPage(pageSize, pageNumber, typeId int, language string) (*models.TradeChannelPage, error) {
	var rawData []models.TradeChannel
	var pageItems []models.ReadChannel

	// Calculate offset for pagination
	offset := (pageNumber - 1) * pageSize

	// Use correct join clause for product and product type, assuming a foreign key like `product_type_id` exists

	result := r.db.Model(&models.TradeChannel{}).
		Preload("Type").
		Preload("Structure").
		Preload("Size").
		Preload("Management").
		Joins("JOIN channel_types ON channel_types.id = trade_channels.type_id").
		Order("trade_channels.id")

	if typeId > 0 {
		result = result.Where("type_id = ?", typeId)
	}

	var count int64
	result.Count(&count)

	result = result.Offset(offset).Limit(pageSize).Find(&rawData)

	if result.Error != nil {
		return nil, result.Error
	}

	// Calculate total count for pagination

	for _, item := range rawData {
		var newItem models.ReadChannel

		newItem.Id = item.Id
		newItem.TypeId = item.TypeId
		newItem.StructureId = item.StructureId
		newItem.SizeId = item.SizeId
		newItem.ManagementId = item.ManagementId
		newItem.Status = item.Status
		newItem.UpdatedAt = item.UpdatedAt.Format("02.01.2006")
		if language == "ru" {
			newItem.Type = item.Type.NameRu
			newItem.Structure = item.Structure.NameRu
			newItem.Size = item.Size.NameRu
			newItem.Management = item.Management.NameRu
		} else {
			newItem.Type = item.Type.NameTm
			newItem.Structure = item.Structure.NameTm
			newItem.Size = item.Size.NameTm
			newItem.Management = item.Management.NameTm
		}

		pageItems = append(pageItems, newItem)
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if pageItems == nil {
		pageItems = []models.ReadChannel{}
	}

	// Create the page object
	page := models.TradeChannelPage{
		Items:      pageItems,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	return &page, nil
}

func (r *TradeChannelPostgres) ChangeStatus(Id int) error {
	var data models.TradeChannel

	result := r.db.Model(&models.TradeChannel{}).Where("id = ?", Id).First(&data)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.TradeChannel{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"updated_at": time.Now(),
		}
		result := r.db.Model(&models.TradeChannel{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *TradeChannelPostgres) Delete(Id int) error {
	result := r.db.Delete(&models.TradeChannel{}, Id)
	return result.Error
}

func (r *TradeChannelPostgres) Update(Id int, channel models.CreateChannel) error {
	updateData := map[string]interface{}{
		"type_id":       channel.TypeId,
		"structure_id":  channel.StructureId,
		"size_id":       channel.SizeId,
		"management_id": channel.ManagementId,
		"status":        channel.Status,
		"updated_at":    time.Now(),
	}
	result := r.db.Model(&models.TradeChannel{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *TradeChannelPostgres) GetTypeForChannel(language string) (*[]models.SelectObject, error) {
	var rawData []models.ChannelType
	var resultData []models.SelectObject

	result := r.db.Model(&models.ChannelType{}).Order("id asc").Find(&rawData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, item := range rawData {
		var newItem models.SelectObject
		newItem.Id = item.Id
		if language == "ru" {
			newItem.Name = item.NameRu
		} else {
			newItem.Name = item.NameTm
		}
		resultData = append(resultData, newItem)
	}

	return &resultData, nil
}

func (r *TradeChannelPostgres) GetStructureForChannel(language string) (*[]models.SelectObject, error) {
	var rawData []models.ChannelStructure
	var resultData []models.SelectObject

	result := r.db.Model(&models.ChannelStructure{}).Order("id asc").Find(&rawData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, item := range rawData {
		var newItem models.SelectObject
		newItem.Id = item.Id
		if language == "ru" {
			newItem.Name = item.NameRu
		} else {
			newItem.Name = item.NameTm
		}
		resultData = append(resultData, newItem)
	}

	return &resultData, nil
}

func (r *TradeChannelPostgres) GetSizeForChannel(language string) (*[]models.SelectObject, error) {
	var rawData []models.ChannelSize
	var resultData []models.SelectObject

	result := r.db.Model(&models.ChannelSize{}).Order("id asc").Find(&rawData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, item := range rawData {
		var newItem models.SelectObject
		newItem.Id = item.Id
		if language == "ru" {
			newItem.Name = item.NameRu
		} else {
			newItem.Name = item.NameTm
		}
		resultData = append(resultData, newItem)
	}

	return &resultData, nil
}

func (r *TradeChannelPostgres) GetManagementForChannel(language string) (*[]models.SelectObject, error) {
	var rawData []models.ChannelManagement
	var resultData []models.SelectObject

	result := r.db.Model(&models.ChannelManagement{}).Order("id asc").Find(&rawData)

	if result.Error != nil {
		return nil, result.Error
	}

	for _, item := range rawData {
		var newItem models.SelectObject
		newItem.Id = item.Id
		if language == "ru" {
			newItem.Name = item.NameRu
		} else {
			newItem.Name = item.NameTm
		}
		resultData = append(resultData, newItem)
	}

	return &resultData, nil
}
