package repository

import (
	"github.com/Saparly-go-dev/tagma/models"
	"gorm.io/gorm"
	"time"
)

type MerchandiserPostgres struct {
	db *gorm.DB
}

func NewMerchandiserPostgres(db *gorm.DB) *MerchandiserPostgres {
	return &MerchandiserPostgres{db: db}
}

func (r *MerchandiserPostgres) Create(merchandiser models.CreateMerchandiser) (int, error) {
	newMerchandiser := models.Merchandiser{
		NameRu:    merchandiser.NameRu,
		NameTm:    merchandiser.NameTm,
		Status:    merchandiser.Status,
		Number:    merchandiser.Number,
		UpdatedAt: time.Now(),
	}

	result := r.db.Create(&newMerchandiser)
	if err := result.Error; err != nil {
		return 0, err
	}

	return newMerchandiser.Id, nil
}

func (r *MerchandiserPostgres) GetAll(pageNumber, pageSize int, name, language string) (*models.MerchandiserPage, error) {
	var merchandisers []models.Merchandiser
	var pagedata []models.ReadMerchandiser

	offset := (pageNumber - 1) * pageSize

	result := r.db.Model(&models.Merchandiser{}).Preload("TradeAgent").Where("name_ru ILIKE ?", "%"+name+"%").
		Or("name_tm ILIKE ?", "%"+name+"%").Order("id").Offset(offset).Limit(pageSize).Find(&merchandisers)

	if result.Error != nil {
		return nil, result.Error
	}

	var count int64
	result2 := r.db.Model(&models.Merchandiser{}).Where("name_ru ILIKE ?", "%"+name+"%").
		Or("name_tm ILIKE ?", "%"+name+"%").Count(&count)

	if result2.Error != nil {
		return nil, result.Error
	}

	pageCount := count / int64(pageSize)
	if count%int64(pageSize) != 0 {
		pageCount++
	}

	for _, merchandiser := range merchandisers {
		var data models.ReadMerchandiser
		var tradeAgent models.TradeAgent
		if len(merchandiser.TradeAgent) > 0 {
			tradeAgent = merchandiser.TradeAgent[0]
		}

		data.Id = merchandiser.Id
		data.Code = tradeAgent.Code
		data.NameRu = merchandiser.NameRu
		data.NameTm = merchandiser.NameTm
		data.Status = merchandiser.Status
		data.Number = merchandiser.Number
		data.UpdatedAt = merchandiser.UpdatedAt.Format("02.01.2006")
		if language == "ru" {
			data.TradeAgent = tradeAgent.NameRu
		} else {
			data.TradeAgent = tradeAgent.NameTm
		}

		pagedata = append(pagedata, data)
	}

	var page = models.MerchandiserPage{
		Items:      pagedata,
		PageCount:  int(pageCount),
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}

	return &page, nil
}
func (r *MerchandiserPostgres) ChangeStatus(Id int) error {
	var merchandiser models.Merchandiser
	result := r.db.First(&merchandiser, Id)

	if err := result.Error; err != nil {
		return err
	}

	// Toggle status
	newStatus := !merchandiser.Status
	updateData := map[string]interface{}{
		"status":     newStatus,
		"updated_at": time.Now(),
	}

	// Use "=" in the Where clause instead of "=="
	result = r.db.Model(&models.Merchandiser{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *MerchandiserPostgres) GetById(Id int) (models.Merchandiser, error) {
	var merchandiser models.Merchandiser
	result := r.db.First(&merchandiser, Id)
	return merchandiser, result.Error
}

func (r *MerchandiserPostgres) Delete(Id int) error {
	result := r.db.Delete(&models.Merchandiser{}, Id)
	return result.Error
}

func (r *MerchandiserPostgres) Update(Id int, merchandiser models.CreateMerchandiser) error {
	updateData := map[string]interface{}{
		"name_ru": merchandiser.NameRu,
		"name_tm": merchandiser.NameTm,
		"status":  merchandiser.Status,
		"number":  merchandiser.Number,
	}

	result := r.db.Model(&models.Merchandiser{}).Where("id = ?", Id).Updates(updateData)
	return result.Error
}

func (r *MerchandiserPostgres) GetMerchandiserForSelect(language, name string) ([]models.SelectObject, error) {
	result := []models.SelectObject{}
	var merchandiser_list []models.Merchandiser

	get_list_request := r.db.Model(&models.Merchandiser{})

	if len(name) > 0 {
		get_list_request = get_list_request.Where("name_ru ILIKE ? OR name_tm ILIKE ?", "%"+name+"%", "%"+name+"%")
	}

	get_list_request = get_list_request.Find(&merchandiser_list)
	if get_list_request.Error != nil {
		return result, get_list_request.Error
	}

	for _, item := range merchandiser_list {
		var data models.SelectObject

		data.Id = item.Id
		if language == "ru" {
			data.Name = item.NameRu
		} else {
			data.Name = item.NameTm
		}

		result = append(result, data)
	}
	return result, get_list_request.Error
}

func (r *MerchandiserPostgres) GetMarketingPageForMobile(day, month, year, userId int, language string) (*[]models.TradePointMobile, error) {
	var userMerchandiser models.UserMerchandiser

	get_merchandiserUser := r.db.Model(&models.UserMerchandiser{}).Where("user_id = ?", userId).First(&userMerchandiser)
	if get_merchandiserUser.Error != nil {
		return nil, get_merchandiserUser.Error
	}

	var merchandiser models.Merchandiser

	get_merchandiser := r.db.Model(&models.Merchandiser{}).Preload("TradeAgent").
		Where("id = ?", userMerchandiser.MerchandiserId).Find(&merchandiser)

	if get_merchandiser.Error != nil {
		return nil, get_merchandiser.Error
	}

	agentId := 0

	for _, item := range merchandiser.TradeAgent {
		if item.Id > 0 {
			agentId = item.Id
			break
		}
	}

	orderRepo := OrderPostgres{db: r.db}

	result, err := orderRepo.GetTradePointsForMarketing(day, month, year, agentId, language)
	if err != nil {
		return nil, err
	}

	return result, err
}
