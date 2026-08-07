package repository

import (
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"time"
)

type GiftPostgres struct {
	db *gorm.DB
}

func NewGiftPostgres(db *gorm.DB) *GiftPostgres {
	return &GiftPostgres{db: db}
}

func (r *GiftPostgres) Create(gift models.CreateGift) error {
	newGift := models.Gift{
		NameRu:     gift.NameRu,
		NameTm:     gift.NameTm,
		GiftTypeId: gift.GiftTypeId,
		Status:     gift.Status,
		CreatedAt:  time.Now(),
	}

	request := r.db.Create(&newGift)
	return request.Error
}

func (r *GiftPostgres) GetPage(pageSize, pageNumber int, name, language string) (*models.GiftPage, error) {
	var gifts []models.Gift
	var items []models.ReadGift

	offset := (pageNumber - 1) * pageSize

	request := r.db.Model(&models.Gift{}).Preload("GiftType").
		Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%").
		Limit(pageSize).Offset(offset).Find(&gifts)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, gift := range gifts {
		var data models.ReadGift
		data.Id = gift.Id
		data.NameRu = gift.NameRu
		data.NameTm = gift.NameTm
		data.GiftTypeId = gift.GiftTypeId
		data.Status = gift.Status
		data.UpdatedAt = gift.CreatedAt.Format("02.01.2006")
		if language == "ru" {
			data.GiftType = gift.GiftType.NameRu
		} else {
			data.GiftType = gift.GiftType.NameTm
		}

		items = append(items, data)
	}

	var count int64

	request2 := r.db.Model(&models.Gift{}).
		Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%").Count(&count)

	if request2.Error != nil {
		return nil, request2.Error
	}

	// Calculate total page count
	pageCount := int(count) / pageSize
	if int(count)%pageSize != 0 {
		pageCount++
	}

	if items == nil {
		items = []models.ReadGift{}
	}

	// Create the page object
	page := models.GiftPage{
		Items:      items,
		PageCount:  pageCount,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	return &page, nil
}

func (r *GiftPostgres) ChangeStatus(Id int) error {
	var data models.Gift
	result := r.db.First(&data, Id)
	if result.Error != nil {
		return result.Error
	}

	if data.Status == true {
		updateData := map[string]interface{}{
			"status":     false,
			"created_at": time.Now(),
		}
		result := r.db.Model(&models.Gift{}).Where("id = ?", Id).Updates(updateData)
		return result.Error

	} else {
		updateData := map[string]interface{}{
			"status":     true,
			"created_at": time.Now(),
		}
		result := r.db.Model(&models.Gift{}).Where("id = ?", Id).Updates(updateData)
		return result.Error
	}

}

func (r *GiftPostgres) GetById(Id int) (*models.Gift, error) {
	var gift models.Gift

	request := r.db.Model(&models.Gift{}).Preload("GiftType").Preload("TradePoint").First(&gift, Id)

	if request.Error != nil {
		return nil, request.Error
	}

	return &gift, nil
}

func (r *GiftPostgres) Delete(Id int) error {
	request := r.db.Delete(&models.Gift{}, Id)
	return request.Error
}

func (r *GiftPostgres) Update(Id int, gift models.CreateGift) error {
	updateData := map[string]interface{}{
		"name_ru":      gift.NameRu,
		"name_tm":      gift.NameTm,
		"status":       gift.Status,
		"gift_type_id": gift.GiftTypeId,
	}

	request := r.db.Model(&models.Gift{}).Where("id = ?", Id).Updates(updateData)
	return request.Error
}

func (r *GiftPostgres) GetGiftTypes(language string) (*[]models.SelectObject, error) {
	var list []models.GiftType
	var result []models.SelectObject
	//fmt.Println("geldi")

	request := r.db.Model(&models.GiftType{}).Where("id > ?", 0).Find(&list)

	if request.Error != nil {
		return nil, request.Error
	}

	for _, item := range list {
		var data models.SelectObject
		data.Id = item.Id
		if language == "ru" {
			data.Name = item.NameRu
		} else {
			data.Name = item.NameTm
		}
		result = append(result, data)
	}

	return &result, nil
}
