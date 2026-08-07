package repository

import (
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"time"
)

type DistrictPostgres struct {
	db *gorm.DB
}

func NewDistrictPostgres(db *gorm.DB) *DistrictPostgres {
	return &DistrictPostgres{db: db}
}

func (r *DistrictPostgres) Create(district models.CreateDistrict) error {
	newDistrict := models.District{
		NameRu:    district.NameRu,
		NameTm:    district.NameTm,
		UpdatedAt: time.Now(),
		CityId:    district.CityId,
		Status:    true,
	}

	request := r.db.Model(&models.District{}).Create(&newDistrict)

	return request.Error
}

func (r *DistrictPostgres) GetPage(pageSize, pageNumber int, name, language string) (*models.DistrictPage, error) {
	var districts []models.District
	var items []models.ReadDistrict

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.District{}).Preload("City").
		Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%").
		Order("id").Offset(offset).Limit(pageSize).Find(&districts)

	if request.Error != nil {
		return &models.DistrictPage{}, request.Error
	}

	for _, item := range districts {
		var data models.ReadDistrict

		data.Id = item.Id
		data.NameRu = item.NameRu
		data.NameTm = item.NameTm
		data.Status = item.Status
		data.UpdatedAt = item.UpdatedAt.Format("02.01.2006")
		data.CityId = item.CityId
		if language == "ru" {
			data.City = item.City.NameRu
		} else {
			data.City = item.City.NameTm
		}

		items = append(items, data)
	}

	var count int64
	request2 := r.db.Model(&models.District{}).
		Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%").Count(&count)

	if request2.Error != nil {
		return &models.DistrictPage{}, request2.Error
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if items == nil {
		items = []models.ReadDistrict{}
	}

	// Create the page object
	page := models.DistrictPage{
		Items:      items,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	return &page, nil
}

func (r *DistrictPostgres) GetById(Id int) (*models.District, error) {
	var result models.District

	request := r.db.Model(&models.District{}).Preload("City").Where("id = ?", Id).First(&result)

	if request.Error != nil {
		return &models.District{}, request.Error
	}

	return &result, request.Error
}

func (r *DistrictPostgres) ChangeStatus(Id int) error {
	var data models.District
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status": false,
		}
		result := r.db.Model(&models.District{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status": true,
		}
		result := r.db.Model(&models.District{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *DistrictPostgres) Delete(Id int) error {
	request := r.db.Model(&models.District{}).Delete(&models.District{}, Id)
	return request.Error
}

func (r *DistrictPostgres) Update(Id int, district models.CreateDistrict) error {
	updateData := map[string]interface{}{
		"name_ru": district.NameRu,
		"name_tm": district.NameTm,
		"city_id": district.CityId,
	}

	request := r.db.Model(&models.District{}).Where("id = ?", Id).Updates(updateData)
	return request.Error
}
