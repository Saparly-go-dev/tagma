package repository

import (
	"errors"
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"time"
)

type FurniturePostgres struct {
	db *gorm.DB
}

func NewFurniturePostgres(db *gorm.DB) *FurniturePostgres {
	return &FurniturePostgres{db: db}
}

func (r *FurniturePostgres) GetMobile(trade_point_id int, language string) ([]models.ReadFurnitureForMobile, error) {
	result := []models.ReadFurnitureForMobile{}

	var data_list []models.TradePointsFurniture

	get_request := r.db.Model(&models.TradePointsFurniture{}).Preload("Furniture").
		Where("trade_point_id = ?", trade_point_id).Find(&data_list)

	if get_request.Error != nil {
		return result, get_request.Error
	}

	for _, element := range data_list {
		var item models.ReadFurnitureForMobile
		item.Count = element.Count

		if language == "ru" {
			item.Name = element.Furniture.NameRu
		} else {
			item.Name = element.Furniture.NameTm
		}
		result = append(result, item)
	}
	return result, nil
}

func (r *FurniturePostgres) Get(name string, pageSize, pageNUmber int) (models.FurniturePage, error) {
	result := models.FurniturePage{}

	var furniture_list []models.Furniture
	var result_items []models.ReadFurniture
	var count int64
	offset := (pageNUmber - 1) * pageSize

	get_request := r.db.Model(&models.Furniture{})

	if len(name) > 0 {
		get_request = get_request.Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%")
	}

	get_request.Count(&count)
	get_request = get_request.Offset(offset).Limit(pageSize).Find(&furniture_list)

	if get_request.Error != nil {
		return result, get_request.Error
	}

	for _, element := range furniture_list {
		var data models.ReadFurniture
		data.Id = element.Id
		data.NameRu = element.NameRu
		data.NameTm = element.NameTm
		data.CreatedAt = element.CreatedAt.Format("02.01.2006")

		result_items = append(result_items, data)
	}

	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	result.Items = result_items
	result.PageSize = pageSize
	result.PageNumber = pageNUmber
	result.PageCount = int(pageCount)
	result.Total = int(count)

	return result, nil
}

func (r *FurniturePostgres) Create(data models.CreateFurniture) error {
	newFurniture := models.Furniture{
		NameRu:    data.NameRu,
		NameTm:    data.NameTm,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	save_request := r.db.Model(&models.Furniture{}).Create(&newFurniture)
	return save_request.Error
}

func (r *FurniturePostgres) Update(id int, data models.CreateFurniture) error {
	updateData := map[string]interface{}{
		"name_ru":    data.NameRu,
		"name_tm":    data.NameTm,
		"updated_at": time.Now(),
	}

	update_request := r.db.Model(&models.Furniture{}).Where("id = ?", id).Updates(updateData)
	return update_request.Error
}

func (r *FurniturePostgres) Delete(id int) error {
	var count int64
	ishave_request := r.db.Model(&models.Furniture{}).Where("id = ?", id).Count(&count)
	if ishave_request.Error != nil || count == 0 {
		return errors.New("Tapylmady")
	}

	delete_request := r.db.Model(&models.Furniture{}).Where("id = ?", id).Delete(&models.Furniture{})
	return delete_request.Error
}

func (r *FurniturePostgres) GetFurnitureList(name, language string) ([]models.SelectObject, error) {
	var list []models.Furniture
	result := []models.SelectObject{}

	get_request := r.db.Model(&models.Furniture{})

	if len(name) > 0 {
		get_request = get_request.Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%")
	}

	get_request = get_request.Find(&list)

	if get_request.Error != nil {
		return result, get_request.Error
	}

	for _, element := range list {
		var data models.SelectObject
		data.Id = element.Id
		if language == "ru" {
			data.Name = element.NameRu
		} else {
			data.Name = element.NameTm
		}

		result = append(result, data)
	}

	return result, nil
}

func (r *FurniturePostgres) GetTradePointFurnitures(trade_point_id int, language string) ([]models.TradePointFurnitureRead, error) {
	var db_list []models.TradePointsFurniture
	result := []models.TradePointFurnitureRead{}

	get_request := r.db.Model(&models.TradePointsFurniture{}).Preload("Furniture").
		Where("trade_point_id = ?", trade_point_id).Find(&db_list)

	if get_request.Error != nil {
		return result, get_request.Error
	}

	for _, element := range db_list {
		var data models.TradePointFurnitureRead
		data.Id = element.Id
		data.Count = element.Count
		data.TradePointId = element.TradePointId
		data.FurnitureId = element.FurnitureId
		if language == "ru" {
			data.Name = element.Furniture.NameRu
		} else {
			data.Name = element.Furniture.NameTm
		}

		result = append(result, data)
	}

	return result, nil
}

func (r *FurniturePostgres) CreateTradePointFurniture(trade_point_id, furniture_id, count int) error {
	newFurniture := models.TradePointsFurniture{
		FurnitureId:  furniture_id,
		TradePointId: trade_point_id,
		Count:        count,
		CreatedAt:    time.Now(),
	}

	save_request := r.db.Model(&models.TradePointsFurniture{}).Create(&newFurniture)
	return save_request.Error
}

func (r *FurniturePostgres) UpdateTradePointFurniture(id, furniture_id, count int) error {
	updateData := map[string]interface{}{
		"furniture_id": furniture_id,
		"count":        count,
		"CreatedAt":    time.Now(),
	}

	update_request := r.db.Model(&models.TradePointsFurniture{}).Where("id = ?", id).Updates(updateData)
	return update_request.Error
}

func (r *FurniturePostgres) DeleteTradePointFurniture(id int) error {
	var count int64
	ishave_request := r.db.Model(&models.TradePointsFurniture{}).Where("id = ?", id).Count(&count)
	if ishave_request.Error != nil || count == 0 {
		return errors.New("Tapylmady")
	}

	delete_request := r.db.Model(&models.TradePointsFurniture{}).Where("id = ?", id).Delete(&models.TradePointsFurniture{})
	return delete_request.Error
}
